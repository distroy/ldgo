/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
)

func clearCopyStructIgnoreField(c *context, v reflect.Value, info *copyStructInfo) {
	for _, f := range info.Ignores {
		field := v.Field(f.Index)
		field.Set(reflect.Zero(f.Type))
	}
}

func copyReflectToStruct(c *context, target, source reflect.Value) bool {
	// source, _ = prepareCopySourceReflect(c, source)
	source, _ = indirectSourceReflect(source)

	switch source.Kind() {
	default:
		return false

	case reflect.Invalid:
		tTyp := target.Type()
		target.Set(reflect.Zero(tTyp))

	case reflect.Struct:
		return copyReflectToStructFromStruct(c, target, source)

	case reflect.Map:
		return copyReflectToStructFromMap(c, target, source)

	}

	return true
}

func copyReflectToStructFromStruct(c *context, target, source reflect.Value) bool {
	tTyp := target.Type()
	tInfo := getCopyTypeInfo(tTyp)

	sTyp := source.Type()
	sInfo := getCopyTypeInfo(sTyp)
	if !c.IsDeep && tTyp == sTyp {
		target.Set(source)
		clearCopyStructIgnoreField(c, target, tInfo)
		return true
	}

	for _, sFieldInfo := range sInfo.Fields {
		tFieldInfo := tInfo.Fields[sFieldInfo.Name]
		if tFieldInfo == nil {
			continue
		}

		tField := target.Field(tFieldInfo.Index)
		sField := source.Field(sFieldInfo.Index)

		c.PushField(tFieldInfo.Name)
		copyReflect(c, tField, sField)
		c.PopField()
	}

	return true
}

func copyReflectToStructFromMap(c *context, target, source reflect.Value) bool {
	tTyp := target.Type()
	tInfo := getCopyTypeInfo(tTyp)

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
