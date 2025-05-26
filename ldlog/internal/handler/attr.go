/*
 * Copyright (C) distroy
 */

package handler

import (
	"log/slog"
	"reflect"
	"time"
)

func init() {
	checkTypeEqual(reflect.TypeOf(Attr{}), reflect.TypeOf(slog.Attr{}))
}

func GetAttrPtr(v *slog.Attr) *Attr { return toType[*Attr](v) }
func GetAttr(v slog.Attr) Attr      { return *GetAttrPtr(&v) }
func GetAttrs(v []slog.Attr) []Attr { return toType[[]Attr](v) }

func GetSAttrs(v []Attr) []slog.Attr { return toType[[]slog.Attr](v) }

func String(key, value string) Attr { return GetAttr(slog.String(key, value)) }

func Int64(key string, value int64) Attr { return GetAttr(slog.Int64(key, value)) }
func Int(key string, value int) Attr     { return GetAttr(slog.Int(key, value)) }
func Uint64(key string, v uint64) Attr   { return GetAttr(slog.Uint64(key, v)) }

func Float64(key string, v float64) Attr { return GetAttr(slog.Float64(key, v)) }

func Bool(key string, v bool) Attr { return GetAttr(slog.Bool(key, v)) }

func Time(key string, v time.Time) Attr         { return GetAttr(slog.Time(key, v)) }
func Duration(key string, v time.Duration) Attr { return GetAttr(slog.Duration(key, v)) }

func Group(key string, args ...any) Attr { return GetAttr(slog.Group(key, args...)) }
func Any(key string, value any) Attr     { return GetAttr(slog.Any(key, value)) }

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
