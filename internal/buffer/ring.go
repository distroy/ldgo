/*
 * Copyright (C) distroy
 */

package buffer

import (
	"io"

	"github.com/distroy/ldgo/v2/ldmath"
)

func MakeRing[T any](buf []T) Ring[T] {
	return Ring[T]{
		buf: buf,
	}
}

type Ring[T any] struct {
	buf    []T
	begin  int
	end    int
	closed bool
	full   bool
}

func (b *Ring[T]) Cap() int { return len(b.buf) }
func (b *Ring[T]) Size() int {
	capacity := len(b.buf)
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

func (b *Ring[T]) Closed() bool { return b.closed }
func (b *Ring[T]) Close() error {
	b.closed = true
	return nil
}

func (b *Ring[T]) addEndPos(n int) {
	capacity := len(b.buf)
	b.end += n
	if b.end >= capacity {
		b.end -= capacity
	}
	if b.end == b.begin {
		b.full = true
	}
}

func (b *Ring[T]) addBeginPos(n int) {
	capacity := len(b.buf)
	b.begin += n
	if b.begin >= capacity {
		b.begin -= capacity
	}
	b.full = false
}

func (b *Ring[T]) Write(d []T) (int, error) {
	if b.closed {
		return 0, io.ErrUnexpectedEOF
	}
	if len(d) == 0 || b.full {
		return 0, nil
	}

	capacity := len(b.buf)

	pos := b.end
	if b.end < b.begin {
		n := b.begin - pos
		n = ldmath.Min(n, len(d))
		copy(b.buf[pos:], d[:n])
		b.addEndPos(n)
		return n, nil
	}

	n := capacity - pos
	n = ldmath.Min(n, len(d))
	copy(b.buf[pos:], d[:n])
	if n >= len(d) || b.begin == 0 {
		b.addEndPos(n)
		return n, nil
	}

	n1 := ldmath.Min(b.begin, len(d)-n)
	copy(b.buf[:n1], d[n:n+n1])
	b.addEndPos(n + n1)
	return n + n1, nil
}

func (b *Ring[T]) Read(d []T) (int, error) {
	if b.end == b.begin && !b.full {
		if b.closed {
			return 0, io.EOF
		}
		return 0, nil
	}
	if len(d) == 0 {
		return 0, nil
	}

	capacity := len(b.buf)
	pos := b.begin
	n := len(d)
	if b.begin < b.end {
		n = ldmath.Min(n, b.end-b.begin)
		copy(d, b.buf[pos:pos+n])
		b.addBeginPos(n)
		return n, nil
	}

	n = ldmath.Min(n, capacity-pos)
	copy(d, b.buf[pos:pos+n])
	if n >= len(d) || b.end == 0 {
		b.addBeginPos(n)
		return n, nil
	}

	n1 := len(d) - n
	n1 = ldmath.Min(n1, b.end)
	copy(d[n:], b.buf[:n1])
	b.addBeginPos(n + n1)
	return n + n1, nil
}

func (b *Ring[T]) Put(d T) (bool, error) {
	if b.closed {
		return false, io.ErrUnexpectedEOF
	}
	if b.full {
		return false, nil
	}
	b.buf[b.end] = d
	b.addEndPos(1)
	return true, nil
}

func (b *Ring[T]) Pop() (T, bool, error) {
	if b.begin == b.end && !b.full {
		var d T
		if b.closed {
			return d, false, io.EOF
		}
		return d, false, nil
	}

	d := b.buf[b.begin]
	b.addBeginPos(1)
	return d, true, nil
}
