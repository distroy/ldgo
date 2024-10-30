/*
 * Copyright (C) distroy
 */

package ldbyte

import (
	"io"

	"github.com/distroy/ldgo/v3/ldmath"
)

var (
	_ io.ReadWriteCloser = (*RingBuffer)(nil)
)

func NewRingBuffer(n int) *RingBuffer {
	if n <= 0 {
		n = 1024
	}
	return &RingBuffer{
		data: make([]byte, n, n),
	}
}

type RingBuffer struct {
	data   []byte
	begin  int
	end    int
	closed bool
	full   bool
}

func (b *RingBuffer) Close() error {
	b.closed = true
	return nil
}

func (b *RingBuffer) Size() int {
	capacity := len(b.data)
	if b.begin == b.end {
		if b.full {
			return capacity
		}
		return 0
	}
	if b.end > b.begin {
		return b.end - b.begin
	}
	return b.end + capacity - b.begin
}

func (b *RingBuffer) Write(d []byte) (int, error) {
	if b.closed || len(d) == 0 || b.full {
		return 0, nil
	}

	capacity := len(b.data)
	addEndPos := func(n int) {
		b.end += n
		if b.end >= capacity {
			b.end -= capacity
		}
		if b.end == b.begin {
			b.full = true
		}
	}

	pos := b.end
	if b.end < b.begin {
		n := b.begin - pos
		n = ldmath.Min(n, len(d))
		copy(b.data[pos:], d[:n])
		addEndPos(n)
		return n, nil
	}

	n := capacity - pos
	n = ldmath.Min(n, len(d))
	copy(b.data[pos:], d[:n])
	if n >= len(d) || b.begin == 0 {
		addEndPos(n)
		return n, nil
	}

	n1 := ldmath.Min(b.begin, len(d)-n)
	copy(b.data[:n1], d[n:n+n1])
	addEndPos(n + n1)
	return n + n1, nil
}

func (b *RingBuffer) Read(d []byte) (int, error) {
	if b.end == b.begin && !b.full {
		if b.closed {
			return 0, io.EOF
		}
		return 0, nil
	}
	if len(d) == 0 {
		return 0, nil
	}

	capacity := len(b.data)
	addBeginPos := func(n int) {
		b.begin += n
		if b.begin >= capacity {
			b.begin -= capacity
		}
		b.full = false
	}

	pos := b.begin
	n := len(d)
	if b.begin < b.end {
		n = ldmath.Min(n, b.end-b.begin)
		copy(d, b.data[pos:pos+n])
		addBeginPos(n)
		return n, nil
	}

	n = ldmath.Min(n, capacity-pos)
	copy(d, b.data[pos:pos+n])
	if n >= len(d) || b.end == 0 {
		addBeginPos(n)
		return n, nil
	}

	n1 := len(d) - n
	n1 = ldmath.Min(n1, b.end)
	copy(d[n:], b.data[:n1])
	addBeginPos(n + n1)
	return n + n1, nil
}
