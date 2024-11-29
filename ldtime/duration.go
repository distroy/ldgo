/*
 * Copyright (C) distroy
 */

package ldtime

import (
	"encoding/json"
	"fmt"
	"time"
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
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch vv := v.(type) {
	case float64:
		*d = Duration(time.Duration(vv))
		return nil

	case string:
		var err error
		*d.ptr(), err = time.ParseDuration(vv)
		if err != nil {
			return err
		}
		return nil

	case json.Number:
		if i64, err := vv.Int64(); err == nil {
			*d = Duration(i64)
			return nil
		} else if f64, err := vv.Float64(); err != nil {
			*d = Duration(f64)
			return nil
		} else {
			return err
		}

	default:
		return fmt.Errorf("invalid duration. str:%s", b)
	}
}
