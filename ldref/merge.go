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
// where target or field of target is zero
//   - Merge(*int, int)
//   - Merge(*int, *int)
//   - Merge(*structA, structA)
//   - Merge(*structA, *structA)
//   - Merge(*map, map)
//   - Merge(*map, *map)
func Merge(target, source any, cfg ...*MergeConfig) error {
	return mergeV2(target, source, cfg...)
}
func mergeV1(target, source any, cfg ...*MergeConfig) error {
	c := &mergeContext{
		MergeConfig: &MergeConfig{},
	}
	if len(cfg) > 0 && cfg[0] != nil {
		c.MergeConfig = cfg[0]
	}

	return mergeWithContext(c, target, source, mergeReflect)
}
func mergeV2(target, source any, cfg ...*MergeConfig) error {
	c := &mergeContext{
		MergeConfig: &MergeConfig{},
	}
	if len(cfg) > 0 && cfg[0] != nil {
		c.MergeConfig = cfg[0]
	}

	// log.Printf(" === clone:%v", c.Clone)
	return mergeWithContext(c, target, source, mergeReflectV2)
}

type mergeContext struct {
	*MergeConfig
}

type mergeFuncType = func(c *mergeContext, target, source reflect.Value)

var mergePool = &commFuncPool[reflect.Type, mergeFuncType]{}

func mergeWithContext(c *mergeContext, target, source any, fnMerge mergeFuncType) error {
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
		// log.Printf(" === clone:%v", c.Clone)

		if sVal.IsNil() {
			// do not need to merge
			return nil
		}

		tVal = tVal.Elem()
		sVal = sVal.Elem()

	case tElemType == sTyp ||
		(tElemType.Kind() == reflect.Interface && sTyp.Implements(tElemType)):
		// log.Printf(" === clone:%v", c.Clone)
		tVal = tVal.Elem()
	}

	// log.Printf(" === clone:%v", c.Clone)
	// mergeReflect(c, tVal, sVal)
	fnMerge(c, tVal, sVal)
	return nil
}

func cloneForMerge(c *mergeContext, x reflect.Value) reflect.Value {
	// log.Printf(" === clone:%v", c.Clone)
	v := x
	if c.Clone {
		// v = deepCloneRef(v)
		pfClone, done := getCloneFuncByPool(v.Type(), true)
		done()
		v = (*pfClone)(v)
	}
	return v
}
func getMergeFuncByClone(_ *mergeContext, typ reflect.Type) mergeFuncType {
	// func is saved by pool without clone flag
	pfClone, done := getCloneFuncByPool(typ, true)
	return func(c *mergeContext, target, source reflect.Value) {
		if c.Clone {
			// log.Printf(" === clone:%v", c.Clone)
			done()
			val := (*pfClone)(source)
			target.Set(val)
		} else {
			target.Set(source)
		}
	}
}

var (
	mergeKindMapV1 []mergeFuncType
	mergeKindMapV2 []func(c *mergeContext, typ reflect.Type) mergeFuncType
)

func init() {
	type mergeKindData struct {
		Kind    reflect.Kind
		MergeV1 mergeFuncType
		MergeV2 func(c *mergeContext, typ reflect.Type) mergeFuncType
	}
	s := []mergeKindData{
		{reflect.Invalid, mergeReflectInvalid, getMergeFuncInvalid},
		{reflect.Interface, mergeReflectIface, getMergeFuncIface},
		{reflect.Ptr, mergeReflectPtr, getMergeFuncPtr},
		{reflect.Func, mergeReflectFunc, getMergeFuncChan},
		{reflect.Chan, mergeReflectFunc, getMergeFuncChan},
		{reflect.Map, mergeReflectMap, getMergeFuncMap},
		{reflect.Slice, mergeReflectSlice, getMergeFuncSlice},
		{reflect.Array, mergeReflectArray, getMergeFuncArray},
		{reflect.Struct, mergeReflectStruct, getMergeFuncStruct},
	}
	l := reflect.Invalid
	for _, v := range s {
		if l < v.Kind {
			l = v.Kind
		}
	}
	mergeKindMapV1 = make([]mergeFuncType, l+1)
	mergeKindMapV2 = make([]func(c *mergeContext, typ reflect.Type) mergeFuncType, l+1)
	for _, v := range s {
		mergeKindMapV1[v.Kind] = v.MergeV1
		mergeKindMapV2[v.Kind] = v.MergeV2
	}
}

