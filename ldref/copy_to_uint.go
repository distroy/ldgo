/*
 * Copyright (C) distroy
 */

package ldref

import (
	"math"
	"reflect"
	"strconv"
)

func init() {
	registerCopyFunc(map[copyPair]copyFuncType{
		{To: reflect.Uint, From: reflect.Invalid}:    copyReflectToUintFromInvalid,
		{To: reflect.Uint, From: reflect.Bool}:       copyReflectToUintFromBool,
		{To: reflect.Uint, From: reflect.Complex64}:  copyReflectToUintFromComplex,
		{To: reflect.Uint, From: reflect.Complex128}: copyReflectToUintFromComplex,
		{To: reflect.Uint, From: reflect.Float32}:    copyReflectToUintFromFloat,
		{To: reflect.Uint, From: reflect.Float64}:    copyReflectToUintFromFloat,
		{To: reflect.Uint, From: reflect.Int}:        copyReflectToUintFromInt,
		{To: reflect.Uint, From: reflect.Int8}:       copyReflectToUintFromInt,
		{To: reflect.Uint, From: reflect.Int16}:      copyReflectToUintFromInt,
		{To: reflect.Uint, From: reflect.Int32}:      copyReflectToUintFromInt,
		{To: reflect.Uint, From: reflect.Int64}:      copyReflectToUintFromInt,
		{To: reflect.Uint, From: reflect.Uint}:       copyReflectToUintFromUint,
		{To: reflect.Uint, From: reflect.Uint8}:      copyReflectToUintFromUint,
		{To: reflect.Uint, From: reflect.Uint16}:     copyReflectToUintFromUint,
		{To: reflect.Uint, From: reflect.Uint32}:     copyReflectToUintFromUint,
		{To: reflect.Uint, From: reflect.Uint64}:     copyReflectToUintFromUint,
		{To: reflect.Uint, From: reflect.Uintptr}:    copyReflectToUintFromUint,
		{To: reflect.Uint, From: reflect.String}:     copyReflectToUintFromString,

		{To: reflect.Uint8, From: reflect.Invalid}:    copyReflectToUintFromInvalid,
		{To: reflect.Uint8, From: reflect.Bool}:       copyReflectToUintFromBool,
		{To: reflect.Uint8, From: reflect.Complex64}:  copyReflectToUintFromComplex,
		{To: reflect.Uint8, From: reflect.Complex128}: copyReflectToUintFromComplex,
		{To: reflect.Uint8, From: reflect.Float32}:    copyReflectToUintFromFloat,
		{To: reflect.Uint8, From: reflect.Float64}:    copyReflectToUintFromFloat,
		{To: reflect.Uint8, From: reflect.Int}:        copyReflectToUintFromInt,
		{To: reflect.Uint8, From: reflect.Int8}:       copyReflectToUintFromInt,
		{To: reflect.Uint8, From: reflect.Int16}:      copyReflectToUintFromInt,
		{To: reflect.Uint8, From: reflect.Int32}:      copyReflectToUintFromInt,
		{To: reflect.Uint8, From: reflect.Int64}:      copyReflectToUintFromInt,
		{To: reflect.Uint8, From: reflect.Uint}:       copyReflectToUintFromUint,
		{To: reflect.Uint8, From: reflect.Uint8}:      copyReflectToUintFromUint,
		{To: reflect.Uint8, From: reflect.Uint16}:     copyReflectToUintFromUint,
		{To: reflect.Uint8, From: reflect.Uint32}:     copyReflectToUintFromUint,
		{To: reflect.Uint8, From: reflect.Uint64}:     copyReflectToUintFromUint,
		{To: reflect.Uint8, From: reflect.Uintptr}:    copyReflectToUintFromUint,
		{To: reflect.Uint8, From: reflect.String}:     copyReflectToUintFromString,

		{To: reflect.Uint16, From: reflect.Invalid}:    copyReflectToUintFromInvalid,
		{To: reflect.Uint16, From: reflect.Bool}:       copyReflectToUintFromBool,
		{To: reflect.Uint16, From: reflect.Complex64}:  copyReflectToUintFromComplex,
		{To: reflect.Uint16, From: reflect.Complex128}: copyReflectToUintFromComplex,
		{To: reflect.Uint16, From: reflect.Float32}:    copyReflectToUintFromFloat,
		{To: reflect.Uint16, From: reflect.Float64}:    copyReflectToUintFromFloat,
		{To: reflect.Uint16, From: reflect.Int}:        copyReflectToUintFromInt,
		{To: reflect.Uint16, From: reflect.Int8}:       copyReflectToUintFromInt,
		{To: reflect.Uint16, From: reflect.Int16}:      copyReflectToUintFromInt,
		{To: reflect.Uint16, From: reflect.Int32}:      copyReflectToUintFromInt,
		{To: reflect.Uint16, From: reflect.Int64}:      copyReflectToUintFromInt,
		{To: reflect.Uint16, From: reflect.Uint}:       copyReflectToUintFromUint,
		{To: reflect.Uint16, From: reflect.Uint8}:      copyReflectToUintFromUint,
		{To: reflect.Uint16, From: reflect.Uint16}:     copyReflectToUintFromUint,
		{To: reflect.Uint16, From: reflect.Uint32}:     copyReflectToUintFromUint,
		{To: reflect.Uint16, From: reflect.Uint64}:     copyReflectToUintFromUint,
		{To: reflect.Uint16, From: reflect.Uintptr}:    copyReflectToUintFromUint,
		{To: reflect.Uint16, From: reflect.String}:     copyReflectToUintFromString,

		{To: reflect.Uint32, From: reflect.Invalid}:    copyReflectToUintFromInvalid,
		{To: reflect.Uint32, From: reflect.Bool}:       copyReflectToUintFromBool,
		{To: reflect.Uint32, From: reflect.Complex64}:  copyReflectToUintFromComplex,
		{To: reflect.Uint32, From: reflect.Complex128}: copyReflectToUintFromComplex,
		{To: reflect.Uint32, From: reflect.Float32}:    copyReflectToUintFromFloat,
		{To: reflect.Uint32, From: reflect.Float64}:    copyReflectToUintFromFloat,
		{To: reflect.Uint32, From: reflect.Int}:        copyReflectToUintFromInt,
		{To: reflect.Uint32, From: reflect.Int8}:       copyReflectToUintFromInt,
		{To: reflect.Uint32, From: reflect.Int16}:      copyReflectToUintFromInt,
		{To: reflect.Uint32, From: reflect.Int32}:      copyReflectToUintFromInt,
		{To: reflect.Uint32, From: reflect.Int64}:      copyReflectToUintFromInt,
		{To: reflect.Uint32, From: reflect.Uint}:       copyReflectToUintFromUint,
		{To: reflect.Uint32, From: reflect.Uint8}:      copyReflectToUintFromUint,
		{To: reflect.Uint32, From: reflect.Uint16}:     copyReflectToUintFromUint,
		{To: reflect.Uint32, From: reflect.Uint32}:     copyReflectToUintFromUint,
		{To: reflect.Uint32, From: reflect.Uint64}:     copyReflectToUintFromUint,
		{To: reflect.Uint32, From: reflect.Uintptr}:    copyReflectToUintFromUint,
		{To: reflect.Uint32, From: reflect.String}:     copyReflectToUintFromString,

		{To: reflect.Uint64, From: reflect.Invalid}:    copyReflectToUintFromInvalid,
		{To: reflect.Uint64, From: reflect.Bool}:       copyReflectToUintFromBool,
		{To: reflect.Uint64, From: reflect.Complex64}:  copyReflectToUintFromComplex,
		{To: reflect.Uint64, From: reflect.Complex128}: copyReflectToUintFromComplex,
		{To: reflect.Uint64, From: reflect.Float32}:    copyReflectToUintFromFloat,
		{To: reflect.Uint64, From: reflect.Float64}:    copyReflectToUintFromFloat,
		{To: reflect.Uint64, From: reflect.Int}:        copyReflectToUintFromInt,
		{To: reflect.Uint64, From: reflect.Int8}:       copyReflectToUintFromInt,
		{To: reflect.Uint64, From: reflect.Int16}:      copyReflectToUintFromInt,
		{To: reflect.Uint64, From: reflect.Int32}:      copyReflectToUintFromInt,
		{To: reflect.Uint64, From: reflect.Int64}:      copyReflectToUintFromInt,
		{To: reflect.Uint64, From: reflect.Uint}:       copyReflectToUintFromUint,
		{To: reflect.Uint64, From: reflect.Uint8}:      copyReflectToUintFromUint,
		{To: reflect.Uint64, From: reflect.Uint16}:     copyReflectToUintFromUint,
		{To: reflect.Uint64, From: reflect.Uint32}:     copyReflectToUintFromUint,
		{To: reflect.Uint64, From: reflect.Uint64}:     copyReflectToUintFromUint,
		{To: reflect.Uint64, From: reflect.Uintptr}:    copyReflectToUintFromUint,
		{To: reflect.Uint64, From: reflect.String}:     copyReflectToUintFromString,

		{To: reflect.Uintptr, From: reflect.Invalid}:    copyReflectToUintFromInvalid,
		{To: reflect.Uintptr, From: reflect.Bool}:       copyReflectToUintFromBool,
		{To: reflect.Uintptr, From: reflect.Complex64}:  copyReflectToUintFromComplex,
		{To: reflect.Uintptr, From: reflect.Complex128}: copyReflectToUintFromComplex,
		{To: reflect.Uintptr, From: reflect.Float32}:    copyReflectToUintFromFloat,
		{To: reflect.Uintptr, From: reflect.Float64}:    copyReflectToUintFromFloat,
		{To: reflect.Uintptr, From: reflect.Int}:        copyReflectToUintFromInt,
		{To: reflect.Uintptr, From: reflect.Int8}:       copyReflectToUintFromInt,
		{To: reflect.Uintptr, From: reflect.Int16}:      copyReflectToUintFromInt,
		{To: reflect.Uintptr, From: reflect.Int32}:      copyReflectToUintFromInt,
		{To: reflect.Uintptr, From: reflect.Int64}:      copyReflectToUintFromInt,
		{To: reflect.Uintptr, From: reflect.Uint}:       copyReflectToUintFromUint,
		{To: reflect.Uintptr, From: reflect.Uint8}:      copyReflectToUintFromUint,
		{To: reflect.Uintptr, From: reflect.Uint16}:     copyReflectToUintFromUint,
		{To: reflect.Uintptr, From: reflect.Uint32}:     copyReflectToUintFromUint,
		{To: reflect.Uintptr, From: reflect.Uint64}:     copyReflectToUintFromUint,
		{To: reflect.Uintptr, From: reflect.Uintptr}:    copyReflectToUintFromUint,
		{To: reflect.Uintptr, From: reflect.String}:     copyReflectToUintFromString,
	})
}

func copyReflectToUintFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.SetUint(0)
	return true
}

func copyReflectToUintFromBool(c *copyContext, target, source reflect.Value) bool {
	b := source.Bool()
	target.SetUint(uint64(bool2int(b)))
	return true
}

func copyReflectToUintFromFloat(c *copyContext, target, source reflect.Value) bool {
	n := source.Float()
	target.SetUint(uint64(n))
	if target.OverflowUint(uint64(n)) {
		c.AddErrorf("%s(%f) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToUintFromComplex(c *copyContext, target, source reflect.Value) bool {
	n := source.Complex()
	r := real(n)
	target.SetUint(uint64(r))
	if r > math.MaxUint64 || r < 0 || target.OverflowUint(uint64(r)) {
		c.AddErrorf("%s(%v) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToUintFromInt(c *copyContext, target, source reflect.Value) bool {
	n := source.Int()
	target.SetUint(uint64(n))
	if n < 0 || target.OverflowUint(uint64(n)) {
		c.AddErrorf("%s(%d) overflow", target.Type().String(), n)
	}

	return true
}

func copyReflectToUintFromUint(c *copyContext, target, source reflect.Value) bool {
	n := source.Uint()
	target.SetUint(n)
	if target.OverflowUint(n) {
		c.AddErrorf("%s(%d) overflow", target.Type().String(), n)
	}
	return true
}

func copyReflectToUintFromString(c *copyContext, target, source reflect.Value) bool {
	s := source.String()
	n, err := strconv.ParseUint(s, 10, 64)
	target.SetUint(n)
	if err != nil {
		c.AddErrorf("can not convert to %s, %q", target.Type().String(), s)

	} else if target.OverflowUint(n) {
		c.AddErrorf("%s(%s) overflow", target.Type().String(), s)
	}
	return true
}
