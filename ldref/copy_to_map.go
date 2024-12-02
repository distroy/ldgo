/*
 * Copyright (C) distroy
 */

package ldref

import (
	"fmt"
	"reflect"
)

func init() {
	registerCopyFunc(map[copyPair]copyFuncType{
		{To: reflect.Map, From: reflect.Invalid}: copyReflectToMapFromInvalid,
		{To: reflect.Map, From: reflect.Struct}:  copyReflectToMapFromStruct,
		{To: reflect.Map, From: reflect.Map}:     copyReflectToMapFromMap,
		{To: reflect.Map, From: reflect.Slice}:   copyReflectToMapFromSlice,
		{To: reflect.Map, From: reflect.Array}:   copyReflectToMapFromArray,
	})
}

func isEmptyStruct(typ reflect.Type) bool {
	if typ.Kind() == reflect.Struct && typ.NumField() == 0 {
		return true
	}
	return false
}

func copyReflectToMapFromInvalid(c *copyContext, target, source reflect.Value) bool {
	target.Set(reflect.Zero(target.Type()))
	return true
}

func copyReflectToMapFromArray(c *copyContext, target, source reflect.Value) bool {
	// source = source.Slice(0, source.Len())
	return copyReflectToMapFromSlice(c, target, source)
}

func isStructFieldNilPtr(v reflect.Value) bool {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}

		v = v.Elem()
	}

	return false
}

func copyReflectToMapFromStruct(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	keyTyp := tTyp.Key()
	valTyp := tTyp.Elem()

	if keyTyp.Kind() != reflect.String && keyTyp.Kind() != reflect.Interface {
		return false
	}

	if target.IsNil() {
		if !target.CanAddr() {
			return false
		}
		target.Set(reflect.MakeMap(target.Type()))
	}

	sInfo := getCopyTypeInfo(source.Type(), c.SourceTag)
	for _, sFieldInfo := range sInfo.Fields {
		sField := source.Field(sFieldInfo.Index)
		if isStructFieldNilPtr(sField) {
			continue
		}

		key := reflect.ValueOf(sFieldInfo.Name)
		value := reflect.New(valTyp).Elem()

		c.PushField(sFieldInfo.Name)
		copyReflect(c, value, sField)
		c.PopField()

		target.SetMapIndex(key, value)
	}

	return true
}

func copyReflectToMapFromMap(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	tKeyTyp := tTyp.Key()
	tValTyp := tTyp.Elem()

	sTyp := source.Type()
	sKeyTyp := sTyp.Key()
	sValTyp := sTyp.Elem()

	if !c.Clone && tTyp == sTyp {
		target.Set(source)
		return true
	}

	if !isCopyTypeConvertible(tKeyTyp, sKeyTyp) {
		return false
	}
	if !isCopyTypeConvertible(tValTyp, sValTyp) {
		return false
	}

	// if target.IsNil() {
	target.Set(reflect.MakeMap(target.Type()))
	// }

	for sIt := source.MapRange(); sIt.Next(); {
		sKey := sIt.Key()
		sVal := sIt.Value()

		tKey := reflect.New(tKeyTyp).Elem()
		tVal := reflect.New(tValTyp).Elem()

		c.PushField(fmt.Sprintf("%v(key)", sKey.Interface()))
		keyEnd := copyReflect(c, tKey, sKey)
		c.PopField()
		if !keyEnd {
			continue
		}

		c.PushField(fmt.Sprintf("%v(val)", sKey.Interface()))
		valEnd := copyReflect(c, tVal, sVal)
		c.PopField()
		if !valEnd {
			continue
		}

		target.SetMapIndex(tKey, tVal)
	}

	return true
}

func copyReflectToMapFromSlice(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	tValTyp := tTyp.Elem()

	if isEmptyStruct(tValTyp) {
		return copyReflectToMapFromSliceWithEmptyStructValue(c, target, source)
	}

	return false
	// return copyReflectToMapFromSliceWithIndexBeKey(c, target, source)
}

func copyReflectToMapFromSliceWithIndexBeKey(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	tKeyTyp := tTyp.Key()
	tValTyp := tTyp.Elem()

	sTyp := source.Type()
	sKeyTyp := reflect.TypeOf(int(0))
	sValTyp := sTyp.Elem()

	if !isCopyTypeConvertible(tKeyTyp, sKeyTyp) {
		return false
	}

	if !isCopyTypeConvertible(tValTyp, sValTyp) {
		return false
	}

	// if target.IsNil() {
	target.Set(reflect.MakeMap(target.Type()))
	// }

	for i, l := 0, source.Len(); i < l; i++ {
		sKey := reflect.ValueOf(i)
		sVal := source.Index(i)

		tKey := reflect.New(tKeyTyp).Elem()
		tVal := reflect.New(tValTyp).Elem()

		c.PushField(fmt.Sprintf("%d(key)", i))
		keyEnd := copyReflect(c, tKey, sKey)
		c.PopField()
		if !keyEnd {
			continue
		}

		c.PushField(fmt.Sprintf("%d(val)", i))
		valEnd := copyReflect(c, tVal, sVal)
		c.PopField()
		if !valEnd {
			continue
		}

		target.SetMapIndex(tKey, tVal)
	}

	return true
}

func copyReflectToMapFromSliceWithEmptyStructValue(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	tElemTyp := tTyp.Key()

	sTyp := source.Type()
	sElemTyp := sTyp.Elem()

	if !isCopyTypeConvertible(tElemTyp, sElemTyp) {
		return false
	}

	target.Set(reflect.MakeMap(target.Type()))

	tValue := reflect.ValueOf(struct{}{})
	for i, l := 0, source.Len(); i < l; i++ {
		// sKey := reflect.ValueOf(i)
		sElem := source.Index(i)

		tElem := reflect.New(tElemTyp).Elem()

		c.PushField(fmt.Sprintf("%d", i))
		keyEnd := copyReflect(c, tElem, sElem)
		c.PopField()
		if !keyEnd {
			continue
		}

		target.SetMapIndex(tElem, tValue)
	}

	return true
}
