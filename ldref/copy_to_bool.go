/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"strconv"
	"strings"
)

func init() {
	registerCopyFunc(map[copyPair]copyFuncType{
		{To: reflect.Bool, From: reflect.Invalid}:    copyReflectToBoolFromInvalid,
		{To: reflect.Bool, From: reflect.Bool}:       copyReflectToBoolFromBool,
		{To: reflect.Bool, From: reflect.String}:     copyReflectToBoolFromString,
		{To: reflect.Bool, From: reflect.Complex64}:  copyReflectToBoolFromComplex,
		{To: reflect.Bool, From: reflect.Complex128}: copyReflectToBoolFromComplex,
		{To: reflect.Bool, From: reflect.Float32}:    copyReflectToBoolFromFloat,
		{To: reflect.Bool, From: reflect.Float64}:    copyReflectToBoolFromFloat,
		{To: reflect.Bool, From: reflect.Int}:        copyReflectToBoolFromInt,
		{To: reflect.Bool, From: reflect.Int8}:       copyReflectToBoolFromInt,
		{To: reflect.Bool, From: reflect.Int16}:      copyReflectToBoolFromInt,
		{To: reflect.Bool, From: reflect.Int32}:      copyReflectToBoolFromInt,
		{To: reflect.Bool, From: reflect.Int64}:      copyReflectToBoolFromInt,
		{To: reflect.Bool, From: reflect.Uint}:       copyReflectToBoolFromUint,
		{To: reflect.Bool, From: reflect.Uint8}:      copyReflectToBoolFromUint,
		{To: reflect.Bool, From: reflect.Uint16}:     copyReflectToBoolFromUint,
		{To: reflect.Bool, From: reflect.Uint32}:     copyReflectToBoolFromUint,
		{To: reflect.Bool, From: reflect.Uint64}:     copyReflectToBoolFromUint,
		{To: reflect.Bool, From: reflect.Uintptr}:    copyReflectToBoolFromUint,
	})
}

func bool2int(b bool) int {
	if b {
		return 1
	}
	return 0
}

func copyReflectToBoolFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.SetBool(false)
	return true
}

func copyReflectToBoolFromBool(c *copyContext, target, source reflect.Value) bool {
	b := source.Bool()
	target.SetBool(b)
	return true
}

func copyReflectToBoolFromComplex(c *copyContext, target, source reflect.Value) bool {
	n := source.Complex()
	target.SetBool(n != 0)
	return true
}

func copyReflectToBoolFromFloat(c *copyContext, target, source reflect.Value) bool {
	n := source.Float()
	target.SetBool(n != 0)
	return true
}

func copyReflectToBoolFromInt(c *copyContext, target, source reflect.Value) bool {
	n := source.Int()
	target.SetBool(n != 0)
	return true
}

func copyReflectToBoolFromUint(c *copyContext, target, source reflect.Value) bool {
	n := source.Uint()
	target.SetBool(n != 0)
	return true
}

func copyReflectToBoolFromString(c *copyContext, target, source reflect.Value) bool {
	s := source.String()

	if len(s) <= 5 {
		switch strings.ToLower(s) {
		case "true", "on":
			target.SetBool(true)
			return true

		case "", "false", "off",
			"nil", "null", "undefined":
			target.SetBool(false)
			return true
		}
	}

	n, err := strconv.ParseFloat(s, 64)
	target.SetBool(n != 0)
	if err != nil {
		c.AddErrorf("can not convert to %s, %q", target.Type().String(), s)
	}
	return true
}
