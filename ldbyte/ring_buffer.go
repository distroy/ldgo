/*
 * Copyright (C) distroy
 */

package ldbyte

import (
	"io"

	"github.com/distroy/ldgo/v2/internal/buffer"
)

var (
	_ io.ReadWriteCloser = (*RingBuffer)(nil)
)

func NewRingBuffer(n int) *RingBuffer {
	if n <= 0 {
		n = 1024
	}
	buf := make([]byte, n, n)
	return &RingBuffer{
		buf: buffer.MakeRing(buf),
	}
}

type RingBuffer struct {
	buf buffer.Ring[byte]
}

func (b *RingBuffer) Close() error { return b.buf.Close() }
func (b *RingBuffer) Closed() bool { return b.buf.Closed() }

func (b *RingBuffer) Cap() int  { return b.buf.Cap() }
func (b *RingBuffer) Size() int { return b.buf.Size() }

func (b *RingBuffer) Write(d []byte) (int, error) { return b.buf.Write(d) }
func (b *RingBuffer) Read(d []byte) (int, error)  { return b.buf.Read(d) }

func NewBlockingRingBuffer(n int) *BlockingRingBuffer {
	if n <= 0 {
		n = 1024
	}
	buf := make([]byte, n, n)
	b := &BlockingRingBuffer{
		buf: buffer.BlockingRing[byte]{
			Buf: buffer.MakeRing(buf),
		},
	}
	return b
}

type BlockingRingBuffer struct {
	buf buffer.BlockingRing[byte]
}

func (b *BlockingRingBuffer) Close() error { return b.buf.Close() }
func (b *BlockingRingBuffer) Closed() bool { return b.buf.Closed() }

func (b *BlockingRingBuffer) Cap() int  { return b.buf.Cap() }
func (b *BlockingRingBuffer) Size() int { return b.buf.Size() }

func (b *BlockingRingBuffer) Write(d []byte) (int, error) { return b.buf.Write(d) }
func (b *BlockingRingBuffer) Read(d []byte) (int, error)  { return b.buf.Read(d) }
