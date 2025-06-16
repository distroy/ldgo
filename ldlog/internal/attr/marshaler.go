/*
 * Copyright (C) distroy
 */

package attr

import (
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/distroy/ldgo/v3/ldconv"
	"github.com/distroy/ldgo/v3/ldlog/internal/buffer"
)

type (
	WriterTo      = io.WriterTo
	TextMarshaler = encoding.TextMarshaler
	JsonMarshaler = json.Marshaler
)

type Marshaler interface {
	WriterTo
	TextMarshaler
	JsonMarshaler
}

var (
	_ Marshaler = nil_t{}
	_ Marshaler = complex_t(0)
	_ Marshaler = (*slice_t[int])(nil)
)

func getBuf() *buffer.Buffer { return buffer.NewBuffer() }

func b2s(b []byte) string   { return ldconv.BytesToStrUnsafe(b) }
func s2b(b string) []byte   { return ldconv.StrToBytesUnsafe(b) }
func b64(b []byte) string   { return base64.StdEncoding.EncodeToString(b) }
func quote(s string) string { return strconv.Quote(s) }
func e2s(e error) string {
	if e == nil {
		return nil_t{}.String()
	}
	return e.Error()
}
func writeStringer(w io.Writer, s fmt.Stringer) (int64, error) {
	n, err := w.Write(s2b(s.String()))
	return int64(n), err
}

type nil_t struct{}

func (p nil_t) String() string                     { return "null" }
func (p nil_t) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p nil_t) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p nil_t) WriteTo(w io.Writer) (int64, error) { return writeStringer(w, p) }

type complex_t complex128

func (n complex_t) MarshalJSON() ([]byte, error)       { return s2b(n.String()), nil }
func (n complex_t) MarshalText() ([]byte, error)       { return s2b(n.String()), nil }
func (p complex_t) WriteTo(w io.Writer) (int64, error) { return writeStringer(w, p) }
func (n complex_t) String() string                     { return strconv.FormatComplex(complex128(n), 'f', -1, 128) }

type slice_t[T any] struct {
	data []T
	text func(buf []byte, v T) []byte
}

func (p *slice_t[T]) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p *slice_t[T]) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p *slice_t[T]) WriteTo(w io.Writer) (int64, error) { return writeStringer(w, p) }
func (p *slice_t[T]) String() string {
	buf := getBuf()
	defer buf.Free()
	if p.data == nil {
		buf.WriteString(nil_t{}.String())
		return buf.String()
	}
	buf.WriteByte('[')
	for i, v := range p.data {
		if i > 0 {
			buf.WriteByte(',')
		}
		*buf = p.text(*buf, v)
	}
	buf.WriteByte(']')
	return buf.String()
}
