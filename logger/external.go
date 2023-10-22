package logger

import (
	"io"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strings"
	"sync/atomic"
)

// SetTimeFormat sets the logger time format.
func SetTimeFormat(format string) {
	std.setTimeFormat(format)
}

// SetAdditionalOut sets an additional logger output.
func SetAdditionalOut(out io.Writer) {
	std.setAdditionalOut(out)
}

// SetFormat sets the logger format.
func SetFormat(f Format) {
	std.setFormat(f)
}

// SetLevel sets the logger level.
func SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&std.level), uint32(level))
}

// SetFuncNamePrinting sets whether the logger should print the caller function name.
func SetFuncNamePrinting(on bool) {
	std.setFuncNamePrinting(on)
}

// InLevel returns true if the given level is less than or equal to the current logger level.
func InLevel(level Level) bool {
	return std.inLevel(level)
}

func Fatalf(format string, v ...any) {
	std.fatalf(format, v...)
}

func Fatal(v ...any) {
	std.fatal(v...)
}

func Errorf(format string, v ...any) {
	std.errorf(format, v...)
}

func Error(v ...any) {
	std.error(v...)
}

func Warnf(format string, v ...any) {
	std.warnf(format, v...)
}

func Warn(v ...any) {
	std.warn(v...)
}

func Infof(format string, v ...any) {
	std.infof(format, v...)
}

func Info(v ...any) {
	std.info(v...)
}

func Debugf(format string, v ...any) {
	std.debugf(format, v...)
}

func Debug(v ...any) {
	std.debug(v...)
}

func Tracef(format string, v ...any) {
	std.tracef(format, v...)
}

func Trace(v ...any) {
	std.trace(v...)
}

// DumpHttpRequest dumps the HTTP request and prints out with logFunc.
func DumpHttpRequest(r *http.Request, logLevel Level) {
	dumpFunc := httputil.DumpRequestOut
	if r.URL.Scheme == "" || r.URL.Host == "" {
		dumpFunc = httputil.DumpRequest
	}
	b, err := dumpFunc(r, true)
	if err != nil {
		std.error("REQUEST LOG error: ", err)
		return
	}
	logLevel.Print("REQUEST:\n", string(b))
}

// DumpHttpResponse dumps the HTTP response and prints out with logFunc.
func DumpHttpResponse(r *http.Response, logLevel Level) {
	b, err := httputil.DumpResponse(r, true)
	if err != nil {
		std.error("RESPONSE LOG error: ", err)
		return
	}
	logLevel.Print("RESPONSE:\n", string(b))
}

// Closer calls the Close method, if the closure occurred with an error, it prints it to the log.
func Closer(c io.Closer) {
	if err := c.Close(); err != nil {
		std.error(err)
	}
}

// FunctionInfo returns the name of the function and file, the line number on the calling goroutine's stack.
// The argument skip is the number of stack frames to ascend.
func FunctionInfo(skip int) (name, file string, line int, ok bool) {
	var pc uintptr
	if pc, file, line, ok = runtime.Caller(skip); !ok {
		return
	}
	name = runtime.FuncForPC(pc).Name()
	// Truncate the path to the package.
	name = name[strings.LastIndex(name, "/")+1:]
	// Truncate the package name.
	name = name[strings.Index(name, ".")+1:]
	// If a function is a method truncate the type name.
	if name[0] == '(' {
		name = name[strings.Index(name, ".")+1:]
	}
	// If the function runs anonymous functions, truncate the name of the anonymous function.
	if i := strings.Index(name, "."); i > 0 {
		name = name[:i]
	}
	return
}
