/*
 * Copyright (C) distroy
 */

package ldgo

import (
	"log"
	"runtime/debug"
)

type GoPool interface {
	Go(fn func(done <-chan struct{}))
	GoN(n int, fn func(done <-chan struct{}))
	Run(fn func())
	RunN(n int, fn func())
	Stop()
	Wait()
}

func NewGoPool() GoPool {
	return goPool{
		done: newDoneWait(),
	}
}

func GoPoolGo(fn func(done <-chan struct{})) GoPool { return GoPoolGoN(1, fn) }
func GoPoolRun(fn func()) GoPool                    { return GoPoolRunN(1, fn) }

func GoPoolGoN(n int, fn func(done <-chan struct{})) GoPool {
	pool := NewGoPool()
	pool.GoN(n, fn)
	return pool
}

func GoPoolRunN(n int, fn func()) GoPool {
	pool := NewGoPool()
	pool.RunN(n, fn)
	return pool
}

type goPool struct {
	done *doneWait
}

func (that goPool) Stop() { that.done.Stop() }
func (that goPool) Wait() { that.done.Wait() }

func (that goPool) Go(fn func(done <-chan none)) { that.GoN(1, fn) }
func (that goPool) Run(fn func())                { that.RunN(1, fn) }

func (that goPool) GoN(n int, fn func(done <-chan none)) {
	fnGo := func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err, BytesToStrUnsafe(debug.Stack()))
			}
			that.done.Done()
		}()

		fn(that.done.Chan())
	}
	that.run(n, fnGo)
}

func (that goPool) RunN(n int, fn func()) {
	fnGo := func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err, BytesToStrUnsafe(debug.Stack()))
			}
			that.done.Done()
		}()

		fn()
	}
	that.run(n, fnGo)
}

func (that goPool) run(n int, fn func()) {
	that.done.Add(n)
	for i := 0; i < n; i++ {
		go fn()
	}
}
