/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"strconv"
)

func copyReflectToArray(c *context, target, source reflect.Value) bool {
	return copyReflectToSlice(c, target, source)
}

func copyReflectToSlice(c *context, target, source reflect.Value) bool {
	// source, _ = prepareCopySourceReflect(c, source)
	source, _ = indirectSourceReflect(source)

	switch source.Kind() {
	default:
		return false

	case reflect.Invalid:
		target.Set(reflect.Zero(target.Type()))

	case reflect.String:
		fallthrough
	case reflect.Array, reflect.Slice:
		return copyReflectToSliceFromSlice(c, target, source)
	}

	return true
}

func copyReflectToSliceFromSlice(c *context, target, source reflect.Value) bool {
	if source.Kind() == reflect.String {
		switch target.Interface().(type) {
		case []byte:
			source = reflect.ValueOf([]byte(source.String()))

		case []rune:
			source = reflect.ValueOf([]rune(source.String()))
		}
	}

	tTyp := target.Type()
	sTyp := source.Type()
	if !c.IsDeep && tTyp == sTyp {
		target.Set(source)
		return true
	}

	if source.Kind() == reflect.Invalid {
		target.Set(reflect.Zero(tTyp))
		return true
	}

	if source.Kind() != reflect.Array {
		source = source.Slice(0, source.Len())
	}

	if source.Kind() != reflect.Slice || !isCopyTypeConvertible(tTyp.Elem(), sTyp.Elem()) {
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

	for i := 0; i < l; i++ {
		tItem := target.Index(i)
		sItem := source.Index(i)

		c.PushField(strconv.Itoa(i))
		copyReflect(c, tItem, sItem)
		c.PopField()
	}

	return true
}