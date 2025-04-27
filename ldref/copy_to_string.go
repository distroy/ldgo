/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"strconv"
	"unsafe"

	"github.com/distroy/ldgo/v3/ldconv"
)

func typeNameOfReflect(v reflect.Value) string {
	if v.Kind() == reflect.Invalid {
		return "nil"
	}

	return v.Type().String()
}
func typeNameOfReflectType(t reflect.Type) string {
	if t.Kind() == reflect.Invalid {
		return "nil"
	}

	return t.String()
}
func refTypeOfValue(v reflect.Value) reflect.Type {
	if v.Kind() != reflect.Invalid {
		return v.Type()
	}

	return typeOfNil
}
func refKindOfType(t reflect.Type) reflect.Kind {
	if t == nil {
		return reflect.Invalid
	}
	return t.Kind()
}

func refStructFieldByCopyFieldInfo(v reflect.Value, info *copyFieldInfo) reflect.Value {
	return refStructField(v, info.Index, &info.StructField)
}
func refStructField(v reflect.Value, index int, info *reflect.StructField) reflect.Value {
	if info.IsExported() {
		return v.Field(index)
	}
	addr := unsafe.Pointer(v.Field(index).UnsafeAddr())
	return reflect.NewAt(info.Type, addr).Elem()
}

func getCopyFuncSetZero(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	zero := reflect.Zero(tTyp)
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		target.Set(zero)
		return true
	}
}

func init() {
	registerCopyFunc(map[copyPair]copyFuncType{
		{To: reflect.String, From: reflect.Invalid}: copyReflectToStringFromInvalid,
		{To: reflect.String, From: reflect.Bool}:    copyReflectToStringFromBool,
		{To: reflect.String, From: reflect.Float32}: copyReflectToStringFromFloat,
		{To: reflect.String, From: reflect.Float64}: copyReflectToStringFromFloat,
		{To: reflect.String, From: reflect.Int}:     copyReflectToStringFromInt,
		{To: reflect.String, From: reflect.Int8}:    copyReflectToStringFromInt,
		{To: reflect.String, From: reflect.Int16}:   copyReflectToStringFromInt,
		{To: reflect.String, From: reflect.Int32}:   copyReflectToStringFromInt,
		{To: reflect.String, From: reflect.Int64}:   copyReflectToStringFromInt,
		{To: reflect.String, From: reflect.Uint}:    copyReflectToStringFromUint,
		{To: reflect.String, From: reflect.Uint8}:   copyReflectToStringFromUint,
		{To: reflect.String, From: reflect.Uint16}:  copyReflectToStringFromUint,
		{To: reflect.String, From: reflect.Uint32}:  copyReflectToStringFromUint,
		{To: reflect.String, From: reflect.Uint64}:  copyReflectToStringFromUint,
		{To: reflect.String, From: reflect.Uintptr}: copyReflectToStringFromUint,
		{To: reflect.String, From: reflect.String}:  copyReflectToStringFromString,
		{To: reflect.String, From: reflect.Array}:   copyReflectToStringFromArray,
		{To: reflect.String, From: reflect.Slice}:   copyReflectToStringFromSlice,

		{To: reflect.String, From: reflect.Complex64}:  copyReflectToStringFromComplex,
		{To: reflect.String, From: reflect.Complex128}: copyReflectToStringFromComplex,
	})
	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.String, From: reflect.Invalid}: getCopyFuncToStringFromInvalid,
		{To: reflect.String, From: reflect.Bool}:    getCopyFuncToStringFromBool,
		{To: reflect.String, From: reflect.Float32}: getCopyFuncToStringFromFloat,
		{To: reflect.String, From: reflect.Float64}: getCopyFuncToStringFromFloat,
		{To: reflect.String, From: reflect.Int}:     getCopyFuncToStringFromInt,
		{To: reflect.String, From: reflect.Int8}:    getCopyFuncToStringFromInt,
		{To: reflect.String, From: reflect.Int16}:   getCopyFuncToStringFromInt,
		{To: reflect.String, From: reflect.Int32}:   getCopyFuncToStringFromInt,
		{To: reflect.String, From: reflect.Int64}:   getCopyFuncToStringFromInt,
		{To: reflect.String, From: reflect.Uint}:    getCopyFuncToStringFromUint,
		{To: reflect.String, From: reflect.Uint8}:   getCopyFuncToStringFromUint,
		{To: reflect.String, From: reflect.Uint16}:  getCopyFuncToStringFromUint,
		{To: reflect.String, From: reflect.Uint32}:  getCopyFuncToStringFromUint,
		{To: reflect.String, From: reflect.Uint64}:  getCopyFuncToStringFromUint,
		{To: reflect.String, From: reflect.Uintptr}: getCopyFuncToStringFromUint,
		{To: reflect.String, From: reflect.String}:  getCopyFuncToStringFromString,
		{To: reflect.String, From: reflect.Array}:   getCopyFuncToStringFromArray,
		{To: reflect.String, From: reflect.Slice}:   getCopyFuncToStringFromSlice,

		{To: reflect.String, From: reflect.Complex64}:  getCopyFuncToStringFromComplex,
		{To: reflect.String, From: reflect.Complex128}: getCopyFuncToStringFromComplex,
	})
}

func copyReflectToStringFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.SetString("")
	return true
}
func getCopyFuncToStringFromInvalid(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		target.SetString("")
		return true
	}
}

func copyReflectToStringFromBool(c *copyContext, target, source reflect.Value) bool {
	b := source.Bool()
	if b {
		target.SetString("true")
	} else {
		target.SetString("false")
	}
	return true
}
func getCopyFuncToStringFromBool(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return copyReflectToStringFromBool
}

func copyReflectToStringFromFloat(c *copyContext, target, source reflect.Value) bool {
	n := source.Float()
	target.SetString(strconv.FormatFloat(n, 'f', -1, 64))
	return true
}
func getCopyFuncToStringFromFloat(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		n := source.Float()
		target.SetString(strconv.FormatFloat(n, 'f', -1, 64))
		return true
	}
}

func copyReflectToStringFromComplex(c *copyContext, target, source reflect.Value) bool {
	n := source.Complex()
	target.SetString(strconv.FormatComplex(n, 'f', -1, 128))
	return true
}
func getCopyFuncToStringFromComplex(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		n := source.Complex()
		target.SetString(strconv.FormatComplex(n, 'f', -1, 128))
		return true
	}
}

func copyReflectToStringFromInt(c *copyContext, target, source reflect.Value) bool {
	n := source.Int()
	target.SetString(strconv.FormatInt(n, 10))
	return true
}
func getCopyFuncToStringFromInt(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		n := source.Int()
		target.SetString(strconv.FormatInt(n, 10))
		return true
	}
}

func copyReflectToStringFromUint(c *copyContext, target, source reflect.Value) bool {
	n := source.Uint()
	target.SetString(strconv.FormatUint(n, 10))
	return true
}
func getCopyFuncToStringFromUint(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		n := source.Uint()
		target.SetString(strconv.FormatUint(n, 10))
		return true
	}
}

func copyReflectToStringFromString(c *copyContext, target, source reflect.Value) bool {
	s := source.String()
	target.SetString(s)
	return true
}
func getCopyFuncToStringFromString(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		s := source.String()
		target.SetString(s)
		return true
	}
}

func copyReflectToStringFromArray(c *copyContext, target, source reflect.Value) bool {
	// sVal := source.Slice(0, source.Len())
	sVal := reflectArrayToSlice(source)
	// return copyReflectToStringFromSlice(c, target, sVal)
	switch ss := sVal.Interface().(type) {
	default:
		return false

	case []byte:
		target.SetString(ldconv.BytesToStrUnsafe(ss))

	case []rune:
		target.SetString(string(ss))
	}
	return true
}
func getCopyFuncToStringFromArray(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	switch sTyp.Elem() {
	default:
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }

	case typeOfByteSlice.Elem():
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			sVal := reflectArrayToSlice(source)
			vv := sVal.Interface().([]byte)
			target.SetString(string(vv))
			return true
		}

	case typeOfRuneSlice.Elem():
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			sVal := reflectArrayToSlice(source)
			vv := sVal.Interface().([]rune)
			target.SetString(string(vv))
			return true
		}
	}
}

func copyReflectToStringFromSlice(c *copyContext, target, source reflect.Value) bool {
	switch ss := source.Interface().(type) {
	default:
		return false

	case []byte:
		target.SetString(string(ss))

	case []rune:
		target.SetString(string(ss))
	}
	return true
}
func getCopyFuncToStringFromSlice(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	switch sTyp {
	default:
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }

	case typeOfByteSlice:
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			vv := source.Interface().([]byte)
			target.SetString(string(vv))
			return true
		}

	case typeOfRuneSlice:
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			vv := source.Interface().([]rune)
			target.SetString(string(vv))
			return true
		}
	}
}
