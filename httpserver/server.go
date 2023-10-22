package httpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/easy-techno-lab/proton/coder"
	"github.com/easy-techno-lab/proton/logger"
)

type Formatter interface {
	coder.Coder
	WriteResponse(w http.ResponseWriter, statusCode int, v any)
}

// NewFormatter returns a new Formatter.
func NewFormatter(coder coder.Coder) Formatter {
	return &protoFormatter{Coder: coder}
}

type protoFormatter struct {
	coder.Coder
}

// WriteResponse encodes the value pointed to by v and writes it and statusCode to the stream.
func (f *protoFormatter) WriteResponse(w http.ResponseWriter, statusCode int, v any) {
	if v != nil {
		if w.Header().Get(coder.ContentType) == "" && f.ContentType() != "" {
			w.Header().Set(coder.ContentType, f.ContentType())
		}
		w.WriteHeader(statusCode)
		if err := f.Encode(w, v); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			logger.Error("Can't encode response: ", err)
		}
		return
	}
	w.WriteHeader(statusCode)
}

// Controller is a wrapper around *http.Server to control the server.
//
//	Server — *http.Server, which will be managed.
//	GracefulTimeout — time that is given to the server to shut down gracefully.
type Controller struct {
	Server          *http.Server
	GracefulTimeout time.Duration

	isRan   atomic.Bool
	restart atomic.Bool
}

// Start starts the *http.Server.
// If *tls.Config on the server is non nil, the server listens and serves using tls.
func (c *Controller) Start() error {
	for {
		if err := c.start(); !c.restart.Load() {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		} else if err != nil {
			logger.Error(err)
		}

		logger.Info("Server is restarting")

		c.clone()
		c.restart.Store(false)
	}
}

// Restart restarts the server if necessary.
// For changes to the following parameters to take effect:
//
//	Addr; TLSConfig; TLSNextProto; ConnState; BaseContext; ConnContext,
//
// a server restart is required.
// Other parameters can be changed without restarting the server.
// If the server is not running, the function will be skipped.
func (c *Controller) Restart() {
	if !c.isRan.Load() {
		return
	}

	c.restart.Store(true)

	if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
		logger.Error("Send stop signal: ", err)
	}
}

// Shutdown gracefully shuts down the server.
func (c *Controller) Shutdown() {
	ctx, cancelWithTimeout := context.WithTimeout(context.Background(), c.GracefulTimeout)
	defer cancelWithTimeout()

	if err := c.Server.Shutdown(ctx); err != nil {
		logger.Error("Shutdown server: ", err)
	}
}

func (c *Controller) start() error {
	idle := make(chan struct{})

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(sigint)
		close(sigint)
	}()

	go func() {
		defer logger.Info("Server is shutdown")

		<-sigint

		c.Shutdown()
		close(idle)
	}()

	c.isRan.Store(true)
	defer c.isRan.Store(false)

	if c.Server.TLSConfig != nil {
		logger.Info("HTTPS server listening on ", c.Server.Addr)

		err := c.Server.ListenAndServeTLS("", "")
		return fmt.Errorf("HTTPS server ListenAndServeTLS: %w", err)
	} else {
		logger.Info("HTTP server listening on ", c.Server.Addr)

		err := c.Server.ListenAndServe()
		return fmt.Errorf("HTTP server ListenAndServe: %w", err)
	}
}

// clone clones the server before restarting, since it is impossible to start a stopped server.
func (c *Controller) clone() {
	var tlsConfig *tls.Config

	if c.Server.TLSConfig != nil && len(c.Server.TLSConfig.Certificates) != 0 {
		tlsConfig = c.Server.TLSConfig.Clone()
	}

	c.Server = &http.Server{
		Addr:                         c.Server.Addr, // need to restart
		Handler:                      c.Server.Handler,
		DisableGeneralOptionsHandler: c.Server.DisableGeneralOptionsHandler,
		TLSConfig:                    tlsConfig, // need to restart
		ReadTimeout:                  c.Server.ReadTimeout,
		ReadHeaderTimeout:            c.Server.ReadHeaderTimeout,
		WriteTimeout:                 c.Server.WriteTimeout,
		IdleTimeout:                  c.Server.IdleTimeout,
		MaxHeaderBytes:               c.Server.MaxHeaderBytes,
		TLSNextProto:                 c.Server.TLSNextProto, // need to restart
		ConnState:                    c.Server.ConnState,    // need to restart
		ErrorLog:                     c.Server.ErrorLog,
		BaseContext:                  c.Server.BaseContext, // need to restart
		ConnContext:                  c.Server.ConnContext, // need to restart
	}
}
