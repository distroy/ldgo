/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"unsafe"

	"github.com/distroy/ldgo/v2/lderr"
)

type MergeConfig struct {
	Clone      bool // is clone if target is nil
	MergeArray bool // is merge array. `false` mean only assign target at whole array is zero value
	MergeSlice bool // is merge slice. `false` mean only assign target at slice is nil
}

// Merge will merge the data from source to target
//   - Merge(*int, int)
//   - Merge(*int, *int)
//   - Merge(*structA, structA)
//   - Merge(*structA, *structA)
//   - Merge(*map, map)
//   - Merge(*map, *map)
func Merge(target, source interface{}, cfg ...*MergeConfig) error {
	c := &mergeContext{
		MergeConfig: &MergeConfig{},
	}
	if len(cfg) > 0 && cfg[0] != nil {
		c.MergeConfig = cfg[0]
	}

	return mergeWithContext(c, target, source)
}

type mergeContext struct {
	*MergeConfig
}

func mergeWithContext(c *mergeContext, target, source interface{}) error {
	tVal := valueOf(target)
	sVal := valueOf(source)

	tTyp := tVal.Type()
	sTyp := sVal.Type()

	if tTyp.Kind() != reflect.Ptr {
		return lderr.ErrReflectTargetNotPtr
	}

	if tVal.IsNil() {
		return lderr.ErrReflectTargetNilPtr
	}

	tElemType := tTyp.Elem()
	switch {
	default:
		return lderr.ErrReflectTypeNotEqual

	case tTyp == sTyp ||
		(tElemType.Kind() == reflect.Interface && sTyp.Kind() == reflect.Ptr && sTyp.Elem().Implements(tElemType)):

		if sVal.IsNil() {
			// do not need to merge
			return nil
		}

		tVal = tVal.Elem()
		sVal = sVal.Elem()

	case tElemType == sTyp ||
		(tElemType.Kind() == reflect.Interface && sTyp.Implements(tElemType)):
		tVal = tVal.Elem()
	}

	mergeReflect(c, tVal, sVal)
	return nil
}

func cloneForMerge(c *mergeContext, x reflect.Value) reflect.Value {
	v := x
	if c.Clone {
		v = deepClone(v)
	}
	return v
}

func mergeReflect(c *mergeContext, target, source reflect.Value) {
	switch target.Kind() {
	default:
		mergeReflectNormal(c, target, source)

	case reflect.Invalid:
		break

	case reflect.Interface:
		mergeReflectIface(c, target, source)

	case reflect.Ptr:
		mergeReflectPtr(c, target, source)

	case reflect.Func, reflect.Chan:
		mergeReflectFunc(c, target, source)

	case reflect.Map:
		mergeReflectMap(c, target, source)

	case reflect.Slice:
		mergeReflectSlice(c, target, source)

	case reflect.Array:
		mergeReflectArray(c, target, source)

	case reflect.Struct:
		mergeReflectStruct(c, target, source)
	}
}

func mergeReflectIface(c *mergeContext, target, source reflect.Value) {
	if target.IsNil() {
		source = cloneForMerge(c, source)
		target.Set(source)
		return
	}

	if source.IsNil() {
		return
	}

	// target = reflect.ValueOf(target.Interface())
	// source = reflect.ValueOf(source.Interface())
	// if target.Type() != source.Type() {
	// 	return
	// }

	tDataTyp := reflect.TypeOf(target.Interface())
	source = reflect.ValueOf(source.Interface())
	if tDataTyp != source.Type() {
		return
	}

	tDataVal := reflect.New(tDataTyp).Elem()
	tDataVal.Set(reflect.ValueOf(target.Interface()))
	mergeReflect(c, tDataVal, source)

	// log.Printf(" === %s: %#v", target.Type().String(), target.Interface())
	// log.Printf(" === %s: %#v", tDataVal.Type().String(), tDataVal.Interface())
	target.Set(tDataVal)
}

func mergeReflectPtr(c *mergeContext, target, source reflect.Value) {
	if target.IsNil() {
		source = cloneForMerge(c, source)
		target.Set(source)
		return
	}

	if source.IsNil() {
		return
	}

	mergeReflect(c, target.Elem(), source.Elem())
}

func mergeReflectFunc(c *mergeContext, target, source reflect.Value) {
	if target.IsNil() {
		source = cloneForMerge(c, source)
		target.Set(source)
	}
}

func mergeReflectMap(c *mergeContext, target, source reflect.Value) {
	if target.IsNil() {
		source = cloneForMerge(c, source)
		target.Set(source)
		return
	}

	for it := source.MapRange(); it.Next(); {
		key := it.Key()
		sVal := it.Value()
		if !sVal.IsValid() {
			continue
		}

		tVal := target.MapIndex(key)
		if !tVal.IsValid() {
			sVal = cloneForMerge(c, sVal)
			target.SetMapIndex(key, sVal)
			continue
		}

		tmp := reflect.New(tVal.Type()).Elem()
		tmp.Set(tVal)
		mergeReflect(c, tmp, sVal)
		target.SetMapIndex(key, tmp)
	}
}

func mergeReflectSlice(c *mergeContext, target, source reflect.Value) {
	if source.IsNil() {
		return
	}

	if target.IsNil() {
		source = cloneForMerge(c, source)
		target.Set(source)
		return
	}

	if !c.MergeSlice {
		return
	}

	tLen := target.Len()
	sLen := source.Len()

	resizeSliceReflect(target, sLen)

	for i := 0; i < sLen; i++ {
		tVal := target.Index(i)
		sVal := source.Index(i)

		if i < tLen {
			mergeReflect(c, tVal, sVal)
			continue
		}

		sVal = cloneForMerge(c, sVal)
		tVal.Set(sVal)
	}
}

func mergeReflectArray(c *mergeContext, target, source reflect.Value) {
	if !c.MergeArray {
		if IsValZero(target) {
			source = cloneForMerge(c, source)
			target.Set(source)
		}
		return
	}

	l := source.Len()
	for i := 0; i < l; i++ {
		tVal := target.Index(i)
		sVal := source.Index(i)

		mergeReflect(c, tVal, sVal)
	}
}

func mergeReflectStruct(c *mergeContext, target, source reflect.Value) {
	n := target.NumField()
	for i := 0; i < n; i++ {

		tField := target.Field(i)
		sField := source.Field(i)

		tFieldAddr := unsafe.Pointer(tField.UnsafeAddr())
		tField = reflect.NewAt(tField.Type(), tFieldAddr).Elem()

		mergeReflect(c, tField, sField)
	}
}

func mergeReflectNormal(_ *mergeContext, target, source reflect.Value) {
	if IsValZero(target) {
		target.Set(source)
	}
}
