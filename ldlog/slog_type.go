/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/distroy/ldgo/v3/ldlog/internal/handler"
)

type (
	Level   = slog.Level
	Attr    = slog.Attr
	Value   = slog.Value
	Record  = slog.Record
	Handler = slog.Handler
)

const (
	LevelDebug Level = slog.LevelDebug
	LevelInfo  Level = slog.LevelInfo
	LevelWarn  Level = slog.LevelWarn
	LevelError Level = slog.LevelError
	LevelPanic Level = handler.LevelPanic
)

func rec2err(r *Record) error {
	buf := &strings.Builder{}
	buf.Grow(1024)
	buf.WriteString(r.Message)

	first := true
	r.Attrs(func(a slog.Attr) bool {
		if first {
			buf.WriteByte('|')
			first = false

		} else {
			buf.WriteByte(',')
		}
		buf.WriteString(a.Key)
		buf.WriteByte('=')
		buf.WriteString(a.Value.String())
		return true
	})
	return errors.New(buf.String())
}
