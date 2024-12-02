/*
 * Copyright (C) distroy
 */

package cmp

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"strings"
	"testing"

	"github.com/distroy/ldgo/v2/ldconv"
)

/*
goos: darwin
goarch: amd64
pkg: github.com/distroy/ldgo/v2/internal/cmp
cpu: Intel(R) Core(TM) i7-4770HQ CPU @ 2.20GHz
BenchmarkStringsCompare_128
BenchmarkStringsCompare_128/same
BenchmarkStringsCompare_128/same/strings.Compare
BenchmarkStringsCompare_128/same/strings.Compare-8              863649567                1.356 ns/op
BenchmarkStringsCompare_128/same/bytes.Compare
BenchmarkStringsCompare_128/same/bytes.Compare-8                854837110                1.368 ns/op
BenchmarkStringsCompare_128/same/bytes.Compare_with_convert
BenchmarkStringsCompare_128/same/bytes.Compare_with_convert-8           772947661                1.517 ns/op
BenchmarkStringsCompare_128/same/CompareBytes
BenchmarkStringsCompare_128/same/CompareBytes-8                         822200548                1.453 ns/op
BenchmarkStringsCompare_128/same/CompareString
BenchmarkStringsCompare_128/same/CompareString-8                        641744119                1.847 ns/op
BenchmarkStringsCompare_128/different
BenchmarkStringsCompare_128/different/strings.Compare
BenchmarkStringsCompare_128/different/strings.Compare-8                 315854049                3.799 ns/op
BenchmarkStringsCompare_128/different/bytes.Compare
BenchmarkStringsCompare_128/different/bytes.Compare-8                   499989393                2.422 ns/op
BenchmarkStringsCompare_128/different/bytes.Compare_with_convert
BenchmarkStringsCompare_128/different/bytes.Compare_with_convert-8      452922696                2.667 ns/op
BenchmarkStringsCompare_128/different/CompareBytes
BenchmarkStringsCompare_128/different/CompareBytes-8                    455882289                2.587 ns/op
BenchmarkStringsCompare_128/different/CompareString
BenchmarkStringsCompare_128/different/CompareString-8                   393111260                3.024 ns/op
BenchmarkStringsCompare_4096
BenchmarkStringsCompare_4096/same
BenchmarkStringsCompare_4096/same/strings.Compare
BenchmarkStringsCompare_4096/same/strings.Compare-8                     830410903                1.500 ns/op
BenchmarkStringsCompare_4096/same/bytes.Compare
BenchmarkStringsCompare_4096/same/bytes.Compare-8                       761810070                1.598 ns/op
BenchmarkStringsCompare_4096/same/bytes.Compare_with_convert
BenchmarkStringsCompare_4096/same/bytes.Compare_with_convert-8          693995869                1.723 ns/op
BenchmarkStringsCompare_4096/same/CompareBytes
BenchmarkStringsCompare_4096/same/CompareBytes-8                        709518660                1.673 ns/op
BenchmarkStringsCompare_4096/same/CompareString
BenchmarkStringsCompare_4096/same/CompareString-8                       553333569                2.190 ns/op
BenchmarkStringsCompare_4096/different
BenchmarkStringsCompare_4096/different/strings.Compare
BenchmarkStringsCompare_4096/different/strings.Compare-8                126961375                9.576 ns/op
BenchmarkStringsCompare_4096/different/bytes.Compare
BenchmarkStringsCompare_4096/different/bytes.Compare-8                  200258566                6.062 ns/op
BenchmarkStringsCompare_4096/different/bytes.Compare_with_convert
BenchmarkStringsCompare_4096/different/bytes.Compare_with_convert-8     148459071                8.081 ns/op
BenchmarkStringsCompare_4096/different/CompareBytes
BenchmarkStringsCompare_4096/different/CompareBytes-8                   177549223                6.957 ns/op
BenchmarkStringsCompare_4096/different/CompareString
BenchmarkStringsCompare_4096/different/CompareString-8                  139446196                8.880 ns/op
*/

func randString(n int) string {
	nn := (n + 1) / 2
	b := make([]byte, nn)
	rand.Read(b)
	return hex.EncodeToString(b)[:n]
}

var testStringCompareCases map[int][][2]string

