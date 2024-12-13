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
		{To: reflect.Func, From: reflect.Invalid}: copyReflectToFuncFromInvalid,
		{To: reflect.Func, From: reflect.Func}:    copyReflectToFuncFromFunc,
		// {To: reflect.Func, From: reflect.UnsafePointer}: copyReflectToFuncFromUnsafePointer,
	})
	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.Func, From: reflect.Invalid}: getCopyFuncToFuncFromInvalid,
		{To: reflect.Func, From: reflect.Func}:    getCopyFuncToFuncFromFunc,
		// {To: reflect.Func, From: reflect.UnsafePointer}: copyReflectToFuncFromUnsafePointer,
	})
}

func copyReflectToFuncFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.Set(reflect.Zero(target.Type()))
	return true
}
func getCopyFuncToFuncFromInvalid(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncSetZero(c, tTyp, sTyp)
}

func copyReflectToFuncFromFunc(c *copyContext, target, source reflect.Value) bool {
	if target.Type() != source.Type() {
		return false
	}

	target.Set(source)
	return true
}
func getCopyFuncToFuncFromFunc(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	if tTyp != sTyp {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}

	return func(c *copyContext, target, source reflect.Value) (end bool) {
		target.Set(source)
		return true
	}
}

func copyReflectToFuncFromUnsafePointer(c *copyContext, target, source reflect.Value) bool {
	// log.Printf(" === source:%#x", source.Pointer())
	// log.Printf(" === target:%#x", target.Pointer())

	// sTemp := reflect.New(source.Type()).Elem()
	// sTemp.Set(source)
	// sAddr := unsafe.Pointer(sTemp.UnsafeAddr())
	// sVal := reflect.NewAt(target.Type(), sAddr).Elem()
	// target.Set(sVal)
	// log.Printf(" === source:%#x", source.Pointer())
	// log.Printf(" === target:%#x", target.Pointer())

	tAddr := unsafe.Pointer(target.UnsafeAddr())
	tTemp := reflect.NewAt(source.Type(), tAddr).Elem()
	// log.Printf(" === tTemp addr:%#x", tTemp.UnsafeAddr())
	// log.Printf(" === target addr:%#x", target.UnsafeAddr())
	// log.Printf(" === tTemp type:%s", tTemp.Type().String())
	// log.Printf(" === target type:%s", target.Type().String())
	tTemp.Set(source)
	// log.Printf(" === source:%#x", source.Pointer())
	// log.Printf(" === tTemp:%#x", tTemp.Pointer())
	// log.Printf(" === target:%#x", target.Pointer())
	//
	// tTemp.SetPointer(unsafe.Pointer(source.Pointer()))
	// log.Printf(" === source:%#x", source.Pointer())
	// log.Printf(" === tTemp:%#x", tTemp.Pointer())
	// log.Printf(" === target:%#x", target.Pointer())
	//
	// log.Printf(" === tTemp type:%s", tTemp.Type().String())
	// log.Printf(" === target type:%s", target.Type().String())
	return true
}
