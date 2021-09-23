/*
 * Copyright (C) distroy
 */

package ldtime

import "time"

const (
	fmtDateStr = "2006-01-02"
)

func divisionTimeNum(num *int64, dividend int64) int {
	n := *num % dividend
	*num /= dividend
	return int(n)
}

func GetTopicWeekDay(tm time.Time) int {
	wday := int(tm.Weekday())
	if wday == 0 {
		wday = 7
	}
	return wday
}

// TimeToDateNum format: 20060102 '%Y%m%d'
func TimeToDateNum(t time.Time) int64 {
	year, month, mday := t.Date()

	num := int64(year)
	num = num*100 + int64(month)
	num = num*100 + int64(mday)

	return num
}

// TimeToDateStr format: 2006-01-02 '%Y-%m-%d'
func TimeToDateStr(t time.Time) string {
	return t.Format(fmtDateStr)
}

// DateNumToTime format: 20060102 '%Y%m%d'
func DateNumToTime(dateNum int64, loc ...*time.Location) time.Time {
	nsec := 0
	hour, min, sec := 0, 0, 0
	day := divisionTimeNum(&dateNum, 100)
	month := divisionTimeNum(&dateNum, 100)
	year := int(dateNum)

	tz := time.Local
	if len(loc) != 0 {
		tz = loc[0]
	}
	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, tz)
}

func DateBegin(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, min, sec, nsec := 0, 0, 0, 0
	return time.Date(int(year), time.Month(month), int(day), hour, min, sec, nsec, t.Location())
}

// TimeToNum format: 20060102150405 '%Y%m%d%H%M%S'
func TimeToNum(t time.Time) int64 {
	year, month, mday := t.Date()
	hour, min, sec := t.Hour(), t.Minute(), t.Second()

	num := int64(year)
	num = num*100 + int64(month)
	num = num*100 + int64(mday)
	num = num*100 + int64(hour)
	num = num*100 + int64(min)
	num = num*100 + int64(sec)

	return num
}

// NumToTime format: 20060102150405 '%Y%m%d%H%M%S'
func NumToTime(timeNum int64, loc ...*time.Location) time.Time {
	nsec := 0
	sec := divisionTimeNum(&timeNum, 100)
	min := divisionTimeNum(&timeNum, 100)
	hour := divisionTimeNum(&timeNum, 100)
	mday := divisionTimeNum(&timeNum, 100)
	month := divisionTimeNum(&timeNum, 100)
	year := int(timeNum)

	tz := time.Local
	if len(loc) != 0 {
		tz = loc[0]
	}
	return time.Date(year, time.Month(month), mday, hour, min, sec, nsec, tz)
}
