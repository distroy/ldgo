/*
 * Copyright (C) distroy
 */

package _slogtype

import (
	"fmt"
	"log/slog"
)

var (
	_ slog.Leveler = Level(0)
)

const (
	LevelTrace = Level(-100)
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
	LevelPanic = Level(100)
)

type Level slog.Level

func (l Level) Level() slog.Level { return slog.Level(l) }

func (l Level) String() string {
	str := func(base string, val Level) string {
		if val <= 0 {
			return base
		}
		return fmt.Sprintf("%s%+d", base, val)
	}

	switch {
	case l < LevelDebug:
		return str("TRACE", l-LevelTrace)
	case l < LevelInfo:
		return str("DEBUG", l-LevelDebug)
	case l < LevelWarn:
		return str("INFO", l-LevelInfo)
	case l < LevelError:
		return str("WARN", l-LevelWarn)
	case l < LevelPanic:
		return str("ERROR", l-LevelError)
	default:
		return str("PANIC", l-LevelPanic)
	}
}
