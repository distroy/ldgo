/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"sync/atomic"
	"unsafe"
)

var _ = (*atomic.Value)(nil)

var firstStoreValueInProgress byte

//go:linkname runtime_procPin sync/atomic.runtime_procPin
func runtime_procPin()

//go:linkname runtime_procUnpin sync/atomic.runtime_procUnpin
func runtime_procUnpin()

// A Value provides an atomic load and store of a consistently typed value.
// The zero value for a Value returns nil from Load.
// Once Store has been called, a Value must not be copied.
//
// A Value must not be copied after first use.
type Value struct {
	v interface{}
}

func NewValue(d interface{}) *Value {
	p := &Value{}
	if d != nil {
		p.Store(d)
	}
	return p
}

// ifaceWords is interface{} internal representation.
type ifaceWords struct {
	typ  unsafe.Pointer
	data unsafe.Pointer
}

// Load returns the value set by the most recent Store.
// It returns nil if there has been no call to Store for this Value.
func (v *Value) Load() (val interface{}) {
	vw := (*ifaceWords)(unsafe.Pointer(v))
	typ := atomic.LoadPointer(&vw.typ)
	if typ == nil || typ == unsafe.Pointer(&firstStoreValueInProgress) {
		// First store not yet completed.
		return nil
	}
	data := atomic.LoadPointer(&vw.data)
	nw := (*ifaceWords)(unsafe.Pointer(&val))
	nw.typ = typ
	nw.data = data
	return
}

// Store sets the value of the Value to x.
// All calls to Store for a given Value must use values of the same concrete type.
// Store of an inconsistent type panics, as does Store(nil).
func (v *Value) Store(val interface{}) {
	if val == nil {
		panic("ldatomic: store of nil value into Value")
	}
	vw := (*ifaceWords)(unsafe.Pointer(v))
	nw := (*ifaceWords)(unsafe.Pointer(&val))
	for {
		typ := atomic.LoadPointer(&vw.typ)
		if typ == nil {
			// Attempt to start first store.
			// Disable preemption so that other goroutines can use
			// active spin wait to wait for completion.
			runtime_procPin()
			if !atomic.CompareAndSwapPointer(&vw.typ, nil, unsafe.Pointer(&firstStoreValueInProgress)) {
				runtime_procUnpin()
				continue
			}
			// Complete first store.
			atomic.StorePointer(&vw.data, nw.data)
			atomic.StorePointer(&vw.typ, nw.typ)
			runtime_procUnpin()
			return

		} else if typ == unsafe.Pointer(&firstStoreValueInProgress) {
			// First store in progress. Wait.
			// Since we disable preemption around the first store,
			// we can wait with active spinning.
			continue
		}

		// First store completed. Check type and overwrite data.
		if typ != nw.typ {
			panic("ldatomic: store of inconsistently typed value into Value")
		}
		atomic.StorePointer(&vw.data, nw.data)
		return
	}
}

// Swap stores new into Value and returns the previous value. It returns nil if
// the Value is empty.
//
// All calls to Swap for a given Value must use values of the same concrete
// type. Swap of an inconsistent type panics, as does Swap(nil).
func (v *Value) Swap(new interface{}) (old interface{}) {
	if new == nil {
		panic("ldatomic: swap of nil value into Value")
	}
	vw := (*ifaceWords)(unsafe.Pointer(v))
	nw := (*ifaceWords)(unsafe.Pointer(&new))
	for {
		typ := atomic.LoadPointer(&vw.typ)
		if typ == nil {
			// Attempt to start first store.
			// Disable preemption so that other goroutines can use
			// active spin wait to wait for completion; and so that
			// GC does not see the fake type accidentally.
			runtime_procPin()
			if !atomic.CompareAndSwapPointer(&vw.typ, nil, unsafe.Pointer(&firstStoreValueInProgress)) {
				runtime_procUnpin()
				continue
			}
			// Complete first store.
			atomic.StorePointer(&vw.data, nw.data)
			atomic.StorePointer(&vw.typ, nw.typ)
			runtime_procUnpin()
			return nil

		} else if typ == unsafe.Pointer(&firstStoreValueInProgress) {
			// First store in progress. Wait.
			// Since we disable preemption around the first store,
			// we can wait with active spinning.
			continue
		}

		// First store completed. Check type and overwrite data.
		if typ != nw.typ {
			panic("ldatomic: swap of inconsistently typed value into Value")
		}
		op := (*ifaceWords)(unsafe.Pointer(&old))
		op.typ, op.data = nw.typ, atomic.SwapPointer(&vw.data, nw.data)
		return old
	}
}

// CompareAndSwap executes the compare-and-swap operation for the Value.
//
// All calls to CompareAndSwap for a given Value must use values of the same
// concrete type. CompareAndSwap of an inconsistent type panics, as does
// CompareAndSwap(old, nil).
func (v *Value) CompareAndSwap(old, new interface{}) (swapped bool) {
	if new == nil {
		panic("ldatomic: compare and swap of nil value into Value")
	}
	vw := (*ifaceWords)(unsafe.Pointer(v))
	nw := (*ifaceWords)(unsafe.Pointer(&new))
	ow := (*ifaceWords)(unsafe.Pointer(&old))
	if ow.typ != nil && nw.typ != ow.typ {
		panic("ldatomic: compare and swap of inconsistently typed values")
	}

	var typ unsafe.Pointer
	for {
		typ = atomic.LoadPointer(&vw.typ)
		if typ == nil {
			if old != nil {
				return false
			}
			// Attempt to start first store.
			// Disable preemption so that other goroutines can use
			// active spin wait to wait for completion; and so that
			// GC does not see the fake type accidentally.
			runtime_procPin()
			if !atomic.CompareAndSwapPointer(&vw.typ, nil, unsafe.Pointer(&firstStoreValueInProgress)) {
				runtime_procUnpin()
				continue
			}
			// Complete first store.
			atomic.StorePointer(&vw.data, nw.data)
			atomic.StorePointer(&vw.typ, nw.typ)
			runtime_procUnpin()
			return true

		} else if typ == unsafe.Pointer(&firstStoreValueInProgress) {
			// First store in progress. Wait.
			// Since we disable preemption around the first store,
			// we can wait with active spinning.
			continue
		}

		break
	}

	// First store completed. Check type and overwrite data.
	if typ != nw.typ {
		panic("ldatomic: compare and swap of inconsistently typed value into Value")
	}

	// Compare old and current via runtime equality check.
	// This allows value types to be compared, something
	// not offered by the package functions.
	// CompareAndSwapPointer below only ensures vp.data
	// has not changed since LoadPointer.
	data := atomic.LoadPointer(&vw.data)
	var i interface{}
	(*ifaceWords)(unsafe.Pointer(&i)).typ = typ
	(*ifaceWords)(unsafe.Pointer(&i)).data = data
	if i != old {
		return false
	}
	return atomic.CompareAndSwapPointer(&vw.data, data, nw.data)
}
