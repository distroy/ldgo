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
		{To: reflect.Slice, From: reflect.Invalid}: copyReflectToSliceFromInvalid,
		{To: reflect.Slice, From: reflect.String}:  copyReflectToSliceFromString,
		{To: reflect.Slice, From: reflect.Slice}:   copyReflectToSliceFromSlice,
		{To: reflect.Slice, From: reflect.Array}:   copyReflectToSliceFromArray,
		{To: reflect.Slice, From: reflect.Map}:     copyReflectToSliceFromMap,
	})
}

func copyReflectToSliceFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.Set(reflect.Zero(target.Type()))
	return true
}

func copyReflectToSliceFromString(c *copyContext, target, source reflect.Value) bool {
	switch target.Type() {
	default:
		return false

	case typeOfByteSlice:
		source = reflect.ValueOf([]byte(source.String()))

	case typeOfRuneSlice:
		source = reflect.ValueOf([]rune(source.String()))
	}

	target.Set(source)
	return true
}

func copyReflectToSliceFromArray(c *copyContext, target, source reflect.Value) bool {
	// sVal := source.Slice(0, source.Len())
	// return copyReflectToSliceFromSlice(c, target, sVal)
	return copyReflectToSliceFromSlice(c, target, source)
}

func copyReflectToSliceFromSlice(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	sTyp := source.Type()
	if !c.Clone && tTyp == sTyp {
		target.Set(source)
		return true
	}

	if sTyp.Kind() != reflect.Slice && sTyp.Kind() != reflect.Array {
		return false
	}
	if !isCopyTypeConvertible(tTyp.Elem(), sTyp.Elem()) {
		return false
	}

	l := source.Len()
	if l > target.Len() {
		target.Set(reflect.MakeSlice(tTyp, l, l))
	}

	for i := 0; i < l; i++ {
		tItem := target.Index(i)
		sItem := source.Index(i)

		c.PushField(strconv.Itoa(i))
		copyReflect(c, tItem, sItem)
		c.PopField()
	}

	return true
}

func copyReflectToSliceFromMap(c *copyContext, target, source reflect.Value) bool {
	sTyp := source.Type()

	if isEmptyStruct(sTyp.Elem()) {
		return copyReflectToSliceFromMapWithEmptyStructValue(c, target, source)
	}
	return false
}

func copyReflectToSliceFromMapWithEmptyStructValue(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	sTyp := source.Type()

	if isCopyTypeConvertible(sTyp.Elem(), sTyp.Key()) {
		return false
	}

	l := source.Len()
	if l > target.Len() {
		if target.Kind() == reflect.Array {
			c.AddErrorf("%s has %d elements, can not convert to %s", sTyp.String(), l, tTyp.String())
			l = target.Len()

		} else {
			target.Set(reflect.MakeSlice(tTyp, l, l))
		}
	}

	// sTyp.Comparable()
	i := 0
	for it := source.MapRange(); i < l && it.Next(); i++ {
		tItem := target.Index(i)
		sItem := it.Key()

		c.PushField(strconv.Itoa(i))
		copyReflect(c, tItem, sItem)
		c.PopField()
	}

	return true
}
