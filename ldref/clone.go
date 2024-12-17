/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"sync"
)

func Clone[T any](d T) T { return cloneV1(d) }

func cloneV1[T any](d T) T {
	var i interface{} = d
	if x, ok := i.(reflect.Value); ok {
		var r interface{} = cloneRef(x)
		return r.(T)
	}

	x := cloneRef(reflect.ValueOf(i))
	if x.Kind() == reflect.Invalid {
		var zero T
		return zero
	}

	return x.Interface().(T)
}
func cloneV2[T any](d T) T {
	var i interface{} = d
	if x, ok := i.(reflect.Value); ok {
		var r interface{} = cloneRefV2(x)
		return r.(T)
	}

	x := cloneRefV2(reflect.ValueOf(i))
	if x.Kind() == reflect.Invalid {
		var zero T
		return zero
	}

	return x.Interface().(T)
}
func cloneRef(x0 reflect.Value) reflect.Value {
	if x0.Kind() == reflect.Interface {
		if x0.IsNil() {
			return x0
		}

		x0 = x0.Elem()
	}

	switch x0.Kind() {
	case reflect.Invalid:
		return reflect.Value{}

	case reflect.Struct:
		if isSyncType(x0.Type()) {
			return reflect.Zero(x0.Type())
		}
		return x0

	case reflect.Ptr:
		return clonePtr(x0)

	case reflect.Array:
		return cloneArray(x0)

	case reflect.Slice:
		return cloneSlice(x0)

	case reflect.Map:
		return cloneMap(x0)
	}

	return x0
}
func cloneRefV2(x0 reflect.Value) reflect.Value {
	t := refTypeOfValue(x0)
	pf := getCloneFuncByPool(t, false)
	return (*pf)(x0)
}
func getCloneFunc(t reflect.Type) cloneFuncType {
	switch refKindOfType(t) {
	case reflect.Invalid:
		return func(x0 reflect.Value) reflect.Value { return reflect.Value{} }

	case reflect.Interface:
		return func(x0 reflect.Value) reflect.Value {
			if x0.IsNil() {
				return x0
			}
			return cloneRefV2(x0.Elem())
		}

	case reflect.Struct:
		zero := reflect.Zero(t)
		if isSyncType(t) {
			return func(x0 reflect.Value) reflect.Value { return zero }
		}
		return func(x0 reflect.Value) reflect.Value { return x0 }

	case reflect.Ptr:
		return getCloneFuncPtr(t)

	case reflect.Array:
		return cloneArray

	case reflect.Slice:
		return cloneSlice

	case reflect.Map:
		return cloneMap
	}

	return func(v reflect.Value) reflect.Value { return v }
}

func clonePtr(x0 reflect.Value) reflect.Value {
	if x0.IsNil() {
		return x0
	}

	x1 := reflect.New(x0.Type().Elem())
	x1.Elem().Set(x0.Elem())
	return x1
}
func getCloneFuncPtr(t reflect.Type) cloneFuncType {
	tt := t.Elem()
	zero := reflect.Zero(t)
	return func(x0 reflect.Value) reflect.Value {
		if x0.IsNil() {
			return zero
		}

		x1 := reflect.New(tt)
		x1.Elem().Set(x0.Elem())
		return x1
	}
}

func isSyncType(t reflect.Type) bool {
	if t.Kind() != reflect.Struct {
		return false
	}
	v := reflect.Zero(t)
	switch v.Interface().(type) {
	case sync.Mutex, sync.RWMutex, sync.Cond, sync.WaitGroup:
		return true
	}
	return false
}

func cloneArray(x0 reflect.Value) reflect.Value {
	return x0
}

func cloneSlice(x0 reflect.Value) reflect.Value {
	l0 := x0.Len()
	x1 := reflect.MakeSlice(reflect.SliceOf(x0.Type().Elem()), 0, l0)
	for i := 0; i < l0; i++ {
		v0 := x0.Index(i)
		x1 = reflect.Append(x1, v0)
		// x1.Index(i).Set(v0)
	}
	return x1
}

func cloneMap(x0 reflect.Value) reflect.Value {
	x1 := reflect.MakeMap(x0.Type())
	for it := x0.MapRange(); it.Next(); {
		key := it.Key()
		val := it.Value()
		x1.SetMapIndex(key, val)
	}
	return x1
}
