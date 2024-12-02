/*
 * Copyright (C) distroy
 */

package cmp

import (
	"math"
	"reflect"
	"sort"
	"sync"
)

var (
	comparerTypes = &sync.Map{}
)

type typeInfo struct {
	IsComparer bool
	Compare    reflect.Value
}

type kind int

const (
	kindNil kind = iota
	kindBool
	kindNumber
	kindComplex
	kindString
	kindOrthers
)

func reflectValueOf(v interface{}) reflect.Value {
	if vv, ok := v.(reflect.Value); ok {
		return vv
	}

	vv := reflect.ValueOf(v)
	return vv
}

func CompareReflect(a, b reflect.Value) int {
	if r := compareReflectType(a, b); r != 0 {
		return r
	}

	if r, ok := compareReflectComparer(a, b); ok {
		return r
	}

	switch a.Kind() {
	default:
		// Certain types cannot appear as keys (maps, funcs, slices), but be explicit.
		panic("bad type in compare: " + a.Type().String())

	case reflect.Invalid:
		return 0

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		aa := a.Int()
		return compareReflectNumberLeftInt(aa, b)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		aa := a.Uint()
		return compareReflectNumberLeftUint(aa, b)

	case reflect.Float32, reflect.Float64:
		aa := a.Float()
		return compareReflectNumberLeftFloat(aa, b)

	case reflect.String:
		aa := a.String()
		bb := b.String()
		return CompareString(aa, bb)

	case reflect.Complex64, reflect.Complex128:
		aa := a.Complex()
		bb := b.Complex()
		return CompareComplex(aa, bb)

	case reflect.Bool:
		aa := a.Bool()
		bb := b.Bool()
		return CompareBool(aa, bb)

	case reflect.Ptr:
		return compareReflectPointer(a, b)

	case reflect.UnsafePointer:
		aa := a.Pointer()
		bb := b.Pointer()
		return CompareOrderable(aa, bb)

	case reflect.Chan, reflect.Func:
		aa := a.Pointer()
		bb := b.Pointer()
		return CompareOrderable(aa, bb)

	case reflect.Map:
		return compareReflectMap(a, b)

	case reflect.Struct:
		return compareReflectStruct(a, b)

	case reflect.Slice:
		return compareReflectSlice(a, b)

	case reflect.Array:
		return compareReflectArray(a, b)

	case reflect.Interface:
		return compareReflectIface(a, b)
	}
}

func convertKind(k reflect.Kind) kind {
	switch k {
	case reflect.Invalid:
		return kindNil

	case reflect.Bool:
		return kindBool

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return kindNumber

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return kindNumber

	case reflect.Float32, reflect.Float64:
		return kindNumber

	case reflect.Complex64, reflect.Complex128:
		return kindComplex

	case reflect.String:
		return kindString

	default:
		return kindOrthers
	}
}

func compareReflectType(a, b reflect.Value) int {
	aKind := convertKind(a.Kind())
	bKind := convertKind(b.Kind())
	if aKind != kindOrthers && aKind == bKind {
		return 0
	}

	if aKind != kindOrthers || bKind != kindOrthers {
		return CompareOrderable(aKind, bKind)
	}

	aType := a.Type()
	bType := b.Type()

	if aType == bType {
		return 0
	}

	aName := aType.String()
	bName := bType.String()
	if r := CompareOrderable(len(aName), len(bName)); r != 0 {
		return r
	}

	if aName < bName {
		return -1
	}
	return 1
}

func compareNilReflect(a, b reflect.Value) (int, bool) {
	aNil := a.IsNil()
	bNil := b.IsNil()
	if aNil && bNil {
		return 0, true
	}
	if aNil {
		return -1, true
	}
	if bNil {
		return 1, true
	}
	return 0, false
}

func compareReflectComparer(a, b reflect.Value) (int, bool) {
	aType := a.Type()
	bType := b.Type()
	if aType != bType {
		return 0, false
	}
	typ := getComparerTypeInfo(aType)
	if !typ.IsComparer {
		return 0, false
	}
	ins := [...]reflect.Value{a, b}
	r := typ.Compare.Call(ins[:])[0].Int()
	return int(r), true
}

func getComparerTypeInfo(typ reflect.Type) typeInfo {
	if i, ok := comparerTypes.Load(typ); ok {
		p, _ := i.(typeInfo)
		return p
	}

	res := typeInfo{}
	defer func() {
		comparerTypes.Store(typ, res)
	}()

	method, ok := typ.MethodByName("Compare")
	if !ok {
		return res
	}

	mType := method.Type
	if mType.NumIn() != 2 {
		return res
	}
	if in := mType.In(1); in != typ && (in.Kind() != reflect.Interface || !typ.Implements(in)) {
		return res
	}
	if mType.NumOut() != 1 || mType.Out(0).Kind() != reflect.Int {
		return res
	}

	res.IsComparer = true
	res.Compare = method.Func
	return res
}

