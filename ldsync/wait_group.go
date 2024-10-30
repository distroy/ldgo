/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sync/atomic"
	"unsafe"

	"github.com/distroy/ldgo/v3/internal/race"
)

// A WaitGroup waits for a collection of goroutines to finish.
// The main goroutine calls Add to set the number of
// goroutines to wait for. Then each of the goroutines
// runs and calls Done when finished. At the same time,
// Wait can be used to block until all goroutines have finished.
//
// A WaitGroup must not be copied after first use.
//
// In the terminology of the Go memory model, a call to Done
// “synchronizes before” the return of any Wait call that it unblocks.
type WaitGroup struct {
	noCopy noCopy

	state uint64 // high 32 bits are counter, low 32 bits are waiter count.
	sema  uint32
}

func (wg *WaitGroup) Count() int {
	state := atomic.LoadUint64(&wg.state)
	return int(int32(state >> 32))
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a WaitGroup is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the WaitGroup example.
func (wg *WaitGroup) Add(delta int) {
	if race.Enabled {
		if delta < 0 {
			// Synchronize decrements with Wait.
			race.ReleaseMerge(unsafe.Pointer(wg))
		}
		race.Disable()
		defer race.Enable()
	}
	state := atomic.AddUint64(&wg.state, uint64(delta)<<32)
	v := int32(state >> 32)
	w := uint32(state)
	if race.Enabled && delta > 0 && v == int32(delta) {
		// The first increment must be synchronized with Wait.
		// Need to model this as a read, because there can be
		// several concurrent wg.counter transitions from 0.
		race.Read(unsafe.Pointer(&wg.sema))
	}
	if v < 0 {
		panic("ldsync: negative WaitGroup counter")
	}
	if w != 0 && delta > 0 && v == int32(delta) {
		panic("ldsync: WaitGroup misuse: Add called concurrently with Wait")
	}
	if v > 0 || w == 0 {
		return
	}
	// This goroutine has set counter to 0 when waiters > 0.
	// Now there can't be concurrent mutations of state:
	// - Adds must not happen concurrently with Wait,
	// - Wait does not increment waiters if it sees counter == 0.
	// Still do a cheap sanity check to detect WaitGroup misuse.
	if atomic.LoadUint64(&wg.state) != state {
		panic("ldsync: WaitGroup misuse: Add called concurrently with Wait")
	}
	// Reset waiters count to 0.
	atomic.StoreUint64(&wg.state, 0)
	for ; w != 0; w-- {
		runtime_Semrelease(&wg.sema, false, 0)
	}
}

// Done decrements the WaitGroup counter by one.
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

// Wait blocks until the WaitGroup counter is zero.
func (wg *WaitGroup) Wait() {
	if race.Enabled {
		race.Disable()
	}
	for {
		ok := wg.waitOnce()
		if ok {
			if race.Enabled {
				race.Enable()
				race.Acquire(unsafe.Pointer(wg))
			}
			return
		}
	}
}

func (wg *WaitGroup) waitOnce() bool {
	state := atomic.LoadUint64(&wg.state)
	v := int32(state >> 32)
	w := uint32(state)
	if v == 0 {
		// Counter is 0, no need to wait.
		return true
	}
	// Increment waiters count.
	if atomic.CompareAndSwapUint64(&wg.state, state, state+1) {
		if race.Enabled && w == 0 {
			// Wait must be synchronized with the first Add.
			// Need to model this is as a write to race with the read in Add.
			// As a consequence, can do the write only for the first waiter,
			// otherwise concurrent Waits will race with each other.
			race.Write(unsafe.Pointer(&wg.sema))
		}
		runtime_Semacquire(&wg.sema)
		if atomic.LoadUint64(&wg.state) != 0 {
			panic("ldsync: WaitGroup is reused before previous Wait has returned")
		}
		return true
	}

	return false
}
