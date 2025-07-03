/*
 * Copyright (C) distroy
 */

package _handler

import (
	"log/slog"

	"github.com/distroy/ldgo/v3/ldlog/internal/_logref"
	"github.com/distroy/ldgo/v3/ldlog/internal/_slogtype"
)

func asType[T any](v any, def ...T) T { return _logref.AsType(v, def...) }

type (
	Level  = _slogtype.Level
	Kind   = _slogtype.Kind
	Attr   = _slogtype.Attr
	Value  = _slogtype.Value
	Source = _slogtype.Source
	Record = _slogtype.Record
)

// *** level begin ****

const (
	LevelDebug Level = _slogtype.LevelDebug
	LevelInfo  Level = _slogtype.LevelInfo
	LevelWarn  Level = _slogtype.LevelWarn
	LevelError Level = _slogtype.LevelError
	LevelPanic Level = _slogtype.LevelPanic
)

// *** level end ****

// *** kind begin ****

const (
	KindAny       = _slogtype.KindAny
	KindBool      = _slogtype.KindBool
	KindDuration  = _slogtype.KindDuration
	KindFloat64   = _slogtype.KindFloat64
	KindInt64     = _slogtype.KindInt64
	KindString    = _slogtype.KindString
	KindTime      = _slogtype.KindTime
	KindUint64    = _slogtype.KindUint64
	KindGroup     = _slogtype.KindGroup
	KindLogValuer = _slogtype.KindLogValuer
)

// *** kind end ****

// *** source begin ****

func GetSourcePtr(v *slog.Source) *Source { return _slogtype.GetSourcePtr(v) }
func GetSource(v slog.Source) Source      { return _slogtype.GetSource(v) }

// *** source end ****

// *** value begin ****

func GetValuePtr(v *slog.Value) *Value { return _slogtype.GetValuePtr(v) }
func GetValue(v slog.Value) Value      { return _slogtype.GetValue(v) }

func CountEmptyGroups(as []Attr) int { return _slogtype.CountEmptyGroups(as) }

// *** value end ****

// *** attr begin ****

func GetAttrPtr(v *slog.Attr) *Attr { return _slogtype.GetAttrPtr(v) }
func GetAttr(v slog.Attr) Attr      { return _slogtype.GetAttr(v) }
func GetAttrs(v []slog.Attr) []Attr { return _slogtype.GetAttrs(v) }

func GetSAttrs(v []Attr) []slog.Attr { return _slogtype.GetSAttrs(v) }

// *** attr end ****

// *** record begin ****

func GetRecordPtr(v *slog.Record) *Record { return _slogtype.GetRecordPtr(v) }
func GetRecord(v slog.Record) Record      { return _slogtype.GetRecord(v) }

// *** record end ****
