/*
 * Copyright (C) distroy
 */

package _attr

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/distroy/ldgo/v3/ldlog/internal/_slogtype"
)

type (
	Attr = slog.Attr

	Value     = slog.Value
	LogValuer = slog.LogValuer
)

var (
	_ LogValuer = str_func_log_valuer_t(nil)
)

func attr(k string, v Value) Attr { return Attr{Key: k, Value: v} }
func value(num uint64, any any) slog.Value {
	return _slogtype.GetSValue(_slogtype.NewValue(num, any))
}

func Bool(k string, v bool) Attr { return slog.Bool(k, v) }

// func String(k string, v string) Attr   { return slog.String(k, v) }
// func Int(k string, v int) Attr         { return slog.Int(k, v) }
// func Int64(k string, v int64) Attr     { return slog.Int64(k, v) }
// func Uint64(k string, v uint64) Attr   { return slog.Uint64(k, v) }
// func Float64(k string, v float64) Attr { return slog.Float64(k, v) }

// func Duration(k string, v time.Duration) Attr { return slog.Duration(k, v) }
// func Time(k string, v time.Time) Attr         { return slog.Time(k, v) }
func Duration(k string, v time.Duration) Attr {
	return Reflect(k, str_func_log_valuer_t(v.String))
}
func Time(k string, v time.Time) Attr {
	fn := func() string { return v.Format(TimeLayout) }
	return Reflect(k, str_func_log_valuer_t(fn))
}

// func Group(k string, v ...any) Attr { return slog.Group(k, v...) }
func Group(k string, v ...Attr) Attr {
	if len(v) == 0 {
		return Skip()
	}
	return attr(k, slog.GroupValue(v...))
}

// func Any(k string, v any) Attr                { return attr(k, AnyValue(v)) }
func Any(k string, v any) Attr {
	switch vv := v.(type) {
	case bool:
		return Bool(k, vv)
	case *bool:
		return Boolp(k, vv)
	case []bool:
		return Bools(k, vv)

	case int:
		return Int(k, vv)
	case int8:
		return Int8(k, vv)
	case int16:
		return Int16(k, vv)
	case int32:
		return Int32(k, vv)
	case int64:
		return Int64(k, vv)

	case []int:
		return Ints(k, vv)
	case []int8:
		return Int8s(k, vv)
	case []int16:
		return Int16s(k, vv)
	case []int32:
		return Int32s(k, vv)
	case []int64:
		return Int64s(k, vv)

	case *int:
		return Intp(k, vv)
	case *int8:
		return Int8p(k, vv)
	case *int16:
		return Int16p(k, vv)
	case *int32:
		return Int32p(k, vv)
	case *int64:
		return Int64p(k, vv)

	case uint:
		return Uint(k, vv)
	case uintptr:
		return Uintptr(k, vv)
	case uint8:
		return Uint8(k, vv)
	case uint16:
		return Uint16(k, vv)
	case uint32:
		return Uint32(k, vv)
	case uint64:
		return Uint64(k, vv)

	case []uint:
		return Uints(k, vv)
	case []uintptr:
		return Uintptrs(k, vv)
	// case []uint8:
	// 	return Uint8s(k, vv)
	case []uint16:
		return Uint16s(k, vv)
	case []uint32:
		return Uint32s(k, vv)
	case []uint64:
		return Uint64s(k, vv)

	case *uint:
		return Uintp(k, vv)
	case *uintptr:
		return Uintptrp(k, vv)
	case *uint8:
		return Uint8p(k, vv)
	case *uint16:
		return Uint16p(k, vv)
	case *uint32:
		return Uint32p(k, vv)
	case *uint64:
		return Uint64p(k, vv)

	case float32:
		return Float32(k, vv)
	case float64:
		return Float64(k, vv)

	case []float32:
		return Float32s(k, vv)
	case []float64:
		return Float64s(k, vv)

	case *float32:
		return Float32p(k, vv)
	case *float64:
		return Float64p(k, vv)

	case complex64:
		return Complex64(k, vv)
	case complex128:
		return Complex128(k, vv)

	case []complex64:
		return Complex64s(k, vv)
	case []complex128:
		return Complex128s(k, vv)

	case *complex64:
		return Complex64p(k, vv)
	case *complex128:
		return Complex128p(k, vv)

	case time.Duration:
		return Duration(k, vv)
	case time.Time:
		return Time(k, vv)

	case []time.Duration:
		return Durations(k, vv)
	case []time.Time:
		return Times(k, vv)

	case *time.Duration:
		return Durationp(k, vv)
	case *time.Time:
		return Timep(k, vv)

	case error:
		return NamedError(k, vv)
	case []error:
		return Errors(k, vv)

	case []Attr:
		return Group(k, vv...)
	case slog.Kind:
		return slog.Any(k, vv)
	case Value:
		return attr(k, vv)

	case string:
		return String(k, vv)
	case []string:
		return Strings(k, vv)
	case *string:
		return Stringp(k, vv)

	case []byte:
		return ByteString(k, vv)
	case [][]byte:
		return ByteStrings(k, vv)
	case fmt.Stringer:
		return Stringer(k, vv)

	default:
		return Reflect(k, vv)
	}
}

