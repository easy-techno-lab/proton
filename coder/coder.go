package coder

import (
	"io"

	"github.com/easy-techno-lab/proton/logger"
)

const ContentType = "Content-Type"

// An Encoder encodes and writes values to an output stream.
type Encoder interface {
	Encode(w io.Writer, v any) error
}

type encoder struct {
	f func(v any) ([]byte, error)
}

// NewEncoder returns a new Encoder that writes to w.
func NewEncoder(marshal func(v any) ([]byte, error)) Encoder {
	return &encoder{f: marshal}
}

// Encode encodes the value pointed to by v and writes it to the stream.
// It will panic if encoder function not set.
func (e *encoder) Encode(w io.Writer, v any) error {
	logger.Debugf("Encoder input: %#v", v)

	p, err := e.f(v)
	if err != nil {
		return err
	}

	logger.Debugf("Encoder output: %s", p)

	logger.Tracef("Encoder output: % x, len: %d", p, len(p))

	if _, err = w.Write(p); err != nil {
		return err
	}

	return nil
}

// A Decoder reads and decodes values from an input stream.
type Decoder interface {
	Decode(r io.Reader, v any) error
}

type decoder struct {
	f func(data []byte, v any) error
}

// NewDecoder returns a new Decoder that reads from r.
func NewDecoder(unmarshal func(data []byte, v any) error) Decoder {
	return &decoder{f: unmarshal}
}

// Decode reads the next encoded value from its input and stores it in the value pointed to by v.
// It will panic if decoder function not set.
func (d *decoder) Decode(r io.Reader, v any) error {
	p, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	logger.Tracef("Decoder input: % x, len: %d", p, len(p))

	logger.Debugf("Decoder input: %s", p)

	if err = d.f(p, v); err != nil {
		return err
	}

	logger.Debugf("Decoder output: %#v", v)

	return nil
}

// A Coder is a pair of Encoder and Decoder.
type Coder interface {
	ContentType() string
	Encoder
	Decoder
}

type coder struct {
	t string
	Encoder
	Decoder
}

// NewCoder returns a new Coder.
func NewCoder(contentType string, marshal func(v any) ([]byte, error), unmarshal func(data []byte, v any) error) Coder {
	return &coder{t: contentType, Encoder: NewEncoder(marshal), Decoder: NewDecoder(unmarshal)}
}

// ContentType returns a string value representing the Coder type.
// Use as the ContentType header of HTTP requests.
func (c coder) ContentType() string {
	return c.t
}
