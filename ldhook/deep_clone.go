/*
 * Copyright (C) distroy
 */

package ldhook

import "reflect"

func DeepClone(i interface{}) interface{} {
	x := deepClone(reflect.ValueOf(i))
	if isValueNil(x) {
		return nil
	}
	return x.Interface()
}

func deepClone(x0 reflect.Value) reflect.Value {
	if isValueNil(x0) {
		return reflect.Value{}
	}

	switch x0.Kind() {
	case reflect.Struct:
		return deepCloneStruct(x0)

	case reflect.Ptr:
		x1 := reflect.New(x0.Type().Elem())
		x1.Elem().Set(deepClone(x0.Elem()))
		return x1

	case reflect.Array:
		return deepCloneArray(x0)

	case reflect.Slice:
		return deepCloneSlice(x0)

	case reflect.Map:
		return deepCloneMap(x0)
	}

	return x0
}

func deepCloneStruct(x0 reflect.Value) reflect.Value {
	x1 := reflect.New(x0.Type()).Elem()
	for i, n := 0, x0.NumField(); i < n; i++ {
		f0 := x0.Field(i)
		if isValueNil(f0) {
			continue
		}
		f1 := x1.Field(i)
		f1.Set(deepClone(f0))
	}
	return x1
}

func deepCloneArray(x0 reflect.Value) reflect.Value {
	l0 := x0.Len()
	x1 := reflect.New(reflect.ArrayOf(l0, x0.Type().Elem())).Elem()
	for i := 0; i < l0; i++ {
		v0 := x0.Index(i)
		if isValueNil(v0) {
			continue
		}
		v0 = deepClone(v0)
		x1.Index(i).Set(v0)
	}
	return x1
}

func deepCloneSlice(x0 reflect.Value) reflect.Value {
	l0 := x0.Len()
	x1 := reflect.MakeSlice(x0.Type(), l0, l0)
	for i := 0; i < l0; i++ {
		v0 := x0.Index(i)
		if isValueNil(v0) {
			continue
		}
		v0 = deepClone(v0)
		x1.Index(i).Set(v0)
	}
	return x1
}

func deepCloneMap(x0 reflect.Value) reflect.Value {
	x1 := reflect.MakeMap(x0.Type())
	for it := x0.MapRange(); it.Next(); {
		key := it.Key()
		val := it.Value()
		if isValueNil(key) || isValueNil(val) {
			continue
		}
		key = deepClone(key)
		val = deepClone(val)
		x1.SetMapIndex(key, val)
	}
	return x1
}