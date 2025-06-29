/*
 * Copyright (C) distroy
 */

package attr

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"
)

type (
	Attr = slog.Attr

	Value     = slog.Value
	LogValuer = slog.LogValuer
)

var (
	_ LogValuer = string_t(nil)
)

func Any(k string, v any) Attr                { return slog.Any(k, v) }
func Bool(k string, v bool) Attr              { return slog.Bool(k, v) }
func String(k string, v string) Attr          { return slog.String(k, v) }
func Int(k string, v int) Attr                { return slog.Int(k, v) }
func Int64(k string, v int64) Attr            { return slog.Int64(k, v) }
func Uint64(k string, v uint64) Attr          { return slog.Uint64(k, v) }
func Float64(k string, v float64) Attr        { return slog.Float64(k, v) }
func Time(k string, v time.Time) Attr         { return slog.Time(k, v) }
func Duration(k string, v time.Duration) Attr { return slog.Duration(k, v) }

func Skip() Attr          { return Attr{} }
func nil_f(k string) Attr { return slog.Any(k, nil_t{}) }

func string_f(k string, f func() string) Attr { return slog.Any(k, string_t(f)) }

type string_t func() string

func (f string_t) LogValue() Value { return slog.StringValue(f()) }

func Binary(key string, val []byte) Attr     { return string_f(key, func() string { return b64(val) }) }
func ByteString(key string, val []byte) Attr { return slog.String(key, b2s(val)) }
func Stringer(key string, val fmt.Stringer) Attr {
	if val == nil {
		return nil_f(key)
	}
	return string_f(key, val.String)
}

func Complex64(key string, val complex64) Attr   { return slog.Any(key, complex_t(val)) }
func Complex128(key string, val complex128) Attr { return slog.Any(key, complex_t(val)) }

func Boolp(key string, val *bool) Attr {
	if val == nil {
		return nil_f(key)
	}
	return slog.Bool(key, *val)
}
func Stringp(key string, val *string) Attr {
	if val == nil {
		return slog.Any(key, nil)
	}
	return slog.String(key, *val)
}

func Bools(key string, val []bool) Attr {
	return slog.Any(key, &slice_t[bool]{
		data: val,
		text: func(buf []byte, v bool) []byte { return strconv.AppendBool(buf, v) },
	})
}
func Strings(key string, val []string) Attr {
	return slog.Any(key, &slice_t[string]{
		data: val,
		text: func(buf []byte, v string) []byte { return strconv.AppendQuote(buf, v) },
	})
}
func Stringers[T fmt.Stringer](key string, val []T) Attr {
	return slog.Any(key, &slice_t[T]{
		data: val,
		text: func(buf []byte, v T) []byte { return strconv.AppendQuote(buf, v.String()) },
	})
}
func ByteStrings(key string, val [][]byte) Attr {
	return slog.Any(key, &slice_t[[]byte]{
		data: val,
		text: func(buf []byte, v []byte) []byte { return strconv.AppendQuote(buf, b2s(v)) },
	})
}

func intptr_f[T int | int8 | int16 | int32 | int64](key string, val *T) Attr {
	if val == nil {
		return slog.Any(key, nil)
	}
	return slog.Int64(key, int64(*val))
}

func Intp(key string, val *int) Attr     { return intptr_f(key, val) }
func Int8p(key string, val *int8) Attr   { return intptr_f(key, val) }
func Int16p(key string, val *int16) Attr { return intptr_f(key, val) }
func Int32p(key string, val *int32) Attr { return intptr_f(key, val) }
func Int64p(key string, val *int64) Attr { return intptr_f(key, val) }

func uintptr_f[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](key string, val *T) Attr {
	if val == nil {
		return slog.Any(key, nil)
	}
	return slog.Uint64(key, uint64(*val))
}

