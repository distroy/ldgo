/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"

	"github.com/distroy/ldgo/v2/ldsync"
)

type cloneFuncType func(x0 reflect.Value) reflect.Value
type cloneFuncKey struct {
	Type reflect.Type
	Deep bool
}

var cloneFuncPool = &ldsync.Map[cloneFuncKey, *cloneFuncType]{}

func getCloneFuncByPool(t reflect.Type, deep bool) *cloneFuncType {
	key := cloneFuncKey{
		Type: t,
		Deep: deep,
	}

	pool := cloneFuncPool
	if pfn, _ := pool.Load(key); pfn != nil {
		return pfn
	}

	var fn cloneFuncType
	pf := &fn
	pf, loaded := pool.LoadOrStore(key, pf)
	if loaded {
		return pf
	}

	fnGet := getCloneFunc
	if deep {
		fnGet = getDeepCloneFunc
	}

	*pf = fnGet(t)
	return pf
}
