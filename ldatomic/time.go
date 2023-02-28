/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	storeTimeInProgress = time.Location{}
	zeroTime            = time.Time{}
	internalToUnix      = zeroTime.Unix()
)

type Time struct {
	// d    Pointer
	sec  int64
	nano int32
	loc  unsafe.Pointer
}

func NewTime(d time.Time) *Time {
	sec := d.Unix() - internalToUnix
	nano := d.Nanosecond()
	loc := d.Location()
	// log.Printf("new time. nano:%d, loc:%s, sec:%d, time:%s",
	// 	nano, loc.String(), d.Unix(), d.Format("2006-01-02T15:04:05-0700"))

	p := &Time{}
	p.storeSec(sec)
	p.storeNano(nano)
	p.storeLoc(loc)
	return p
}

func (p *Time) Store(d time.Time) {
	newSec := d.Unix() - internalToUnix
	newNano := d.Nanosecond()
	newLoc := d.Location()
	// log.Printf("store time. nano:%d, loc:%s, time:%s",
	// 	newNano, newLoc.String(), d.Format("2006-01-02T15:04:05-0700"))
	for {
		loc := p.loadLoc()
		if loc == &storeTimeInProgress {
			// Store in progress. Wait.
			// Since we disable preemption around the store,
			// we can wait with active spinning.
			continue
		}

		// Attempt to start store.
		// Disable preemption so that other goroutines can use
		// active spin wait to wait for completion.
		runtime_procPin()

		if !p.cmpAndSwapLoc(loc, &storeTimeInProgress) {
			runtime_procUnpin()
			continue
		}

		p.storeSec(newSec)
		p.storeNano(newNano)
		p.storeLoc(newLoc)
		runtime_procUnpin()

		return
	}
}

func (p *Time) Load() time.Time {
	for {
		loc := p.loadLoc()
		if loc == &storeTimeInProgress {
			continue
		}

		sec := p.loadSec()
		nano := p.loadNano()
		return p.makeTime(sec, nano, loc)
	}
}

func (p *Time) Swap(new time.Time) (old time.Time) {
	newSec := new.Unix() - internalToUnix
	newNano := new.Nanosecond()
	newLoc := new.Location()

	for {
		oldLoc := p.loadLoc()
		if oldLoc == &storeTimeInProgress {
			// Store in progress. Wait.
			// Since we disable preemption around the store,
			// we can wait with active spinning.
			continue
		}

		// Attempt to start store.
		// Disable preemption so that other goroutines can use
		// active spin wait to wait for completion.
		runtime_procPin()

		if !p.cmpAndSwapLoc(oldLoc, &storeTimeInProgress) {
			runtime_procUnpin()
			continue
		}

		oldSec := p.swapSec(newSec)
		oldNano := p.swapNano(newNano)
		p.storeLoc(newLoc)
		runtime_procUnpin()

		return p.makeTime(oldSec, oldNano, oldLoc)
	}
}

func (p *Time) CompareAndSwap(old, new time.Time) (swapped bool) {
	newSec := new.Unix() - internalToUnix
	newNano := new.Nanosecond()
	newLoc := new.Location()

	oldSec := old.Unix() - internalToUnix
	oldNano := old.Nanosecond()
	oldLoc := old.Location()

	for {
		pLoc := p.loadLoc()
		if pLoc == &storeTimeInProgress {
			continue

		} else if pLoc != oldLoc {
			return false
		}

		runtime_procPin()
		if !p.cmpAndSwapLoc(pLoc, &storeTimeInProgress) {
			runtime_procUnpin()
			return false
		}

		if !p.cmpAndSwapNano(oldNano, newNano) {
			p.storeLoc(pLoc)
			runtime_procUnpin()
			return false
		}

		if !p.cmpAndSwapSec(oldSec, newSec) {
			p.storeNano(oldNano)
			p.storeLoc(pLoc)
			runtime_procUnpin()
			return false
		}

		p.storeLoc(newLoc)
		runtime_procUnpin()
		return true
	}
}

func (p *Time) Add(d time.Duration) (new time.Time) {
	return p.MustChange(func(old time.Time) (new time.Time) {
		return old.Add(d)
	})
}

func (p *Time) AddDate(years int, months int, days int) (new time.Time) {
	return p.MustChange(func(old time.Time) (new time.Time) {
		return old.AddDate(years, months, days)
	})
}

func (p *Time) MustChange(change func(old time.Time) (new time.Time)) (new time.Time) {
	for {
		new, swapped := p.Change(change)
		if swapped {
			return new
		}
	}
}

func (p *Time) Change(change func(old time.Time) (new time.Time)) (new time.Time, swapped bool) {
	old := p.Load()
	new = change(old)
	return new, p.CompareAndSwap(old, new)
}

func (p *Time) makeTime(sec int64, nano int, loc *time.Location) time.Time {
	sec += internalToUnix
	if loc == nil {
		loc = time.UTC
	}

	t := time.Unix(sec, int64(nano))
	t = t.In(loc)

	// log.Printf("make time. sec:%d, nano:%d, loc:%s, time:%s",
	// 	sec, nano, loc.String(), t.Format("2006-01-02T15:04:05-0700"))

	return t
}

func (p *Time) storeSec(d int64)              { atomic.StoreInt64(&p.sec, d) }
func (p *Time) loadSec() int64                { return atomic.LoadInt64(&p.sec) }
func (p *Time) swapSec(new int64) (old int64) { return atomic.SwapInt64(&p.sec, new) }
func (p *Time) cmpAndSwapSec(old, new int64) (swapped bool) {
	return atomic.CompareAndSwapInt64(&p.sec, old, new)
}

func (p *Time) storeNano(d int)            { atomic.StoreInt32(&p.nano, int32(d)) }
func (p *Time) loadNano() int              { return int(atomic.LoadInt32(&p.nano)) }
func (p *Time) swapNano(new int) (old int) { return int(atomic.SwapInt32(&p.nano, int32(new))) }
func (p *Time) cmpAndSwapNano(old, new int) (swapped bool) {
	return atomic.CompareAndSwapInt32(&p.nano, int32(old), int32(new))
}

func (p *Time) storeLoc(d *time.Location) { atomic.StorePointer(&p.loc, unsafe.Pointer(d)) }
func (p *Time) loadLoc() *time.Location   { return (*time.Location)(atomic.LoadPointer(&p.loc)) }
func (p *Time) swapLoc(new *time.Location) (old *time.Location) {
	return (*time.Location)(atomic.SwapPointer(&p.loc, unsafe.Pointer(new)))
}
func (p *Time) cmpAndSwapLoc(old, new *time.Location) (swapped bool) {
	return atomic.CompareAndSwapPointer(&p.loc, unsafe.Pointer(old), unsafe.Pointer(new))
}
