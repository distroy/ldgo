/*
 * Copyright (C) distroy
 */

package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	buf := [16]byte{}
	postion := [16]uint64{}

	for i := range buf {
		buf[i] = byte(i)
	}

	for i := 0; i < 16; i++ {
		rand.Shuffle(len(buf), func(i, j int) { buf[i], buf[j] = buf[j], buf[i] })
		fmt.Fprintf(os.Stdout, "%#v\n", buf)

		var n uint64
		for _, v := range buf {
			n = (n << 4) | uint64(v)
		}
		postion[i] = n
	}

	for _, v := range postion {
		fmt.Fprintf(os.Stdout, "%#0.16x\n", v)
	}
}
