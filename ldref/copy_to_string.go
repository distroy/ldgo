/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"strconv"

	"github.com/distroy/ldgo/v2/ldconv"
)

func typeNameOfReflect(v reflect.Value) string {
	if v.Kind() == reflect.Invalid {
		return "nil"
	}

	return v.Type().String()
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
}

func copyReflectToStringFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.SetString("")
	return true
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

func copyReflectToStringFromFloat(c *copyContext, target, source reflect.Value) bool {
	n := source.Float()
	target.SetString(strconv.FormatFloat(n, 'f', -1, 64))
	return true
}

func copyReflectToStringFromComplex(c *copyContext, target, source reflect.Value) bool {
	n := source.Complex()
	target.SetString(strconv.FormatComplex(n, 'f', -1, 128))
	return true
}

func copyReflectToStringFromInt(c *copyContext, target, source reflect.Value) bool {
	n := source.Int()
	target.SetString(strconv.FormatInt(n, 10))
	return true
}

func copyReflectToStringFromUint(c *copyContext, target, source reflect.Value) bool {
	n := source.Uint()
	target.SetString(strconv.FormatUint(n, 10))
	return true
}

func copyReflectToStringFromString(c *copyContext, target, source reflect.Value) bool {
	s := source.String()
	target.SetString(s)
	return true
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
