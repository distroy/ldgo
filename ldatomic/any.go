/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"reflect"
	"sync/atomic"
	"unsafe"
)

func NewAny[T any](d T) *Any[T] {
	v := &Any[T]{}
	v.Store(d)
	return v
}

type Any[T any] struct {
	v any
}

func (v *Any[T]) Store(d T)          { v.store(v.pack(d)) }
func (v *Any[T]) Load() T            { return v.load().Data }
func (v *Any[T]) Swap(new T) (old T) { return v.swap(v.pack(new)).Data }
func (v *Any[T]) CompareAndSwap(old, new T) (swapped bool) {
	return v.compareAndSwap(v.pack(old), v.pack(new))
}

func (v Any[T]) MarshalJSON() ([]byte, error)  { return marshalJSON[T](&v) }
func (v *Any[T]) UnmarshalJSON(b []byte) error { return unmarshalJSON[T](v, b) }

func (v *Any[T]) pack(d T) anyData[T] { return anyData[T]{Data: d} }

func (v *Any[T]) load() anyData[T] {
	vp := (*efaceWords)(unsafe.Pointer(v))
	typ := atomic.LoadPointer(&vp.typ)
	if typ == nil {
		// First store not yet completed.
		return anyData[T]{}
	}

	var val any
	data := atomic.LoadPointer(&vp.data)
	vlp := (*efaceWords)(unsafe.Pointer(&val))
	vlp.typ = typ
	vlp.data = data
	return val.(anyData[T])
}

func (v *Any[T]) store(val anyData[T]) {
	vli := any(val)
	vlp := (*efaceWords)(unsafe.Pointer(&vli))

	vp := (*efaceWords)(unsafe.Pointer(v))

	typ := atomic.LoadPointer(&vp.typ)
	if typ == nil {
		// Complete first store.
		runtime_procPin()
		atomic.StorePointer(&vp.data, vlp.data)
		atomic.StorePointer(&vp.typ, vlp.typ)
		runtime_procUnpin()
		return
	}

	atomic.StorePointer(&vp.data, vlp.data)
}

func (v *Any[T]) swap(new anyData[T]) (old anyData[T]) {
	ni := any(new)
	np := (*efaceWords)(unsafe.Pointer(&ni))

	vp := (*efaceWords)(unsafe.Pointer(v))
	typ := atomic.LoadPointer(&vp.typ)
	if typ == nil {
		runtime_procPin()
		if atomic.CompareAndSwapPointer(&vp.data, nil, np.data) {
			// Complete first store.
			atomic.StorePointer(&vp.typ, np.typ)
			runtime_procUnpin()
			return anyData[T]{}
		}
		runtime_procUnpin()
	}

	var oi any
	op := (*efaceWords)(unsafe.Pointer(&oi))
	op.typ, op.data = np.typ, atomic.SwapPointer(&vp.data, np.data)
	return oi.(anyData[T])
}

func (v *Any[T]) compareAndSwap(old, new anyData[T]) (swapped bool) {
	var zero anyData[T]
	// zi := any(zero)
	// zp := (*efaceWords)(unsafe.Pointer(&zi))

	// oi := any(old)
	// op := (*efaceWords)(unsafe.Pointer(&oi))

	ni := any(new)
	np := (*efaceWords)(unsafe.Pointer(&ni))

	vp := (*efaceWords)(unsafe.Pointer(v))

	typ := atomic.LoadPointer(&vp.typ)
	if typ == nil {
		if !v.equal(old, zero) {
			return false
		}
		// if oi != zero {
		// 	return false
		// }
		runtime_procPin()
		if !atomic.CompareAndSwapPointer(&vp.data, nil, np.data) {
			runtime_procUnpin()
			return false
		}
		// Complete first store.
		atomic.StorePointer(&vp.typ, np.typ)
		runtime_procUnpin()
		return true
	}

	data := atomic.LoadPointer(&vp.data)
	var i any
	(*efaceWords)(unsafe.Pointer(&i)).typ = typ
	(*efaceWords)(unsafe.Pointer(&i)).data = data
	// if i != (any)(old) {
	// 	return false
	// }
	if !v.equal(i.(anyData[T]), old) {
		return false
	}
	return atomic.CompareAndSwapPointer(&vp.data, data, np.data)
}

func (v *Any[T]) equal(a, b anyData[T]) bool {
	aa := reflect.ValueOf(&a.Data).Elem()
	bb := reflect.ValueOf(&b.Data).Elem()
	return equal(aa, bb)
}

// efaceWords is any internal representation.
type efaceWords struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}

type anyData[T any] struct {
	Data T
}

//go:linkname runtime_procPin sync/atomic.runtime_procPin
func runtime_procPin() int

//go:linkname runtime_procUnpin sync/atomic.runtime_procUnpin
func runtime_procUnpin()
