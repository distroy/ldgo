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
		{To: reflect.Ptr, From: reflect.Invalid}:       copyReflectToPtrFromInvalid,
		{To: reflect.Ptr, From: reflect.Ptr}:           copyReflectToPtrFromPtr,
		{To: reflect.Ptr, From: reflect.UnsafePointer}: copyReflectToPtrFromUnsafePointer,
		{To: reflect.Ptr, From: reflect.Interface}:     copyReflectToPtrFromIface,

		{To: reflect.Ptr, From: reflect.Func}:       copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Map}:        copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Slice}:      copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Array}:      copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Struct}:     copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Bool}:       copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.String}:     copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Float32}:    copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Float64}:    copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Complex64}:  copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Complex128}: copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int}:        copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int8}:       copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int16}:      copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int32}:      copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int64}:      copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint}:       copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint8}:      copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint16}:     copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint32}:     copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint64}:     copyReflectToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uintptr}:    copyReflectToPtrFromOthers,
	})
	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.Ptr, From: reflect.Invalid}:       getCopyFuncToPtrFromInvalid,
		{To: reflect.Ptr, From: reflect.Ptr}:           getCopyFuncToPtrFromPtr,
		{To: reflect.Ptr, From: reflect.UnsafePointer}: getCopyFuncToPtrFromUnsafePointer,
		{To: reflect.Ptr, From: reflect.Interface}:     getCopyFuncToPtrFromIface,

		{To: reflect.Ptr, From: reflect.Func}:       getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Map}:        getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Slice}:      getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Array}:      getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Struct}:     getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Bool}:       getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.String}:     getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Float32}:    getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Float64}:    getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Complex64}:  getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Complex128}: getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int}:        getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int8}:       getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int16}:      getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int32}:      getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Int64}:      getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint}:       getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint8}:      getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint16}:     getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint32}:     getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uint64}:     getCopyFuncToPtrFromOthers,
		{To: reflect.Ptr, From: reflect.Uintptr}:    getCopyFuncToPtrFromOthers,
	})
}

func copyReflectToPtrFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.Set(reflect.Zero(target.Type()))
	return true
}
func getCopyFuncToPtrFromInvalid(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncSetZero(c, tTyp, sTyp)
}

func copyReflectToPtrFromPtr(c *copyContext, target, source reflect.Value) bool {
	if source.IsNil() {
		target.SetZero()
	}
	if !c.Clone {
		if target.Type() == source.Type() {
			target.Set(source)
			return true

		} else if target.Type() == source.Type().Elem() {
			target.Set(source.Elem())
			return true
		}
	}
	// sVal, sPtrLvl := indirectCopySource(source)
	sVal, _ := indirectCopySource(source)
	if sVal.Kind() == reflect.Ptr {
		target.Set(reflect.Zero(target.Type()))
		return true
	}

	// tTyp, tPtrLvl := indirectType(target.Type())
	// if !c.Clone && tTyp == sVal.Type() {
	// 	tVal := target
	// 	sVal := source
	// 	for i := 0; i+sPtrLvl < tPtrLvl; i++ {
	// 		if tVal.IsNil() {
	// 			tVal.Set(reflect.New(tVal.Type().Elem()))
	// 		}
	// 		tVal = tVal.Elem()
	// 	}
	// 	for i := 0; i+tPtrLvl < sPtrLvl; i++ {
	// 		sVal = sVal.Elem()
	// 	}
	//
	// 	tVal.Set(sVal)
	// 	return true
	// }

	tVal, _ := indirectCopyTarget(target)
	return copyReflect(c, tVal, sVal)
}
func getCopyFuncToPtrFromPtr(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	if !c.Clone && tTyp == sTyp {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			target.Set(source)
			return true
		}
	}
	if !c.Clone && tTyp == sTyp {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			target.Set(source)
			return true
		}
	}
	if !c.Clone && tTyp == sTyp.Elem() {
		tZero := reflect.Zero(tTyp)
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			if source.IsNil() {
				target.Set(tZero)
				return true
			}
			target.Set(source.Elem())
			return true
		}
	}

	tElemType := tTyp
	sElemType := sTyp
	for tElemType.Kind() == reflect.Ptr {
		tElemType = tElemType.Elem()
	}
	for sElemType.Kind() == reflect.Ptr {
		sElemType = sElemType.Elem()
	}

	pfe, done := getCopyFunc(c, tElemType, sElemType)
	tZero := reflect.Zero(tTyp)
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		sVal, _ := indirectCopySource(source)
		if sVal.Kind() == reflect.Ptr {
			target.Set(tZero)
			return true
		}
		tVal, _ := indirectCopyTarget(target)
		done()
		(*pfe)(c, tVal, sVal)
		return true
	}
}

func copyReflectToPtrFromUnsafePointer(c *copyContext, target, source reflect.Value) bool {
	tAddr := unsafe.Pointer(target.UnsafeAddr())
	tTemp := reflect.NewAt(source.Type(), tAddr).Elem()
	tTemp.Set(source)
	return true
}
func getCopyFuncToPtrFromUnsafePointer(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return copyReflectToPtrFromUnsafePointer
}

func copyReflectToPtrFromIface(c *copyContext, target, source reflect.Value) bool {
	sVal := reflect.ValueOf(source.Interface())
	tVal := target
	return copyReflect(c, tVal, sVal)
}
func getCopyFuncToPtrFromIface(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		// sVal := reflect.ValueOf(source.Interface())
		sVal := source.Elem()
		tVal := target
		return copyReflectV2(c, tVal, sVal)
	}
}

func copyReflectToPtrFromOthers(c *copyContext, target, source reflect.Value) bool {
	sVal := source
	tVal, _ := indirectCopyTarget(target)
	return copyReflect(c, tVal, sVal)
}
func getCopyFuncToPtrFromOthers(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	// base types
	tt := tTyp
	for tt.Kind() == reflect.Ptr {
		tt = tt.Elem()
	}
	pfe, done := getCopyFunc(c, tt, sTyp)
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		sVal := source
		tVal, _ := indirectCopyTarget(target)
		done()
		return (*pfe)(c, tVal, sVal)
	}
}