func testStrToBytesUnsafe(s string) []byte {
	return ldconv.StrToBytesUnsafe(s)

	// type bytesFromString reflect.SliceHeader
	// sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// bh := &bytesFromString{
	// 	Data: sh.Data,
	// 	Len:  sh.Len,
	// 	Cap:  0,
	// }
	// return *(*[]byte)(unsafe.Pointer(bh))
}

func getGestGetStringCompareCases(strLen int) [][2]string {
	if testStringCompareCases == nil {
		testStringCompareCases = make(map[int][][2]string)
	}
	const size = 1 << 10
	// const strLen = 1024
	if testStringCompareCases[strLen] == nil {
		cases := make([][2]string, 0, size)
		for i := 0; i < size; i++ {
			cases = append(cases, [2]string{
				randString(strLen), randString(strLen),
			})
		}

		testStringCompareCases[strLen] = cases
	}

	return testStringCompareCases[strLen]
}

func benchmarkStrCompare_same[T any](b *testing.B, cases [][2]T, cmp func(a, b T) int) {
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		size := len(cases)
		pos := 0
		for p.Next() {
			v := cases[pos]
			// pos = (pos + 1) % size
			if pos++; pos >= size {
				pos -= size
			}

			s0 := v[0]
			s1 := v[0]
			// strings.Compare(s0, s1)
			cmp(s0, s1)
		}
	})
}

func benchmarkStrCompare_different[T any](b *testing.B, cases [][2]T, cmp func(a, b T) int) {
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		size := len(cases)
		pos := 0
		for p.Next() {
			v := cases[pos]
			// pos = (pos + 1) % size
			if pos++; pos >= size {
				pos -= size
			}

			s0 := v[0]
			s1 := v[1]
			// strings.Compare(s0, s1)
			cmp(s0, s1)
		}
	})
}

func benchmarkStrCompare(b *testing.B, strCases [][2]string) {
	caseCount := len(strCases)
	byteCases := make([][2][]byte, 0, caseCount)
	for _, v := range strCases {
		byteCases = append(byteCases, [2][]byte{
			testStrToBytesUnsafe(v[0]), testStrToBytesUnsafe(v[1]),
		})
	}

	b.Run("same", func(b *testing.B) {
		b.Run("strings.Compare", func(b *testing.B) {
			benchmarkStrCompare_same(b, strCases, strings.Compare)
		})
		b.Run("bytes.Compare", func(b *testing.B) {
			benchmarkStrCompare_same(b, byteCases, bytes.Compare)
		})
		b.Run("bytes.Compare with convert", func(b *testing.B) {
			benchmarkStrCompare_same(b, strCases, func(a, b string) int {
				return bytes.Compare(testStrToBytesUnsafe(a), testStrToBytesUnsafe(b))
			})
		})
		b.Run("CompareBytes", func(b *testing.B) {
			benchmarkStrCompare_same(b, byteCases, CompareBytes[[]byte])
		})
		b.Run("CompareString", func(b *testing.B) {
			benchmarkStrCompare_same(b, strCases, CompareString[string])
		})
	})
	b.Run("different", func(b *testing.B) {
		b.Run("strings.Compare", func(b *testing.B) {
			benchmarkStrCompare_different(b, strCases, strings.Compare)
		})
		b.Run("bytes.Compare", func(b *testing.B) {
			benchmarkStrCompare_different(b, byteCases, bytes.Compare)
		})
		b.Run("bytes.Compare with convert", func(b *testing.B) {
			benchmarkStrCompare_different(b, strCases, func(a, b string) int {
				return bytes.Compare(testStrToBytesUnsafe(a), testStrToBytesUnsafe(b))
			})
		})
		b.Run("CompareBytes", func(b *testing.B) {
			benchmarkStrCompare_different(b, byteCases, CompareBytes[[]byte])
		})
		b.Run("CompareString", func(b *testing.B) {
			benchmarkStrCompare_different(b, strCases, CompareString[string])
		})
	})
}

func BenchmarkStringsCompare_128(b *testing.B) {
	cases := getGestGetStringCompareCases(128)
	benchmarkStrCompare(b, cases)
}

//	func BenchmarkStringsCompare_1024(b *testing.B) {
//		cases := getGestGetStringCompareCases(1024)
//		benchmarkStrCompare(b, cases)
//	}

func BenchmarkStringsCompare_4096(b *testing.B) {
	cases := getGestGetStringCompareCases(4096)
	benchmarkStrCompare(b, cases)
}
