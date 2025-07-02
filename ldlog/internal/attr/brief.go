/*
 * Copyright (C) distroy
 */

package attr

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/distroy/ldgo/v3/ldlog/internal/buffer"
)

const (
	tagLen   = "<len>"
	tagType  = "<type>"
	tagBrief = "<brief>"

	minBriefStringLen = 10
	minBriefArrayLen  = 1
	minBriefMapLen    = 3
)

var (
	briefStringLen = 100
	briefArrayLen  = 1
	briefMapLen    = 10
)

func SetBriefStringLen(n int) { briefStringLen = max(n, minBriefStringLen) }
func SetBriefArrayLen(n int)  { briefArrayLen = max(n, minBriefArrayLen) }
func SetBriefMapLen(n int)    { briefMapLen = max(n, minBriefMapLen) }

var (
	_ Marshaler = (*brief_stringer_t)(nil)
	_ Marshaler = (*brief_slice_t[int])(nil)
)

func writeBrief(b *buffer.Buffer, l int, typ string, f func(b *buffer.Buffer)) {
	writeKey := func(b *buffer.Buffer, key string) {
		b.AppendByte('"')
		b.AppendString(key)
		b.AppendByte('"')
	}

	b.AppendByte('{')

	writeKey(b, tagLen)
	b.AppendByte(':')
	b.AppendInt(int64(l))
	b.AppendByte(',')

	writeKey(b, tagType)
	b.AppendByte(':')
	writeKey(b, typ)
	b.AppendByte(',')

	writeKey(b, tagBrief)
	b.AppendByte(':')
	// addQuote(b, s[:n])
	f(b)

	b.AppendByte('}')
}

func addBriefStr(b *buffer.Buffer, s string) {
	n := briefStringLen
	l := len(s)
	if l <= n {
		b.AppendQuote(s)
		return
	}
	writeBrief(b, l, "string", func(b *buffer.Buffer) {
		b.AppendQuote(s[:n])
	})
}

func addBriefSlice[T any](b *buffer.Buffer, s *brief_slice_t[T]) {
	n := briefArrayLen
	l := len(s.data)
	if l <= n {
		x := &slice_t[T]{
			data: s.data,
			text: s.text,
		}
		x.WriteTo(b)
		return
	}

	writeBrief(b, l, "array", func(b *buffer.Buffer) {
		x := &slice_t[T]{
			data: s.data[:n],
			text: s.text,
		}
		x.WriteTo(b)
	})
}

type string_t string

func (p string_t) String() string { return string(p) }

type brief_stringer_t struct {
	val fmt.Stringer
}

func (p brief_stringer_t) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p brief_stringer_t) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p brief_stringer_t) WriteTo(w io.Writer) (int64, error) { return writeTo(w, p) }
func (p brief_stringer_t) WriteToBuffer(buf *buffer.Buffer)   { addBriefStr(buf, p.val.String()) }
func (p brief_stringer_t) String() string {
	buf := getBuf()
	defer buf.Free()
	p.WriteToBuffer(buf)
	return buf.String()
}

type brief_slice_t[T any] slice_t[T]

func (p *brief_slice_t[T]) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p *brief_slice_t[T]) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p *brief_slice_t[T]) WriteTo(w io.Writer) (int64, error) { return writeTo(w, p) }
func (p *brief_slice_t[T]) WriteToBuffer(buf *buffer.Buffer)   { addBriefSlice(buf, p) }
func (p *brief_slice_t[T]) String() string {
	buf := getBuf()
	defer buf.Free()
	p.WriteToBuffer(buf)
	return buf.String()
}

func BriefString(key, val string) Attr            { return slog.Any(key, brief_stringer_t{string_t(val)}) }
func BriefByteString(key string, val []byte) Attr { return BriefString(key, b2s(val)) }
func BriefStringer(key string, val fmt.Stringer) Attr {
	if val == nil {
		return nil_f(key)
	}
	return slog.Any(key, brief_stringer_t{val})
}

func BriefStringp(key string, val *string) Attr {
	if val == nil {
		return nil_f(key)
	}
	return BriefString(key, *val)
}
func BriefStrings(key string, val []string) Attr {
	return slog.Any(key, &brief_slice_t[string]{
		data: val,
		text: func(buf *buffer.Buffer, v string) { addBriefStr(buf, v) },
	})
}
func BriefByteStrings(key string, val [][]byte) Attr {
	return slog.Any(key, &brief_slice_t[[]byte]{
		data: val,
		text: func(buf *buffer.Buffer, v []byte) { addBriefStr(buf, b2s(v)) },
	})
}
func BriefStringers[T fmt.Stringer](key string, val []T) Attr {
	return slog.Any(key, &brief_slice_t[T]{
		data: val,
		text: func(buf *buffer.Buffer, v T) { addBriefStr(buf, v.String()) },
	})
}
