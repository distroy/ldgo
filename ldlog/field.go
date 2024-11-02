/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"fmt"
	"time"

	"github.com/distroy/ldgo/v3/ldlog/internal/field"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Field = zapcore.Field

	ObjectMarshaler = zapcore.ObjectMarshaler
	ArrayMarshaler  = zapcore.ArrayMarshaler
)

func Bool(key string, val bool) Field             { return zap.Bool(key, val) }
func String(key string, val string) Field         { return zap.String(key, val) }
func Stringer(key string, val fmt.Stringer) Field { return zap.Stringer(key, val) }
func Binary(key string, val []byte) Field         { return zap.Binary(key, val) }
func ByteString(key string, val []byte) Field     { return zap.ByteString(key, val) }

func Int(key string, val int) Field     { return zap.Int(key, val) }
func Int8(key string, val int8) Field   { return zap.Int8(key, val) }
func Int16(key string, val int16) Field { return zap.Int16(key, val) }
func Int32(key string, val int32) Field { return zap.Int32(key, val) }
func Int64(key string, val int64) Field { return zap.Int64(key, val) }

func Uint(key string, val uint) Field       { return zap.Uint(key, val) }
func Uint8(key string, val uint8) Field     { return zap.Uint8(key, val) }
func Uint16(key string, val uint16) Field   { return zap.Uint16(key, val) }
func Uint32(key string, val uint32) Field   { return zap.Uint32(key, val) }
func Uint64(key string, val uint64) Field   { return zap.Uint64(key, val) }
func Uintptr(key string, val uintptr) Field { return zap.Uintptr(key, val) }

func Float32(key string, val float32) Field { return zap.Float32(key, val) }
func Float64(key string, val float64) Field { return zap.Float64(key, val) }

func Complex64(key string, val complex64) Field   { return zap.Complex64(key, val) }
func Complex128(key string, val complex128) Field { return zap.Complex128(key, val) }

func Boolp(key string, val *bool) Field                   { return zap.Boolp(key, val) }
func Bools(key string, val []bool) Field                  { return zap.Bools(key, val) }
func Stringp(key string, val *string) Field               { return zap.Stringp(key, val) }
func Strings(key string, val []string) Field              { return zap.Strings(key, val) }
func Stringers[T fmt.Stringer](key string, val []T) Field { return zap.Stringers(key, val) }
func ByteStrings(key string, val [][]byte) Field          { return zap.ByteStrings(key, val) }

func Intp(key string, val *int) Field     { return zap.Intp(key, val) }
func Int8p(key string, val *int8) Field   { return zap.Int8p(key, val) }
func Int16p(key string, val *int16) Field { return zap.Int16p(key, val) }
func Int32p(key string, val *int32) Field { return zap.Int32p(key, val) }
func Int64p(key string, val *int64) Field { return zap.Int64p(key, val) }

func Uintp(key string, val *uint) Field       { return zap.Uintp(key, val) }
func Uint8p(key string, val *uint8) Field     { return zap.Uint8p(key, val) }
func Uint16p(key string, val *uint16) Field   { return zap.Uint16p(key, val) }
func Uint32p(key string, val *uint32) Field   { return zap.Uint32p(key, val) }
func Uint64p(key string, val *uint64) Field   { return zap.Uint64p(key, val) }
func Uintptrp(key string, val *uintptr) Field { return zap.Uintptrp(key, val) }

func Float32p(key string, val *float32) Field { return zap.Float32p(key, val) }
func Float64p(key string, val *float64) Field { return zap.Float64p(key, val) }

func Complex64p(key string, val *complex64) Field   { return zap.Complex64p(key, val) }
func Complex128p(key string, val *complex128) Field { return zap.Complex128p(key, val) }

func Ints(key string, val []int) Field     { return zap.Ints(key, val) }
func Int8s(key string, val []int8) Field   { return zap.Int8s(key, val) }
func Int16s(key string, val []int16) Field { return zap.Int16s(key, val) }
func Int32s(key string, val []int32) Field { return zap.Int32s(key, val) }
func Int64s(key string, val []int64) Field { return zap.Int64s(key, val) }

func Uints(key string, val []uint) Field       { return zap.Uints(key, val) }
func Uint8s(key string, val []uint8) Field     { return zap.Uint8s(key, val) }
func Uint16s(key string, val []uint16) Field   { return zap.Uint16s(key, val) }
func Uint32s(key string, val []uint32) Field   { return zap.Uint32s(key, val) }
func Uint64s(key string, val []uint64) Field   { return zap.Uint64s(key, val) }
func Uintptrs(key string, val []uintptr) Field { return zap.Uintptrs(key, val) }

func Float32s(key string, val []float32) Field { return zap.Float32s(key, val) }
func Float64s(key string, val []float64) Field { return zap.Float64s(key, val) }

func Complex64s(key string, val []complex64) Field   { return zap.Complex64s(key, val) }
func Complex128s(key string, val []complex128) Field { return zap.Complex128s(key, val) }

func Time(key string, val time.Time) Field         { return zap.Time(key, val) }
func Duration(key string, val time.Duration) Field { return zap.Duration(key, val) }

func Timep(key string, val *time.Time) Field         { return zap.Timep(key, val) }
func Durationp(key string, val *time.Duration) Field { return zap.Durationp(key, val) }

func Times(key string, val []time.Time) Field         { return zap.Times(key, val) }
func Durations(key string, val []time.Duration) Field { return zap.Durations(key, val) }

func Error(err error) Field                  { return zap.Error(err) }
func Errors(key string, err []error) Field   { return zap.Errors(key, err) }
func NamedError(key string, err error) Field { return zap.NamedError(key, err) }

func Stack(key string) Field               { return zap.StackSkip(key, 1) }
func StackSkip(key string, skip int) Field { return zap.StackSkip(key, skip+1) }

func Namespace(key string) Field       { return zap.Namespace(key) }
func Inline(val ObjectMarshaler) Field { return zap.Inline(val) }

func Array(key string, val ArrayMarshaler) Field           { return zap.Array(key, val) }
func Object(key string, val ObjectMarshaler) Field         { return zap.Object(key, val) }
func Objects[T ObjectMarshaler](key string, val []T) Field { return zap.Objects(key, val) }

func Any(key string, val interface{}) Field     { return zap.Any(key, val) }
func Reflect(key string, val interface{}) Field { return zap.Reflect(key, val) }

func Skip() Field { return zap.Skip() }

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~uintptr
}

func Integer[Int integer](key string, val Int) Field {
	if val >= 0 {
		return Uint64(key, uint64(val))
	}
	return Int64(key, int64(val))
}

func BriefString(key string, val string) Field         { return field.BriefString(key, val) }
func BriefByteString(key string, val []byte) Field     { return field.BriefByteString(key, val) }
func BriefStringer(key string, val fmt.Stringer) Field { return field.BriefStringer(key, val) }

func BriefStringp(key string, val *string) Field               { return field.BriefStringp(key, val) }
func BriefStrings(key string, val []string) Field              { return field.BriefStrings(key, val) }
func BriefByteStrings(key string, val [][]byte) Field          { return field.BriefByteStrings(key, val) }
func BriefStringers[T fmt.Stringer](key string, val []T) Field { return field.BriefStringers(key, val) }

func BriefReflect(key string, val interface{}) Field { return field.BriefReflect(key, val) }
