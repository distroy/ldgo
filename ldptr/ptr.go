/*
 * Copyright (C) distroy
 */

package ldptr

func Get[T any](p *T, def ...T) T {
	if p != nil {
		return *p
	}
	if len(def) > 0 {
		return def[0]
	}
	var v T
	return v
}

func New[T any](d T) *T { return &d }

func NewByPtr[T any](d *T, def ...T) *T {
	if d == nil {
		if len(def) > 0 {
			return &def[0]
		}
		return nil
	}
	cp := *d
	return &cp
}
