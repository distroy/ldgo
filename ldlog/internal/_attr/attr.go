/*
 * Copyright (C) distroy
 */

package _attr

import (
	"fmt"
	"log/slog"
	"time"
)

type (
	Attr = slog.Attr
)

func attr(k string, v Value) Attr { return Attr{Key: k, Value: v} }

func Any(k string, v any) Attr                { return attr(k, AnyValue(v)) }
func Bool(k string, v bool) Attr              { return attr(k, BoolValue(v)) }
func String(k string, v string) Attr          { return attr(k, StringValue(v)) }
func Int(k string, v int) Attr                { return attr(k, IntValue(v)) }
func Int64(k string, v int64) Attr            { return attr(k, Int64Value(v)) }
func Uint64(k string, v uint64) Attr          { return attr(k, Uint64Value(v)) }
func Float64(k string, v float64) Attr        { return attr(k, Float64Value(v)) }
func Duration(k string, v time.Duration) Attr { return attr(k, DurationValue(v)) }
func Time(k string, v time.Time) Attr         { return attr(k, TimeValue(v)) }

func Reflect(k string, v any) Attr { return attr(k, ReflectValue(v)) }

func Skip() Attr        { return Attr{} }
func Nil(k string) Attr { return attr(k, NilValue()) }

func Binary(k string, v []byte) Attr         { return attr(k, StrFnValue(func() string { return b64(v) })) }
func ByteString(k string, v []byte) Attr     { return String(k, b2s(v)) }
func Stringer(k string, v fmt.Stringer) Attr { return attr(k, StringerValue(v)) }

func Complex64(key string, val complex64) Attr   { return attr(key, Complex64Value(val)) }
func Complex128(key string, val complex128) Attr { return attr(key, Complex128Value(val)) }

func Boolp(key string, val *bool) Attr {
	if val == nil {
		return Nil(key)
	}
	return Bool(key, *val)
}
func Stringp(key string, val *string) Attr {
	if val == nil {
		return Nil(key)
	}
	return String(key, *val)
}

func Bools(key string, val []bool) Attr {
	return Reflect(key, &slice_t[bool]{
		data: val,
		text: func(buf *Buffer, v bool) { buf.AppendBool(v) },
	})
}
func Strings(key string, val []string) Attr {
	return Reflect(key, &slice_t[string]{
		data: val,
		text: func(buf *Buffer, v string) { buf.AppendQuote(v) },
	})
}
func Stringers[T fmt.Stringer](key string, val []T) Attr {
	return Reflect(key, &slice_t[T]{
		data: val,
		text: func(buf *Buffer, v T) { buf.AppendQuote(v.String()) },
	})
}
func ByteStrings(key string, val [][]byte) Attr {
	return Reflect(key, &slice_t[[]byte]{
		data: val,
		text: func(buf *Buffer, v []byte) { buf.AppendQuote(b2s(v)) },
	})
}

func intptr_f[T int | int8 | int16 | int32 | int64](key string, val *T) Attr {
	if val == nil {
		return Nil(key)
	}
	return Int64(key, int64(*val))
}

func Intp(key string, val *int) Attr     { return intptr_f(key, val) }
func Int8p(key string, val *int8) Attr   { return intptr_f(key, val) }
func Int16p(key string, val *int16) Attr { return intptr_f(key, val) }
func Int32p(key string, val *int32) Attr { return intptr_f(key, val) }
func Int64p(key string, val *int64) Attr { return intptr_f(key, val) }

func uintptr_f[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](key string, val *T) Attr {
	if val == nil {
		return Nil(key)
	}
	return Uint64(key, uint64(*val))
}

func Uintp(key string, val *uint) Attr       { return uintptr_f(key, val) }
func Uint8p(key string, val *uint8) Attr     { return uintptr_f(key, val) }
func Uint16p(key string, val *uint16) Attr   { return uintptr_f(key, val) }
func Uint32p(key string, val *uint32) Attr   { return uintptr_f(key, val) }
func Uint64p(key string, val *uint64) Attr   { return uintptr_f(key, val) }
func Uintptrp(key string, val *uintptr) Attr { return uintptr_f(key, val) }

func float_f[T float32 | float64](key string, val *T) Attr {
	if val == nil {
		return Nil(key)
	}
	return Float64(key, float64(*val))
}

