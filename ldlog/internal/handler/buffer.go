/*
 * Copyright (C) distroy
 */

package handler

import (
	"time"

	"github.com/distroy/ldgo/v3/ldsync"
)

// Buffer is a byte buffer.
//
// This implementation is adapted from the unexported type buffer
// in go/src/fmt/print.go.
type Buffer []byte

// Having an initial size gives a dramatic speedup.
var bufPool = ldsync.GetPool(func() []byte { return make([]byte, 0, 1024) })

func newBuffer() *Buffer {
	buf := bufPool.Get()
	return (*Buffer)(&buf)
}

func (b *Buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(*(*[]byte)(b))
	}
}

func (b *Buffer) Reset() {
	b.SetLen(0)
}

func (b *Buffer) WriteTime(t time.Time, layout string) (int, error) {
	if layout == "" {
		layout = "2006-01-02T15:04:05.000Z0700"
	}
	l0 := len(*b)
	*b = t.AppendFormat(*b, layout)
	l1 := len(*b)
	return l1 - l0, nil
}

func (b *Buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *Buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

func (b *Buffer) String() string {
	return string(*b)
}

func (b *Buffer) Len() int {
	return len(*b)
}

func (b *Buffer) SetLen(n int) {
	*b = (*b)[:n]
}
