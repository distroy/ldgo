/*
 * Copyright (C) distroy
 */

package attr

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"

	"github.com/distroy/ldgo/v3/ldlog/internal/buffer"
	"github.com/distroy/ldgo/v3/ldptr"
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

func addBriefStr(b *buffer.Buffer, s string) {
	n := briefStringLen
	l := len(s)
	if l <= n {
		b.AppendString(quote(s))
		return
	}
	b.AppendByte('{')

	b.AppendString(tagLen)
	b.AppendByte(':')
	b.AppendInt(int64(l))
	b.AppendByte(',')

	b.AppendString(tagType)
	b.AppendByte(':')
	b.AppendString("string")
	b.AppendByte(',')

	b.AppendString(tagBrief)
	b.AppendByte(':')
	b.AppendString(quote(s[:n]))

	b.AppendByte('}')
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

	b.AppendByte('{')

	b.AppendString(tagLen)
	b.AppendByte(':')
	b.AppendInt(int64(l))
	b.AppendByte(',')

	b.AppendString(tagType)
	b.AppendByte(':')
	b.AppendString("array")
	b.AppendByte(',')

	x := &slice_t[T]{
		data: s.data[:n],
		text: s.text,
	}
	b.AppendString(tagBrief)
	b.AppendByte(':')
	x.WriteTo(b)

	b.AppendByte('}')
}

type stringer_t string

func (p stringer_t) String() string { return string(p) }

type StringPtr string

func (p *StringPtr) String() string { return string(ldptr.Get(p)) }

type brief_stringer_t struct {
	val fmt.Stringer
}

func (p brief_stringer_t) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p brief_stringer_t) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p brief_stringer_t) WriteTo(w io.Writer) (int64, error) { return writeStringer(w, p) }
func (p brief_stringer_t) String() string {
	buf := getBuf()
	defer buf.Free()
	addBriefStr(buf, p.val.String())
	return buf.String()
}

type StringArray interface {
	Len() int
	Get(idx int) string
}

type brief_slice_t[T any] slice_t[T]

func (p *brief_slice_t[T]) MarshalJSON() ([]byte, error)       { return s2b(p.String()), nil }
func (p *brief_slice_t[T]) MarshalText() ([]byte, error)       { return s2b(p.String()), nil }
func (p *brief_slice_t[T]) WriteTo(w io.Writer) (int64, error) { return writeStringer(w, p) }
func (p *brief_slice_t[T]) String() string {
	buf := getBuf()
	defer buf.Free()
	addBriefSlice(buf, p)
	return buf.String()
}

func BriefString(key, val string) Attr            { return slog.Any(key, brief_stringer_t{stringer_t(val)}) }
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
		text: func(buf []byte, v string) []byte { return strconv.AppendQuote(buf, v) },
	})
}
func BriefByteStrings(key string, val [][]byte) Attr {
	return slog.Any(key, &brief_slice_t[[]byte]{
		data: val,
		text: func(buf []byte, v []byte) []byte { return strconv.AppendQuote(buf, b2s(v)) },
	})
}
func BriefStringers[T fmt.Stringer](key string, val []T) Attr {
	return slog.Any(key, &brief_slice_t[T]{
		data: val,
		text: func(buf []byte, v T) []byte { return strconv.AppendQuote(buf, v.String()) },
	})
}
