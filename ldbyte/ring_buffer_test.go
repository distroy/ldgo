/*
 * Copyright (C) distroy
 */

package ldbyte

import (
	"io"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestNewRingBuffer(t *testing.T) {
	io.Pipe()

	convey.Convey(t.Name(), t, func(c convey.C) {
		b := NewRingBuffer(16)
		c.Convey("new", func(c convey.C) {
			c.So(b, convey.ShouldNotBeNil)
			c.So(b.Size(), convey.ShouldEqual, 0)
		})

		// read empty
		d := make([]byte, 32)
		n, err := b.Read(d)
		c.Convey("read empty", func(c convey.C) {
			c.So(n, convey.ShouldEqual, 0)
			c.So(err, convey.ShouldBeNil)
			c.So(b.Size(), convey.ShouldEqual, 0)
		})

		c.Convey("from zero", func(c convey.C) {
			n, err = b.Write([]byte("0123456789"))
			c.Convey("first write 10", func(c convey.C) {
				c.So(n, convey.ShouldEqual, 10)
				c.So(err, convey.ShouldBeNil)
				c.So(b.Size(), convey.ShouldEqual, 10)
			})

			n, err = b.Write([]byte("0123456789"))
			c.Convey("second write 10 -- no enough buffer", func(c convey.C) {
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
			c.Convey("first write 10", func(c convey.C) {
				c.So(n, convey.ShouldEqual, 10)
				c.So(err, convey.ShouldBeNil)
				c.So(b.Size(), convey.ShouldEqual, 10)
			})

			n, err = b.Write([]byte("0123456789"))
			c.Convey("second write 10 -- no enough buffer", func(c convey.C) {
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
			c.Convey("third write 10", func(c convey.C) {
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
}
