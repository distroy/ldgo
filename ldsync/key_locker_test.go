/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"sync"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestLockerPool(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		kl := NewLockerPool()
		key0 := "123"
		key1 := 123

		l00 := kl.Get(key0)
		l10 := kl.Get(key1)
		l01 := kl.Get(key0)
		l11 := kl.Get(key1)

		c.So(l00.TryLock(), convey.ShouldBeTrue)
		c.So(l10.TryLock(), convey.ShouldBeTrue)
		c.So(l00.TryLock(), convey.ShouldBeTrue)
		c.So(l10.TryLock(), convey.ShouldBeTrue)
		c.So(l01.TryLock(), convey.ShouldBeFalse)
		c.So(l11.TryLock(), convey.ShouldBeFalse)

		var key1UnlockTime, key1LockTime time.Time
		var key2UnlockTime, key2LockTime time.Time

		wg := &sync.WaitGroup{}
		wg.Add(4)
		go func() {
			defer wg.Done()
			l01.Lock()
			key1LockTime = time.Now()
		}()
		go func() {
			defer wg.Done()
			time.Sleep(1 * time.Millisecond)
			l00.Unlock()
			key1UnlockTime = time.Now()
		}()
		go func() {
			defer wg.Done()
			l11.Lock()
			key2LockTime = time.Now()
		}()
		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			l10.Unlock()
			key2UnlockTime = time.Now()
		}()
		wg.Wait()

		c.So(key1UnlockTime, convey.ShouldHappenBefore, key1LockTime)
		c.So(key2UnlockTime, convey.ShouldHappenBefore, key2LockTime)
		c.So(key1LockTime, convey.ShouldHappenBefore, key2LockTime)
		c.So(key1UnlockTime, convey.ShouldHappenBefore, key2UnlockTime)

		l01.Unlock()
		l11.Unlock()
		l11.Unlock() // double unlock

		c.So(l00.TryLock(), convey.ShouldBeTrue)
		c.So(l10.TryLock(), convey.ShouldBeTrue)
		c.So(l01.TryLock(), convey.ShouldBeFalse)
		c.So(l11.TryLock(), convey.ShouldBeFalse)
	})
}
