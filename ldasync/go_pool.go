/*
 * Copyright (C) distroy
 */

package ldasync

import (
	"log"
	"runtime/debug"

	"github.com/distroy/ldgo/v3/ldsync"
)

func Go(fn func()) *GoPool { return GoN(1, fn) }
func GoN(n int, fn func()) *GoPool {
	p := &GoPool{}
	p.GoN(n, fn)
	return p
}

type GoPool struct {
	wg ldsync.WaitGroup
}

func (p *GoPool) Wait()        { p.wg.Wait() }
func (p *GoPool) Count() int   { return p.wg.Count() }
func (p *GoPool) Go(fn func()) { p.GoN(1, fn) }
func (p *GoPool) GoN(n int, fn func()) {
	if n <= 0 {
		n = 1
	}

	fnGo := func() {
		defer func() {
			p.wg.Done()

			if err := recover(); err != nil {
				buf := debug.Stack()

				// log.Println(err, ldconv.BytesToStrUnsafe(buf))
				log.Printf("[go pool] go func panic. err:%v, stack:\n%s", err, buf)
			}
		}()

		fn()
	}

	p.wg.Add(n)
	for i := 0; i < n; i++ {
		go fnGo()
	}
}
