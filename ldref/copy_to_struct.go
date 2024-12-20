/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
)

func init() {
	registerCopyFunc(map[copyPair]copyFuncType{
		// {To: reflect.Struct, From: reflect.Invalid}: copyReflectToStructFromInvalid,
		{To: reflect.Struct, From: reflect.Struct}: copyReflectToStructFromStruct,
		{To: reflect.Struct, From: reflect.Map}:    copyReflectToStructFromMap,
	})

	registerGetCopyFunc(map[copyPair]getCopyFuncType{
		// {To: reflect.Struct, From: reflect.Invalid}: copyReflectToStructFromInvalid,
		{To: reflect.Struct, From: reflect.Struct}: getCopyFuncToStructFromStruct,
		{To: reflect.Struct, From: reflect.Map}:    getCopyFuncToStructFromMap,
	})
}

func clearCopyStructIgnoreField(c *copyContext, v reflect.Value, info *copyStructValue) {
	for _, f := range info.Ignores {
		field := v.Field(f.Index)
		field.Set(f.TypeZero)
	}
}

func copyReflectToStructFromStruct(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	tInfo := getCopyTypeInfo(tTyp, c.TargetTag)

	sTyp := source.Type()
	sInfo := getCopyTypeInfo(sTyp, c.SourceTag)
	if !c.Clone && tTyp == sTyp && c.TargetTag == c.SourceTag {
		target.Set(source)
		clearCopyStructIgnoreField(c, target, tInfo)
		return true
	}

	if !source.CanAddr() {
		if tTyp == sTyp && c.TargetTag == c.SourceTag {
			target.Set(source)
			source = target
			clearCopyStructIgnoreField(c, target, tInfo)

		} else {
			tmp := reflect.New(sTyp).Elem()
			tmp.Set(source)
			source = tmp
		}
	}

	for _, sFieldInfo := range sInfo.Fields {
		tFieldInfo := tInfo.Fields[sFieldInfo.Name]
		if tFieldInfo == nil {
			continue
		}

		// tFieldAddr := unsafe.Pointer(target.Field(tFieldInfo.Index).UnsafeAddr())
		// tField := reflect.NewAt(tFieldInfo.Type, tFieldAddr).Elem()
		// sFieldAddr := unsafe.Pointer(source.Field(sFieldInfo.Index).UnsafeAddr())
		// sField := reflect.NewAt(sFieldInfo.Type, sFieldAddr).Elem()
		tField := refStructFieldByCopyFieldInfo(target, tFieldInfo)
		sField := refStructFieldByCopyFieldInfo(source, sFieldInfo)

		c.PushField(tFieldInfo.Name)
		copyReflect(c, tField, sField)
		c.PopField()
	}

	return true
}

func getCopyFuncToStructFromStruct(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	tInfo := getCopyTypeInfo(tTyp, c.TargetTag)
	sInfo := getCopyTypeInfo(sTyp, c.SourceTag)
	if !c.Clone && tTyp == sTyp && c.TargetTag == c.SourceTag {
		return func(c *copyContext, target, source reflect.Value) (end bool) {
			target.Set(source)
			clearCopyStructIgnoreField(c, target, tInfo)
			return true
		}
	}

	copyFields := make([]copyFuncType, 0, len(sInfo.Fields))
	for _, _v := range sInfo.Fields {
		sFieldInfo := _v
		tFieldInfo := tInfo.Fields[sFieldInfo.Name]
		if tFieldInfo == nil {
			continue
		}

		pff, done := getCopyFunc(c, tFieldInfo.Type, sFieldInfo.Type)
		copyFields = append(copyFields, func(c *copyContext, target, source reflect.Value) (end bool) {
			// tFieldAddr := unsafe.Pointer(target.Field(tFieldInfo.Index).UnsafeAddr())
			// tField := reflect.NewAt(tFieldInfo.Type, tFieldAddr).Elem()
			// sFieldAddr := unsafe.Pointer(source.Field(sFieldInfo.Index).UnsafeAddr())
			// sField := reflect.NewAt(sFieldInfo.Type, sFieldAddr).Elem()
			tField := refStructFieldByCopyFieldInfo(target, tFieldInfo)
			sField := refStructFieldByCopyFieldInfo(source, sFieldInfo)

			done()
			c.PushField(tFieldInfo.Name)
			// copyReflect(c, tField, sField)
			end = (*pff)(c, tField, sField)
			c.PopField()

			return end
		})
	}

	return func(c *copyContext, target, source reflect.Value) (end bool) {
		if !source.CanAddr() {
			if tTyp == sTyp && c.TargetTag == c.SourceTag {
				target.Set(source)
				source = target
				clearCopyStructIgnoreField(c, target, tInfo)

			} else {
				tmp := reflect.New(sTyp).Elem()
				tmp.Set(source)
				source = tmp
			}
		}

		for _, fnCopy := range copyFields {
			fnCopy(c, target, source)
		}

		return true
	}
}

func copyReflectToStructFromMap(c *copyContext, target, source reflect.Value) bool {
	tTyp := target.Type()
	tInfo := getCopyTypeInfo(tTyp, c.TargetTag)

	sTyp := source.Type()
	if sTyp.Key().Kind() != reflect.String {
		return false
	}

	it := source.MapRange()
	for it.Next() {
		key := it.Key().String()
		tFieldInfo := tInfo.Fields[key]
		if tFieldInfo == nil {
			continue
		}

		tField := target.Field(tFieldInfo.Index)
		value := it.Value()

		c.PushField(tFieldInfo.Name)
		copyReflect(c, tField, value)
		c.PopField()
	}
	return true
}

func getCopyFuncToStructFromMap(c *copyContext, tTyp, sTyp reflect.Type) copyFuncType {
	// tTyp := target.Type()
	tInfo := getCopyTypeInfo(tTyp, c.TargetTag)

	// sTyp := source.Type()
	if sTyp.Key().Kind() != reflect.String {
		return func(c *copyContext, target, source reflect.Value) (end bool) { return false }
	}

	fnFieldCopies := make([]copyFuncType, 0, len(tInfo.Fields))
	for _, v := range tInfo.Fields {
		tFieldInfo := v
		keyVal := reflect.ValueOf(tFieldInfo.Name)
		pff, done := getCopyFunc(c, tFieldInfo.Type, sTyp.Elem())
		fnFieldCopies = append(fnFieldCopies, func(c *copyContext, target, source reflect.Value) (end bool) {
			sField := source.MapIndex(keyVal)

			// tFieldAddr := unsafe.Pointer(target.Field(tFieldInfo.Index).UnsafeAddr())
			// tField := reflect.NewAt(tFieldInfo.Type, tFieldAddr).Elem()
			tField := refStructFieldByCopyFieldInfo(target, tFieldInfo)
			if !sField.IsValid() {
				tField.Set(tFieldInfo.TypeZero)
				return true
			}

			done()
			c.PushField(tFieldInfo.Name)
			ok := (*pff)(c, tField, sField)
			c.PopField()
			return ok
		})
	}

	return func(c *copyContext, target, source reflect.Value) (end bool) {
		for _, fn := range fnFieldCopies {
			fn(c, target, source)
		}
		return true
	}
}
