/*
 * Copyright (C) distroy
 */

package ldatomic

import (
	"time"

	"github.com/distroy/ldgo/v2/internal/timeinternal"
)

type Duration int64

func NewDuration(d time.Duration) *Duration {
	v := Duration(d)
	return &v
}

func (p *Duration) get() *Int64 { return (*Int64)(p) }

func (p *Duration) Store(d time.Duration) { p.get().Store(int64(d)) }
func (p *Duration) Load() time.Duration   { return time.Duration(p.get().Load()) }
func (p *Duration) Swap(old time.Duration) (new time.Duration) {
	return time.Duration(p.get().Swap(int64(old)))
}
func (p *Duration) CompareAndSwap(old, new time.Duration) (swapped bool) {
	return p.get().CompareAndSwap(int64(old), int64(new))
}
func (p *Duration) Add(delta time.Duration) (new time.Duration) {
	return time.Duration(p.get().Add(int64(delta)))
}
func (p *Duration) Sub(delta time.Duration) (new time.Duration) {
	return time.Duration(p.get().Sub(int64(delta)))
}

func (p Duration) MarshalJSON() ([]byte, error) { return timeinternal.DurationMarshalJson(p.Load()) }
func (p *Duration) UnmarshalJSON(b []byte) error {
	dur, err := timeinternal.DurationUnmarshalJson(b)
	if err == nil {
		p.Store(dur)
	}
	return err
}
