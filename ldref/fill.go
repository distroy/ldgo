/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"unsafe"

	"github.com/distroy/ldgo/v3/lderr"
)

type FillConfig struct {
	Deep bool
}

func DeepFill(v interface{}) error {
	cfg := &FillConfig{
		Deep: true,
	}
	return Fill(v, cfg)
}

// Fill will fill the fields of struct, if field is nil pointer
func Fill(v interface{}, cfg ...*FillConfig) error {
	c := &fillContext{
		context:    getContext(),
		FillConfig: &FillConfig{},
	}
	putContext(c.context)

	if len(cfg) > 0 && cfg[0] != nil {
		c.FillConfig = cfg[0]
	}

	return fillWithContext(c, v)
}

type fillContext struct {
	*context
	*FillConfig
}

func fillWithContext(c *fillContext, v interface{}) error {
	vv := valueOf(v)
	if vv.Kind() == reflect.Ptr {
		if vv.IsNil() {
			return lderr.ErrReflectTargetNilPtr
		}

		vv = vv.Elem()
	}

	if !vv.CanAddr() {
		return lderr.ErrReflectTargetNilPtr
	}

	if vv.Kind() == reflect.Struct {
		fillReflectStruct(c, vv)
	} else {
		fillReflect(c, vv)
	}
	return c.Error()
}

func fillReflectStruct(c *fillContext, v reflect.Value) {
	t := v.Type()
	num := v.NumField()
	for i := 0; i < num; i++ {
		fVal := v.Field(i)

		fTyp := t.Field(i)
		name := fTyp.Name

		addr := unsafe.Pointer(fVal.UnsafeAddr())
		fVal = reflect.NewAt(fVal.Type(), addr).Elem()

		c.PushField(name)
		fillReflect(c, fVal)
		c.PopField()

		// log.Printf("%s(%s): %v", name, fVal.Type().String(), fVal.Interface())
	}
}

func fillReflectNotDeep(c *fillContext, v reflect.Value) {
	var fnNew func(typ reflect.Type) reflect.Value
	switch v.Kind() {
	// default:
	// 	return

	case reflect.Ptr:
		fnNew = func(typ reflect.Type) reflect.Value { return reflect.New(typ.Elem()) }
	case reflect.Map:
		fnNew = func(typ reflect.Type) reflect.Value { return reflect.MakeMap(typ) }
	case reflect.Slice:
		fnNew = func(typ reflect.Type) reflect.Value { return reflect.MakeSlice(typ, 0, 0) }
	case reflect.Chan:
		fnNew = func(typ reflect.Type) reflect.Value { return reflect.MakeChan(typ, 0) }
	}

	if fnNew != nil && v.IsNil() {
		v.Set(fnNew(v.Type()))
	}
}

func fillReflect(c *fillContext, v reflect.Value) {
	if !c.Deep {
		fillReflectNotDeep(c, v)
		return
	}

	vv := v
	for ; vv.Kind() == reflect.Ptr; vv = vv.Elem() {
		if !vv.IsNil() {
			continue
		}

		vv.Set(reflect.New(vv.Type().Elem()))
	}

	fillReflectNotDeep(c, vv)

	if vv.Kind() == reflect.Struct {
		fillReflectStruct(c, vv)
	}
}
