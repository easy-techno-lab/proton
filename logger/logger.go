package logger

import (
	"io"
	"os"
	"sync"
	"sync/atomic"
)

const timeFormatDefault = "2006/01/02 15:04:05.000"

var std = &logger{
	timeFormat: timeFormatDefault,
	format:     FormatText,
	level:      LevelTrace,
	funcName:   true,
}

type logger struct {
	mu         sync.RWMutex
	timeFormat string
	out        io.Writer
	format     Format
	level      Level
	funcName   bool
}

func (l *logger) setTimeFormat(format string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.timeFormat = format
}

func (l *logger) setAdditionalOut(out io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = out
}

func (l *logger) setFormat(f Format) {
	l.mu.Lock()
	defer l.mu.Unlock()
	switch f {
	case FormatJson:
		l.format = FormatJson
	default:
		l.format = FormatText
	}
}

func (l *logger) setLevel(level Level) {
	atomic.StoreUint32((*uint32)(&l.level), uint32(level))
}

func (l *logger) setFuncNamePrinting(on bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.funcName = on
}

func (l *logger) inLevel(lvl Level) bool {
	return Level(atomic.LoadUint32((*uint32)(&l.level))) >= lvl
}

func (l *logger) fatalf(format string, v ...any) {
	l.newPrinter(LevelFatal).printf(format, v...)
	os.Exit(1)
}

func (l *logger) fatal(v ...any) {
	l.newPrinter(LevelFatal).print(v...)
	os.Exit(1)
}

func (l *logger) errorf(format string, v ...any) {
	if l.inLevel(LevelError) {
		l.newPrinter(LevelError).printf(format, v...)
	}
}

func (l *logger) error(v ...any) {
	if l.inLevel(LevelError) {
		l.newPrinter(LevelError).print(v...)
	}
}

func (l *logger) warnf(format string, v ...any) {
	if l.inLevel(LevelWarn) {
		l.newPrinter(LevelWarn).printf(format, v...)
	}
}

func (l *logger) warn(v ...any) {
	if l.inLevel(LevelWarn) {
		l.newPrinter(LevelWarn).print(v...)
	}
}

func (l *logger) infof(format string, v ...any) {
	if l.inLevel(LevelInfo) {
		l.newPrinter(LevelInfo).printf(format, v...)
	}
}

func (l *logger) info(v ...any) {
	if l.inLevel(LevelInfo) {
		l.newPrinter(LevelInfo).print(v...)
	}
}

func (l *logger) debugf(format string, v ...any) {
	if l.inLevel(LevelDebug) {
		l.newPrinter(LevelDebug).printf(format, v...)
	}
}

func (l *logger) debug(v ...any) {
	if l.inLevel(LevelDebug) {
		l.newPrinter(LevelDebug).print(v...)
	}
}

func (l *logger) tracef(format string, v ...any) {
	if l.inLevel(LevelTrace) {
		l.newPrinter(LevelTrace).printf(format, v...)
	}
}

func (l *logger) trace(v ...any) {
	if l.inLevel(LevelTrace) {
		l.newPrinter(LevelTrace).print(v...)
	}
}
