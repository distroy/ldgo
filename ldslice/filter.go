/*
 * Copyright (C) distroy
 */

package ldslice

import (
	"fmt"
	"reflect"
)

// Filter filters the slice with filter
// slice type must be slice or array
// filter type must be:
//
//	func (v TypeOfSliceElement) bool
func Filter(slice interface{}, filter interface{}) int {
	sliceVal := reflect.ValueOf(slice)
	sliceTyp := sliceVal.Type()

	if sliceTyp.Kind() != reflect.Slice && sliceTyp.Kind() != reflect.Array {
		panic(fmt.Sprintf("the slice type must be slice or array. type:%s", sliceTyp))
	}
	elemTyp := sliceVal.Type().Elem()

	filterVal := reflect.ValueOf(filter)
	filterTyp := filterVal.Type()
	if filterTyp.Kind() != reflect.Func {
		panic(fmt.Sprintf("the filter must be func. type:%s", filterTyp))
	}

	if filterTyp.NumIn() != 1 {
		panic(fmt.Sprintf("the filter must have 1 input parameter. type:%s", filterTyp))
	}
	if typ := filterTyp.In(0); !(typ == elemTyp || (typ.Kind() == reflect.Interface && elemTyp.Implements(typ))) {
		panic(fmt.Sprintf("the parameter of filter must be or interface for %s", elemTyp))
	}

	if filterTyp.NumOut() != 1 {
		panic(fmt.Sprintf("the filter must have 1 return value. type:%s", filterTyp))
	}
	if typ := filterTyp.Out(0); typ.Kind() != reflect.Bool {
		panic(fmt.Sprintf("the return value of filter must be bool. type:%s", typ))
	}

	return filterSlice(sliceVal, filterVal)
}

func filterSlice(slice, filter reflect.Value) int {
	i := 0
	j := slice.Len()

	for i < j {
		var vi, vj reflect.Value

		for ; i < j; i++ {
			vi = slice.Index(i)
			if !filter.Call([]reflect.Value{vi})[0].Bool() {
				break
			}
		}

		for ; i < j; j-- {
			vj = slice.Index(j - 1)
			if filter.Call([]reflect.Value{vj})[0].Bool() {
				break
			}
		}

		if i < j-1 {
			tmp := vi.Interface()
			vi.Set(vj)
			vj.Set(reflect.ValueOf(tmp))
			i++
		}
	}

	return i
}