func isNormailTypeForMerge(kind reflect.Kind) bool {
	m := mergeKindMapV2
	return int(kind) >= len(m) || m[kind] == nil
}

func mergeReflect(c *mergeContext, target, source reflect.Value) {
	// log.Printf(" === type:%s", target.Type())
	kind := target.Kind()
	fn := mergeReflectNormal
	if m := mergeKindMapV1; int(kind) <= len(m) {
		tmp := m[kind]
		if tmp != nil {
			fn = tmp
		}
	}
	fn(c, target, source)
}
func mergeReflectV2(c *mergeContext, target, source reflect.Value) {
	// log.Printf(" === type:%s", target.Type())
	pf, done := getMergeFuncByPool(c, target.Type())
	done()
	(*pf)(c, target, source)
}
func getMergeFuncByPool(c *mergeContext, typ reflect.Type) (*mergeFuncType, func()) {
	pool := mergePool
	return pool.Get(typ, func() mergeFuncType {
		return getMergeFunc(c, typ)
	})
}
func getMergeFunc(c *mergeContext, typ reflect.Type) mergeFuncType {
	// log.Printf(" === type:%s", typ)
	kind := typ.Kind()
	fn := getMergeFuncNormal
	if m := mergeKindMapV2; int(kind) <= len(m) {
		tmp := m[kind]
		if tmp != nil {
			fn = tmp
		}
	}
	return fn(c, typ)
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
func getMergeFuncIface(c *mergeContext, typ reflect.Type) mergeFuncType {
	return func(c *mergeContext, target, source reflect.Value) {
		if target.IsNil() {
			source = cloneForMerge(c, source)
			target.Set(source)
			return

		} else if source.IsNil() {
			return
		}

		tDataVal := target.Elem()
		tDataTyp := target.Elem().Type()
		if !tDataVal.CanAddr() {
			val := reflect.New(tDataTyp).Elem()
			val.Set(tDataVal)
			tDataVal = val
		}
		source = source.Elem()
		if tDataTyp != source.Type() {
			return
		}

		mergeReflectV2(c, tDataVal, source)
		target.Set(tDataVal)
	}
}

func mergeReflectPtr(c *mergeContext, target, source reflect.Value) {
	// log.Printf(" === clone:%v", c.Clone)
	if target.IsNil() {
		source = cloneForMerge(c, source)
		target.Set(source)
		return
	}

	if source.IsNil() {
		return
	}

	target = target.Elem()
	source = source.Elem()
	kind := target.Kind()
	if m := mergeKindMapV1; int(kind) <= len(m) {
		fn := m[kind]
		if fn != nil {
			fn(c, target, source)
		}
	}
	// target.Set(source)
}
func getMergeFuncPtr(c *mergeContext, typ reflect.Type) mergeFuncType {
	// log.Printf(" === clone:%v", c.Clone)
	tElem := typ.Elem()
	pfClone := getMergeFuncByClone(c, typ)
	if isNormailTypeForMerge(tElem.Kind()) {
		return func(c *mergeContext, target, source reflect.Value) {
			// log.Printf(" === type:%s", target.Type())
			if source.IsNil() {
				// log.Printf(" === clone:%v", c.Clone)
				return
			}

			if target.IsNil() {
				// log.Printf(" === clone:%v", c.Clone)
				pfClone(c, target, source)
				return
			}
			// log.Printf(" === clone:%v", c.Clone)

			// target = target.Elem()
			// source = source.Elem()
			// // log.Printf(" === clone:%v", c.Clone)
			// target.Set(source)
		}
	}
	pfElem, dElem := getMergeFuncByPool(c, tElem)
	return func(c *mergeContext, target, source reflect.Value) {
		// log.Printf(" === type:%s", target.Type())
		if source.IsNil() {
			// log.Printf(" === clone:%v", c.Clone)
			return
		}

		if target.IsNil() {
			// log.Printf(" === clone:%v", c.Clone)
			pfClone(c, target, source)
			return
		}
		// log.Printf(" === clone:%v", c.Clone)

		target = target.Elem()
		source = source.Elem()

		dElem()
		(*pfElem)(c, target, source)
	}
}

func mergeReflectFunc(c *mergeContext, target, source reflect.Value) {
	if target.IsNil() {
		source = cloneForMerge(c, source)
		target.Set(source)
	}
}
func getMergeFuncChan(c *mergeContext, typ reflect.Type) mergeFuncType {
	pf := getMergeFuncByClone(c, typ)
	return func(c *mergeContext, target, source reflect.Value) {
		if target.IsNil() {
			pf(c, target, source)
		}
	}
}

func mergeReflectMap(c *mergeContext, target, source reflect.Value) {
	if source.IsNil() {
		return
	}

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
func getMergeFuncMap(c *mergeContext, typ reflect.Type) mergeFuncType {
	pfClone := getMergeFuncByClone(c, typ)
	pfElem, dElem := getMergeFuncByPool(c, typ.Elem())
	return func(c *mergeContext, target, source reflect.Value) {
		if source.IsNil() {
			return
		}

		if target.IsNil() {
			pfClone(c, target, source)
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
			dElem()
			(*pfElem)(c, tmp, sVal)
			target.SetMapIndex(key, tmp)
		}
	}
}

func mergeReflectSlice(c *mergeContext, target, source reflect.Value) {
	if source.IsNil() || source.Len() == 0 {
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
func getMergeFuncSlice(c *mergeContext, typ reflect.Type) mergeFuncType {
	pfClone := getMergeFuncByClone(c, typ)
	pfElem, dElem := getMergeFuncByPool(c, typ.Elem())
	pfCloneElem := getMergeFuncByClone(c, typ.Elem())
	return func(c *mergeContext, target, source reflect.Value) {
		switch {
		case source.IsNil() || source.Len() == 0:
			return

		case target.IsNil():
			pfClone(c, target, source)
			return

		case !c.MergeSlice:
			return
		}

		tLen := target.Len()
		sLen := source.Len()

		resizeSliceReflect(target, sLen)

		for i := 0; i < sLen; i++ {
			tVal := target.Index(i)
			sVal := source.Index(i)

			if i < tLen {
				dElem()
				(*pfElem)(c, tVal, sVal)
				continue
			}

			// sVal = cloneForMerge(c, sVal)
			// tVal.Set(sVal)
			pfCloneElem(c, tVal, sVal)
		}
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
func getMergeFuncArray(c *mergeContext, typ reflect.Type) mergeFuncType {
	pfClone := getMergeFuncByClone(c, typ)
	pfElem, dElem := getMergeFuncByPool(c, typ.Elem())
	// pfCloneElem := getMergeFuncByClone(c, typ.Elem())
	return func(c *mergeContext, target, source reflect.Value) {
		if !c.MergeArray {
			if IsValZero(target) {
				// source = cloneForMerge(c, source)
				// target.Set(source)
				pfClone(c, target, source)
			}
			return
		}

		l := source.Len()
		for i := 0; i < l; i++ {
			tVal := target.Index(i)
			sVal := source.Index(i)

			// mergeReflect(c, tVal, sVal)
			dElem()
			(*pfElem)(c, tVal, sVal)
		}
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
func getMergeFuncStruct(c *mergeContext, typ reflect.Type) mergeFuncType {
	// pfClone := getMergeFuncByClone(c, typ)
	fFields := make([]mergeFuncType, 0, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		idx := i
		field := typ.Field(idx)
		pfField, dField := getMergeFuncByPool(c, field.Type)

		fFields = append(fFields, func(c *mergeContext, target, source reflect.Value) {
			tField := refStructField(target, idx, &field)
			sField := source.Field(idx)

			dField()
			(*pfField)(c, tField, sField)
		})
	}
	// pfCloneElem := getMergeFuncByClone(c, typ.Elem())
	return func(c *mergeContext, target, source reflect.Value) {
		for _, f := range fFields {
			f(c, target, source)
		}
	}
}

func mergeReflectInvalid(_ *mergeContext, _, _ reflect.Value) {}
func getMergeFuncInvalid(_ *mergeContext, _ reflect.Type) mergeFuncType {
	return mergeReflectInvalid
}

func mergeReflectNormal(_ *mergeContext, target, source reflect.Value) {
	if IsValZero(target) {
		target.Set(source)
	}
}
func getMergeFuncNormal(_ *mergeContext, _ reflect.Type) mergeFuncType {
	return mergeReflectNormal
}