func Reflect(k string, v any) Attr { return attr(k, value(0, v)) }

func Skip() Attr        { return Attr{} }
func Nil(k string) Attr { return Reflect(k, nil_t{}) }

type str_func_log_valuer_t func() string

func (f str_func_log_valuer_t) LogValue() Value { return slog.StringValue(f()) }

func String(k string, v string) Attr     { return slog.String(k, v) }
func ByteString(k string, v []byte) Attr { return String(k, b2s(v)) }
func Stringer(k string, v fmt.Stringer) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, str_func_log_valuer_t(v.String))
}
func Binary(k string, v []byte) Attr {
	fn := func() string { return b64(v) }
	return Reflect(k, str_func_log_valuer_t(fn))
}

func Complex64(k string, v complex64) Attr   { return Reflect(k, complex64_t(v)) }
func Complex128(k string, v complex128) Attr { return Reflect(k, complex128_t(v)) }

func Boolp(k string, v *bool) Attr {
	if v == nil {
		return Nil(k)
	}
	return Bool(k, *v)
}
func Stringp(k string, v *string) Attr {
	if v == nil {
		return Nil(k)
	}
	return String(k, *v)
}

func Bools(k string, v []bool) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[bool]{
		data: v,
		text: func(buf *Buffer, v bool) { buf.AppendBool(v) },
	})
}
func Strings(k string, v []string) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[string]{
		data: v,
		text: func(buf *Buffer, v string) { buf.AppendQuote(v) },
	})
}
func Stringers[T fmt.Stringer](k string, v []T) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[T]{
		data: v,
		text: func(buf *Buffer, v T) { buf.AppendQuote(v.String()) },
	})
}
func ByteStrings(k string, v [][]byte) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[[]byte]{
		data: v,
		text: func(buf *Buffer, v []byte) { buf.AppendQuote(b2s(v)) },
	})
}

func Int(k string, v int) Attr     { return slog.Int(k, v) }
func Int8(k string, v int8) Attr   { return Int64(k, int64(v)) }
func Int16(k string, v int16) Attr { return Int64(k, int64(v)) }
func Int32(k string, v int32) Attr { return Int64(k, int64(v)) }
func Int64(k string, v int64) Attr { return slog.Int64(k, v) }

func Uint(k string, v uint) Attr       { return Uint64(k, uint64(v)) }
func Uintptr(k string, v uintptr) Attr { return Uint64(k, uint64(v)) }
func Uint8(k string, v uint8) Attr     { return Uint64(k, uint64(v)) }
func Uint16(k string, v uint16) Attr   { return Uint64(k, uint64(v)) }
func Uint32(k string, v uint32) Attr   { return Uint64(k, uint64(v)) }
func Uint64(k string, v uint64) Attr   { return slog.Uint64(k, v) }

func Float32(k string, v float32) Attr { return Float64(k, float64(v)) }
func Float64(k string, v float64) Attr { return slog.Float64(k, v) }

func intpAttr[T int | int8 | int16 | int32 | int64](k string, v *T) Attr {
	if v == nil {
		return Nil(k)
	}
	return Int64(k, int64(*v))
}

func Intp(k string, v *int) Attr     { return intpAttr(k, v) }
func Int8p(k string, v *int8) Attr   { return intpAttr(k, v) }
func Int16p(k string, v *int16) Attr { return intpAttr(k, v) }
func Int32p(k string, v *int32) Attr { return intpAttr(k, v) }
func Int64p(k string, v *int64) Attr { return intpAttr(k, v) }

func uintpAttr[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](k string, v *T) Attr {
	if v == nil {
		return Nil(k)
	}
	return Uint64(k, uint64(*v))
}

