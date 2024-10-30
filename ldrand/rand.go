/*
 * Copyright (C) distroy
 */

package ldrand

import (
	"time"

	"github.com/distroy/ldgo/v3/ldconv"
)

var (
	randStringLetters = ldconv.StrToBytesUnsafe("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	globalRand        = New(NewFastSource(time.Now().UnixNano()))
)

func Int() int     { return globalRand.Int() }
func Int31() int32 { return globalRand.Int31() }
func Int63() int64 { return globalRand.Int63() }

func Intn(n int) int       { return globalRand.Intn(n) }
func Int31n(n int32) int32 { return globalRand.Int31n(n) }
func Int63n(n int64) int64 { return globalRand.Int63n(n) }

func Uint() uint     { return globalRand.Uint() }
func Uint32() uint32 { return globalRand.Uint32() }
func Uint64() uint64 { return globalRand.Uint64() }

// Float64 returns, as a float64, a pseudo-random number in [0.0,1.0).
func Float64() float64 { return globalRand.Float64() }

// Float32 returns, as a float32, a pseudo-random number in [0.0,1.0).
func Float32() float32 { return globalRand.Float32() }

// NormFloat64 returns a normally distributed float64 in the range
// [-math.MaxFloat64, +math.MaxFloat64] with
// standard normal distribution (mean = 0, stddev = 1)
// from the default Source.
// To produce a different normal distribution, callers can
// adjust the output using:
//
//	sample = NormFloat64() * desiredStdDev + desiredMean
func NormFloat64() float64 { return globalRand.NormFloat64() }

// ExpFloat64 returns an exponentially distributed float64 in the range
// (0, +math.MaxFloat64] with an exponential distribution whose rate parameter
// (lambda) is 1 and whose mean is 1/lambda (1) from the default Source.
// To produce a distribution with a different rate parameter,
// callers can adjust the output using:
//
//	sample = ExpFloat64() / desiredRateParameter
func ExpFloat64() float64 { return globalRand.ExpFloat64() }

func Read(p []byte) (int, error)         { return globalRand.Read(p) }
func Shuffle(n int, swap func(i, j int)) { globalRand.Shuffle(n, swap) }
func Perm(n int) []int                   { return globalRand.Perm(n) }

func String(n int) string {
	letters := randStringLetters

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[Intn(len(letters))]
	}

	return ldconv.BytesToStrUnsafe(b)
}

func Bytes(n int) []byte {
	if n <= 0 {
		return nil
	}
	b := make([]byte, n)
	globalRand.Read(b) // nolint
	return b
}
