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
	Value     = slog.Value
	LogValuer = slog.LogValuer
)

var (
	_ LogValuer = string_value_t(nil)
)

func value(num uint64, any any) slog.Value {
	return _slogtype.GetSValue(_slogtype.NewValue(num, any))
}

func BoolValue(v bool) Value         { return slog.BoolValue(v) }
func StringValue(value string) Value { return slog.StringValue(value) }
func Float64Value(v float64) Value   { return slog.Float64Value(v) }

func IntValue(v int) Value       { return slog.IntValue(v) }
func Int64Value(v int64) Value   { return slog.Int64Value(v) }
func Uint64Value(v uint64) Value { return slog.Uint64Value(v) }

// func TimeValue(v time.Time) Value         { return slog.TimeValue(v) }
func TimeValue(v time.Time) Value {
	return value(0, string_value_t(func() string { return v.Format(TimeLayout) }))
}

// func DurationValue(v time.Duration) Value { return slog.DurationValue(v) }
func DurationValue(v time.Duration) Value {
	if v == 0 {
		return IntValue(0)
	}
	return value(0, string_value_t(v.String))
}

func GroupValue(as ...Attr) Value { return slog.GroupValue(as...) }

// func AnyValue(v any) Value { return slog.AnyValue(v) }
func AnyValue(v any) Value {
	switch v := v.(type) {
	case string:
		return StringValue(v)
	case int:
		return IntValue(v)
	case uint:
		return Uint64Value(uint64(v))
	case int64:
		return Int64Value(v)
	case uint64:
		return Uint64Value(v)
	case bool:
		return BoolValue(v)
	case time.Duration:
		return DurationValue(v)
	case time.Time:
		return TimeValue(v)
	case uint8:
		return Uint64Value(uint64(v))
	case uint16:
		return Uint64Value(uint64(v))
	case uint32:
		return Uint64Value(uint64(v))
	case uintptr:
		return Uint64Value(uint64(v))
	case int8:
		return Int64Value(int64(v))
	case int16:
		return Int64Value(int64(v))
	case int32:
		return Int64Value(int64(v))
	case float64:
		return Float64Value(v)
	case float32:
		return Float64Value(float64(v))
	case complex64:
		return Complex64Value(v)
	case complex128:
		return Complex128Value(v)
	case []Attr:
		return GroupValue(v...)
	case slog.Kind:
		return slog.AnyValue(v)
	case Value:
		return v
	case []byte:
		return StringValue(b2s(v))
	case fmt.Stringer:
		return StringerValue(v)
	case error:
		return ErrorValue(v)
	default:
		return ReflectValue(v)
	}
}

// *** all codes below are customized ***

func ReflectValue(v any) Value { return value(0, v) }

func NilValue() Value { return ReflectValue(nil_t{}) }

type string_value_t func() string

func (f string_value_t) LogValue() Value { return StringValue(f()) }

func StrFnValue(v func() string) Value { return ReflectValue(string_value_t(v)) }
func ByteStringValue(v []byte) Value   { return StringValue(b2s(v)) }
func StringerValue(v fmt.Stringer) Value {
	if v == nil {
		return NilValue()
	}
	return ReflectValue(string_value_t(v.String))
}

func Complex64Value(v complex64) Value   { return ReflectValue(complex64_t(v)) }
func Complex128Value(v complex128) Value { return ReflectValue(complex128_t(v)) }

func ErrorValue(v error) Value {
	if v == nil {
		return NilValue()
	}
	return ReflectValue(string_value_t(v.Error))
}
