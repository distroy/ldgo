/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/distroy/ldgo/v3/ldlog/internal/field"
)

func Bool(key string, val bool) Attr             { return slog.Bool(key, val) }
func String(key string, val string) Attr         { return slog.String(key, val) }
func Stringer(key string, val fmt.Stringer) Attr { return slog.Stringer(key, val) }
func Binary(key string, val []byte) Attr         { return slog.Binary(key, val) }
func ByteString(key string, val []byte) Attr     { return slog.ByteString(key, val) }

func Int(key string, val int) Attr     { return slog.Int(key, val) }
func Int8(key string, val int8) Attr   { return slog.Int64(key, int64(val)) }
func Int16(key string, val int16) Attr { return slog.Int64(key, int64(val)) }
func Int32(key string, val int32) Attr { return slog.Int64(key, int64(val)) }
func Int64(key string, val int64) Attr { return slog.Int64(key, val) }

func Uint(key string, val uint) Attr       { return slog.Uint64(key, uint64(val)) }
func Uint8(key string, val uint8) Attr     { return slog.Uint64(key, uint64(val)) }
func Uint16(key string, val uint16) Attr   { return slog.Uint64(key, uint64(val)) }
func Uint32(key string, val uint32) Attr   { return slog.Uint64(key, uint64(val)) }
func Uint64(key string, val uint64) Attr   { return slog.Uint64(key, val) }
func Uintptr(key string, val uintptr) Attr { return slog.Uint64(key, uint64(val)) }

func Float32(key string, val float32) Attr { return slog.Float64(key, float64(val)) }
func Float64(key string, val float64) Attr { return slog.Float64(key, val) }

func Complex64(key string, val complex64) Attr   { return slog.Complex64(key, val) }
func Complex128(key string, val complex128) Attr { return slog.Complex128(key, val) }

func Boolp(key string, val *bool) Attr                   { return slog.Boolp(key, val) }
func Bools(key string, val []bool) Attr                  { return slog.Bools(key, val) }
func Stringp(key string, val *string) Attr               { return slog.Stringp(key, val) }
func Strings(key string, val []string) Attr              { return slog.Strings(key, val) }
func Stringers[T fmt.Stringer](key string, val []T) Attr { return slog.Stringers(key, val) }
func ByteStrings(key string, val [][]byte) Attr          { return slog.ByteStrings(key, val) }

func Intp(key string, val *int) Attr     { return slog.Intp(key, val) }
func Int8p(key string, val *int8) Attr   { return slog.Int8p(key, val) }
func Int16p(key string, val *int16) Attr { return slog.Int16p(key, val) }
func Int32p(key string, val *int32) Attr { return slog.Int32p(key, val) }
func Int64p(key string, val *int64) Attr { return slog.Int64p(key, val) }

func Uintp(key string, val *uint) Attr       { return slog.Uintp(key, val) }
func Uint8p(key string, val *uint8) Attr     { return slog.Uint8p(key, val) }
func Uint16p(key string, val *uint16) Attr   { return slog.Uint16p(key, val) }
func Uint32p(key string, val *uint32) Attr   { return slog.Uint32p(key, val) }
func Uint64p(key string, val *uint64) Attr   { return slog.Uint64p(key, val) }
func Uintptrp(key string, val *uintptr) Attr { return slog.Uintptrp(key, val) }

func Float32p(key string, val *float32) Attr { return slog.Float32p(key, val) }
func Float64p(key string, val *float64) Attr { return slog.Float64p(key, val) }

func Complex64p(key string, val *complex64) Attr   { return slog.Complex64p(key, val) }
func Complex128p(key string, val *complex128) Attr { return slog.Complex128p(key, val) }

func Ints(key string, val []int) Attr     { return slog.Ints(key, val) }
func Int8s(key string, val []int8) Attr   { return slog.Int8s(key, val) }
func Int16s(key string, val []int16) Attr { return slog.Int16s(key, val) }
func Int32s(key string, val []int32) Attr { return slog.Int32s(key, val) }
func Int64s(key string, val []int64) Attr { return slog.Int64s(key, val) }

func Uints(key string, val []uint) Attr       { return slog.Uints(key, val) }
func Uint8s(key string, val []uint8) Attr     { return slog.Uint8s(key, val) }
func Uint16s(key string, val []uint16) Attr   { return slog.Uint16s(key, val) }
func Uint32s(key string, val []uint32) Attr   { return slog.Uint32s(key, val) }
func Uint64s(key string, val []uint64) Attr   { return slog.Uint64s(key, val) }
func Uintptrs(key string, val []uintptr) Attr { return slog.Uintptrs(key, val) }

func Float32s(key string, val []float32) Attr { return slog.Float32s(key, val) }
func Float64s(key string, val []float64) Attr { return slog.Float64s(key, val) }

func Complex64s(key string, val []complex64) Attr   { return slog.Complex64s(key, val) }
func Complex128s(key string, val []complex128) Attr { return slog.Complex128s(key, val) }

func Time(key string, val time.Time) Attr         { return slog.Time(key, val) }
func Duration(key string, val time.Duration) Attr { return slog.Duration(key, val) }

func Timep(key string, val *time.Time) Attr         { return slog.Timep(key, val) }
func Durationp(key string, val *time.Duration) Attr { return slog.Durationp(key, val) }

func Times(key string, val []time.Time) Attr         { return slog.Times(key, val) }
func Durations(key string, val []time.Duration) Attr { return slog.Durations(key, val) }

func Error(err error) Attr                  { return slog.Error(err) }
func Errors(key string, err []error) Attr   { return slog.Errors(key, err) }
func NamedError(key string, err error) Attr { return slog.NamedError(key, err) }

func Stack(key string) Attr               { return slog.StackSkip(key, 1) }
func StackSkip(key string, skip int) Attr { return slog.StackSkip(key, skip+1) }

func Namespace(key string) Attr       { return slog.Namespace(key) }
func Inline(val ObjectMarshaler) Attr { return slog.Inline(val) }

func Array(key string, val ArrayMarshaler) Attr           { return slog.Array(key, val) }
func Object(key string, val ObjectMarshaler) Attr         { return slog.Object(key, val) }
func Objects[T ObjectMarshaler](key string, val []T) Attr { return slog.Objects(key, val) }

func Any(key string, val interface{}) Attr     { return slog.Any(key, val) }
func Reflect(key string, val interface{}) Attr { return slog.Reflect(key, val) }

func Skip() Attr { return Attr{} }

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~uintptr
}

func Integer[Int integer](key string, val Int) Attr {
	if val >= 0 {
		return Uint64(key, uint64(val))
	}
	return Int64(key, int64(val))
}

func BriefString(key string, val string) Attr         { return field.BriefString(key, val) }
func BriefByteString(key string, val []byte) Attr     { return field.BriefByteString(key, val) }
func BriefStringer(key string, val fmt.Stringer) Attr { return field.BriefStringer(key, val) }

func BriefStringp(key string, val *string) Attr               { return field.BriefStringp(key, val) }
func BriefStrings(key string, val []string) Attr              { return field.BriefStrings(key, val) }
func BriefByteStrings(key string, val [][]byte) Attr          { return field.BriefByteStrings(key, val) }
func BriefStringers[T fmt.Stringer](key string, val []T) Attr { return field.BriefStringers(key, val) }

func BriefReflect(key string, val interface{}) Attr { return field.BriefReflect(key, val) }
