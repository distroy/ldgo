/*
 * Copyright (C) distroy
 */

package ldref

import (
	"encoding/json"
	"testing"

	"github.com/distroy/ldgo/v2/ldref/internal/copybenchstruct1"
	"github.com/distroy/ldgo/v2/ldref/internal/copybenchstruct2"
)

/*
goos: darwin
goarch: amd64
pkg: github.com/distroy/ldgo/v2/ldref
cpu: VirtualApple @ 2.50GHz
Benchmark_copyV1
Benchmark_copyV1-10                20092             59424 ns/op
Benchmark_copyV2
Benchmark_copyV2-10                42955             28423 ns/op
Benchmark_deepCopyV1
Benchmark_deepCopyV1-10            12066             94852 ns/op
Benchmark_deepCopyV2
Benchmark_deepCopyV2-10            23590             46978 ns/op
Benchmark_jsonCopy
Benchmark_jsonCopy-10               5764            202620 ns/op
PASS
ok      github.com/distroy/ldgo/v2/ldref        19.954s
*/

func benchPrepareObjects(n int) []*copybenchstruct1.ItemCardData {
	obj := &copybenchstruct1.ItemCardData{}
	json.Unmarshal(copybenchstruct1.JSON_STING, obj)
	res := make([]*copybenchstruct1.ItemCardData, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, DeepClone(obj))
	}
	return res
}

func benchCopyFunc(b *testing.B, copyFunc func(target, source interface{}, cfg ...*CopyConfig) error) {
	size := 1024
	srcs := benchPrepareObjects(size)
	{
		var (
			target = &copybenchstruct2.ItemCardData{}
			source = srcs[0]
		)
		copyFunc(target, source)
	}

	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		count := 0
		for p.Next() {
			var (
				index  = count
				target = &copybenchstruct2.ItemCardData{}
				source = srcs[index&(size-1)]
			)
			count++
			copyFunc(target, source)
		}
	})
	b.StopTimer()
}

func Benchmark_copyV1(b *testing.B) { benchCopyFunc(b, copyV1) }
func Benchmark_copyV2(b *testing.B) { benchCopyFunc(b, copyV2) }

func Benchmark_deepCopyV1(b *testing.B) { benchCopyFunc(b, deepCopyV1) }
func Benchmark_deepCopyV2(b *testing.B) { benchCopyFunc(b, deepCopyV2) }

func Benchmark_jsonCopy(b *testing.B) {
	benchCopyFunc(b, func(target, source interface{}, cfg ...*CopyConfig) error {
		raw, err := json.Marshal(source)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(raw, target); err != nil {
			return err
		}
		return nil
	})
}
