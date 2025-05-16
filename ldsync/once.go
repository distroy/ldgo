/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	mutex sync.Mutex
	done  uint32
}

func (o *Once) get() once {
	return once{
		mutex: &o.mutex,
		done:  &o.done,
	}
}

func (o *Once) Done() bool               { return o.get().Done() }
func (o *Once) Reset()                   { o.get().Reset() }
func (o *Once) Do(fn func() error) error { return o.get().Do(fn) }

type once struct {
	done  *uint32
	mutex *sync.Mutex
}

func (o once) Done() bool { return atomic.LoadUint32(o.done) != 0 }

func (o once) Reset() {
	o.mutex.Lock()
	atomic.StoreUint32(o.done, 0)
	o.mutex.Unlock()
}

func (o once) Do(fn func() error) error {
	if atomic.LoadUint32(o.done) != 0 {
		return nil
	}

	o.mutex.Lock()
	defer o.mutex.Unlock()

	if atomic.LoadUint32(o.done) != 0 {
		return nil
	}

	err := fn()
	if err != nil {
		return err
	}

	atomic.StoreUint32(o.done, 1)
	return nil
}
