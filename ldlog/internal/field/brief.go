/*
 * Copyright (C) distroy
 */

package field

import (
	"fmt"

	"github.com/distroy/ldgo/v2/ldconv"
	"github.com/distroy/ldgo/v2/ldmath"
	"github.com/distroy/ldgo/v2/ldptr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Field = zapcore.Field

	ObjectEncoder   = zapcore.ObjectEncoder
	ObjectMarshaler = zapcore.ObjectMarshaler
	ArrayEncoder    = zapcore.ArrayEncoder
	ArrayMarshaler  = zapcore.ArrayMarshaler
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

func bytes2str(d []byte) string { return ldconv.BytesToStrUnsafe(d) }

func SetBriefStringLen(n int) { briefStringLen = ldmath.MaxInt(n, minBriefStringLen) }
func SetBriefArrayLen(n int)  { briefArrayLen = ldmath.MaxInt(n, minBriefArrayLen) }
func SetBriefMapLen(n int)    { briefMapLen = ldmath.MaxInt(n, minBriefMapLen) }

func AddStr2Log(enc ObjectEncoder, k, s string) error {
	oe := &objectEncoder{
		Key: k,
		Enc: enc,
	}
	return AppendStr2Log(oe, s)
}

func AppendStr2Log(enc ArrayEncoder, v string) error {
	n := briefStringLen
	l := len(v)
	if l <= n {
		enc.AppendString(v)
		return nil
	}
	return enc.AppendObject(&briefStringerType{Len: l, Str: v[:n]})
}

func AddStrs2Log(enc ObjectEncoder, k string, v StringArray) error {
	oe := &objectEncoder{
		Key: k,
		Enc: enc,
	}
	return AppendStrs2Log(oe, v)
}

func AppendStrs2Log(enc ArrayEncoder, v StringArray) error {
	n := briefArrayLen
	l := v.Len()
	if l <= n {
		return enc.AppendArray(&briefStringArrayType{Len: l, Val: v})
	}
	return enc.AppendObject(&briefStringArrayType{Len: n, Val: v})
}

type String string

func (p String) String() string { return string(p) }

type StringPtr string

func (p *StringPtr) String() string { return string(ldptr.Get(p)) }

type BriefStringerType struct {
	Key string
	Val fmt.Stringer
}

func (p *BriefStringerType) MarshalLogObject(enc ObjectEncoder) error {
	return AddStr2Log(enc, p.Key, p.Val.String())
}

type briefStringerType struct {
	Len int
	Str string
}

func (p *briefStringerType) MarshalLogObject(enc ObjectEncoder) error {
	enc.AddInt(tagLen, p.Len)
	enc.AddString(tagType, "string")
	enc.AddString(tagBrief, p.Str)
	return nil
}

type StringArray interface {
	Len() int
	Get(idx int) string
}

type BriefStringersType struct {
	Key string
	Val StringArray
}

func (p *BriefStringersType) MarshalLogObject(enc ObjectEncoder) error {
	return AddStrs2Log(enc, p.Key, p.Val)
}

type briefStringArrayType struct {
	Len int
	Val StringArray
}

func (p *briefStringArrayType) MarshalLogArray(enc ArrayEncoder) error {
	for i := 0; i < p.Len; i++ {
		v := p.Val.Get(i)
		err := AppendStr2Log(enc, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *briefStringArrayType) MarshalLogObject(enc ObjectEncoder) error {
	enc.AddInt(tagLen, p.Val.Len())
	enc.AddString(tagType, "array")
	return enc.AddArray(tagBrief, p)
}

type Strings []string

func (p Strings) Len() int         { return len(p) }
func (p Strings) Get(i int) string { return p[i] }

type ByteStrings [][]byte

func (p ByteStrings) Len() int         { return len(p) }
func (p ByteStrings) Get(i int) string { return bytes2str(p[i]) }

type Stringers[T fmt.Stringer] []T

func (p Stringers[T]) Len() int         { return len(p) }
func (p Stringers[T]) Get(i int) string { return p[i].String() }

func BriefString(key, val string) Field            { return BriefStringer(key, String(val)) }
func BriefByteString(key string, val []byte) Field { return BriefString(key, bytes2str(val)) }
func BriefStringer(key string, val fmt.Stringer) Field {
	return zap.Inline(&BriefStringerType{Key: key, Val: val})
}

func BriefStringp(key string, val *string) Field      { return BriefStringer(key, (*StringPtr)(val)) }
func BriefStrings(key string, val []string) Field     { return briefStrings(key, Strings(val)) }
func BriefByteStrings(key string, val [][]byte) Field { return briefStrings(key, ByteStrings(val)) }
func BriefStringers[T fmt.Stringer](key string, val []T) Field {
	return briefStrings(key, Stringers[T](val))
}

func briefStrings(key string, val StringArray) Field {
	return zap.Inline(&BriefStringersType{Key: key, Val: val})
}
