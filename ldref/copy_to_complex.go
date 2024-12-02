/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"strconv"
)

func init() {
	registerCopyFunc(map[copyPair]copyFuncType{
		{To: reflect.Complex64, From: reflect.Invalid}:    copyReflectToComplexFromInvalid,
		{To: reflect.Complex64, From: reflect.Bool}:       copyReflectToComplexFromBool,
		{To: reflect.Complex64, From: reflect.Complex64}:  copyReflectToComplexFromComplex,
		{To: reflect.Complex64, From: reflect.Complex128}: copyReflectToComplexFromComplex,
		{To: reflect.Complex64, From: reflect.Float32}:    copyReflectToComplexFromFloat,
		{To: reflect.Complex64, From: reflect.Float64}:    copyReflectToComplexFromFloat,
		{To: reflect.Complex64, From: reflect.Int}:        copyReflectToComplexFromInt,
		{To: reflect.Complex64, From: reflect.Int8}:       copyReflectToComplexFromInt,
		{To: reflect.Complex64, From: reflect.Int16}:      copyReflectToComplexFromInt,
		{To: reflect.Complex64, From: reflect.Int32}:      copyReflectToComplexFromInt,
		{To: reflect.Complex64, From: reflect.Int64}:      copyReflectToComplexFromInt,
		{To: reflect.Complex64, From: reflect.Uint}:       copyReflectToComplexFromUint,
		{To: reflect.Complex64, From: reflect.Uint8}:      copyReflectToComplexFromUint,
		{To: reflect.Complex64, From: reflect.Uint16}:     copyReflectToComplexFromUint,
		{To: reflect.Complex64, From: reflect.Uint32}:     copyReflectToComplexFromUint,
		{To: reflect.Complex64, From: reflect.Uint64}:     copyReflectToComplexFromUint,
		{To: reflect.Complex64, From: reflect.Uintptr}:    copyReflectToComplexFromUint,
		{To: reflect.Complex64, From: reflect.String}:     copyReflectToComplexFromString,

		{To: reflect.Complex128, From: reflect.Invalid}:    copyReflectToComplexFromInvalid,
		{To: reflect.Complex128, From: reflect.Bool}:       copyReflectToComplexFromBool,
		{To: reflect.Complex128, From: reflect.Complex64}:  copyReflectToComplexFromComplex,
		{To: reflect.Complex128, From: reflect.Complex128}: copyReflectToComplexFromComplex,
		{To: reflect.Complex128, From: reflect.Float32}:    copyReflectToComplexFromFloat,
		{To: reflect.Complex128, From: reflect.Float64}:    copyReflectToComplexFromFloat,
		{To: reflect.Complex128, From: reflect.Int}:        copyReflectToComplexFromInt,
		{To: reflect.Complex128, From: reflect.Int8}:       copyReflectToComplexFromInt,
		{To: reflect.Complex128, From: reflect.Int16}:      copyReflectToComplexFromInt,
		{To: reflect.Complex128, From: reflect.Int32}:      copyReflectToComplexFromInt,
		{To: reflect.Complex128, From: reflect.Int64}:      copyReflectToComplexFromInt,
		{To: reflect.Complex128, From: reflect.Uint}:       copyReflectToComplexFromUint,
		{To: reflect.Complex128, From: reflect.Uint8}:      copyReflectToComplexFromUint,
		{To: reflect.Complex128, From: reflect.Uint16}:     copyReflectToComplexFromUint,
		{To: reflect.Complex128, From: reflect.Uint32}:     copyReflectToComplexFromUint,
		{To: reflect.Complex128, From: reflect.Uint64}:     copyReflectToComplexFromUint,
		{To: reflect.Complex128, From: reflect.Uintptr}:    copyReflectToComplexFromUint,
		{To: reflect.Complex128, From: reflect.String}:     copyReflectToComplexFromString,
	})
}

func copyReflectToComplexFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.SetComplex(0)
	return true
}

func copyReflectToComplexFromBool(c *copyContext, target, source reflect.Value) bool {
	b := source.Bool()
	target.SetComplex(complex(float64(bool2int(b)), 0))
	return true
}

func copyReflectToComplexFromComplex(c *copyContext, target, source reflect.Value) bool {
	n := source.Complex()
	target.SetComplex(n)
	if target.OverflowComplex(n) {
		c.AddErrorf("%s(%v) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToComplexFromFloat(c *copyContext, target, source reflect.Value) bool {
	n := source.Float()
	x := complex(n, 0)
	target.SetComplex(x)
	if target.OverflowComplex(x) {
		c.AddErrorf("%s(%f) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToComplexFromInt(c *copyContext, target, source reflect.Value) bool {
	n := source.Int()
	x := complex(float64(n), 0)
	target.SetComplex(x)
	if target.OverflowComplex(x) {
		c.AddErrorf("%s(%d) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToComplexFromUint(c *copyContext, target, source reflect.Value) bool {
	n := source.Uint()
	x := complex(float64(n), 0)
	target.SetComplex(x)
	if target.OverflowComplex(x) {
		c.AddErrorf("%s(%d) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToComplexFromString(c *copyContext, target, source reflect.Value) bool {
	s := source.String()
	x, err := strconv.ParseComplex(s, 128)
	if err != nil {
		c.AddErrorf("can not convert to %s, %q", target.Type().String(), s)
		return false
	}
	target.SetComplex(x)
	if target.OverflowComplex(x) {
		c.AddErrorf("%s(%s) overflow", target.Type().String(), s)
	}
	return true
}
