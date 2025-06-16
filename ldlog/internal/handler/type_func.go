/*
 * Copyright (C) distroy
 */

package handler

import (
	"fmt"
	"reflect"
	"unsafe"
)

func toType[T, S any](v S) T { return *(*T)(unsafe.Pointer(&v)) }

func asType[T any](v any, def ...T) T {
	vv, ok := v.(T)
	if ok {
		return vv
	}
	if len(def) > 0 {
		return def[0]
	}
	return vv
}

func checkTypeEqual(this, that reflect.Type) {
	if !isTypeEqual(this, that) {
		panic(fmt.Errorf("%s not not compatible with %s", this.String(), that.String()))
	}
}

func isTypeEqual(t0, t1 reflect.Type) bool {
	if t0 == t1 {
		return true
	}
	if t0.Kind() != t1.Kind() {
		return false
	}

	switch t0.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Bool, reflect.String, reflect.Float32, reflect.Float64:
		return true
	case reflect.Complex64, reflect.Complex128:
		return true
	case reflect.UnsafePointer:
		return true

	case reflect.Array:
		return t0.Len() == t1.Len() && isTypeEqual(t0.Elem(), t1.Elem())

	case reflect.Slice, reflect.Chan, reflect.Pointer:
		return isTypeEqual(t0.Elem(), t1.Elem())

	case reflect.Map:
		return isTypeEqual(t0.Key(), t1.Key()) && isTypeEqual(t0.Elem(), t1.Elem())

	case reflect.Func:
		return isTypeEqualForFunc(t0, t1)

	case reflect.Struct:
		return isTypeEqualForStruct(t0, t1)

	case reflect.Interface:
		return t0.Implements(t1) && t1.Implements(t0)
	}
	return false
}

func isTypeEqualForFunc(t0, t1 reflect.Type) bool {
	nout0 := t0.NumOut()
	nout1 := t1.NumOut()
	nin0 := t0.NumIn()
	nin1 := t1.NumIn()
	if nout0 != nout1 || nin0 != nin1 {
		return false
	}
	for i := range nout0 {
		in0 := t0.In(i)
		in1 := t1.In(i)
		if !isTypeEqual(in0, in1) {
			return false
		}
	}
	for i := range nin0 {
		out0 := t0.Out(i)
		out1 := t1.Out(i)
		if !isTypeEqual(out0, out1) {
			return false
		}
	}
	return true
}

func isTypeEqualForStruct(t0, t1 reflect.Type) bool {
	n0 := t0.NumField()
	n1 := t1.NumField()
	if n0 != n1 {
		return false
	}

	for i := range n0 {
		f0 := t0.Field(i)
		f1 := t1.Field(i)
		if !isTypeEqual(f0.Type, f1.Type) {
			return false
		}
	}
	return true
}
