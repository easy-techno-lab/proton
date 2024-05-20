package logger

import (
	"fmt"
	"strings"
)

type Format uint8

const (
	FormatUnknown Format = iota
	FormatText
	FormatJson
)

func (f Format) IsValid() bool {
	switch f {
	case FormatText, FormatJson:
		return true
	default:
		return false
	}
}

func (f Format) String() string {
	switch f {
	case FormatText:
		return "TEXT"
	case FormatJson:
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
		return FormatJson, nil
	default:
		return FormatUnknown, fmt.Errorf("not a valid logger Format: %s", f)
	}
}
