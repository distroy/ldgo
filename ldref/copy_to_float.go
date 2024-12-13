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
	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.Float32, From: reflect.Invalid}:    getCopyFuncToFloatFromInvalid,
		{To: reflect.Float32, From: reflect.Bool}:       getCopyFuncToFloatFromBool,
		{To: reflect.Float32, From: reflect.Complex64}:  getCopyFuncToFloatFromComplex,
		{To: reflect.Float32, From: reflect.Complex128}: getCopyFuncToFloatFromComplex,
		{To: reflect.Float32, From: reflect.Float32}:    getCopyFuncToFloatFromFloat,
		{To: reflect.Float32, From: reflect.Float64}:    getCopyFuncToFloatFromFloat,
		{To: reflect.Float32, From: reflect.Int}:        getCopyFuncToFloatFromInt,
		{To: reflect.Float32, From: reflect.Int8}:       getCopyFuncToFloatFromInt,
		{To: reflect.Float32, From: reflect.Int16}:      getCopyFuncToFloatFromInt,
		{To: reflect.Float32, From: reflect.Int32}:      getCopyFuncToFloatFromInt,
		{To: reflect.Float32, From: reflect.Int64}:      getCopyFuncToFloatFromInt,
		{To: reflect.Float32, From: reflect.Uint}:       getCopyFuncToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uint8}:      getCopyFuncToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uint16}:     getCopyFuncToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uint32}:     getCopyFuncToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uint64}:     getCopyFuncToFloatFromUint,
		{To: reflect.Float32, From: reflect.Uintptr}:    getCopyFuncToFloatFromUint,
		{To: reflect.Float32, From: reflect.String}:     getCopyFuncToFloatFromString,

		{To: reflect.Float64, From: reflect.Invalid}:    getCopyFuncToFloatFromInvalid,
		{To: reflect.Float64, From: reflect.Bool}:       getCopyFuncToFloatFromBool,
		{To: reflect.Float64, From: reflect.Complex64}:  getCopyFuncToFloatFromComplex,
		{To: reflect.Float64, From: reflect.Complex128}: getCopyFuncToFloatFromComplex,
		{To: reflect.Float64, From: reflect.Float32}:    getCopyFuncToFloatFromFloat,
		{To: reflect.Float64, From: reflect.Float64}:    getCopyFuncToFloatFromFloat,
		{To: reflect.Float64, From: reflect.Int}:        getCopyFuncToFloatFromInt,
		{To: reflect.Float64, From: reflect.Int8}:       getCopyFuncToFloatFromInt,
		{To: reflect.Float64, From: reflect.Int16}:      getCopyFuncToFloatFromInt,
		{To: reflect.Float64, From: reflect.Int32}:      getCopyFuncToFloatFromInt,
		{To: reflect.Float64, From: reflect.Int64}:      getCopyFuncToFloatFromInt,
		{To: reflect.Float64, From: reflect.Uint}:       getCopyFuncToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uint8}:      getCopyFuncToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uint16}:     getCopyFuncToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uint32}:     getCopyFuncToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uint64}:     getCopyFuncToFloatFromUint,
		{To: reflect.Float64, From: reflect.Uintptr}:    getCopyFuncToFloatFromUint,
		{To: reflect.Float64, From: reflect.String}:     getCopyFuncToFloatFromString,
	})
}

func copyReflectToFloatFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.SetFloat(0)
	return true
}
func getCopyFuncToFloatFromInvalid(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		target.SetFloat(0)
		return true
	}
}

func copyReflectToFloatFromBool(c *copyContext, target, source reflect.Value) bool {
	b := source.Bool()
	target.SetFloat(float64(bool2int(b)))
	return true
}
func getCopyFuncToFloatFromBool(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return copyReflectToFloatFromBool
}

func copyReflectToFloatFromFloat(c *copyContext, target, source reflect.Value) bool {
	n := source.Float()
	target.SetFloat(n)
	if target.OverflowFloat(n) {
		c.AddErrorf("%s(%f) overflow", target.Type().String(), n)
	}
	return true
}
func getCopyFuncToFloatFromFloat(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return copyReflectToFloatFromFloat
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
func getCopyFuncToFloatFromComplex(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return copyReflectToFloatFromComplex
}

func copyReflectToFloatFromInt(c *copyContext, target, source reflect.Value) bool {
	n := source.Int()
	target.SetFloat(float64(n))
	if target.OverflowFloat(float64(n)) {
		c.AddErrorf("%s(%d) overflow", target.Type().String(), n)
	}
	return true
}
func getCopyFuncToFloatFromInt(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return copyReflectToFloatFromInt
}

func copyReflectToFloatFromUint(c *copyContext, target, source reflect.Value) bool {
	n := source.Uint()
	target.SetFloat(float64(n))
	if target.OverflowFloat(float64(n)) {
		c.AddErrorf("%s(%d) overflow", target.Type().String(), n)
	}
	return true
}
func getCopyFuncToFloatFromUint(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return copyReflectToFloatFromUint
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
func getCopyFuncToFloatFromString(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return copyReflectToFloatFromString
}
