/*
 * Copyright (C) distroy
 */

package ldref

import (
	"reflect"
	"sync"
)

type commFuncValue[F any] struct {
	Done func()
	Func F
}

type commFuncPool[K any, F any] struct {
	data sync.Map
}

func (m *commFuncPool[K, F]) Get(key any, fnGet func() F) (*F, func()) {
	{
		i, _ := m.data.Load(key)
		v, _ := i.(*commFuncValue[F])
		if v != nil {
			return &v.Func, v.Done
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	v := &commFuncValue[F]{
		Done: wg.Wait,
	}
	i, loaded := m.data.LoadOrStore(key, v)
	if loaded {
		v := i.(*commFuncValue[F])
		return &v.Func, v.Done
	}

	v.Func = fnGet()
	wg.Done()
	return &v.Func, v.Done
}

type cloneFuncType func(x0 reflect.Value) reflect.Value
type cloneFuncKey struct {
	Type reflect.Type
	Deep bool
}

var cloneFuncPool = &commFuncPool[cloneFuncKey, cloneFuncType]{}

func getCloneFuncByPool(t reflect.Type, deep bool) (*cloneFuncType, func()) {
	key := cloneFuncKey{
		Type: t,
		Deep: deep,
	}

	pool := cloneFuncPool
	return pool.Get(key, func() cloneFuncType {
		if deep {
			return getDeepCloneFunc(t)
		}
		return getCloneFunc(t)
	})
}
