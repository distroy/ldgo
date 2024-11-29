/*
 * Copyright (C) distroy
 */

package ldtime

import (
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

type Duration time.Duration

func (d Duration) Duration() time.Duration { return d.get() }
func (d Duration) get() time.Duration      { return time.Duration(d) }
func (d *Duration) ptr() *time.Duration    { return (*time.Duration)(d) }

func (d Duration) Abs() Duration                { return Duration(d.get().Abs()) }
func (d Duration) Hours() float64               { return d.get().Hours() }
func (d Duration) Microseconds() int64          { return d.get().Microseconds() }
func (d Duration) Milliseconds() int64          { return d.get().Milliseconds() }
func (d Duration) Minutes() float64             { return d.get().Minutes() }
func (d Duration) Nanoseconds() int64           { return d.get().Nanoseconds() }
func (d Duration) Round(m Duration) Duration    { return Duration(d.get().Round(m.get())) }
func (d Duration) Seconds() float64             { return d.get().Seconds() }
func (d Duration) String() string               { return d.get().String() }
func (d Duration) Truncate(m Duration) Duration { return Duration(d.get().Truncate(m.get())) }

func (d Duration) MarshalJSON() ([]byte, error) {
	s := d.String()
	buf := make([]byte, 0, len(s)+4)
	buf = strconv.AppendQuote(buf, s)
	return buf, nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return fmt.Errorf("unexpected end of JSON input")
	}

	str := unsafe.String(unsafe.SliceData(b), len(b))
	switch b[0] {
	case '"', '\'', '`':
		vv, err := strconv.Unquote(str)
		if err != nil {
			return fmt.Errorf("invalid duration: %s", b)
		}
		*d.ptr(), err = time.ParseDuration(vv)
		if err != nil {
			return fmt.Errorf("invalid duration: %s", b)
		}
		return nil

	}

	if i64, err := strconv.ParseInt(str, 10, 64); err == nil {
		*d = Duration(i64)
		return nil
	}
	// if f64, err := strconv.ParseFloat(str, 64); err != nil {
	// 	*d = Duration(f64)
	// 	return nil
	// }
	return fmt.Errorf("invalid duration: %s", b)
}
