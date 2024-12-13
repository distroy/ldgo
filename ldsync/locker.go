/*
 * Copyright (C) distroy
 */

package ldsync

import "sync"

var (
	_ Locker = (*AutoLocker)(nil)
)

type Locker interface {
	sync.Locker

	TryLock() bool
}

func DiscardLocker() Locker { return discardLocker{} }

type discardLocker struct{}

func (_ discardLocker) Lock()         {}
func (_ discardLocker) Unlock()       {}
func (_ discardLocker) TryLock() bool { return false }

func AutoLock(l sync.Locker) *AutoLocker {
	switch v := l.(type) {
	case *AutoLocker:
		return v
	}
	return &AutoLocker{Locker: l}
}

type AutoLocker struct {
	Locker sync.Locker

	locked bool
}

func (l *AutoLocker) Lock() {
	if l.locked {
		return
	}
	l.Locker.Lock()
	l.locked = true
}

func (l *AutoLocker) Unlock() {
	if !l.locked || l.Locker == nil {
		return
	}
	l.Locker.Unlock()
	l.locked = false
}

func (l *AutoLocker) TryLock() bool {
	if l.locked {
		return true
	}
	locker, _ := l.Locker.(Locker)
	if locker == nil {
		return false
	}
	ok := locker.TryLock()
	l.locked = ok
	return ok
}
