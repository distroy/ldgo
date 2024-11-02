/*
 * Copyright (C) distroy
 */

package ldasync

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestErrGroup(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("nomarl", func(c convey.C) {
			n := int32(0)
			fn := func() error {
				time.Sleep(time.Millisecond * 10)
				atomic.AddInt32(&n, 1)
				return nil
			}

			p := &ErrGroup{}
			p.Start(10)

			for i := 0; i < 10; i++ {
				p.Async() <- fn
			}

			c.So(atomic.LoadInt32(&n), convey.ShouldEqual, 0)

			p.Close()
			err := p.Wait()
			c.So(err, convey.ShouldBeNil)

			c.So(atomic.LoadInt32(&n), convey.ShouldEqual, 10)
		})

		c.Convey("panic", func(c convey.C) {
			fn := func() error {
				panic(11)
			}

			var errByOn error

			p := &ErrGroup{}
			p.Start(10)
			p.OnError(func(err error) {
				errByOn = err
			})

			p.Reset(5)
			p.Async() <- fn

			p.Close()
			err := p.Wait()
			c.So(err, convey.ShouldNotBeNil)
			c.So(err, convey.ShouldEqual, errByOn)
		})

		c.Convey("size", func(c convey.C) {
			// var seq int32
			fn := func() error {
				// id := atomic.AddInt32(&seq, 1)
				// t.Logf("go[%d] begin", id)
				time.Sleep(time.Second * 1)
				// t.Logf("go[%d] end", id)
				return nil
			}
			p := &ErrGroup{}

			// p.Start(1)
			c.So(p.Capacity(), convey.ShouldEqual, 0)

			p.Async() <- fn
			c.So(p.Capacity(), convey.ShouldEqual, 1)

			p.Async() <- fn
			c.So(p.Capacity(), convey.ShouldEqual, 1)
			c.So(p.Running(), convey.ShouldEqual, 1)

			p.Reset(2)
			time.Sleep(time.Millisecond)
			c.So(p.Capacity(), convey.ShouldEqual, 2)
			c.So(p.Running(), convey.ShouldEqual, 2)

			p.Reset(1)
			c.So(p.Capacity(), convey.ShouldEqual, 1)
			c.So(p.Running(), convey.ShouldEqual, 2)

			p.Close()
			err := p.Wait()
			c.So(err, convey.ShouldBeNil)

			c.So(p.Running(), convey.ShouldEqual, 0)
		})
	})
}
