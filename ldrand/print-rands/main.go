/*
 * Copyright (C) distroy
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/distroy/ldgo/v3/ldrand"
)

func printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
	fmt.Fprintf(os.Stdout, "\n")
}

func main() {
	seed := flag.Int64("seed", time.Now().UnixNano(), "")
	count := flag.Int("count", 10, "")
	mod := flag.Int("mod", 0, "")
	flag.Parse()

	r := ldrand.New(ldrand.NewFastSource(*seed))

	switch *mod {
	case 0:
		printf("uint64:")
		for i := 0; i < *count; i++ {
			printf("%#0.16x", r.Uint64())
		}
	case -1:
		printf("uint32:")
		for i := 0; i < *count; i++ {
			printf("%#0.8x", r.Uint32())
		}

	default:
		printf("mod %d:", *mod)
		for i := 0; i < *count; i++ {
			printf("%d", r.Intn(*mod))
		}
	}
}
