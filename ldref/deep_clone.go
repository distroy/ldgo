/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"unsafe"
)

func DeepClone[T any](d T) T { return deepCloneV2(d) }

func deepCloneV1[T any](d T) T {
	var i interface{} = d
	if x, ok := i.(reflect.Value); ok {
		var r interface{} = deepCloneRef(x)
		return r.(T)
	}

	x := deepCloneRef(reflect.ValueOf(i))
	if x.Kind() == reflect.Invalid {
		var zero T
		return zero
	}

	return x.Interface().(T)
}
func deepCloneV2[T any](d T) T {
	var i interface{} = d
	if x, ok := i.(reflect.Value); ok {
		var r interface{} = deepCloneRefV2(x)
		return r.(T)
	}

	x := deepCloneRefV2(reflect.ValueOf(i))
	if x.Kind() == reflect.Invalid {
		var zero T
		return zero
	}

	return x.Interface().(T)
}

func deepCloneRef(x0 reflect.Value) reflect.Value {
	if x0.Kind() == reflect.Interface {
		if x0.IsNil() {
			return x0
		}

		x0 = x0.Elem()
	}

	switch x0.Kind() {
	case reflect.Invalid:
		return x0

	case reflect.Struct:
		return deepCloneStruct(x0)

	case reflect.Ptr:
		return deepClonePtr(x0)

	case reflect.Array:
		return deepCloneArray(x0)

	case reflect.Slice:
		return deepCloneSlice(x0)

	case reflect.Map:
		return deepCloneMap(x0)
	}

	return x0
}
func getDeepCloneFunc(t reflect.Type) cloneFuncType {
	switch refKindOfType(t) {
	case reflect.Invalid:
		return func(x0 reflect.Value) reflect.Value { return x0 }

	case reflect.Interface:
		return func(x0 reflect.Value) reflect.Value {
			if x0.IsNil() {
				return x0
			}
			return deepCloneRefV2(x0.Elem())
		}

	case reflect.Struct:
		return getDeepCloneFuncStruct(t)

	case reflect.Ptr:
		return getDeepCloneFuncPtr(t)

	case reflect.Array:
		return getDeepCloneFuncArray(t)

	case reflect.Slice:
		return getDeepCloneFuncSlice(t)

	case reflect.Map:
		return getDeepCloneFuncMap(t)
	}

	return func(x0 reflect.Value) reflect.Value { return x0 }
}
func deepCloneRefV2(x0 reflect.Value) reflect.Value {
	t := refTypeOfValue(x0)
	pf := getCloneFuncByPool(t, true)
	return (*pf)(x0)
}

func deepClonePtr(x0 reflect.Value) reflect.Value {
	if x0.IsNil() {
		return x0
	}

	x1 := reflect.New(x0.Type().Elem())

	x1.Elem().Set(deepCloneRef(x0.Elem()))
	return x1
}
func getDeepCloneFuncPtr(t reflect.Type) cloneFuncType {
	tt := t.Elem()
	pf := getCloneFuncByPool(tt, true)
	return func(x0 reflect.Value) reflect.Value {
		if x0.IsNil() {
			return x0
		}

		x1 := reflect.New(tt)

		x1.Elem().Set((*pf)(x0.Elem()))
		return x1
	}
}

