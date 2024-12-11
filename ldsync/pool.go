/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sync"
)

type Pool[T any] struct {
	d sync.Pool

	// New optionally specifies a function to generate
	// a value when Get would otherwise return nil.
	// It may not be changed concurrently with calls to Get.
	New func() T
}

func (p *Pool[T]) Get() T {
	i := p.d.Get()
	if i != nil {
		return i.(T)
	}
	if p.New != nil {
		return p.New()
	}
	var v T
	return v
}

func (p *Pool[T]) Put(v T) {
	p.d.Put(v)
}

var pools = &sync.Map{}

func GetPool[T any](fnNew func() T) *Pool[T] {
	var key T
	// typ := reflect.TypeOf(d)
	if i, _ := pools.Load(key); i != nil {
		return i.(*Pool[T])
	}

	p := &Pool[T]{New: fnNew}
	i, _ := pools.LoadOrStore(key, p)
	return i.(*Pool[T])
}
