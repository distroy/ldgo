/*
 * Copyright (C) distroy
 */

package _attr

import (
	"fmt"
	"runtime"
)

func stack(skip, nFrames int) string {
	pcs := make([]uintptr, nFrames+1)
	n := runtime.Callers(skip+1, pcs)
	if n == 0 {
		return "(no stack)"
	}
	frames := runtime.CallersFrames(pcs[:n])
	// var b strings.Builder
	b := getBuf()
	defer b.Free()
	i := 0
	for {
		frame, more := frames.Next()
		fmt.Fprintf(b, "called from %s (%s:%d)\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
		i++
		if i >= nFrames {
			fmt.Fprintf(b, "(rest of stack elided)\n")
			break
		}
	}
	return b.String()
}
