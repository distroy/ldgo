/*
 * Copyright (C) distroy
 */

package _slogtype

import "log/slog"

type Level = slog.Level

const (
	LevelDebug Level = slog.LevelDebug
	LevelInfo  Level = slog.LevelInfo
	LevelWarn  Level = slog.LevelWarn
	LevelError Level = slog.LevelError
	LevelPanic Level = 100
)