func Float32p(key string, val *float32) Attr { return float_f(key, val) }
func Float64p(key string, val *float64) Attr { return float_f(key, val) }

func complexptr_f[T complex64 | complex128](key string, val *T) Attr {
	if val == nil {
		return Nil(key)
	}
	return Complex128(key, complex128(*val))
}

func Complex64p(key string, val *complex64) Attr   { return complexptr_f(key, val) }
func Complex128p(key string, val *complex128) Attr { return complexptr_f(key, val) }

func ints_f[T int | int8 | int16 | int32 | int64](key string, val []T) Attr {
	if val == nil {
		return Nil(key)
	}
	return Reflect(key, &slice_t[T]{
		data: val,
		text: func(buf *Buffer, v T) { buf.AppendInt(int64(v)) },
	})
}

func Ints(key string, val []int) Attr     { return ints_f(key, val) }
func Int8s(key string, val []int8) Attr   { return ints_f(key, val) }
func Int16s(key string, val []int16) Attr { return ints_f(key, val) }
func Int32s(key string, val []int32) Attr { return ints_f(key, val) }
func Int64s(key string, val []int64) Attr { return ints_f(key, val) }

func uints_f[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](key string, val []T) Attr {
	if val == nil {
		return Nil(key)
	}
	return Reflect(key, &slice_t[T]{
		data: val,
		text: func(buf *Buffer, v T) { buf.AppendUint(uint64(v)) },
	})
}
func Uints(key string, val []uint) Attr       { return uints_f(key, val) }
func Uint8s(key string, val []uint8) Attr     { return uints_f(key, val) }
func Uint16s(key string, val []uint16) Attr   { return uints_f(key, val) }
func Uint32s(key string, val []uint32) Attr   { return uints_f(key, val) }
func Uint64s(key string, val []uint64) Attr   { return uints_f(key, val) }
func Uintptrs(key string, val []uintptr) Attr { return uints_f(key, val) }

func floats_f[T float32 | float64](key string, val []T) Attr {
	if val == nil {
		return Nil(key)
	}
	return Reflect(key, &slice_t[T]{
		data: val,
		text: func(buf *Buffer, v T) { buf.AppendFloat(float64(v), 64) },
	})
}

func Float32s(key string, val []float32) Attr { return floats_f(key, val) }
func Float64s(key string, val []float64) Attr { return floats_f(key, val) }

func complexs_f[T complex64 | complex128](key string, val []T) Attr {
	if val == nil {
		return Nil(key)
	}
	return Reflect(key, &slice_t[T]{
		data: val,
		text: func(buf *Buffer, v T) { buf.AppendString(complex128_t(v).String()) },
	})
}

func Complex64s(key string, val []complex64) Attr   { return complexs_f(key, val) }
func Complex128s(key string, val []complex128) Attr { return complexs_f(key, val) }

func Timep(key string, val *time.Time) Attr {
	if val == nil {
		return Nil(key)
	}
	return Time(key, *val)
}
func Durationp(key string, val *time.Duration) Attr {
	if val == nil {
		return Nil(key)
	}
	return Duration(key, *val)
}

func Durations(key string, val []time.Duration) Attr { return Stringers(key, val) }
func Times(key string, val []time.Time) Attr {
	// return Stringers(key, val)
	return Reflect(key, &slice_t[time.Time]{
		data: val,
		text: func(buf *Buffer, v time.Time) {
			buf.AppendByte('"')
			buf.AppendTime(v, "")
			buf.AppendByte('"')
		},
	})
}

func Error(err error) Attr { return NamedError("error", err) }
func NamedError(k string, v error) Attr {
	if v == nil {
		return Nil(k)
	}
	return attr(k, StrFnValue(v.Error))
}
func Errors(key string, err []error) Attr {
	return Reflect(key, &slice_t[error]{
		data: err,
		text: func(buf *Buffer, v error) {
			if v == nil {
				buf.AppendString(nil_t{}.String())
				return
			}
			buf.AppendQuote(v.Error())
		},
	})
}

func Stack(key string) Attr               { return StackSkip(key, 1) }
func StackSkip(key string, skip int) Attr { return String(key, stack(skip+1, 10)) }
