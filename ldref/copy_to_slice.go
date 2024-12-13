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
	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.Slice, From: reflect.Invalid}: getCopyFuncToSliceFromInvalid,
		{To: reflect.Slice, From: reflect.String}:  getCopyFuncToSliceFromString,
		{To: reflect.Slice, From: reflect.Slice}:   getCopyFuncToSliceFromSlice,
		{To: reflect.Slice, From: reflect.Array}:   getCopyFuncToSliceFromArray,
		{To: reflect.Slice, From: reflect.Map}:     getCopyFuncToSliceFromMap,
	})
}

func copyReflectToSliceFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.Set(reflect.Zero(target.Type()))
	return true
}
func getCopyFuncToSliceFromInvalid(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncSetZero(c, tTyp, sTyp)
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
func getCopyFuncToSliceFromString(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	switch tTyp {
	default:
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }

	case typeOfByteSlice:
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			source = reflect.ValueOf([]byte(source.String()))
			target.Set(source)
			return true
		}

	case typeOfRuneSlice:
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			source = reflect.ValueOf([]rune(source.String()))
			target.Set(source)
			return true
		}
	}
}

func copyReflectToSliceFromArray(c *copyContext, target, source reflect.Value) bool {
	// sVal := source.Slice(0, source.Len())
	// return copyReflectToSliceFromSlice(c, target, sVal)
	return copyReflectToSliceFromSlice(c, target, source)
}
func getCopyFuncToSliceFromArray(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	// sVal := source.Slice(0, source.Len())
	// return copyReflectToSliceFromSlice(c, target, sVal)
	return getCopyFuncToSliceFromSlice(c, tTyp, sTyp)
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
func getCopyFuncToSliceFromSlice(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	if !c.Clone && tTyp == sTyp {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			target.Set(source)
			return true
		}
	}

	if sTyp.Kind() != reflect.Slice && sTyp.Kind() != reflect.Array {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}
	if !isCopyTypeConvertible(tTyp.Elem(), sTyp.Elem()) {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}

	fnElemCopy := getCopyFunc(c, tTyp.Elem(), sTyp.Elem())
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		l := source.Len()
		if l > target.Len() {
			target.Set(reflect.MakeSlice(tTyp, l, l))
		}

		for i := 0; i < l; i++ {
			tItem := target.Index(i)
			sItem := source.Index(i)

			c.PushField(strconv.Itoa(i))
			fnElemCopy(c, tItem, sItem)
			c.PopField()
		}

		return true
	}
}

func copyReflectToSliceFromMap(c *copyContext, target, source reflect.Value) bool {
	sTyp := source.Type()

	if isEmptyStruct(sTyp.Elem()) {
		return copyReflectToSliceFromMapWithEmptyStructValue(c, target, source)
	}
	return false
}
func getCopyFuncToSliceFromMap(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	// sTyp := source.Type()

	if isEmptyStruct(sTyp.Elem()) {
		// return copyReflectToSliceFromMapWithEmptyStructValue(c, target, source)
		return getCopyFuncToArrayFromMapWithEmptyStructValue(c, tTyp, sTyp)
	}
	return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
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
func getCopyFuncToSliceFromMapWithEmptyStructValue(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	// tTyp := target.Type()
	// sTyp := source.Type()

	if isCopyTypeConvertibleV2(sTyp.Elem(), sTyp.Key()) {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}

	fnElemCopy := getCopyFunc(c, tTyp.Elem(), sTyp.Key())
	return func(c *copyContext, target, source reflect.Value) (end bool) {
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
			fnElemCopy(c, tItem, sItem)
			c.PopField()
		}

		return true
	}
}
