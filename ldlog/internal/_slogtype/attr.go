/*
 * Copyright (C) distroy
 */

package _slogtype

import (
	"log/slog"
	"reflect"
)

func init() {
	checkTypeEqual(reflect.TypeOf(Attr{}), reflect.TypeOf(slog.Attr{}))
}

func GetAttrPtr(v *slog.Attr) *Attr { return toType[*Attr](v) }
func GetAttr(v slog.Attr) Attr      { return *GetAttrPtr(&v) }
func GetAttrs(v []slog.Attr) []Attr { return toType[[]Attr](v) }

func GetSAttrs(v []Attr) []slog.Attr { return toType[[]slog.Attr](v) }

type Attr struct {
	Key   string
	Value Value
}

func (a *Attr) Get() *slog.Attr { return toType[*slog.Attr](a) }

func (a *Attr) Equal(b Attr) bool { return a.Get().Equal(*b.Get()) }
func (a *Attr) String() string    { return a.Get().String() }
func (a *Attr) IsEmpty() bool {
	return a.Key == "" && a.Value.num == 0 && a.Value.any == nil
}
