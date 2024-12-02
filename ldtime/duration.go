/*
 * Copyright (C) distroy
 */

package ldtime

import (
	"time"

	"github.com/distroy/ldgo/v2/internal/timeinternal"
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

func (d Duration) MarshalJSON() ([]byte, error) { return timeinternal.DurationMarshalJSON(d.get()) }
func (d *Duration) UnmarshalJSON(b []byte) error {
	dur, err := timeinternal.DurationUnmarshalJSON(b)
	if err == nil {
		*d.ptr() = dur
	}
	return err
}
