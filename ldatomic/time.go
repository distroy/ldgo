/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"time"
	"unsafe"
)

var (
	storeTimeInProgress = time.Location{}
	zeroTime            = time.Time{}
	internalToUnix      = zeroTime.Unix()
)

type Time struct {
	sec  Int64
	nano Int
	loc  Pointer
}

func NewTime(d time.Time) *Time {
	sec := d.Unix() - internalToUnix
	nano := d.Nanosecond()
	loc := d.Location()
	// log.Printf("new time. nano:%d, loc:%s, sec:%d, time:%s",
	// 	nano, loc.String(), d.Unix(), d.Format("2006-01-02T15:04:05-0700"))

	p := &Time{}
	p.sec.Store(sec)
	p.nano.Store(nano)
	p.loc.Store(unsafe.Pointer(loc))
	return p
}

func (p *Time) Store(d time.Time) {
	newSec := d.Unix() - internalToUnix
	newNano := d.Nanosecond()
	newLoc := d.Location()
	// log.Printf("store time. nano:%d, loc:%s, time:%s",
	// 	newNano, newLoc.String(), d.Format("2006-01-02T15:04:05-0700"))
	for {
		loc := (*time.Location)(p.loc.Load())
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

		if !p.loc.CompareAndSwap(unsafe.Pointer(loc), unsafe.Pointer(&storeTimeInProgress)) {
			runtime_procUnpin()
			continue
		}

		p.sec.Store(newSec)
		p.nano.Store(newNano)
		p.loc.Store(unsafe.Pointer(newLoc))
		runtime_procUnpin()

		return
	}
}

func (p *Time) Load() time.Time {
	for {
		loc := (*time.Location)(p.loc.Load())
		if loc == &storeTimeInProgress {
			continue
		}

		sec := p.sec.Load()
		nano := p.nano.Load()
		return p.makeTime(sec, nano, loc)
	}
}

func (p *Time) Swap(new time.Time) (old time.Time) {
	newSec := new.Unix() - internalToUnix
	newNano := new.Nanosecond()
	newLoc := new.Location()

	for {
		oldLoc := (*time.Location)(p.loc.Load())
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

		if !p.loc.CompareAndSwap(unsafe.Pointer(oldLoc), unsafe.Pointer(&storeTimeInProgress)) {
			runtime_procUnpin()
			continue
		}

		oldSec := p.sec.Swap(newSec)
		oldNano := p.nano.Swap(newNano)
		p.loc.Store(unsafe.Pointer(newLoc))
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
		pLoc := (*time.Location)(p.loc.Load())
		if pLoc == &storeTimeInProgress {
			continue

		} else if pLoc != oldLoc {
			return false
		}

		runtime_procPin()
		if !p.loc.CompareAndSwap(unsafe.Pointer(pLoc), unsafe.Pointer(&storeTimeInProgress)) {
			runtime_procUnpin()
			return false
		}

		if !p.nano.CompareAndSwap(oldNano, newNano) {
			p.loc.Store(unsafe.Pointer(pLoc))
			runtime_procUnpin()
			return false
		}

		if !p.sec.CompareAndSwap(oldSec, newSec) {
			p.nano.Store(oldNano)
			p.loc.Store(unsafe.Pointer(pLoc))
			runtime_procUnpin()
			return false
		}

		p.loc.Store(unsafe.Pointer(newLoc))
		runtime_procUnpin()
		return true
	}
}

func (p *Time) Add(d time.Duration) (new time.Time) {
	sec := d / time.Second
	nano := d % time.Second

	for {
		loc := (*time.Location)(p.loc.Load())
		if loc == &storeTimeInProgress {
			continue
		}

		runtime_procPin()

		if !p.loc.CompareAndSwap(unsafe.Pointer(loc), unsafe.Pointer(&storeTimeInProgress)) {
			runtime_procUnpin()
			continue
		}

		newNano := p.nano.Add(int(nano))
		if newNano < 0 {
			newNano = p.nano.Add(int(time.Second))
			sec--

		} else if newNano >= int(time.Second) {
			newNano = p.nano.Sub(int(time.Second))
			sec++
		}

		newSec := p.sec.Add(int64(sec))

		p.loc.Store(unsafe.Pointer(loc))
		runtime_procUnpin()

		return p.makeTime(newSec, newNano, loc)
	}
}

func (p *Time) Sub(d time.Time) time.Duration {
	return p.Load().Sub(d)
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
