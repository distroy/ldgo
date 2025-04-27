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
		{To: reflect.Array, From: reflect.Invalid}: copyReflectToArrayFromInvalid,
		{To: reflect.Array, From: reflect.String}:  copyReflectToArrayFromString,
		{To: reflect.Array, From: reflect.Slice}:   copyReflectToArrayFromSlice,
		{To: reflect.Array, From: reflect.Array}:   copyReflectToArrayFromArray,
		{To: reflect.Array, From: reflect.Map}:     copyReflectToArrayFromMapWithEmptyStructValue,
	})

	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.Array, From: reflect.Invalid}: getCopyFuncToArrayFromInvalid,
		{To: reflect.Array, From: reflect.String}:  getCopyFuncToArrayFromString,
		{To: reflect.Array, From: reflect.Slice}:   getCopyFuncToArrayFromSlice,
		{To: reflect.Array, From: reflect.Array}:   getCopyFuncToArrayFromArray,
		{To: reflect.Array, From: reflect.Map}:     getCopyFuncToArrayFromMapWithEmptyStructValue,
	})
}

func copyReflectToArrayFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.Set(reflect.Zero(target.Type()))
	return true
}

func getCopyFuncToArrayFromInvalid(c *copyContext, tType, sType reflect.Type) copyFuncType {
	return getCopyFuncSetZero(c, tType, sType)
}

func copyReflectToArrayFromString(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()

	sVal := source
	switch tTyp.Elem().Kind() {
	default:
		return false

	// case typeOfByteSlice:
	case reflect.Uint8:
		sVal = reflect.ValueOf([]byte(sVal.String()))

	case reflect.Int32:
		sVal = reflect.ValueOf([]rune(sVal.String()))
	}
	sTyp := source.Type()

	l := sVal.Len()
	if l > target.Len() {
		c.AddErrorf("%s has %d elements, can not convert to %s",
			sTyp.String(), l, tTyp.String())
		l = target.Len()
	}
	for i := 0; i < l; i++ {
		tItem := target.Index(i)
		sItem := sVal.Index(i)
		tItem.Set(sItem)
	}

	return true
}

func getCopyFuncToArrayFromString(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	// var fnElem *copyFuncValue
	var fnValConvert func(v reflect.Value) reflect.Value
	switch tTyp.Elem().Kind() {
	default:
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			return false
		}

	// case typeOfByteSlice:
	case reflect.Uint8:
		fnValConvert = func(v reflect.Value) reflect.Value {
			return reflect.ValueOf([]byte(v.String()))
		}

	case reflect.Int32:
		fnValConvert = func(v reflect.Value) reflect.Value {
			return reflect.ValueOf([]rune(v.String()))
		}
	}

	return func(c *copyContext, target, source reflect.Value) (end bool) {
		sVal := source
		sVal = fnValConvert(sVal)

		l := sVal.Len()
		if l > target.Len() {
			c.AddErrorf("%s has %d elements, can not convert to %s",
				sTyp.String(), l, tTyp.String())
			l = target.Len()
		}
		for i := 0; i < l; i++ {
			tItem := target.Index(i)
			sItem := sVal.Index(i)
			tItem.Set(sItem)
		}
		return true
	}
}

func copyReflectToArrayFromArray(c *copyContext, target, source reflect.Value) bool {
	// sVal := source.Slice(0, source.Len())
	// return copyReflectToArrayFromSlice(c, target, sVal)
	return copyReflectToArrayFromSlice(c, target, source)
}

func getCopyFuncToArrayFromArray(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncToArrayFromSlice(c, tTyp, sTyp)
}

func copyReflectToArrayFromSlice(c *copyContext, target, source reflect.Value) bool {
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
		c.AddErrorf("%s has %d elements, can not convert to %s", sTyp.String(), l, tTyp.String())
		l = target.Len()
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

func getCopyFuncToArrayFromSlice(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	if !c.Clone && tTyp == sTyp {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			target.Set(source)
			return true
		}
	}

	if sTyp.Kind() != reflect.Slice && sTyp.Kind() != reflect.Array {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			return false
		}
	}
	if !isCopyTypeConvertibleV2(tTyp.Elem(), sTyp.Elem()) {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			return false
		}
	}

	pfe, done := getCopyFunc(c, tTyp.Elem(), sTyp.Elem())

	return func(c *copyContext, target, source reflect.Value) (end bool) {
		l := source.Len()
		if l > target.Len() {
			c.AddErrorf("%s has %d elements, can not convert to %s", sTyp.String(), l, tTyp.String())
			l = target.Len()
		}

		for i := 0; i < l; i++ {
			tItem := target.Index(i)
			sItem := source.Index(i)

			c.PushField(strconv.Itoa(i))
			// copyReflect(c, tItem, sItem)
			done()
			(*pfe)(c, tItem, sItem)
			c.PopField()
		}

		return true
	}
}

func copyReflectToArrayFromMapWithEmptyStructValue(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	sTyp := source.Type()

	if !isEmptyStruct(sTyp.Elem()) {
		return false
	}

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

func getCopyFuncToArrayFromMapWithEmptyStructValue(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	if !isEmptyStruct(sTyp.Elem()) {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			return false
		}
	}

	if isCopyTypeConvertible(sTyp.Elem(), sTyp.Key()) {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			return false
		}
	}

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
			copyReflect(c, tItem, sItem)
			c.PopField()
		}

		return true
	}
}
