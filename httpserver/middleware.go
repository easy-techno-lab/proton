package httpserver

import (
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"time"

	"github.com/easy-techno-lab/proton/logger"
)

// MiddlewareSequencer chains middleware functions in a chain.
func MiddlewareSequencer(baseHandler http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for _, f := range mws {
		baseHandler = f(baseHandler)
	}
	return baseHandler
}

// Timer measures the time taken by http.HandlerFunc.
func Timer(logLevel logger.Level) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if logger.InLevel(logLevel) {
				defer func(start time.Time) {
					logLevel.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
				}(time.Now())
			}
			next.ServeHTTP(w, r)
		})
	}
}

// PanicCatcher handles panics in http.HandlerFunc.
func PanicCatcher(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				if logger.InLevel(logger.LevelError) {
					logger.Errorf("%s\n%s", rec, debug.Stack())
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// DumpHttp dumps the HTTP request and response, and prints out with logFunc.
func DumpHttp(logLevel logger.Level) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if logger.InLevel(logLevel) {
				logger.DumpHttpRequest(r, logLevel)

				recorder := httptest.NewRecorder()

				next.ServeHTTP(recorder, r)

				for key, values := range recorder.Header() {
					w.Header().Del(key)
					for _, value := range values {
						w.Header().Set(key, value)
					}
				}

				w.WriteHeader(recorder.Code)

				response := recorder.Result()
				response.ContentLength, _ = recorder.Body.WriteTo(w)

				logger.DumpHttpResponse(response, logLevel)

				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
