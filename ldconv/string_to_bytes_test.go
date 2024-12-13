/*
 * Copyright (C) distroy
 */

package ldconv

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

/*
goos: linux
goarch: amd64
pkg: github.com/distroy/ldgo/v2/ldconv
cpu: 12th Gen Intel(R) Core(TM) i7-12650H
Benchmark_BytesToStr
Benchmark_BytesToStr/len-10
Benchmark_BytesToStr/len-10-16          187010206                6.479 ns/op
Benchmark_BytesToStr/len-100
Benchmark_BytesToStr/len-100-16         46684989                26.70 ns/op
Benchmark_BytesToStr/len-500
Benchmark_BytesToStr/len-500-16         12696064                89.12 ns/op
Benchmark_BytesToStr/len-1000
Benchmark_BytesToStr/len-1000-16         7481668               168.9 ns/op
Benchmark_BytesToStrUnsafe
Benchmark_BytesToStrUnsafe/len-10
Benchmark_BytesToStrUnsafe/len-10-16            1000000000               0.2030 ns/op
Benchmark_BytesToStrUnsafe/len-100
Benchmark_BytesToStrUnsafe/len-100-16           1000000000               0.1998 ns/op
Benchmark_BytesToStrUnsafe/len-500
Benchmark_BytesToStrUnsafe/len-500-16           1000000000               0.1994 ns/op
Benchmark_BytesToStrUnsafe/len-1000
Benchmark_BytesToStrUnsafe/len-1000-16          1000000000               0.2012 ns/op
Benchmark_StrToBytes
Benchmark_StrToBytes/len-10
Benchmark_StrToBytes/len-10-16                  167892522                6.870 ns/op
Benchmark_StrToBytes/len-100
Benchmark_StrToBytes/len-100-16                 46696956                25.16 ns/op
Benchmark_StrToBytes/len-500
Benchmark_StrToBytes/len-500-16                 10698717                99.79 ns/op
Benchmark_StrToBytes/len-1000
Benchmark_StrToBytes/len-1000-16                 6156220               188.1 ns/op
Benchmark_StrToBytesUnsafe
Benchmark_StrToBytesUnsafe/len-10
Benchmark_StrToBytesUnsafe/len-10-16            1000000000               0.1957 ns/op
Benchmark_StrToBytesUnsafe/len-100
Benchmark_StrToBytesUnsafe/len-100-16           1000000000               0.1922 ns/op
Benchmark_StrToBytesUnsafe/len-500
Benchmark_StrToBytesUnsafe/len-500-16           1000000000               0.1953 ns/op
Benchmark_StrToBytesUnsafe/len-1000
Benchmark_StrToBytesUnsafe/len-1000-16          1000000000               0.1987 ns/op
*/

var testSizesForConvBetweenStrAndBytes = []int{
	10, 100, 500, 1000,
}

func testRandString(n int) string {
	buf := make([]byte, n)
	rand.Read(buf)
	return string(buf)
}

func Test_BytesToStrUnsafe(t *testing.T) {
	fn := BytesToStrUnsafe

	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("nil", func(c convey.C) {
			c.So(fn(nil), convey.ShouldEqual, "")
		})
		c.Convey("empty", func(c convey.C) {
			c.So(fn([]byte{}), convey.ShouldEqual, "")
		})
		c.Convey("conv and change", func(c convey.C) {
			b := []byte("123456789")
			s := fn(b)
			c.So(s, convey.ShouldEqual, "123456789")
			copy(b, []byte("abcd"))
			c.So(s, convey.ShouldEqual, "abcd56789")
		})
	})
}

func Test_StrToBytesUnsafe(t *testing.T) {
	fn := StrToBytesUnsafe

	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(fn(""), convey.ShouldResemble, []byte(nil))
		c.So(len(fn("abc")), convey.ShouldEqual, 3)
		c.So(cap(fn("abc")), convey.ShouldBeIn, []int{0, 3})
	})
}

func Test_strToBytesUnsafeV1(t *testing.T) {
	fn := strToBytesUnsafeV1

	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(fn(""), convey.ShouldResemble, []byte(nil))
		c.So(len(fn("abc")), convey.ShouldEqual, 3)
		c.SkipSo(cap(fn("abc")), convey.ShouldEqual, 0)
	})
}

func Test_strToBytesUnsafeV2(t *testing.T) {
	fn := strToBytesUnsafeV2

	convey.Convey(t.Name(), t, func(c convey.C) {
		c.So(fn(""), convey.ShouldResemble, []byte(nil))
		c.So(len(fn("abc")), convey.ShouldEqual, 3)
		c.So(cap(fn("abc")), convey.ShouldEqual, 3)
	})
}

func benchmarkBytesToStrWithSize(b *testing.B, fn func(d []byte) string, size int) {
	b.Run(fmt.Sprintf("len-%d", size), func(b *testing.B) {
		str := []byte(testRandString(size))

		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				fn(str)
			}
		})
	})
}

func benchmarkBytesToStr(b *testing.B, fn func(d []byte) string) {
	for _, n := range testSizesForConvBetweenStrAndBytes {
		benchmarkBytesToStrWithSize(b, fn, n)
	}
}

func Benchmark_BytesToStr(b *testing.B) {
	benchmarkBytesToStr(b, BytesToStr)
}
func Benchmark_BytesToStrUnsafe(b *testing.B) {
	benchmarkBytesToStr(b, BytesToStrUnsafe)
}
func Benchmark_bytesToStrUnsafeV1(b *testing.B) {
	benchmarkBytesToStr(b, bytesToStrUnsafeV1)
}

func benchmarkStrToBytesWithSize(b *testing.B, fn func(s string) []byte, size int) {
	b.Run(fmt.Sprintf("len-%d", size), func(b *testing.B) {
		str := testRandString(size)

		b.ResetTimer()
		b.RunParallel(func(p *testing.PB) {
			for p.Next() {
				fn(str)
			}
		})
	})
}
func benchmarkStrToBytes(b *testing.B, fn func(s string) []byte) {
	for _, n := range testSizesForConvBetweenStrAndBytes {
		benchmarkStrToBytesWithSize(b, fn, n)
	}
}

func Benchmark_StrToBytes(b *testing.B) {
	benchmarkStrToBytes(b, StrToBytes)
}
func Benchmark_StrToBytesUnsafe(b *testing.B) {
	benchmarkStrToBytes(b, StrToBytesUnsafe)
}
func Benchmark_strToBytesUnsafeV1(b *testing.B) {
	benchmarkStrToBytes(b, strToBytesUnsafeV1)
}
func Benchmark_strToBytesUnsafeV2(b *testing.B) {
	benchmarkStrToBytes(b, strToBytesUnsafeV2)
}
