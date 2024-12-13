/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"unsafe"
)

func init() {
	registerCopyFunc(map[copyPair]copyFuncType{
		{To: reflect.UnsafePointer, From: reflect.Invalid}:       copyReflectToUnsafePointerFromInvalid,
		{To: reflect.UnsafePointer, From: reflect.UnsafePointer}: copyReflectToUnsafePointerFromUnsafePointer,
		{To: reflect.UnsafePointer, From: reflect.Ptr}:           copyReflectToUnsafePointerFromPtr,
		// {To: reflect.UnsafePointer, From: reflect.Func}:          copyReflectToUnsafePointerFromPtr,
	})
	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.UnsafePointer, From: reflect.Invalid}:       getCopyFuncToUnsafePointerFromInvalid,
		{To: reflect.UnsafePointer, From: reflect.UnsafePointer}: getCopyFuncToUnsafePointerFromUnsafePointer,
		{To: reflect.UnsafePointer, From: reflect.Ptr}:           getCopyFuncToUnsafePointerFromPtr,
		// {To: reflect.UnsafePointer, From: reflect.Func}:          getCopyFuncToUnsafePointerFromPtr,
	})
}

func copyReflectToUnsafePointerFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.Set(reflect.Zero(target.Type()))
	return true
}
func getCopyFuncToUnsafePointerFromInvalid(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncSetZero(c, tTyp, sTyp)
}

func copyReflectToUnsafePointerFromUnsafePointer(c *copyContext, target, source reflect.Value) bool {
	target.Set(source)
	return true
}
func getCopyFuncToUnsafePointerFromUnsafePointer(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		target.Set(source)
		return true
	}
}

func copyReflectToUnsafePointerFromPtr(c *copyContext, target, source reflect.Value) bool {
	target.SetPointer(unsafe.Pointer(source.Pointer()))
	return true
}
func getCopyFuncToUnsafePointerFromPtr(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		target.SetPointer(unsafe.Pointer(source.Pointer()))
		return true
	}
}
