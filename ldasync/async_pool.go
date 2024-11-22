/*
 * Copyright (C) distroy
 */

package ldasync

import (
	"runtime/debug"
)

type AsyncPool struct {
	asyncBase[func()]
}

func NewAsyncPool(concurrency int) *AsyncPool {
	p := &AsyncPool{}
	p.Start(concurrency)
	return p
}

func (p *AsyncPool) Start(concurrency int) { p.asyncBase.start(concurrency, p.doWithRecover) }
func (p *AsyncPool) Reset(concurrency int) { p.asyncBase.reset(concurrency, p.doWithRecover) }

func (p *AsyncPool) Capacity() int { return p.asyncBase.getCap() }
func (p *AsyncPool) Running() int  { return p.asyncBase.getLen() }

func (p *AsyncPool) SetLogger(l Logger) { p.asyncBase.setLogger(l) }

func (p *AsyncPool) init() {
	p.asyncBase.init(p.doWithRecover)
}

func (p *AsyncPool) Wait()  { p.asyncBase.wait() }
func (p *AsyncPool) Close() { p.asyncBase.close() }

func (p *AsyncPool) Async() chan<- func() {
	p.init()
	return p.asyncBase.async()
}

func (p *AsyncPool) doWithRecover(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			buf := debug.Stack()
			p.getLogger().Printf("[async pool] do async func panic. err:%v, stack:\n%s", err, buf)
		}
	}()

	fn()
}