func deepCloneStruct(x0 reflect.Value) reflect.Value {
	if isSyncType(x0.Type()) {
		return reflect.Zero(x0.Type())
	}

	x1 := reflect.New(x0.Type()).Elem()
	if !x0.CanAddr() {
		x1.Set(x0)
		x0 = x1
	}

	for i, n := 0, x0.NumField(); i < n; i++ {
		f0 := x0.Field(i)
		f1 := x1.Field(i)
		// f1.Set(deepClone(f0))

		a0 := unsafe.Pointer(f0.UnsafeAddr())
		o0 := reflect.NewAt(f0.Type(), a0).Elem()

		a1 := unsafe.Pointer(f1.UnsafeAddr())
		o1 := reflect.NewAt(f1.Type(), a1).Elem()
		o1.Set(deepCloneRef(o0))
	}
	return x1
}
func getDeepCloneFuncStruct(t reflect.Type) cloneFuncType {
	if isSyncType(t) {
		zero := reflect.Zero(t)
		return func(x0 reflect.Value) reflect.Value { return zero }
	}

	n := t.NumField()
	fnFields := make([]func(x0, x1 reflect.Value), 0, n)
	for _i := 0; _i < n; _i++ {
		i := _i
		sf := t.Field(i)
		pf := getCloneFuncByPool(sf.Type, true)
		fnFields = append(fnFields, func(x0, x1 reflect.Value) {
			f0 := refStructField(x0, i, &sf)
			f1 := refStructField(x1, i, &sf)

			f1.Set((*pf)(f0))
		})
	}

	return func(x0 reflect.Value) reflect.Value {
		x1 := reflect.New(t).Elem()
		if !x0.CanAddr() {
			x1.Set(x0)
			x0 = x1
		}

		for _, fn := range fnFields {
			fn(x0, x1)
		}
		return x1
	}
}

func deepCloneArray(x0 reflect.Value) reflect.Value {
	l0 := x0.Len()
	x1 := reflect.New(reflect.ArrayOf(l0, x0.Type().Elem())).Elem()
	for i := 0; i < l0; i++ {
		v0 := x0.Index(i)
		v0 = deepCloneRef(v0)
		x1.Index(i).Set(v0)
	}
	return x1
}
func getDeepCloneFuncArray(t reflect.Type) cloneFuncType {
	pf := getCloneFuncByPool(t.Elem(), true)
	return func(x0 reflect.Value) reflect.Value {
		l0 := x0.Len()
		x1 := reflect.New(t).Elem()
		for i := 0; i < l0; i++ {
			v0 := x0.Index(i)
			v1 := (*pf)(v0)
			x1.Index(i).Set(v1)
		}
		return x1
	}
}

func deepCloneSlice(x0 reflect.Value) reflect.Value {
	if x0.IsNil() {
		return x0
	}
	l0 := x0.Len()
	x1 := reflect.MakeSlice(reflect.SliceOf(x0.Type().Elem()), l0, l0)
	for i := 0; i < l0; i++ {
		v0 := x0.Index(i)
		v0 = deepCloneRef(v0)
		x1.Index(i).Set(v0)
	}
	return x1
}
func getDeepCloneFuncSlice(t reflect.Type) cloneFuncType {
	tt := t.Elem()
	pf := getCloneFuncByPool(tt, true)
	return func(x0 reflect.Value) reflect.Value {
		if x0.IsNil() {
			return x0
		}
		l0 := x0.Len()
		// x1 := reflect.MakeSlice(reflect.SliceOf(x0.Type().Elem()), l0, l0)
		x1 := reflect.MakeSlice(t, l0, l0)
		for i := 0; i < l0; i++ {
			v0 := x0.Index(i)
			v1 := (*pf)(v0)
			x1.Index(i).Set(v1)
		}
		return x1
	}
}

func deepCloneMap(x0 reflect.Value) reflect.Value {
	x1 := reflect.MakeMap(x0.Type())
	for it := x0.MapRange(); it.Next(); {
		key := it.Key()
		val := it.Value()

		key = deepCloneRef(key)
		val = deepCloneRef(val)
		x1.SetMapIndex(key, val)
	}
	return x1
}
func getDeepCloneFuncMap(t reflect.Type) cloneFuncType {
	tk := t.Key()
	tv := t.Elem()
	pfk := getCloneFuncByPool(tk, true)
	pfv := getCloneFuncByPool(tv, true)
	return func(x0 reflect.Value) reflect.Value {
		x1 := reflect.MakeMap(t)
		for it := x0.MapRange(); it.Next(); {
			key := it.Key()
			val := it.Value()

			key = (*pfk)(key)
			val = (*pfv)(val)
			x1.SetMapIndex(key, val)
		}
		return x1
	}
}
