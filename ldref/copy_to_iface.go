/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
)

func init() {
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
}

func copyReflectToIfaceFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.Set(reflect.Zero(target.Type()))
	return true
}

func copyReflectToIfaceFromPtr(c *copyContext, target, source reflect.Value) bool {
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

func copyReflectToIfaceFromIface(c *copyContext, target, source reflect.Value) bool {
	if source.IsNil() {
		val := reflect.Zero(typeOfIface)
		target.Set(val)
		return true
	}

	return copyReflectWithIndirect(c, target, source.Elem())
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
		sVal = deepClone(sVal)
	}

	target.Set(sVal)
	return true
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

func copyReflectToIfaceFromMap(c *copyContext, target, source reflect.Value) bool {
	return copyReflectToIfaceFromComplexKinds(c, target, source, nil)
}

func copyReflectToIfaceFromBaseKinds(c *copyContext, target, source reflect.Value) bool {
	return copyReflectToIfaceFromComplexKinds(c, target, source, nil)
}
