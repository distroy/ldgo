/*
 * Copyright (C) distroy
 */

package buffer

import (
	"slices"
	"strconv"
	"time"

	"github.com/distroy/ldgo/v3/ldsync"
)

const smallBufferSize = 64

// Buffer is a byte buffer.
//
// This implementation is adapted from the unexported type buffer
// in go/src/fmt/print.go.
type Buffer []byte

// Having an initial size gives a dramatic speedup.
var bufPool = ldsync.GetPool(func() []byte { return make([]byte, 0, 1024) })

func NewBuffer() *Buffer {
	buf := bufPool.Get()
	buf = buf[:0]
	return (*Buffer)(&buf)
}

func (b *Buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		bufPool.Put(*(*[]byte)(b))
	}
}

func (b *Buffer) Grow(n int) {
	if n <= 0 {
		return
	}
	*b = slices.Grow(*b, n)
}

// TrimNewline trims any final "\n" byte from the end of the buffer.
func (b *Buffer) TrimNewline() {
	// *b = bytes.TrimSuffix(*b, []byte{'\n'})
	if i := len(*b) - 1; i >= 0 {
		if (*b)[i] == '\n' {
			*b = (*b)[:i]
		}
	}
}

func (b *Buffer) Reset() { b.SetLen(0) }

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

func (b *Buffer) Bytes() []byte  { return *b }
func (b *Buffer) String() string { return string(*b) }

func (b *Buffer) Len() int     { return len(*b) }
func (b *Buffer) Cap() int     { return cap(*b) }
func (b *Buffer) SetLen(n int) { *b = (*b)[:n] }

func (b *Buffer) AppendBool(v bool)     { *b = strconv.AppendBool(*b, v) }
func (b *Buffer) AppendByte(v byte)     { *b = append(*b, v) }
func (b *Buffer) AppendBytes(v []byte)  { *b = append(*b, v...) }
func (b *Buffer) AppendString(v string) { *b = append(*b, v...) }

func (b *Buffer) AppendInt(v int64)   { *b = strconv.AppendInt(*b, v, 10) }
func (b *Buffer) AppendUint(v uint64) { *b = strconv.AppendUint(*b, v, 10) }
func (b *Buffer) AppendFloat(v float64, bitSize int) {
	*b = strconv.AppendFloat(*b, v, 'f', -1, bitSize)
}
func (b *Buffer) AppendComplex(v complex128, bitSize int) {
	s := strconv.FormatComplex(v, 'f', -1, bitSize)
	*b = append(*b, s...)
}

func (b *Buffer) AppendTime(t time.Time, layout string) {
	if layout == "" {
		b.AppendString(strconv.Quote(t.String()))
		return
	}
	b.AppendString(strconv.Quote(t.Format(layout)))
}
