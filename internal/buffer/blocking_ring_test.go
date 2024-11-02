/*
 * Copyright (C) distroy
 */

package buffer

import (
	"io"
	"sync"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestBlockingRing(t *testing.T) {
	// io.Pipe()

	convey.Convey(t.Name(), t, func(c convey.C) {
		var (
			wg  = &sync.WaitGroup{}
			tt0 time.Time
			tt1 time.Time
			tt2 time.Time
		)

		fGo := func(f func()) {
			wg.Add(1)
			go func() {
				defer func() {
					wg.Done()
				}()
				f()
			}()
		}
		fNew := func(n int) *BlockingRing[byte] {
			b := &BlockingRing[byte]{
				Buf: MakeRing(make([]byte, n, n)),
			}
			b.Init()
			return b
		}
		fSleep := func() { time.Sleep(10 * time.Millisecond) }

		c.Convey("single", func(c convey.C) {
			b := fNew(4)

			c.Convey("new", func(c convey.C) {
				c.So(b, convey.ShouldNotBeNil)
				c.So(b.Size(), convey.ShouldEqual, 0)
			})

			c.Convey("pop before put", func(c convey.C) {
				fGo(func() {
					d, err := b.Pop()
					c.So(err, convey.ShouldBeNil)
					c.So(d, convey.ShouldBeIn, []byte{'1', '2'})
					tt1 = time.Now()
				})
				fGo(func() {
					d, err := b.Pop()
					c.So(err, convey.ShouldBeNil)
					c.So(d, convey.ShouldBeIn, []byte{'1', '2'})
					tt2 = time.Now()
				})

				fSleep()

				tt0 = time.Now()
				c.So(b.Put('1'), convey.ShouldBeNil)
				c.So(b.Put('2'), convey.ShouldBeNil)

				wg.Wait()

				c.So(b.Size(), convey.ShouldEqual, 0)
				c.So(tt0, convey.ShouldHappenBefore, tt1)
				c.So(tt0, convey.ShouldHappenBefore, tt2)
			})

			c.Convey("put when full, then pop", func(c convey.C) {
				b.Put('1')
				b.Put('2')
				b.Put('3')
				b.Put('4')
				c.So(b.Cap(), convey.ShouldEqual, 4)
				c.So(b.Size(), convey.ShouldEqual, 4)

				fGo(func() {
					err := b.Put('5')
					c.So(err, convey.ShouldBeNil)
					tt1 = time.Now()
				})
				fGo(func() {
					err := b.Put('6')
					c.So(err, convey.ShouldBeNil)
					tt2 = time.Now()
				})

				fSleep()
				tt0 = time.Now()

				d, err := b.Pop()
				c.So(err, convey.ShouldBeNil)
				c.So(d, convey.ShouldEqual, '1')

				d, err = b.Pop()
				c.So(err, convey.ShouldBeNil)
				c.So(d, convey.ShouldEqual, '2')

				wg.Wait()
				c.So(b.Size(), convey.ShouldEqual, 4)
				c.So(tt0, convey.ShouldHappenBefore, tt1)
				c.So(tt0, convey.ShouldHappenBefore, tt2)

				c.So(b.Close(), convey.ShouldBeNil)
				c.So(b.Closed(), convey.ShouldEqual, true)

				d, err = b.Pop()
				c.So(err, convey.ShouldBeNil)
				c.So(d, convey.ShouldEqual, '3')

				d, err = b.Pop()
				c.So(err, convey.ShouldBeNil)
				c.So(d, convey.ShouldEqual, '4')

				d, err = b.Pop()
				c.So(err, convey.ShouldBeNil)
				c.So(d, convey.ShouldBeIn, []byte{'5', '6'})

				d, err = b.Pop()
				c.So(err, convey.ShouldBeNil)
				c.So(d, convey.ShouldBeIn, []byte{'5', '6'})

				c.So(b.Size(), convey.ShouldEqual, 0)

				d, err = b.Pop()
				c.So(err, convey.ShouldEqual, io.EOF)
				c.So(d, convey.ShouldEqual, 0)
			})

			c.Convey("pop & put after close", func(c convey.C) {
				c.So(b.Close(), convey.ShouldBeNil)
				c.So(b.Closed(), convey.ShouldEqual, true)

				err := b.Put('1')
				c.So(err, convey.ShouldEqual, io.ErrUnexpectedEOF)

				d, err := b.Pop()
				c.So(err, convey.ShouldEqual, io.EOF)
				c.So(d, convey.ShouldEqual, 0)
			})
		})

		c.Convey("multi", func(c convey.C) {
			b := fNew(16)
			c.Convey("new", func(c convey.C) {
				c.So(b, convey.ShouldNotBeNil)
				c.So(b.Size(), convey.ShouldEqual, 0)
			})

			c.Convey("read before write", func(c convey.C) {
				fGo(func() {
					d := make([]byte, 4)
					n, err := b.Read(d)
					c.So(err, convey.ShouldBeNil)
					c.So(n, convey.ShouldEqual, 4)
					c.So(string(d), convey.ShouldBeIn, []string{"0123", "4567"})
					tt1 = time.Now()
				})
				fGo(func() {
					d := make([]byte, 4)
					n, err := b.Read(d)
					c.So(err, convey.ShouldBeNil)
					c.So(n, convey.ShouldEqual, 4)
					c.So(string(d), convey.ShouldBeIn, []string{"0123", "4567"})
					tt2 = time.Now()
				})

				fSleep()
				tt0 = time.Now()

				n, err := b.Write([]byte("0123456789012345"))
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 16)

				wg.Wait()
				c.So(b.Size(), convey.ShouldEqual, 8)

				c.So(tt0, convey.ShouldHappenBefore, tt1)
				c.So(tt0, convey.ShouldHappenBefore, tt2)
			})

			c.Convey("write when full, then read", func(c convey.C) {
				n, err := b.Write([]byte("0123456789012345"))
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 16)
				c.So(b.Cap(), convey.ShouldEqual, 16)
				c.So(b.Size(), convey.ShouldEqual, 16)

				fGo(func() {
					n, err := b.Write([]byte("0123"))
					c.So(err, convey.ShouldBeNil)
					c.So(n, convey.ShouldEqual, 4)
					tt1 = time.Now()
				})
				fGo(func() {
					n, err := b.Write([]byte("4567"))
					c.So(err, convey.ShouldBeNil)
					c.So(n, convey.ShouldEqual, 4)
					tt2 = time.Now()
				})

				fSleep()
				tt0 = time.Now()

				d := make([]byte, 32)
				n, err = b.Read(d)
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 16)
				c.So(string(d[:n]), convey.ShouldEqual, "0123456789012345")

				wg.Wait()
				c.So(b.Size(), convey.ShouldEqual, 8)
				c.So(tt0, convey.ShouldHappenBefore, tt1)
				c.So(tt0, convey.ShouldHappenBefore, tt2)

				c.So(b.Close(), convey.ShouldBeNil)
				c.So(b.Closed(), convey.ShouldEqual, true)

				n, err = b.Read(d)
				c.So(err, convey.ShouldBeNil)
				c.So(n, convey.ShouldEqual, 8)
				c.So(string(d[:n]), convey.ShouldBeIn, []string{"01234567", "45670123"})
			})

			c.Convey("read & write after close", func(c convey.C) {
				c.So(b.Close(), convey.ShouldBeNil)
				c.So(b.Closed(), convey.ShouldEqual, true)

				n, err := b.Write([]byte("1"))
				c.So(err, convey.ShouldEqual, io.ErrUnexpectedEOF)
				c.So(n, convey.ShouldEqual, 0)

				d := make([]byte, 4)
				n, err = b.Read(d)
				c.So(err, convey.ShouldEqual, io.EOF)
				c.So(n, convey.ShouldEqual, 0)
			})
		})
	})
}
