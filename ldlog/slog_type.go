/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"log/slog"

	"github.com/distroy/ldgo/v3/ldlog/internal/_slogtype"
)

type (
	Level   = slog.Level
	Attr    = slog.Attr
	Value   = slog.Value
	Record  = slog.Record
	Handler = slog.Handler
)

const (
	LevelDebug Level = _slogtype.LevelDebug
	LevelInfo  Level = _slogtype.LevelInfo
	LevelWarn  Level = _slogtype.LevelWarn
	LevelError Level = _slogtype.LevelError
	LevelPanic Level = _slogtype.LevelPanic
)
