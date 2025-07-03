/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"fmt"
	"time"

	"github.com/distroy/ldgo/v3/ldlog/internal/_attr"
)

func Skip() Attr { return _attr.Skip() }

func Bool(key string, val bool) Attr             { return _attr.Bool(key, val) }
func String(key string, val string) Attr         { return _attr.String(key, val) }
func Stringer(key string, val fmt.Stringer) Attr { return _attr.Stringer(key, val) }
func Binary(key string, val []byte) Attr         { return _attr.Binary(key, val) }
func ByteString(key string, val []byte) Attr     { return _attr.ByteString(key, val) }

func Int(key string, val int) Attr     { return _attr.Int(key, val) }
func Int8(key string, val int8) Attr   { return _attr.Int64(key, int64(val)) }
func Int16(key string, val int16) Attr { return _attr.Int64(key, int64(val)) }
func Int32(key string, val int32) Attr { return _attr.Int64(key, int64(val)) }
func Int64(key string, val int64) Attr { return _attr.Int64(key, val) }

func Uint(key string, val uint) Attr       { return _attr.Uint64(key, uint64(val)) }
func Uint8(key string, val uint8) Attr     { return _attr.Uint64(key, uint64(val)) }
func Uint16(key string, val uint16) Attr   { return _attr.Uint64(key, uint64(val)) }
func Uint32(key string, val uint32) Attr   { return _attr.Uint64(key, uint64(val)) }
func Uint64(key string, val uint64) Attr   { return _attr.Uint64(key, val) }
func Uintptr(key string, val uintptr) Attr { return _attr.Uint64(key, uint64(val)) }

func Float32(key string, val float32) Attr { return _attr.Float64(key, float64(val)) }
func Float64(key string, val float64) Attr { return _attr.Float64(key, val) }

func Complex64(key string, val complex64) Attr   { return _attr.Complex64(key, val) }
func Complex128(key string, val complex128) Attr { return _attr.Complex128(key, val) }

func Boolp(key string, val *bool) Attr                   { return _attr.Boolp(key, val) }
func Bools(key string, val []bool) Attr                  { return _attr.Bools(key, val) }
func Stringp(key string, val *string) Attr               { return _attr.Stringp(key, val) }
func Strings(key string, val []string) Attr              { return _attr.Strings(key, val) }
func Stringers[T fmt.Stringer](key string, val []T) Attr { return _attr.Stringers(key, val) }
func ByteStrings(key string, val [][]byte) Attr          { return _attr.ByteStrings(key, val) }

func Intp(key string, val *int) Attr     { return _attr.Intp(key, val) }
func Int8p(key string, val *int8) Attr   { return _attr.Int8p(key, val) }
func Int16p(key string, val *int16) Attr { return _attr.Int16p(key, val) }
func Int32p(key string, val *int32) Attr { return _attr.Int32p(key, val) }
func Int64p(key string, val *int64) Attr { return _attr.Int64p(key, val) }

func Uintp(key string, val *uint) Attr       { return _attr.Uintp(key, val) }
func Uint8p(key string, val *uint8) Attr     { return _attr.Uint8p(key, val) }
func Uint16p(key string, val *uint16) Attr   { return _attr.Uint16p(key, val) }
func Uint32p(key string, val *uint32) Attr   { return _attr.Uint32p(key, val) }
func Uint64p(key string, val *uint64) Attr   { return _attr.Uint64p(key, val) }
func Uintptrp(key string, val *uintptr) Attr { return _attr.Uintptrp(key, val) }

func Float32p(key string, val *float32) Attr { return _attr.Float32p(key, val) }
func Float64p(key string, val *float64) Attr { return _attr.Float64p(key, val) }

func Complex64p(key string, val *complex64) Attr   { return _attr.Complex64p(key, val) }
func Complex128p(key string, val *complex128) Attr { return _attr.Complex128p(key, val) }

func Ints(key string, val []int) Attr     { return _attr.Ints(key, val) }
func Int8s(key string, val []int8) Attr   { return _attr.Int8s(key, val) }
func Int16s(key string, val []int16) Attr { return _attr.Int16s(key, val) }
func Int32s(key string, val []int32) Attr { return _attr.Int32s(key, val) }
func Int64s(key string, val []int64) Attr { return _attr.Int64s(key, val) }

func Uints(key string, val []uint) Attr       { return _attr.Uints(key, val) }
func Uint8s(key string, val []uint8) Attr     { return _attr.Uint8s(key, val) }
func Uint16s(key string, val []uint16) Attr   { return _attr.Uint16s(key, val) }
func Uint32s(key string, val []uint32) Attr   { return _attr.Uint32s(key, val) }
func Uint64s(key string, val []uint64) Attr   { return _attr.Uint64s(key, val) }
func Uintptrs(key string, val []uintptr) Attr { return _attr.Uintptrs(key, val) }

func Float32s(key string, val []float32) Attr { return _attr.Float32s(key, val) }
func Float64s(key string, val []float64) Attr { return _attr.Float64s(key, val) }

func Complex64s(key string, val []complex64) Attr   { return _attr.Complex64s(key, val) }
func Complex128s(key string, val []complex128) Attr { return _attr.Complex128s(key, val) }

func Time(key string, val time.Time) Attr         { return _attr.Time(key, val) }
func Duration(key string, val time.Duration) Attr { return _attr.Duration(key, val) }

func Timep(key string, val *time.Time) Attr         { return _attr.Timep(key, val) }
func Durationp(key string, val *time.Duration) Attr { return _attr.Durationp(key, val) }

func Times(key string, val []time.Time) Attr         { return _attr.Times(key, val) }
func Durations(key string, val []time.Duration) Attr { return _attr.Durations(key, val) }

func Error(err error) Attr                  { return _attr.Error(err) }
func Errors(key string, err []error) Attr   { return _attr.Errors(key, err) }
func NamedError(key string, err error) Attr { return _attr.NamedError(key, err) }

func Stack(key string) Attr               { return _attr.StackSkip(key, 1) }
func StackSkip(key string, skip int) Attr { return _attr.StackSkip(key, skip+1) }

func Any(key string, val any) Attr     { return _attr.Any(key, val) }
func Reflect(key string, val any) Attr { return _attr.Any(key, val) }

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

func BriefString(key string, val string) Attr         { return _attr.BriefString(key, val) }
func BriefByteString(key string, val []byte) Attr     { return _attr.BriefByteString(key, val) }
func BriefStringer(key string, val fmt.Stringer) Attr { return _attr.BriefStringer(key, val) }

func BriefStringp(key string, val *string) Attr               { return _attr.BriefStringp(key, val) }
func BriefStrings(key string, val []string) Attr              { return _attr.BriefStrings(key, val) }
func BriefByteStrings(key string, val [][]byte) Attr          { return _attr.BriefByteStrings(key, val) }
func BriefStringers[T fmt.Stringer](key string, val []T) Attr { return _attr.BriefStringers(key, val) }

func BriefReflect(key string, val any) Attr { return _attr.BriefReflect(key, val) }
