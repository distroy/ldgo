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
		{To: reflect.Float32, From: reflect.Invalid}:    copyReflectToFloatFromInvalid,
		{To: reflect.Float32, From: reflect.Bool}:       copyReflectToFloatFromBool,
		{To: reflect.Float32, From: reflect.Complex64}:  copyReflectToFloatFromComplex,
		{To: reflect.Float32, From: reflect.Complex128}: copyReflectToFloatFromComplex,
		{To: reflect.Float32, From: reflect.Float32}:    copyReflectToFloatFromFloat,
		{To: reflect.Float32, From: reflect.Float64}:    copyReflectToFloatFromFloat,
		{To: reflect.Float32, From: reflect.Int}:        copyReflectToFloatFromInt,
		{To: reflect.Float32, From: reflect.Int8}:       copyReflectToFloatFromInt,
		{To: reflect.Float32, From: reflect.Int16}:      copyReflectToFloatFromInt,
		{To: reflect.Float32, From: reflect.Int32}:      copyReflectToFloatFromInt,
		{To: reflect.Float32, From: reflect.Int64}:      copyReflectToFloatFromInt,
		{To: reflect.Float32, From: reflect.Uint}:       copyReflectToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uint8}:      copyReflectToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uint16}:     copyReflectToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uint32}:     copyReflectToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uint64}:     copyReflectToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uintptr}:    copyReflectToFloatFromUint,
		{To: reflect.Float32, From: reflect.String}:     copyReflectToFloatFromString,

		{To: reflect.Float64, From: reflect.Invalid}:    copyReflectToFloatFromInvalid,
		{To: reflect.Float64, From: reflect.Bool}:       copyReflectToFloatFromBool,
		{To: reflect.Float64, From: reflect.Complex64}:  copyReflectToFloatFromComplex,
		{To: reflect.Float64, From: reflect.Complex128}: copyReflectToFloatFromComplex,
		{To: reflect.Float64, From: reflect.Float32}:    copyReflectToFloatFromFloat,
		{To: reflect.Float64, From: reflect.Float64}:    copyReflectToFloatFromFloat,
		{To: reflect.Float64, From: reflect.Int}:        copyReflectToFloatFromInt,
		{To: reflect.Float64, From: reflect.Int8}:       copyReflectToFloatFromInt,
		{To: reflect.Float64, From: reflect.Int16}:      copyReflectToFloatFromInt,
		{To: reflect.Float64, From: reflect.Int32}:      copyReflectToFloatFromInt,
		{To: reflect.Float64, From: reflect.Int64}:      copyReflectToFloatFromInt,
		{To: reflect.Float64, From: reflect.Uint}:       copyReflectToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uint8}:      copyReflectToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uint16}:     copyReflectToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uint32}:     copyReflectToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uint64}:     copyReflectToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uintptr}:    copyReflectToFloatFromUint,
		{To: reflect.Float64, From: reflect.String}:     copyReflectToFloatFromString,
	})
}

func copyReflectToFloatFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.SetFloat(0)
	return true
}

func copyReflectToFloatFromBool(c *copyContext, target, source reflect.Value) bool {
	b := source.Bool()
	target.SetFloat(float64(bool2int(b)))
	return true
}

func copyReflectToFloatFromFloat(c *copyContext, target, source reflect.Value) bool {
	n := source.Float()
	target.SetFloat(n)
	if target.OverflowFloat(n) {
		c.AddErrorf("%s(%f) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToFloatFromComplex(c *copyContext, target, source reflect.Value) bool {
	n := source.Complex()
	r := real(n)
	target.SetFloat(r)
	if target.OverflowFloat(r) {
		c.AddErrorf("%s(%v) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToFloatFromInt(c *copyContext, target, source reflect.Value) bool {
	n := source.Int()
	target.SetFloat(float64(n))
	if target.OverflowFloat(float64(n)) {
		c.AddErrorf("%s(%d) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToFloatFromUint(c *copyContext, target, source reflect.Value) bool {
	n := source.Uint()
	target.SetFloat(float64(n))
	if target.OverflowFloat(float64(n)) {
		c.AddErrorf("%s(%d) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToFloatFromString(c *copyContext, target, source reflect.Value) bool {
	s := source.String()
	n, err := strconv.ParseFloat(s, 64)
	target.SetFloat(n)
	if err != nil {
		c.AddErrorf("can not convert to %s, %q", target.Type().String(), s)

	} else if target.OverflowFloat(n) {
		c.AddErrorf("%s(%s) overflow", target.Type().String(), s)
	}
	return true
}