func Uintp(key string, val *uint) Attr       { return uintptr_f(key, val) }
func Uint8p(key string, val *uint8) Attr     { return uintptr_f(key, val) }
func Uint16p(key string, val *uint16) Attr   { return uintptr_f(key, val) }
func Uint32p(key string, val *uint32) Attr   { return uintptr_f(key, val) }
func Uint64p(key string, val *uint64) Attr   { return uintptr_f(key, val) }
func Uintptrp(key string, val *uintptr) Attr { return uintptr_f(key, val) }

func float_f[T float32 | float64](key string, val *T) Attr {
	if val == nil {
		return slog.Any(key, nil)
	}
	return slog.Float64(key, float64(*val))
}

func Float32p(key string, val *float32) Attr { return float_f(key, val) }
func Float64p(key string, val *float64) Attr { return float_f(key, val) }

func complexptr_f[T complex64 | complex128](key string, val *T) Attr {
	if val == nil {
		return slog.Any(key, nil)
	}
	return Complex128(key, complex128(*val))
}

func Complex64p(key string, val *complex64) Attr   { return complexptr_f(key, val) }
func Complex128p(key string, val *complex128) Attr { return complexptr_f(key, val) }

func ints_f[T int | int8 | int16 | int32 | int64](key string, val []T) Attr {
	if val == nil {
		return nil_f(key)
	}
	return slog.Any(key, &slice_t[T]{
		data: val,
		text: func(buf []byte, v T) []byte { return strconv.AppendInt(buf, int64(v), 10) },
	})
}

func Ints(key string, val []int) Attr     { return ints_f(key, val) }
func Int8s(key string, val []int8) Attr   { return ints_f(key, val) }
func Int16s(key string, val []int16) Attr { return ints_f(key, val) }
func Int32s(key string, val []int32) Attr { return ints_f(key, val) }
func Int64s(key string, val []int64) Attr { return ints_f(key, val) }

func uints_f[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](key string, val []T) Attr {
	if val == nil {
		return nil_f(key)
	}
	return slog.Any(key, &slice_t[T]{
		data: val,
		text: func(buf []byte, v T) []byte { return strconv.AppendUint(buf, uint64(v), 10) },
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
		return nil_f(key)
	}
	return slog.Any(key, &slice_t[T]{
		data: val,
		text: func(buf []byte, v T) []byte { return strconv.AppendFloat(buf, float64(v), 'g', -1, 64) },
	})
}

func Float32s(key string, val []float32) Attr { return floats_f(key, val) }
func Float64s(key string, val []float64) Attr { return floats_f(key, val) }

func complexs_f[T complex64 | complex128](key string, val []T) Attr {
	if val == nil {
		return nil_f(key)
	}
	return slog.Any(key, &slice_t[T]{
		data: val,
		text: func(buf []byte, v T) []byte { return append(buf, complex_t(v).String()...) },
	})
}

func Complex64s(key string, val []complex64) Attr   { return complexs_f(key, val) }
func Complex128s(key string, val []complex128) Attr { return complexs_f(key, val) }

func Timep(key string, val *time.Time) Attr {
	if val == nil {
		return nil_f(key)
	}
	return slog.Time(key, *val)
}
func Durationp(key string, val *time.Duration) Attr {
	if val == nil {
		return nil_f(key)
	}
	return slog.Duration(key, *val)
}

func Times(key string, val []time.Time) Attr         { return Stringers(key, val) }
func Durations(key string, val []time.Duration) Attr { return Stringers(key, val) }

func Error(err error) Attr              { return NamedError("error", err) }
func NamedError(k string, e error) Attr { return string_f(k, func() string { return e2s(e) }) }
func Errors(key string, err []error) Attr {
	return slog.Any(key, &slice_t[error]{
		data: err,
		text: func(buf []byte, v error) []byte { return strconv.AppendQuote(buf, e2s(v)) },
	})
}

func Stack(key string) Attr               { return StackSkip(key, 1) }
func StackSkip(key string, skip int) Attr { return slog.String(key, stack(skip+1, 10)) }