func Uintp(k string, v *uint) Attr       { return uintpAttr(k, v) }
func Uint8p(k string, v *uint8) Attr     { return uintpAttr(k, v) }
func Uint16p(k string, v *uint16) Attr   { return uintpAttr(k, v) }
func Uint32p(k string, v *uint32) Attr   { return uintpAttr(k, v) }
func Uint64p(k string, v *uint64) Attr   { return uintpAttr(k, v) }
func Uintptrp(k string, v *uintptr) Attr { return uintpAttr(k, v) }

func floatpAttr[T float32 | float64](k string, v *T) Attr {
	if v == nil {
		return Nil(k)
	}
	return Float64(k, float64(*v))
}

func Float32p(k string, v *float32) Attr { return floatpAttr(k, v) }
func Float64p(k string, v *float64) Attr { return floatpAttr(k, v) }

func complexpAttr[T complex64 | complex128](k string, v *T) Attr {
	if v == nil {
		return Nil(k)
	}
	return Complex128(k, complex128(*v))
}

func Complex64p(k string, v *complex64) Attr   { return complexpAttr(k, v) }
func Complex128p(k string, v *complex128) Attr { return complexpAttr(k, v) }

func intsAttr[T int | int8 | int16 | int32 | int64](k string, v []T) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[T]{
		data: v,
		text: func(buf *Buffer, v T) { buf.AppendInt(int64(v)) },
	})
}

func Ints(k string, v []int) Attr     { return intsAttr(k, v) }
func Int8s(k string, v []int8) Attr   { return intsAttr(k, v) }
func Int16s(k string, v []int16) Attr { return intsAttr(k, v) }
func Int32s(k string, v []int32) Attr { return intsAttr(k, v) }
func Int64s(k string, v []int64) Attr { return intsAttr(k, v) }

func uintsAttr[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](k string, v []T) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[T]{
		data: v,
		text: func(buf *Buffer, v T) { buf.AppendUint(uint64(v)) },
	})
}
func Uints(k string, v []uint) Attr       { return uintsAttr(k, v) }
func Uint8s(k string, v []uint8) Attr     { return uintsAttr(k, v) }
func Uint16s(k string, v []uint16) Attr   { return uintsAttr(k, v) }
func Uint32s(k string, v []uint32) Attr   { return uintsAttr(k, v) }
func Uint64s(k string, v []uint64) Attr   { return uintsAttr(k, v) }
func Uintptrs(k string, v []uintptr) Attr { return uintsAttr(k, v) }

func floatsAttr[T float32 | float64](k string, v []T) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[T]{
		data: v,
		text: func(buf *Buffer, v T) { buf.AppendFloat(float64(v), 64) },
	})
}

func Float32s(k string, v []float32) Attr { return floatsAttr(k, v) }
func Float64s(k string, v []float64) Attr { return floatsAttr(k, v) }

func complexsAttr[T complex64 | complex128](k string, v []T) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[T]{
		data: v,
		text: func(buf *Buffer, v T) { buf.AppendString(complex128_t(v).String()) },
	})
}

func Complex64s(k string, v []complex64) Attr   { return complexsAttr(k, v) }
func Complex128s(k string, v []complex128) Attr { return complexsAttr(k, v) }

func Timep(k string, v *time.Time) Attr {
	if v == nil {
		return Nil(k)
	}
	return Time(k, *v)
}
func Durationp(k string, v *time.Duration) Attr {
	if v == nil {
		return Nil(k)
	}
	return Duration(k, *v)
}

func Durations(k string, v []time.Duration) Attr { return Stringers(k, v) }
func Times(k string, v []time.Time) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[time.Time]{
		data: v,
		text: func(buf *Buffer, v time.Time) {
			buf.AppendByte('"')
			buf.AppendTime(v, "")
			buf.AppendByte('"')
		},
	})
}

func Error(v error) Attr { return NamedError("error", v) }
func NamedError(k string, v error) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, str_func_log_valuer_t(v.Error))
}
func Errors(k string, v []error) Attr {
	if v == nil {
		return Nil(k)
	}
	return Reflect(k, &slice_t[error]{
		data: v,
		text: func(buf *Buffer, v error) {
			if v == nil {
				buf.AppendString(nil_t{}.String())
				return
			}
			buf.AppendQuote(v.Error())
		},
	})
}

func Stack(k string) Attr               { return StackSkip(k, 1) }
func StackSkip(k string, skip int) Attr { return String(k, stack(skip+1, 10)) }
