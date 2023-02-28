/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

var (
	testTimezone = testFixedTimezone("Asia/Bejing", +(8 * time.Hour))
)

func testFixedTimezone(name string, tm time.Duration) *time.Location {
	sec := int(tm / time.Second)
	return time.FixedZone(name, sec)
}

func TestTime(t *testing.T) {
	const format = "2006-01-02T15:04:05-0700"
	now := time.Unix(1644479966, 0).In(testTimezone)
	convey.Convey(t.Name(), t, func() {
		convey.So(NewTime(zeroTime).Load().Format(format), convey.ShouldEqual, zeroTime.Format(format))

		p := Time{}
		convey.So(p.Load().Format(format), convey.ShouldEqual, zeroTime.Format(format))

		p.Store(now)
		convey.So(p.Load().Format(format), convey.ShouldEqual, "2022-02-10T15:59:26+0800")

		p.Add(60 * time.Second)
		convey.So(p.Load().Format(format), convey.ShouldEqual, "2022-02-10T16:00:26+0800")

		p.AddDate(-1, 1, 3)
		convey.So(p.Load().Format(format), convey.ShouldEqual, "2021-03-13T16:00:26+0800")

		convey.So(p.Swap(now).Format(format), convey.ShouldEqual, "2021-03-13T16:00:26+0800")
		convey.So(p.Load().Format(format), convey.ShouldEqual, "2022-02-10T15:59:26+0800")
	})
}
