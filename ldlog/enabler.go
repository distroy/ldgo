/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"time"

	"github.com/distroy/ldgo/v2/ldrand"
	"go.uber.org/zap/zapcore"
)

type Enabler interface {
	Enable(lvl zapcore.Level, skip ...int) bool
}

// RateEnabler: enable based on probability(rate).
//   - Rate should be in [0, 1.0].
//   - Always enable log levels higher than error.
func RateEnabler(rate float64) Enabler {
	if rate <= 0 {
		return falseEnabler{}
	}
	if rate >= 1 {
		return defaultEnabler{}
	}
	return rateEnabler{rate: rate}
}

// IntervalEnabler: enable based on time interval.
//   - Calculate the time interval separately at each invocation location.
//   - Always enable log levels higher than error.
func IntervalEnabler(dur time.Duration) Enabler {
	if dur <= 0 {
		return defaultEnabler{}
	}
	return intervalEnabler{interval: dur}
}

// EnablerByInterval: get enabler by interval.
//   - Calculate the time interval separately at the get timing.
//   - Always enable log levels higher than error.
func EnablerByInterval(dur time.Duration, skip int) Enabler {
	enabler := IntervalEnabler(dur)
	if enabler.Enable(LevelInfo, skip+1) {
		return defaultEnabler{}
	}
	return falseEnabler{}
}

// EnablerByNameAndInterval: get enabler by name and interval.
//   - Calculate the time interval separately at the get timing.
//   - Always enable log levels higher than error.
func EnablerByNameAndInterval(name string, dur time.Duration) Enabler {
	if dur <= 0 {
		return defaultEnabler{}
	}
	i := getIntervalerByKey(name)
	if i == nil || !i.hit(dur) {
		return falseEnabler{}
	}
	return defaultEnabler{}
}

type defaultEnabler struct{}

func (p defaultEnabler) Enable(lvl zapcore.Level, skips ...int) bool { return true }

type falseEnabler struct{}

func (p falseEnabler) Enable(lvl zapcore.Level, skips ...int) bool { return lvl >= zapcore.ErrorLevel }

type rateEnabler struct {
	rate float64
}

func (p rateEnabler) Enable(lvl zapcore.Level, skips ...int) bool {
	if (falseEnabler{}).Enable(lvl, skips...) {
		return true
	}
	rate := p.rate
	if rate <= 0 {
		return false
	}
	if rate >= 1 {
		return true
	}
	return ldrand.Float64() < rate
}

type intervalEnabler struct {
	interval time.Duration
}

func (p intervalEnabler) Enable(lvl zapcore.Level, skips ...int) bool {
	if (falseEnabler{}).Enable(lvl, skips...) {
		return true
	}
	skip := 1
	if len(skips) > 0 && skips[0] > 0 {
		skip = skips[0]
	}
	return hitInterval(p.interval, skip+1)
}
