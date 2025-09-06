/*
 * Copyright (C) distroy
 */

package _attr

import (
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/distroy/ldgo/v3/ldconv"
	"github.com/distroy/ldgo/v3/ldlog/internal/_buffer"
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
	_ Marshaler = complex64_t(0)
	_ Marshaler = complex128_t(0)
	_ Marshaler = (*slice_t[int])(nil)
)

const TimeLayout = _buffer.TimeLayout

type Buffer = _buffer.Buffer

func getBuf() *Buffer { return _buffer.NewBuffer() }

func b2s(b []byte) string { return ldconv.BytesToStrUnsafe(b) }
func s2b(b string) []byte { return ldconv.StrToBytesUnsafe(b) }
func b64(b []byte) string { return base64.StdEncoding.EncodeToString(b) }

type BufferWriter interface {
	fmt.Stringer
	WriteToBuffer(buf *Buffer)
}

func writeTo(w io.Writer, s BufferWriter) (int64, error) {
	if buf, _ := w.(*Buffer); buf != nil {
		l0 := buf.Len()
		s.WriteToBuffer(buf)
		return int64(buf.Len() - l0), nil
	}
	n, err := w.Write(s2b(s.String()))
	return int64(n), err
}

type nil_t struct{}

func (p nil_t) String() string                     { return "null" }
func (p nil_t) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p nil_t) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p nil_t) WriteTo(w io.Writer) (int64, error) { return writeTo(w, p) }
func (p nil_t) WriteToBuffer(w *Buffer)            { w.WriteString(p.String()) }

type complex64_t complex64

func (n complex64_t) MarshalJSON() ([]byte, error)       { return s2b(n.String()), nil }
func (n complex64_t) MarshalText() ([]byte, error)       { return s2b(n.String()), nil }
func (p complex64_t) WriteTo(w io.Writer) (int64, error) { return writeTo(w, p) }
func (n complex64_t) WriteToBuffer(b *Buffer)            { complex128_t(n).WriteToBuffer(b) }
func (n complex64_t) String() string                     { return complex128_t(n).String() }

type complex128_t complex128

func (n complex128_t) MarshalJSON() ([]byte, error)       { return s2b(n.String()), nil }
func (n complex128_t) MarshalText() ([]byte, error)       { return s2b(n.String()), nil }
func (p complex128_t) WriteTo(w io.Writer) (int64, error) { return writeTo(w, p) }
func (n complex128_t) WriteToBuffer(b *Buffer) {
	s := strconv.FormatComplex(complex128(n), 'f', -1, 128)
	l := len(s) - 1
	if s[0] == '(' && s[l] == ')' {
		s = s[1 : len(s)-1]
	}
	b.AppendQuote(s)
}
func (n complex128_t) String() string {
	b := getBuf()
	defer b.Free()
	n.WriteToBuffer(b)
	return b.String()
}

type slice_t[T any] struct {
	data []T
	text func(buf *Buffer, v T)
}

func (p *slice_t[T]) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p *slice_t[T]) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p *slice_t[T]) WriteTo(w io.Writer) (int64, error) { return writeTo(w, p) }
func (p *slice_t[T]) String() string {
	buf := getBuf()
	defer buf.Free()
	p.WriteToBuffer(buf)
	return buf.String()
}
func (p *slice_t[T]) WriteToBuffer(buf *Buffer) {
	if p.data == nil {
		buf.WriteString(nil_t{}.String())
		return
	}
	buf.WriteByte('[')
	for i, v := range p.data {
		if i > 0 {
			buf.WriteByte(',')
		}
		p.text(buf, v)
	}
	buf.WriteByte(']')
	return
}
