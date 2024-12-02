/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"runtime"
	"time"

	"github.com/distroy/ldgo/v2/ldatomic"
	"github.com/distroy/ldgo/v2/ldsync"
	"go.uber.org/zap"
)

const (
	ctxCallerCacheSizeMax = 100 * 10000
)

var (
	ctxCallerCache = &ldsync.Map[any, *intervaler]{}
)

type caller struct {
	File string
	Line int
}

type intervaler struct {
	LastLogTime ldatomic.Time // for log at intervals
}

func (i *intervaler) hit(dur time.Duration) bool {
	last := i.LastLogTime.Load()
	next := last.Add(dur)
	now := time.Now()
	if now.After(next) {
		return i.LastLogTime.CompareAndSwap(last, now)
	}

	return false
}

func hitInterval(dur time.Duration, skip int) bool {
	if dur <= 0 {
		return true
	}
	i := getIntervaler(skip + 1)
	if i == nil {
		return false
	}

	return i.hit(dur)
}

func getIntervaler(skip int) *intervaler {
	_, file, line, _ := runtime.Caller(skip + 1)
	key := caller{
		File: file,
		Line: line,
	}
	c := getIntervalerByKey(key)
	if c == nil {
		Default().Error("reached the limit of interval callers, drop it",
			zap.String("file", file), zap.Int("line", line))
	}
	return c
}

func getIntervalerByKey(key interface{}) *intervaler {
	cache := ctxCallerCache
	i, _ := cache.Load(key)
	if i != nil {
		return i
	}

	if cache.Size() >= ctxCallerCacheSizeMax {
		return nil
	}

	i = &intervaler{}
	i.LastLogTime.Store(time.Time{})
	i, loaded := cache.LoadOrStore(key, i)
	if loaded {
		return i
	}

	n := cache.Size()
	if n <= ctxCallerCacheSizeMax {
		return i
	}

	cache.Del(key)
	return nil
}
