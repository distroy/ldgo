/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"fmt"
	"time"

	"github.com/distroy/ldgo/v3/ldlog/internal/attr"
)

func Skip() Attr { return attr.Skip() }

func Bool(key string, val bool) Attr             { return attr.Bool(key, val) }
func String(key string, val string) Attr         { return attr.String(key, val) }
func Stringer(key string, val fmt.Stringer) Attr { return attr.Stringer(key, val) }
func Binary(key string, val []byte) Attr         { return attr.Binary(key, val) }
func ByteString(key string, val []byte) Attr     { return attr.ByteString(key, val) }

func Int(key string, val int) Attr     { return attr.Int(key, val) }
func Int8(key string, val int8) Attr   { return attr.Int64(key, int64(val)) }
func Int16(key string, val int16) Attr { return attr.Int64(key, int64(val)) }
func Int32(key string, val int32) Attr { return attr.Int64(key, int64(val)) }
func Int64(key string, val int64) Attr { return attr.Int64(key, val) }

func Uint(key string, val uint) Attr       { return attr.Uint64(key, uint64(val)) }
func Uint8(key string, val uint8) Attr     { return attr.Uint64(key, uint64(val)) }
func Uint16(key string, val uint16) Attr   { return attr.Uint64(key, uint64(val)) }
func Uint32(key string, val uint32) Attr   { return attr.Uint64(key, uint64(val)) }
func Uint64(key string, val uint64) Attr   { return attr.Uint64(key, val) }
func Uintptr(key string, val uintptr) Attr { return attr.Uint64(key, uint64(val)) }

func Float32(key string, val float32) Attr { return attr.Float64(key, float64(val)) }
func Float64(key string, val float64) Attr { return attr.Float64(key, val) }

func Complex64(key string, val complex64) Attr   { return attr.Complex64(key, val) }
func Complex128(key string, val complex128) Attr { return attr.Complex128(key, val) }

func Boolp(key string, val *bool) Attr                   { return attr.Boolp(key, val) }
func Bools(key string, val []bool) Attr                  { return attr.Bools(key, val) }
func Stringp(key string, val *string) Attr               { return attr.Stringp(key, val) }
func Strings(key string, val []string) Attr              { return attr.Strings(key, val) }
func Stringers[T fmt.Stringer](key string, val []T) Attr { return attr.Stringers(key, val) }
func ByteStrings(key string, val [][]byte) Attr          { return attr.ByteStrings(key, val) }

func Intp(key string, val *int) Attr     { return attr.Intp(key, val) }
func Int8p(key string, val *int8) Attr   { return attr.Int8p(key, val) }
func Int16p(key string, val *int16) Attr { return attr.Int16p(key, val) }
func Int32p(key string, val *int32) Attr { return attr.Int32p(key, val) }
func Int64p(key string, val *int64) Attr { return attr.Int64p(key, val) }

func Uintp(key string, val *uint) Attr       { return attr.Uintp(key, val) }
func Uint8p(key string, val *uint8) Attr     { return attr.Uint8p(key, val) }
func Uint16p(key string, val *uint16) Attr   { return attr.Uint16p(key, val) }
func Uint32p(key string, val *uint32) Attr   { return attr.Uint32p(key, val) }
func Uint64p(key string, val *uint64) Attr   { return attr.Uint64p(key, val) }
func Uintptrp(key string, val *uintptr) Attr { return attr.Uintptrp(key, val) }

func Float32p(key string, val *float32) Attr { return attr.Float32p(key, val) }
func Float64p(key string, val *float64) Attr { return attr.Float64p(key, val) }

func Complex64p(key string, val *complex64) Attr   { return attr.Complex64p(key, val) }
func Complex128p(key string, val *complex128) Attr { return attr.Complex128p(key, val) }

func Ints(key string, val []int) Attr     { return attr.Ints(key, val) }
func Int8s(key string, val []int8) Attr   { return attr.Int8s(key, val) }
func Int16s(key string, val []int16) Attr { return attr.Int16s(key, val) }
func Int32s(key string, val []int32) Attr { return attr.Int32s(key, val) }
func Int64s(key string, val []int64) Attr { return attr.Int64s(key, val) }

func Uints(key string, val []uint) Attr       { return attr.Uints(key, val) }
func Uint8s(key string, val []uint8) Attr     { return attr.Uint8s(key, val) }
func Uint16s(key string, val []uint16) Attr   { return attr.Uint16s(key, val) }
func Uint32s(key string, val []uint32) Attr   { return attr.Uint32s(key, val) }
func Uint64s(key string, val []uint64) Attr   { return attr.Uint64s(key, val) }
func Uintptrs(key string, val []uintptr) Attr { return attr.Uintptrs(key, val) }

func Float32s(key string, val []float32) Attr { return attr.Float32s(key, val) }
func Float64s(key string, val []float64) Attr { return attr.Float64s(key, val) }

func Complex64s(key string, val []complex64) Attr   { return attr.Complex64s(key, val) }
func Complex128s(key string, val []complex128) Attr { return attr.Complex128s(key, val) }

func Time(key string, val time.Time) Attr         { return attr.Time(key, val) }
func Duration(key string, val time.Duration) Attr { return attr.Duration(key, val) }

func Timep(key string, val *time.Time) Attr         { return attr.Timep(key, val) }
func Durationp(key string, val *time.Duration) Attr { return attr.Durationp(key, val) }

func Times(key string, val []time.Time) Attr         { return attr.Times(key, val) }
func Durations(key string, val []time.Duration) Attr { return attr.Durations(key, val) }

func Error(err error) Attr                  { return attr.Error(err) }
func Errors(key string, err []error) Attr   { return attr.Errors(key, err) }
func NamedError(key string, err error) Attr { return attr.NamedError(key, err) }

func Stack(key string) Attr               { return attr.StackSkip(key, 1) }
func StackSkip(key string, skip int) Attr { return attr.StackSkip(key, skip+1) }

func Any(key string, val any) Attr     { return attr.Any(key, val) }
func Reflect(key string, val any) Attr { return attr.Any(key, val) }

func Integer[Int ~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr,
](
	key string, val Int,
) Attr {
	if val >= 0 {
		return Uint64(key, uint64(val))
	}
	return Int64(key, int64(val))
}

func Float[T ~float32 | ~float64](key string, val T) Attr { return Float64(key, float64(val)) }

func BriefString(key string, val string) Attr         { return attr.BriefString(key, val) }
func BriefByteString(key string, val []byte) Attr     { return attr.BriefByteString(key, val) }
func BriefStringer(key string, val fmt.Stringer) Attr { return attr.BriefStringer(key, val) }

func BriefStringp(key string, val *string) Attr               { return attr.BriefStringp(key, val) }
func BriefStrings(key string, val []string) Attr              { return attr.BriefStrings(key, val) }
func BriefByteStrings(key string, val [][]byte) Attr          { return attr.BriefByteStrings(key, val) }
func BriefStringers[T fmt.Stringer](key string, val []T) Attr { return attr.BriefStringers(key, val) }

func BriefReflect(key string, val any) Attr { return attr.BriefReflect(key, val) }
