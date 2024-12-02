/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sync"

	"github.com/distroy/ldgo/v2/ldatomic"
)

func NewMap[K comparable, V any](m map[K]V) *Map[K, V] {
	p := &Map[K, V]{}
	for k, v := range m {
		p.Set(k, v)
	}
	return p
}

type Map[K comparable, V any] struct {
	data sync.Map
	size ldatomic.Int
}

func (p *Map[K, V]) Size() int {
	return p.size.Load()
}

func (p *Map[K, V]) Get(key K) V {
	i, _ := p.data.Load(key)
	v, _ := i.(V)
	return v
}

func (p *Map[K, V]) Set(key K, new V) {
	p.swap(key, new)
}

func (p *Map[K, V]) Load(key K) (V, bool) {
	i, ok := p.data.Load(key)
	v, _ := i.(V)
	return v, ok
}

func (p *Map[K, V]) Swap(key K, val V) (previous V, loaded bool) {
	i, loaded := p.swap(key, val)
	old, _ := i.(V)
	return old, loaded
}

func (p *Map[K, V]) swap(key K, val V) (previous any, loaded bool) {
	i, loaded := p.data.Swap(key, val)
	if !loaded {
		p.size.Add(1)
	}
	return i, loaded
}

func (p *Map[K, V]) Has(key K) bool {
	_, ok := p.data.Load(key)
	return ok
}

func (p *Map[K, V]) Del(key K) {
	p.loadAndDelete(key)
}

func (p *Map[K, V]) loadAndDelete(key K) (val any, loaded bool) {
	i, loaded := p.data.LoadAndDelete(key)
	if loaded {
		p.size.Add(-1)
	}
	return i, loaded
}

func (p *Map[K, V]) loadOrStore(key K, val V) (actual any, loaded bool) {
	i, loaded := p.data.LoadOrStore(key, val)
	if !loaded {
		p.size.Add(1)
	}
	return i, loaded
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (p *Map[K, V]) Range(f func(key K, val V) bool) {
	p.data.Range(func(key, value any) bool {
		k, _ := key.(K)
		v, _ := value.(V)
		return f(k, v)
	})
}

func (p *Map[K, V]) Map() map[K]V {
	m := make(map[K]V, p.Size())
	p.Range(func(key K, val V) bool {
		m[key] = val
		return true
	})
	return m
}

func (p *Map[K, V]) Keys() []K {
	keys := make([]K, 0, p.Size())
	p.Range(func(key K, val V) bool {
		keys = append(keys, key)
		return true
	})
	return keys
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (p *Map[K, V]) LoadAndDelete(key K) (val V, loaded bool) {
	i, loaded := p.loadAndDelete(key)
	v, _ := i.(V)
	return v, loaded
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (p *Map[K, V]) LoadOrStore(key K, val V) (actual V, loaded bool) {
	i, loaded := p.loadOrStore(key, val)
	v, _ := i.(V)
	return v, loaded
}
