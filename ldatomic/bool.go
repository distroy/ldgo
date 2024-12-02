/*
 * Copyright (C) distroy
 */

package ldatomic

type Bool uint32

func NewBool(d bool) *Bool {
	v := Bool(0)
	v.Store(d)
	return &v
}

func (p *Bool) get() *Uint32 { return (*Uint32)(p) }

func (p *Bool) Store(d bool)             { p.get().Store(p.toInt(d)) }
func (p *Bool) Load() bool               { return p.toBool(p.get().Load()) }
func (p *Bool) Swap(new bool) (old bool) { return p.toBool(p.get().Swap(p.toInt(new))) }
func (p *Bool) CompareAndSwap(old, new bool) (swapped bool) {
	return p.get().CompareAndSwap(p.toInt(old), p.toInt(new))
}

func (p *Bool) Toggle() (old bool) {
	for {
		old := p.Load()
		if p.CompareAndSwap(old, !old) {
			return old
		}
	}
}

func (p *Bool) toInt(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}
func (p *Bool) toBool(d uint32) bool {
	return d != 0
}

func (v Bool) MarshalJSON() ([]byte, error)  { return marshalJSON[bool](&v) }
func (v *Bool) UnmarshalJSON(b []byte) error { return unmarshalJSON[bool](v, b) }
