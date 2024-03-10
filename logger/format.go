package logger

import (
	"fmt"
	"strings"
)

type Format uint8

const (
	FormatUnknown Format = iota
	FormatText
	FormatJSON
)

func (f Format) IsValid() bool {
	switch f {
	case FormatText, FormatJSON:
		return true
	default:
		return false
	}
}

func (f Format) String() string {
	switch f {
	case FormatText:
		return "TEXT"
	case FormatJSON:
		return "JSON"
	default:
		return "UNKNOWN"
	}
}

// ParseFormat takes a string format and returns the logger Format constant.
func ParseFormat(f string) (Format, error) {
	switch strings.ToLower(f) {
	case "text":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	default:
		return FormatUnknown, fmt.Errorf("not a valid logger Format: %s", f)
	}
}
