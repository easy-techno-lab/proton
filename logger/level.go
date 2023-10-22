package logger

import (
	"fmt"
	"strings"
)

type Level uint32

const (
	LevelUnknown Level = iota
	LevelFatal
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
)

func (lvl Level) IsValid() bool {
	switch lvl {
	case LevelFatal, LevelError, LevelWarn, LevelInfo, LevelDebug, LevelTrace:
		return true
	}
	return false
}

func (lvl Level) String() string {
	switch lvl {
	case LevelFatal:
		return "FATAL"
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	case LevelTrace:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

// Printf returns a logger printf function for current logger level.
func (lvl Level) Printf(format string, v ...any) {
	switch lvl {
	case LevelFatal:
		std.fatalf(format, v...)
	case LevelError:
		std.errorf(format, v...)
	case LevelWarn:
		std.warnf(format, v...)
	case LevelInfo:
		std.infof(format, v...)
	case LevelDebug:
		std.debugf(format, v...)
	default:
		std.tracef(format, v...)
	}
}

// Print returns a logger print function for current logger level.
func (lvl Level) Print(v ...any) {
	switch lvl {
	case LevelFatal:
		std.fatal(v...)
	case LevelError:
		std.error(v...)
	case LevelWarn:
		std.warn(v...)
	case LevelInfo:
		std.info(v...)
	case LevelDebug:
		std.debug(v...)
	default:
		std.trace(v...)
	}
}

// ParseLevel takes a string level and returns the logger Level constant.
func ParseLevel(lvl string) (Level, error) {
	switch strings.ToLower(lvl) {
	case "fatal":
		return LevelFatal, nil
	case "error":
		return LevelError, nil
	case "warn", "warning":
		return LevelWarn, nil
	case "info":
		return LevelInfo, nil
	case "debug":
		return LevelDebug, nil
	case "trace":
		return LevelTrace, nil
	default:
		return LevelUnknown, fmt.Errorf("not a valid logger Level: %s", lvl)
	}
}
