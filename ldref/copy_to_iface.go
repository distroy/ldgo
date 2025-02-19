/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
)

func init() {
	registerCopyFunc(map[copyPair]copyFuncType{
		{To: reflect.Ptr, From: reflect.Interface}: copyReflectFromIface,

		{To: reflect.Struct, From: reflect.Interface}: copyReflectFromIface,
		{To: reflect.Map, From: reflect.Interface}:    copyReflectFromIface,
		{To: reflect.Slice, From: reflect.Interface}:  copyReflectFromIface,
		{To: reflect.Array, From: reflect.Interface}:  copyReflectFromIface,

		{To: reflect.Bool, From: reflect.Interface}:       copyReflectFromIface,
		{To: reflect.String, From: reflect.Interface}:     copyReflectFromIface,
		{To: reflect.Float32, From: reflect.Interface}:    copyReflectFromIface,
		{To: reflect.Float64, From: reflect.Interface}:    copyReflectFromIface,
		{To: reflect.Complex64, From: reflect.Interface}:  copyReflectFromIface,
		{To: reflect.Complex128, From: reflect.Interface}: copyReflectFromIface,
		{To: reflect.Int, From: reflect.Interface}:        copyReflectFromIface,
		{To: reflect.Int8, From: reflect.Interface}:       copyReflectFromIface,
		{To: reflect.Int16, From: reflect.Interface}:      copyReflectFromIface,
		{To: reflect.Int32, From: reflect.Interface}:      copyReflectFromIface,
		{To: reflect.Int64, From: reflect.Interface}:      copyReflectFromIface,
		{To: reflect.Uint, From: reflect.Interface}:       copyReflectFromIface,
		{To: reflect.Uint8, From: reflect.Interface}:      copyReflectFromIface,
		{To: reflect.Uint16, From: reflect.Interface}:     copyReflectFromIface,
		{To: reflect.Uint32, From: reflect.Interface}:     copyReflectFromIface,
		{To: reflect.Uint64, From: reflect.Interface}:     copyReflectFromIface,
		{To: reflect.Uintptr, From: reflect.Interface}:    copyReflectFromIface,

		{To: reflect.UnsafePointer, From: reflect.Interface}: copyReflectFromIface,
		{To: reflect.Func, From: reflect.Interface}:          copyReflectFromIface,
		{To: reflect.Chan, From: reflect.Interface}:          copyReflectFromIface,
	})
	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.Ptr, From: reflect.Interface}: getCopyFuncFromIface,

		{To: reflect.Struct, From: reflect.Interface}: getCopyFuncFromIface,
		{To: reflect.Map, From: reflect.Interface}:    getCopyFuncFromIface,
		{To: reflect.Slice, From: reflect.Interface}:  getCopyFuncFromIface,
		{To: reflect.Array, From: reflect.Interface}:  getCopyFuncFromIface,

		{To: reflect.Bool, From: reflect.Interface}:       getCopyFuncFromIface,
		{To: reflect.String, From: reflect.Interface}:     getCopyFuncFromIface,
		{To: reflect.Float32, From: reflect.Interface}:    getCopyFuncFromIface,
		{To: reflect.Float64, From: reflect.Interface}:    getCopyFuncFromIface,
		{To: reflect.Complex64, From: reflect.Interface}:  getCopyFuncFromIface,
		{To: reflect.Complex128, From: reflect.Interface}: getCopyFuncFromIface,
		{To: reflect.Int, From: reflect.Interface}:        getCopyFuncFromIface,
		{To: reflect.Int8, From: reflect.Interface}:       getCopyFuncFromIface,
		{To: reflect.Int16, From: reflect.Interface}:      getCopyFuncFromIface,
		{To: reflect.Int32, From: reflect.Interface}:      getCopyFuncFromIface,
		{To: reflect.Int64, From: reflect.Interface}:      getCopyFuncFromIface,
		{To: reflect.Uint, From: reflect.Interface}:       getCopyFuncFromIface,
		{To: reflect.Uint8, From: reflect.Interface}:      getCopyFuncFromIface,
		{To: reflect.Uint16, From: reflect.Interface}:     getCopyFuncFromIface,
		{To: reflect.Uint32, From: reflect.Interface}:     getCopyFuncFromIface,
		{To: reflect.Uint64, From: reflect.Interface}:     getCopyFuncFromIface,
		{To: reflect.Uintptr, From: reflect.Interface}:    getCopyFuncFromIface,

		{To: reflect.UnsafePointer, From: reflect.Interface}: getCopyFuncFromIface,
		{To: reflect.Func, From: reflect.Interface}:          getCopyFuncFromIface,
		{To: reflect.Chan, From: reflect.Interface}:          getCopyFuncFromIface,
	})

	registerCopyFunc(map[copyPair]copyFuncType{
		{To: reflect.Interface, From: reflect.Invalid}:   copyReflectToIfaceFromInvalid,
		{To: reflect.Interface, From: reflect.Interface}: copyReflectToIfaceFromIface,
		{To: reflect.Interface, From: reflect.Ptr}:       copyReflectToIfaceFromPtr,

		{To: reflect.Interface, From: reflect.Struct}: copyReflectToIfaceFromStruct,
		{To: reflect.Interface, From: reflect.Map}:    copyReflectToIfaceFromMap,
		{To: reflect.Interface, From: reflect.Slice}:  copyReflectToIfaceFromSlice,
		{To: reflect.Interface, From: reflect.Array}:  copyReflectToIfaceFromArray,

		{To: reflect.Interface, From: reflect.Bool}:       copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.String}:     copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Float32}:    copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Float64}:    copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Complex64}:  copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Complex128}: copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int}:        copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int8}:       copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int16}:      copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int32}:      copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int64}:      copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint}:       copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint8}:      copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint16}:     copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint32}:     copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint64}:     copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uintptr}:    copyReflectToIfaceFromBaseKinds,

		{To: reflect.Interface, From: reflect.UnsafePointer}: copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Func}:          copyReflectToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Chan}:          copyReflectToIfaceFromBaseKinds,
	})
	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.Interface, From: reflect.Invalid}:   getCopyFuncToIfaceFromInvalid,
		{To: reflect.Interface, From: reflect.Interface}: getCopyFuncToIfaceFromIface,
		{To: reflect.Interface, From: reflect.Ptr}:       getCopyFuncToIfaceFromPtr,

		{To: reflect.Interface, From: reflect.Struct}: getCopyFuncToIfaceFromStruct,
		{To: reflect.Interface, From: reflect.Map}:    getCopyFuncToIfaceFromMap,
		{To: reflect.Interface, From: reflect.Slice}:  getCopyFuncToIfaceFromSlice,
		{To: reflect.Interface, From: reflect.Array}:  getCopyFuncToIfaceFromArray,

		{To: reflect.Interface, From: reflect.Bool}:       getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.String}:     getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Float32}:    getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Float64}:    getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Complex64}:  getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Complex128}: getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int}:        getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int8}:       getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int16}:      getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int32}:      getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Int64}:      getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint}:       getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint8}:      getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint16}:     getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint32}:     getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uint64}:     getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Uintptr}:    getCopyFuncToIfaceFromBaseKinds,

		{To: reflect.Interface, From: reflect.UnsafePointer}: getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Func}:          getCopyFuncToIfaceFromBaseKinds,
		{To: reflect.Interface, From: reflect.Chan}:          getCopyFuncToIfaceFromBaseKinds,
	})
}

func copyReflectToIfaceFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.Set(reflect.Zero(target.Type()))
	return true
}
func getCopyFuncToIfaceFromInvalid(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncSetZero(c, tTyp, sTyp)
}

func copyReflectToIfaceFromPtr(c *copyContext, target, source reflect.Value) bool {
	if !source.Type().Implements(target.Type()) {
		return false
	}

	val, _ := indirectCopySource(source)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		target.Set(source)
		return true
	}

	switch val.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array, reflect.Interface:
		return copyReflectWithIndirect(c, target, source.Elem())
	}

	return copyReflectToIfaceFromComplexKinds(c, target, source, nil)
}
func getCopyFuncToIfaceFromPtr(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	if !sTyp.Implements(tTyp) {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}

	sElemTyp := sTyp
	for sElemTyp.Kind() == reflect.Ptr {
		sElemTyp = sElemTyp.Elem()
	}

	switch sElemTyp.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array, reflect.Interface:
		pfe, done := getCopyFuncIndirect(c, tTyp, sTyp.Elem())
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			val, _ := indirectCopySource(source)
			if val.Kind() == reflect.Ptr && val.IsNil() {
				target.Set(source)
				return true
			}

			// return copyReflectToIfaceFromComplexKinds(c, target, source, nil)
			done()
			return (*pfe)(c, target, source.Elem())
		}
	}

	fnElemCopy := getCopyFuncToIfaceFromComplexKinds(c, tTyp, sTyp, nil)
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		val, _ := indirectCopySource(source)
		if val.Kind() == reflect.Ptr && val.IsNil() {
			target.Set(source)
			return true
		}

		// return copyReflectToIfaceFromComplexKinds(c, target, source, nil)
		return fnElemCopy(c, target, source)
	}
}

