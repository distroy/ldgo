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
	xor := [16]uint64{}

	for i := range buf {
		buf[i] = byte(i)
	}

	for i := 0; i < 16; i++ {
		rand.Shuffle(len(buf), func(i, j int) { buf[i], buf[j] = buf[j], buf[i] })
		for j := range xor {
			xor[j] = (xor[j] << 4) | uint64(buf[j])
		}
	}

	for _, v := range xor {
		fmt.Fprintf(os.Stdout, "%#0.16x\n", v)
	}
}
