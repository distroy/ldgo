/*
 * Copyright (C) distroy
 */

package ldasync

import (
	"sync"
	"sync/atomic"

	"github.com/distroy/ldgo/v3/ldsync"
)

const (
	defaultConcurrency = 1
	defaultChannelSize = 1
)

type asyncBase[T any] struct {
	wg       sync.WaitGroup
	mu       sync.Mutex
	nodes    *asyncNode
	capacity int32
	count    int32
	ch       chan T
}

func (p *asyncBase[T]) async() chan T { return p.ch }
func (p *asyncBase[T]) getLen() int   { return int(atomic.LoadInt32(&p.count)) }
func (p *asyncBase[T]) getCap() int   { return int(atomic.LoadInt32(&p.capacity)) }
func (p *asyncBase[T]) setCap(n int)  { atomic.StoreInt32(&p.capacity, int32(n)) }

func (p *asyncBase[T]) start(concurrency int, run func(fn T)) {
	n := concurrency
	if n <= 0 {
		n = defaultConcurrency
	}
	p.reset(n, run)
}

func (p *asyncBase[T]) reset(concurrency int, run func(fn T)) {
	n := concurrency
	if n <= 0 {
		return
	}
	if p.getCap() == n {
		return
	}
	p.mu.Lock()
	p.startWithoutLock(n, run)
	p.mu.Unlock()
}

func (p *asyncBase[T]) init(run func(fn T)) {
	if p.getCap() > 0 {
		return
	}

	p.mu.Lock()
	if p.getCap() <= 0 {
		p.startWithoutLock(defaultConcurrency, run)
	}
	p.mu.Unlock()
}

func (p *asyncBase[T]) startWithoutLock(n int, run func(fn T)) {
	num := p.getCap()
	if num == 0 {
		p.ch = make(chan T, defaultChannelSize)
	}

	if n == num {
		return
	}
	if n > num {
		p.incrNodesTo(n, run)
		return
	}
	p.decrNodesTo(n)
}

func (p *asyncBase[T]) incrNodesTo(to int, run func(fn T)) {
	delta := to - p.getCap()
	p.setCap(to)

	p.wg.Add(delta)
	for i := 0; i < delta; i++ {
		node := getAsyncNode()

		node.next = p.nodes
		p.nodes = node

		go p.main(node, run)
	}
}

func (p *asyncBase[T]) decrNodesTo(to int) {
	delta := p.getCap() - to
	p.setCap(to)

	for i := 0; i < delta; i++ {
		node := p.nodes
		if node == nil {
			return
		}
		p.nodes = node.next
		node.next = nil
		close(node.ch)
	}
}

func (p *asyncBase[T]) wait() { p.wg.Wait() }
func (p *asyncBase[T]) close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	num := p.getCap()
	if num == 0 {
		return
	}

	p.setCap(0)
	p.nodes = nil
	close(p.ch)
}

func (p *asyncBase[T]) main(node *asyncNode, run func(fn T)) {
	defer p.wg.Done()

	for {
		select {
		case <-node.ch:
			putAsyncNode(node)
			return

		case fn, ok := <-p.ch:
			if !ok {
				putAsyncNode(node)
				return
			}
			atomic.AddInt32(&p.count, 1)
			// p.doWithRecover(fn)
			run(fn)
			atomic.AddInt32(&p.count, -1)
		}
	}
}

var (
	asyncNodePool = &ldsync.Pool[*asyncNode]{}
)

type asyncNode struct {
	ch   chan struct{}
	next *asyncNode
}

func getAsyncNode() *asyncNode {
	p := asyncNodePool

	n := p.Get()
	if n == nil {
		n = &asyncNode{}
	}

	n.ch = make(chan struct{})
	n.next = nil
	return n
}

func putAsyncNode(n *asyncNode) {
	p := asyncNodePool
	if n == nil {
		return
	}
	p.Put(n)
}
