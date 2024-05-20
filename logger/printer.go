package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Proper usage of a sync.Pool requires each entry to have approximately
// the same memory cost. To obtain this property when the stored type
// contains a variably-sized buffer, we add a hard limit on the maximum buffer
// to place back in the pool.
//
// See https://golang.org/issue/23199
const maxSize = 1 << 16 // 64KiB

var printerPool = &sync.Pool{
	New: func() any {
		return &printer{
			buf: make([]byte, 0, 512),
		}
	},
}

func putPrinter(p *printer) {
	if cap(p.buf) > maxSize {
		return
	}
	printerPool.Put(p)
}

func (l *logger) newPrinter(lvl Level) *printer {
	p := printerPool.Get().(*printer)

	p.buf = enc.color(p.buf[:0], lvl)
	p.level = lvl

	l.mu.RLock()
	p.timeFormat = l.timeFormat
	p.out = l.out
	p.format = l.format
	p.funcName = l.funcName
	l.mu.RUnlock()

	return p
}

type printer struct {
	timeFormat string
	out        io.Writer
	buf        []byte
	level      Level
	format     Format
	funcName   bool
}

func (p *printer) printf(format string, v ...any) {
	p.mid(fmt.Sprintf(format, v...))
}

func (p *printer) print(v ...any) {
	p.mid(fmt.Sprint(v...))
}

func (p *printer) mid(msg string) {
	switch p.format {
	case FormatJson:
		p.json(msg)
	default:
		p.text(msg)
	}

	if p.out != nil {
		_, _ = p.out.Write(enc.lineFeed(p.buf[5:]))
	}

	p.buf = enc.lineFeed(enc.value(p.buf, "\x1b[0m"))
	_, _ = os.Stderr.Write(p.buf)

	putPrinter(p)
}

func (p *printer) text(msg string) {
	if p.timeFormat != "" {
		p.buf = enc.space(enc.time(p.buf, p.timeFormat))
	}

	p.buf = enc.space(enc.value(p.buf, p.level.String()))
	p.buf = enc.space(enc.value(enc.routine(enc.value(p.buf, "[")), "]"))

	if p.funcName || std.inLevel(LevelTrace) {
		name, file, line, _ := FunctionInfo(6)

		if std.inLevel(LevelTrace) {
			p.buf = enc.space(enc.caller(p.buf, file, line))
		}

		p.buf = enc.space(enc.value(enc.value(p.buf, name), "()"))
	}

	p.buf = enc.value(p.buf, msg)
}

func (p *printer) json(msg string) {
	p.buf = enc.startMarker(p.buf)

	if p.timeFormat != "" {
		p.buf = enc.quotedTime(enc.key(p.buf, nt), p.timeFormat)
	}

	p.buf = enc.string(enc.key(p.buf, nl), p.level.String())
	p.buf = enc.routine(enc.key(p.buf, nr))

	if p.funcName || std.inLevel(LevelTrace) {
		name, file, line, _ := FunctionInfo(6)

		if std.inLevel(LevelTrace) {
			p.buf = enc.value(enc.key(p.buf, nc), "\"")
			p.buf = enc.value(enc.caller(p.buf, file, line), "\"")
		}

		p.buf = enc.string(enc.key(p.buf, nf), name)
	}

	p.buf = enc.string(enc.key(p.buf, nm), msg)
	p.buf = enc.endMarker(p.buf)
}

const (
	nt = "time"
	nl = "level"
	nr = "routine"
	nc = "caller"
	nf = "func"
	nm = "message"
)

var enc = encoder{buf: make([]byte, 32)}

type encoder struct {
	buf []byte
}

func (e encoder) color(dst []byte, lvl Level) []byte {
	switch lvl {
	case LevelFatal:
		return enc.value(dst, "\x1b[95m")
	case LevelError:
		return enc.value(dst, "\x1b[91m")
	case LevelWarn:
		return enc.value(dst, "\x1b[93m")
	case LevelInfo:
		return enc.value(dst, "\x1b[92m")
	case LevelDebug:
		return enc.value(dst, "\x1b[94m")
	default:
		return enc.value(dst, "\x1b[96m")
	}
}

func (e encoder) routine(dst []byte) []byte {
	runtime.Stack(e.buf, false)
	b := e.buf[10:]
	return append(dst, b[:bytes.IndexByte(b, ' ')]...)
}

func (e encoder) caller(dst []byte, file string, line int) []byte {
	return append(append(append(dst, file...), ':'), strconv.Itoa(line)...)
}

func (e encoder) key(dst []byte, key string) []byte {
	if dst[len(dst)-1] != '{' {
		dst = append(dst, ',')
	}
	return append(e.string(dst, key), ':')
}

func (encoder) value(dst []byte, val string) []byte {
	return append(dst, val...)
}

func (encoder) string(dst []byte, str string) []byte {
	return append(append(append(dst, '"'), str...), '"')
}

func (encoder) time(dst []byte, fmt string) []byte {
	return time.Now().AppendFormat(dst, fmt)
}

func (encoder) quotedTime(dst []byte, fmt string) []byte {
	return append(time.Now().AppendFormat(append(dst, '"'), fmt), '"')
}

func (encoder) startMarker(dst []byte) []byte {
	return append(dst, '{')
}

func (encoder) space(dst []byte) []byte {
	return append(dst, ' ')
}

func (encoder) endMarker(dst []byte) []byte {
	return append(dst, '}')
}

func (encoder) lineFeed(dst []byte) []byte {
	return append(dst, '\n')
}