func copyReflectFromIface(c *copyContext, target, source reflect.Value) bool {
	if source.IsNil() {
		val := reflect.Zero(target.Type())
		target.Set(val)
		return true
	}

	return copyReflectWithIndirect(c, target, source.Elem())
}
func copyReflectToIfaceFromIface(c *copyContext, target, source reflect.Value) bool {
	return copyReflectFromIface(c, target, source)
}
func getCopyFuncFromIface(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	tZero := reflect.Zero(tTyp)
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		if source.IsNil() {
			target.Set(tZero)
			return true
		}

		return copyReflectWithIndirectV2(c, target, source.Elem())
	}
}
func getCopyFuncToIfaceFromIface(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncFromIface(c, tTyp, sTyp)
}

func copyReflectToIfaceFromComplexKinds(
	c *copyContext, target, source reflect.Value,
	fCopyAny func(c *copyContext, target, source reflect.Value) bool,
) bool {
	tTyp := target.Type()
	sTyp := source.Type()

	// log.Printf(" === target type: %s", tTyp.String())
	// log.Printf(" === source type: %s", sTyp.String())
	// log.Printf(" === copy func %x", reflect.ValueOf(fCopyAny).Pointer())

	if fCopyAny != nil && tTyp.NumMethod() == 0 {
		return fCopyAny(c, target, source)
	}

	if !sTyp.Implements(tTyp) {
		c.AddErrorf("%s can not copy to %s", typeNameOfReflect(source), typeNameOfReflect(target))
		return false
	}

	sVal := source
	if c.Clone {
		sVal = deepCloneRef(sVal)
	}

	target.Set(sVal)
	return true
}
func getCopyFuncToIfaceFromComplexKinds(
	c *copyContext, tTyp, sTyp reflect.Type,
	fnAny func(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType,
) copyFuncType {
	// tTyp := target.Type()
	// sTyp := source.Type()

	// log.Printf(" === target type: %s", tTyp.String())
	// log.Printf(" === source type: %s", sTyp.String())
	// log.Printf(" === copy func %x", reflect.ValueOf(fCopyAny).Pointer())

	if fnAny != nil && tTyp.NumMethod() == 0 {
		// return fCopyAny(c, target, source)
		return fnAny(c, tTyp, sTyp)
	}

	if !sTyp.Implements(tTyp) {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			c.AddErrorf("%s can not copy to %s", typeNameOfReflect(source), typeNameOfReflect(target))
			return false
		}
	}
	if c.Clone {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			sVal := source
			sVal = deepCloneRef(sVal)

			target.Set(sVal)
			return true
		}
	}

	return func(c *copyContext, target, source reflect.Value) (end bool) {
		sVal := source
		target.Set(sVal)
		return true
	}
}

