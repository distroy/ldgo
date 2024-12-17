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
	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		{To: reflect.Map, From: reflect.Invalid}: getCopyFuncToMapFromInvalid,
		{To: reflect.Map, From: reflect.Struct}:  getCopyFuncToMapFromStruct,
		{To: reflect.Map, From: reflect.Map}:     getCopyFuncToMapFromMap,
		{To: reflect.Map, From: reflect.Slice}:   getCopyFuncToMapFromSlice,
		{To: reflect.Map, From: reflect.Array}:   getCopyFuncToMapFromArray,
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
func getCopyFuncToMapFromInvalid(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncSetZero(c, tTyp, sTyp)
}

func copyReflectToMapFromArray(c *copyContext, target, source reflect.Value) bool {
	// source = source.Slice(0, source.Len())
	return copyReflectToMapFromSlice(c, target, source)
}
func getCopyFuncToMapFromArray(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	return getCopyFuncToMapFromSlice(c, tTyp, sTyp)
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
func getCopyFuncToMapFromStruct(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	// tTyp := target.Type()
	keyTyp := tTyp.Key()
	valTyp := tTyp.Elem()

	if keyTyp.Kind() != reflect.String && keyTyp.Kind() != reflect.Interface {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}

	return func(c *copyContext, target, source reflect.Value) (end bool) {
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
			copyReflectV2(c, value, sField)
			c.PopField()

			target.SetMapIndex(key, value)
		}

		return true
	}
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

	if !isCopyTypeConvertible(tKeyTyp, sKeyTyp) || !isCopyTypeConvertible(tValTyp, sValTyp) {
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
func getCopyFuncToMapFromMap(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	if !c.Clone && tTyp == sTyp {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			target.Set(source)
			return true
		}
	}

	tKeyTyp := tTyp.Key()
	tValTyp := tTyp.Elem()

	sKeyTyp := sTyp.Key()
	sValTyp := sTyp.Elem()

	if !isCopyTypeConvertibleV2(tKeyTyp, sKeyTyp) || !isCopyTypeConvertibleV2(tValTyp, sValTyp) {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}

	pfk := getCopyFunc(c, tKeyTyp, sKeyTyp)
	pfv := getCopyFunc(c, tValTyp, sValTyp)
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		// if target.IsNil() {
		target.Set(reflect.MakeMap(tTyp))
		// }

		for sIt := source.MapRange(); sIt.Next(); {
			sKey := sIt.Key()
			sVal := sIt.Value()

			tKey := reflect.New(tKeyTyp).Elem()
			tVal := reflect.New(tValTyp).Elem()

			c.PushField(fmt.Sprintf("%v(key)", sKey.Interface()))
			keyEnd := (*pfk)(c, tKey, sKey)
			c.PopField()
			if !keyEnd {
				continue
			}

			c.PushField(fmt.Sprintf("%v(val)", sKey.Interface()))
			valEnd := (*pfv)(c, tVal, sVal)
			c.PopField()
			if !valEnd {
				continue
			}

			target.SetMapIndex(tKey, tVal)
		}

		return true
	}
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
func getCopyFuncToMapFromSlice(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	// tTyp := target.Type()
	tValTyp := tTyp.Elem()

	if isEmptyStruct(tValTyp) {
		// return copyReflectToMapFromSliceWithEmptyStructValue(c, target, source)
		return getCopyFuncToMapFromSliceWithEmptyStructValue(c, tTyp, sTyp)
	}

	return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
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
func getCopyFuncToMapFromSliceWithIndexBeKey(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	// tTyp := target.Type()
	tKeyTyp := tTyp.Key()
	tValTyp := tTyp.Elem()

	// sTyp := source.Type()
	sKeyTyp := reflect.TypeOf(int(0))
	sValTyp := sTyp.Elem()

	if !isCopyTypeConvertibleV2(tKeyTyp, sKeyTyp) || !isCopyTypeConvertibleV2(tValTyp, sValTyp) {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}

	pfk := getCopyFunc(c, tKeyTyp, sKeyTyp)
	pfv := getCopyFunc(c, tValTyp, sValTyp)
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		target.Set(reflect.MakeMap(target.Type()))

		for i, l := 0, source.Len(); i < l; i++ {
			sKey := reflect.ValueOf(i)
			sVal := source.Index(i)

			tKey := reflect.New(tKeyTyp).Elem()
			tVal := reflect.New(tValTyp).Elem()

			c.PushField(fmt.Sprintf("%d(key)", i))
			keyEnd := (*pfk)(c, tKey, sKey)
			c.PopField()
			if !keyEnd {
				continue
			}

			c.PushField(fmt.Sprintf("%d(val)", i))
			valEnd := (*pfv)(c, tVal, sVal)
			c.PopField()
			if !valEnd {
				continue
			}

			target.SetMapIndex(tKey, tVal)
		}

		return true
	}
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
func getCopyFuncToMapFromSliceWithEmptyStructValue(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	tElemTyp := tTyp.Key()
	sElemTyp := sTyp.Elem()

	if !isCopyTypeConvertibleV2(tElemTyp, sElemTyp) {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}

	pfe := getCopyFunc(c, tElemTyp, sElemTyp)
	return func(c *copyContext, target, source reflect.Value) (end bool) {
		target.Set(reflect.MakeMap(target.Type()))

		tValue := reflect.ValueOf(struct{}{})
		for i, l := 0, source.Len(); i < l; i++ {
			// sKey := reflect.ValueOf(i)
			sElem := source.Index(i)

			tElem := reflect.New(tElemTyp).Elem()

			c.PushField(fmt.Sprintf("%d", i))
			keyEnd := (*pfe)(c, tElem, sElem)
			c.PopField()
			if !keyEnd {
				continue
			}

			target.SetMapIndex(tElem, tValue)
		}

		return true
	}
}
