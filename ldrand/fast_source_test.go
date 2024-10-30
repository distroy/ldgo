/*
 * Copyright (C) distroy
 */

package ldrand

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/distroy/ldgo/v3/ldmath"
	"github.com/smartystreets/goconvey/convey"
)

/*
 * pkg: github.com/distroy/ldgo/v2/ldrand
 * cpu: Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz
 * BenchmarkRandGo
 * BenchmarkRandGo-12      18746797                63.71 ns/op
 * BenchmarkRand
 * BenchmarkRand-12        68977040                17.42 ns/op
 */
func BenchmarkRandGo(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Intn(100)
		}
	})
}

func BenchmarkRand(b *testing.B) {
	r := New(NewFastSource(rand.Int63()))
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Intn(100)
		}
	})
}

type testFastSource struct {
	Mod, Scale, Diff int
}

func maxInt(a []int) int { return ldmath.Max(a[0], a[1:]...) }
func minInt(a []int) int { return ldmath.Min(a[0], a[1:]...) }

func diffRatio(a []int) float64 {
	sum := ldmath.SumInt(a...)
	cnt := int64(len(a))
	avg := (sum + cnt/2) / cnt

	diff := int64(0)
	for _, n := range a {
		diff += ldmath.Abs(avg - int64(n))
	}

	return float64(diff) / float64(sum)
}

func (t *testFastSource) Test() {
	var (
		mod   = t.Mod
		scale = t.Scale
		diff  = t.Diff
	)
	name := fmt.Sprintf("mod=%d,scale=%d,diff=%d", mod, scale, diff)
	convey.Convey(name, func() {
		r := New(NewFastSource(time.Now().UnixNano()))

		counts := make([]int, mod)
		for i := 0; i < mod*scale; i++ {
			// x := r.Int() % mod
			x := r.Intn(mod)
			counts[x]++
		}

		min := minInt(counts)
		max := maxInt(counts)
		ratio := diffRatio(counts)

		log.Printf("diff:%d, ratio:%.04g, min:%d, max:%d", max-min, ratio, min, max)
		convey.So(max-min, convey.ShouldBeLessThan, diff)
	})
}

func Test_fastSource_ProbabilityOfOverall(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		(&testFastSource{
			Mod:   100,
			Scale: 1000 * 200,
			Diff:  1000 * 4,
		}).Test()
		(&testFastSource{
			Mod:   16,
			Scale: 1000 * 200,
			Diff:  1000 * 4,
		}).Test()
		(&testFastSource{
			Mod:   256,
			Scale: 1000 * 200,
			Diff:  1000 * 4,
		}).Test()
	})
}

func Test_fastSource_ProbabilityOfVery4Bits(t *testing.T) {
	r := New(NewFastSource(time.Now().UnixNano()))

	convey.Convey(t.Name(), t, func() {
		convey.Convey("check the probability of very 4 bits", func() {
			const (
				scale = 1000 * 100
				diff  = 1000 * 4
			)

			countsPer4Bits := [16][16]int{}
			for i := 0; i < scale*16; i++ {
				v := r.Uint64()
				for i := range countsPer4Bits {
					countsPer4Bits[i][v&0xf]++
					v = v >> 4
				}
			}

			log.Printf("")
			for i, v := range countsPer4Bits {
				min := minInt(v[:])
				max := maxInt(v[:])
				ratio := diffRatio(v[:])
				log.Printf("postion:%d, diff:%d, ratio:%.04g, min:%d, max:%d",
					i, max-min, ratio, min, max)
				convey.So(max-min, convey.ShouldBeLessThan, diff)
			}
		})
	})
}
func testGetDiffThresholdBySliceFunc(diff1, diff2 int) func(idx int, slice interface{}) int {
	return func(idx int, slice interface{}) int {
		sliceV := reflect.ValueOf(slice)
		if idx == sliceV.Len()-1 {
			return diff2
		}
		return diff1
	}
}

