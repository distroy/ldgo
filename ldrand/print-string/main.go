/*
 * Copyright (C) distroy
 */

package main

import (
	"encoding/hex"
	"fmt"

	"github.com/distroy/ldgo/v3/ldflag"
	"github.com/distroy/ldgo/v3/ldmath"
	"github.com/distroy/ldgo/v3/ldrand"
)

type Flags struct {
	Count int `ldflag:"default:10"`
	Size  int `ldflag:"default:16"`
}

func main() {
	flags := &Flags{}
	ldflag.MustParse(flags)

	flags.Count = ldmath.MaxInt(flags.Count, 1)
	flags.Size = ldmath.MaxInt(flags.Size, 1)

	buf := make([]byte, flags.Size)

	for i := 0; i < flags.Count; i++ {
		ldrand.Read(buf)
		fmt.Println(hex.EncodeToString(buf))
	}
}
