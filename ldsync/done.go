/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sync"
	"sync/atomic"
)

type none = struct{}

type Done struct {
	done   chan none
	mutex  sync.Mutex
	inited uint32
	stoped uint32
}

func (p *Done) Stop() {
	p.init()
	stoped := once{mutex: &p.mutex, done: &p.stoped}
	stoped.Do(func() error {
		close(p.done)
		return nil
	})
}

func (p *Done) Chan() <-chan none {
	p.init()
	return p.done
}

func (p *Done) init() {
	inited := once{mutex: &p.mutex, done: &p.inited}
	inited.Do(func() error {
		p.done = make(chan struct{})
		return nil
	})
}

func (p *Done) Reset() {
	p.mutex.Lock()
	atomic.StoreUint32(&p.inited, 0)
	atomic.StoreUint32(&p.stoped, 0)
	p.mutex.Unlock()
}