func copyReflectToIfaceFromStruct(c *copyContext, target, source reflect.Value) bool {
	return copyReflectToIfaceFromComplexKinds(c, target, source, func(c *copyContext, target, source reflect.Value) bool {
		val := target
		if target.IsNil() || target.Elem().Type() != typeOfMapStrIface {
			val = reflect.MakeMap(typeOfMapStrIface)
		} else {
			val = val.Elem()
		}
		ok := copyReflectToMapFromStruct(c, val, source)
		target.Set(val)
		return ok
	})
}
func getCopyFuncToIfaceFromStruct(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncToIfaceFromComplexKinds(c, tTyp, sTyp, func(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
		tValTyp := typeOfMapStrIface
		fnCopy := getCopyFuncToMapFromStruct(c, tValTyp, sTyp)
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			val := target
			if target.IsNil() || target.Elem().Type() != tValTyp {
				val = reflect.MakeMap(tValTyp)
			} else {
				val = val.Elem()
			}
			ok := fnCopy(c, val, source)
			target.Set(val)
			return ok
		}
	})
}

func copyReflectToIfaceFromSlice(c *copyContext, target, source reflect.Value) bool {
	return copyReflectToIfaceFromComplexKinds(c, target, source, func(c *copyContext, target, source reflect.Value) bool {
		val := target
		if target.IsNil() || target.Elem().Type() != typeOfIfaceSlice {
			l := source.Len()
			val = reflect.MakeSlice(typeOfIfaceSlice, l, l)

		} else {
			val = val.Elem()
		}
		ok := copyReflectToSliceFromSlice(c, val, source)
		target.Set(val)
		return ok
	})
}
func getCopyFuncToIfaceFromSlice(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncToIfaceFromComplexKinds(c, tTyp, sTyp, func(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
		tValTyp := typeOfIfaceSlice
		fnCopy := getCopyFuncToSliceFromSlice(c, tValTyp, sTyp)
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			val := target
			if target.IsNil() || target.Elem().Type() != tValTyp {
				l := source.Len()
				val = reflect.MakeSlice(tValTyp, l, l)

			} else {
				val = val.Elem()
			}
			ok := fnCopy(c, val, source)
			target.Set(val)
			return ok
		}
	})
}

func copyReflectToIfaceFromArray(c *copyContext, target, source reflect.Value) bool {
	return copyReflectToIfaceFromComplexKinds(c, target, source, func(c *copyContext, target, source reflect.Value) bool {
		val := target
		if target.IsNil() || target.Elem().Type() != typeOfIfaceSlice {
			l := source.Len()
			val = reflect.MakeSlice(typeOfIfaceSlice, l, l)

		} else {
			val = val.Elem()
		}
		ok := copyReflectToSliceFromArray(c, val, source)
		target.Set(val)
		return ok
	})
}
func getCopyFuncToIfaceFromArray(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncToIfaceFromComplexKinds(c, tTyp, sTyp, func(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
		tValTyp := typeOfIfaceSlice
		fnCopy := getCopyFuncToSliceFromArray(c, tValTyp, sTyp)
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			val := target
			if target.IsNil() || target.Elem().Type() != tValTyp {
				l := source.Len()
				val = reflect.MakeSlice(tValTyp, l, l)

			} else {
				val = val.Elem()
			}
			ok := fnCopy(c, val, source)
			target.Set(val)
			return ok
		}
	})
}

func copyReflectToIfaceFromMap(c *copyContext, target, source reflect.Value) bool {
	return copyReflectToIfaceFromComplexKinds(c, target, source, nil)
}
func getCopyFuncToIfaceFromMap(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncToIfaceFromComplexKinds(c, tTyp, sTyp, nil)
}

func copyReflectToIfaceFromBaseKinds(c *copyContext, target, source reflect.Value) bool {
	return copyReflectToIfaceFromComplexKinds(c, target, source, nil)
}
func getCopyFuncToIfaceFromBaseKinds(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncToIfaceFromComplexKinds(c, tTyp, sTyp, nil)
}
