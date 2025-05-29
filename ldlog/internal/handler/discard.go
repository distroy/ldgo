/*
 * Copyright (C) distroy
 */

package handler

import (
	"context"
	"log/slog"
)

var (
	_ slog.Handler = (*Discard)(nil)
)

type Discard struct{}

func (h Discard) IsDiscard() bool                           { return true }
func (h Discard) Enabled(context.Context, slog.Level) bool  { return false }
func (h Discard) Handle(context.Context, slog.Record) error { return nil }
func (h Discard) WithAttrs(attrs []slog.Attr) slog.Handler  { return h }
func (h Discard) WithGroup(name string) slog.Handler        { return h }