func Test_fastSource_ProbabilityOfVery4BitsWithPreviousNumber(t *testing.T) {
	r := New(NewFastSource(time.Now().UnixNano()))
	const (
		scale = 1000 * 100 * 16 * 16
		diff1 = 5000
		diff2 = 6500
	)

	getDiffThreshold := testGetDiffThresholdBySliceFunc(diff1, diff2)

	convey.Convey(t.Name(), t, func() {
		convey.Convey("check the probability of very 4 bits with previous number", func() {
			countsPer4BitsWithPrev := [16][16][16]int{}
			var prevNum uint64
			for i := 0; i < scale; i++ {
				v := r.Uint64()
				p := prevNum
				prevNum = v
				for i := range countsPer4BitsWithPrev {
					countsPer4BitsWithPrev[i][p&0xf][v&0xf]++
					v = v >> 4
					p = p >> 4
				}
			}

			log.Printf("")
			for i, v := range countsPer4BitsWithPrev {
				for j, w := range v {
					min := minInt(w[:])
					max := maxInt(w[:])
					ratio := diffRatio(w[:])
					log.Printf("postion:%d, prev:%d, diff:%d, ratio:%.04g, min:%d, max:%d, v:%v",
						i, j, max-min, ratio, min, max, w[:])

					diff := getDiffThreshold(i, countsPer4BitsWithPrev[:])
					// sometimes it will be failed
					convey.So(max-min, convey.ShouldBeLessThan, diff)
				}
			}
		})
	})
}

func Test_fastSource_ProbabilityOfVeryByte(t *testing.T) {
	r := New(NewFastSource(time.Now().UnixNano()))

	convey.Convey(t.Name(), t, func() {
		convey.Convey("check the probability of very byte", func() {
			const (
				scale = 1000 * 100
				diff  = 1000 * 5
			)

			countsPer4Bits := [8][256]int{}
			for i := 0; i < scale*256; i++ {
				v := r.Uint64()
				for i := range countsPer4Bits {
					countsPer4Bits[i][v&0xff]++
					v = v >> 8
				}
			}

			log.Printf("")
			for i, v := range countsPer4Bits {
				min := minInt(v[:])
				max := maxInt(v[:])
				ratio := diffRatio(v[:])
				log.Printf("postion:%d, diff:%d, ratio:%.04g, min:%d, max:%d",
					i, max-min, ratio, min, max)
				convey.So(max-min, convey.ShouldBeLessThan, diff)
			}
		})
	})
}

func Test_fastSource_Repeated(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		r := New(NewFastSource(time.Now().UnixNano()))

		convey.Convey("check result if repeated", func() {
			const times = 100 * 10000

			m := make(map[uint64]struct{}, times)
			for i := 0; i < times; i++ {
				x := r.Uint64()
				if _, ok := m[x]; ok {
					t.Fatalf("number repeated. %d", x)
				}
				m[x] = struct{}{}
			}
		})
	})
}

func testGreatestCommonDivisor(a, b uint64) uint64 {
	for b != 0 {
		a %= b
		a, b = b, a
	}

	return a
}

func Test_testGreatestCommonDivisor(t *testing.T) {
	type Case struct {
		A, B, R uint64
	}
	cases := []*Case{
		{A: 0, B: 2, R: 2},
		{A: 2, B: 0, R: 2},
		{A: 2, B: 100, R: 2},
		{A: 2, B: 3, R: 1},
		{A: 10, B: 6, R: 2},
		{A: 6, B: 10, R: 2},
	}
	convey.Convey(t.Name(), t, func(c convey.C) {
		for _, v := range cases {
			name := fmt.Sprintf("a:%d, b:%d", v.A, v.B)
			c.Convey(name, func(c convey.C) {
				r := testGreatestCommonDivisor(v.A, v.B)
				c.So(r, convey.ShouldEqual, v.R)
			})
		}
	})
}

func Test__fastSourceStep(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		r := testGreatestCommonDivisor(math.MaxUint64, fastSourceStep)
		c.So(r, convey.ShouldEqual, 1)
	})
}