func compareReflectPointer(a, b reflect.Value) int {
	aa := a.Pointer()
	bb := b.Pointer()
	if r := CompareOrderable(aa, bb); r == 0 || aa == 0 || bb == 0 {
		return r
	}
	return CompareReflect(a.Elem(), b.Elem())
}

func compareReflectStruct(a, b reflect.Value) int {
	for i := 0; i < a.NumField(); i++ {
		aa := a.Field(i)
		bb := b.Field(i)
		// aa, bb = reflect.ValueOf(aa.Interface()), reflect.ValueOf(bb.Interface())
		if r := CompareReflect(aa, bb); r != 0 {
			return r
		}
	}
	return 0
}

func compareReflectArray(a, b reflect.Value) int {
	for i := 0; i < a.Len(); i++ {
		aa := a.Index(i)
		bb := b.Index(i)
		if r := CompareReflect(aa, bb); r != 0 {
			return r
		}
	}
	return 0
}

func compareReflectMap(a, b reflect.Value) int {
	// if r := CompareOrderable(a.Len(), b.Len()); r != 0 {
	// 	return r
	// }

	aKeys := a.MapKeys()
	bKeys := b.MapKeys()
	sort.Sort(sortedReflects(aKeys))
	sort.Sort(sortedReflects(bKeys))

	al := a.Len()
	bl := b.Len()
	l := al
	if l > bl {
		l = bl
	}

	for i := 0; i < l; i++ {
		aKey, bKey := aKeys[i], bKeys[i]
		if r := CompareReflect(aKey, bKey); r != 0 {
			return -r
		}

		aVal, bVal := a.MapIndex(aKey), b.MapIndex(bKey)
		if r := CompareReflect(aVal, bVal); r != 0 {
			return r
		}
	}

	// return 0
	return CompareOrderable(al, bl)
}

func compareReflectSlice(a, b reflect.Value) int {
	al := a.Len()
	bl := b.Len()
	if al == 0 || bl == 0 || a.Index(0).UnsafeAddr() == b.Index(0).UnsafeAddr() {
		return CompareOrderable(al, bl)
	}
	if a.Kind() == reflect.Uint8 {
		return CompareBytes(a.Bytes(), b.Bytes())
	}

	l := al
	if l > bl {
		l = bl
	}
	for i := 0; i < l; i++ {
		aa, bb := a.Index(i), b.Index(i)
		if r := CompareReflect(aa, bb); r != 0 {
			return r
		}
	}
	return CompareOrderable(al, bl)
}

func compareReflectIface(a, b reflect.Value) int {
	if r, ok := compareNilReflect(a, b); ok {
		return r
	}
	return CompareReflect(a.Elem(), b.Elem())
}

func compareReflectNumberLeftInt(aa int64, b reflect.Value) int {
	switch b.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		bb := b.Uint()
		return compareIntAndUint(aa, bb)

	case reflect.Float32, reflect.Float64:
		bb := b.Float()
		return compareIntAndFloat(aa, bb)
	}

	bb := b.Int()
	return CompareOrderable(aa, bb)
}

func compareReflectNumberLeftUint(aa uint64, b reflect.Value) int {
	switch b.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bb := b.Int()
		return -compareIntAndUint(bb, aa)

	case reflect.Float32, reflect.Float64:
		bb := b.Float()
		return compareUintAndFloat(aa, bb)
	}

	bb := b.Uint()
	return CompareOrderable(aa, bb)
}

func compareReflectNumberLeftFloat(aa float64, b reflect.Value) int {
	switch b.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bb := b.Int()
		return -compareIntAndFloat(bb, aa)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		bb := b.Uint()
		return -compareUintAndFloat(bb, aa)
	}

	bb := b.Float()
	return CompareOrderable(aa, bb)
}

func compareIntAndFloat(a int64, b float64) int {
	if b > math.MaxInt64 {
		return -1
	} else if b < math.MinInt64 {
		return 1
	}

	bb := int64(b)
	r := CompareOrderable(a, bb)
	if r != 0 {
		return r
	}

	d := b - float64(bb)
	if d == 0 {
		return r
	} else if d > 0 {
		return -1
	}
	return 1
}

func compareIntAndUint(a int64, b uint64) int {
	if a < 0 {
		return -1
	}
	return CompareOrderable(uint64(a), b)
}

func compareUintAndFloat(a uint64, b float64) int {
	if b > math.MaxUint64 {
		return -1
	} else if b < 0 {
		return 1
	}

	bb := uint64(b)
	r := CompareOrderable(a, bb)
	if r != 0 {
		return r
	}

	d := b - float64(bb)
	if d == 0 {
		return r
	} else if d > 0 {
		return -1
	}
	return 1
}
