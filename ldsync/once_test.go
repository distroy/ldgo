/*
 * Copyright (C) distroy
 */

package ldsync

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestOnce(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("concurrency", func(c convey.C) {
			times := int32(0)

			once := &Once{}
			c.So(once.Done(), convey.ShouldBeFalse)

			wg := &WaitGroup{}
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					once.Do(func() error {
						time.Sleep(time.Microsecond)
						atomic.AddInt32(&times, 1)
						return nil
					})
				}()
			}
			wg.Wait()

			c.So(once.Done(), convey.ShouldBeTrue)
			c.So(times, convey.ShouldEqual, 1)

			once.Reset()
			c.So(once.Done(), convey.ShouldBeFalse)

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					once.Do(func() error {
						time.Sleep(time.Microsecond)
						atomic.AddInt32(&times, 1)
						return nil
					})
				}()
			}
			wg.Wait()

			c.So(once.Done(), convey.ShouldBeTrue)
			c.So(times, convey.ShouldEqual, 2)
		})

		c.Convey("error", func(c convey.C) {
			var err error
			once := &Once{}
			err = once.Do(func() error { return fmt.Errorf("111") })
			c.So(err, convey.ShouldNotBeNil)
			c.So(once.Done(), convey.ShouldBeFalse)

			err = once.Do(func() error { return fmt.Errorf("111") })
			c.So(err, convey.ShouldNotBeNil)
			c.So(once.Done(), convey.ShouldBeFalse)

			err = once.Do(func() error { return nil })
			c.So(err, convey.ShouldBeNil)
			c.So(once.Done(), convey.ShouldBeTrue)
		})
	})
}
