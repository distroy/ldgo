/*
 * Copyright (C) distroy
 */

package buffer

import (
	"io"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestRing(t *testing.T) {
	// io.Pipe()

	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("single", func(c convey.C) {
			b := MakeRing(make([]byte, 3))
			c.Convey("new", func(c convey.C) {
				c.So(b, convey.ShouldNotBeNil)
				c.So(b.Size(), convey.ShouldEqual, 0)
			})

			d, ok, err := b.Pop()
			c.Convey("pop empty", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, false)
				c.So(d, convey.ShouldEqual, 0)
				c.So(b.Size(), convey.ShouldEqual, 0)
			})

			ok, err = b.Put('1')
			c.Convey("put '1`'", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, true)
				c.So(b.Size(), convey.ShouldEqual, 1)
			})

			ok, err = b.Put('2')
			c.Convey("put '2'", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, true)
				c.So(b.Size(), convey.ShouldEqual, 2)
			})

			ok, err = b.Put('3')
			c.Convey("put '3'", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, true)
				c.So(b.Size(), convey.ShouldEqual, 3)
			})

			ok, err = b.Put('4')
			c.Convey("put '4' when full", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, false)
				c.So(b.Size(), convey.ShouldEqual, 3)
			})

			d, ok, err = b.Pop()
			c.Convey("pop '1'", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, true)
				c.So(d, convey.ShouldEqual, '1')
				c.So(b.Size(), convey.ShouldEqual, 2)
			})

			d, ok, err = b.Pop()
			c.Convey("pop '2'", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, true)
				c.So(d, convey.ShouldEqual, '2')
				c.So(b.Size(), convey.ShouldEqual, 1)
			})

			ok, err = b.Put('5')
			c.Convey("put '5'", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, true)
				c.So(b.Size(), convey.ShouldEqual, 2)
			})

			err = b.Close()
			c.Convey("close", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
			})

			ok, err = b.Put('5')
			c.Convey("put after close", func(c convey.C) {
				c.So(err, convey.ShouldEqual, io.ErrUnexpectedEOF)
				c.So(ok, convey.ShouldEqual, false)
				c.So(b.Size(), convey.ShouldEqual, 2)
			})

			d, ok, err = b.Pop()
			c.Convey("pop '3' when closed", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, true)
				c.So(d, convey.ShouldEqual, '3')
				c.So(b.Size(), convey.ShouldEqual, 1)
			})

			d, ok, err = b.Pop()
			c.Convey("pop '5' when closed", func(c convey.C) {
				c.So(err, convey.ShouldBeNil)
				c.So(ok, convey.ShouldEqual, true)
				c.So(d, convey.ShouldEqual, '5')
				c.So(b.Size(), convey.ShouldEqual, 0)
			})

			d, ok, err = b.Pop()
			c.Convey("pop when closed & empty", func(c convey.C) {
				c.So(err, convey.ShouldEqual, io.EOF)
				c.So(ok, convey.ShouldEqual, false)
				c.So(d, convey.ShouldEqual, 0)
				c.So(b.Size(), convey.ShouldEqual, 0)
			})
		})

		c.Convey("multi", func(c convey.C) {
			b := MakeRing(make([]byte, 16))
			c.Convey("new", func(c convey.C) {
				c.So(b, convey.ShouldNotBeNil)
				c.So(b.Size(), convey.ShouldEqual, 0)
			})

			d := make([]byte, 32)
			n, err := b.Read(d)
			c.Convey("read empty", func(c convey.C) {
				c.So(n, convey.ShouldEqual, 0)
				c.So(err, convey.ShouldBeNil)
				c.So(b.Size(), convey.ShouldEqual, 0)
			})

			c.Convey("from zero", func(c convey.C) {
				n, err = b.Write([]byte("0123456789"))
				c.Convey("1st write 10", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 10)
					c.So(err, convey.ShouldBeNil)
					c.So(b.Size(), convey.ShouldEqual, 10)
				})

				n, err = b.Write([]byte("0123456789"))
				c.Convey("2nd write 10 -- no enough buffer", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 6)
					c.So(err, convey.ShouldBeNil)
					c.So(b.Size(), convey.ShouldEqual, 16)
				})

				n, err = b.Read(d)
				c.Convey("read all", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 16)
					c.So(err, convey.ShouldBeNil)
					c.So(string(d[:n]), convey.ShouldEqual, "0123456789012345")
					c.So(b.Size(), convey.ShouldEqual, 0)
				})
			})

			c.Convey("from middle", func(c convey.C) {
				n, err = b.Write([]byte("0123456789"))
				c.Convey("write for prepare", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 10)
					c.So(err, convey.ShouldBeNil)
					c.So(b.Size(), convey.ShouldEqual, 10)
				})

				n, err = b.Read(d)
				c.Convey("read for prepare", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 10)
					c.So(err, convey.ShouldBeNil)
					c.So(string(d[:n]), convey.ShouldEqual, "0123456789")
					c.So(b.Size(), convey.ShouldEqual, 0)
				})

				n, err = b.Write([]byte("0123456789"))
				c.Convey("1st write 10", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 10)
					c.So(err, convey.ShouldBeNil)
					c.So(b.Size(), convey.ShouldEqual, 10)
				})

				n, err = b.Write([]byte("0123456789"))
				c.Convey("2nd write 10 -- no enough buffer", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 6)
					c.So(err, convey.ShouldBeNil)
					c.So(b.Size(), convey.ShouldEqual, 16)
				})

				n, err = b.Write([]byte("0123456789"))
				c.Convey("write after full", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 0)
					c.So(err, convey.ShouldBeNil)
					c.So(b.Size(), convey.ShouldEqual, 16)
				})

				n, err = b.Read(d)
				c.Convey("read all buffer", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 16)
					c.So(err, convey.ShouldBeNil)
					c.So(string(d[:n]), convey.ShouldEqual, "0123456789012345")
				})

				n, err = b.Write([]byte("0123456789"))
				c.Convey("3rd write 10", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 10)
					c.So(err, convey.ShouldBeNil)
					c.So(b.Size(), convey.ShouldEqual, 10)
				})

				err = b.Close()
				c.Convey("close", func(c convey.C) {
					c.So(err, convey.ShouldBeNil)
					c.So(b.Size(), convey.ShouldEqual, 10)
				})

				n, err = b.Read(d[:5])
				c.Convey("read part after close", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 5)
					c.So(err, convey.ShouldBeNil)
					c.So(string(d[:n]), convey.ShouldEqual, "01234")
				})

				n, err = b.Read(d)
				c.Convey("read all after close", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 5)
					c.So(err, convey.ShouldBeNil)
					c.So(string(d[:n]), convey.ShouldEqual, "56789")
				})

				// empty & read after close
				n, err = b.Read(d)
				c.Convey("read after close & empty", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 0)
					c.So(err, convey.ShouldEqual, io.EOF)
				})

				// write after close
				n, err = b.Write([]byte("0123456789"))
				c.Convey("write after close & empty", func(c convey.C) {
					c.So(n, convey.ShouldEqual, 0)
					c.So(err, convey.ShouldEqual, io.ErrUnexpectedEOF)
				})
			})
		})
	})
}
