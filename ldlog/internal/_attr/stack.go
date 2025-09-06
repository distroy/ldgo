/*
 * Copyright (C) distroy
 */

package _attr

import (
	"fmt"
	"io"
	"runtime"
)

var (
	_ Marshaler = (*callers_t)(nil)
)

type callers_t struct {
	data []uintptr
}

func (c callers_t) MarshalJSON() ([]byte, error)       { return s2b(c.String()), nil }
func (c callers_t) MarshalText() ([]byte, error)       { return s2b(c.String()), nil }
func (c callers_t) WriteTo(w io.Writer) (int64, error) { return writeTo(w, c) }
func (c callers_t) WriteToBuffer(b *Buffer) {
	if len(c.data) == 0 {
		b.AppendString(`"(no stack)"`)
		return
	}

	pcs := c.data
	if n := len(pcs); pcs[n-1] == 0 {
		pcs = pcs[:n-1]
	}
	frames := runtime.CallersFrames(pcs)

	b.WriteByte('[')
	for i := 0; ; i++ {
		if i != 0 {
			b.WriteByte(',')
		}
		frame, more := frames.Next()
		b.AppendString(fmt.Sprintf(`"called from %s (%s:%d)"`, frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	if len(pcs) < len(c.data) {
		b.WriteByte(',')
		b.AppendString(`"(rest of stack elided)"`)
	}
	b.WriteByte(']')
}
func (c callers_t) String() string {
	b := getBuf()
	defer b.Free()
	c.WriteToBuffer(b)
	// log.Printf("%s", b.Bytes())
	return b.String()
}

func stack(skip, nFrames int) callers_t {
	pcs := make([]uintptr, nFrames+1)
	n := runtime.Callers(skip+2, pcs)
	if n == 0 {
		return callers_t{}
	}
	pcs = pcs[:n]
	if n > nFrames {
		pcs[nFrames] = 0
	}
	return callers_t{
		data: pcs[:n],
	}
}
