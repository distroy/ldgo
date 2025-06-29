/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"context"
	"math"
)

func newDiscard() *Logger {
	h := discardHandler{}
	return newLogger(newCore(h))
}

var _ logHandler = (*discardHandler)(nil)

type discardHandler struct{}

func (_ discardHandler) Enabled(context.Context, Level) bool  { return false }
func (_ discardHandler) Handle(context.Context, Record) error { return nil }
func (h discardHandler) WithAttrs(attrs []Attr) Handler       { return h }
func (h discardHandler) WithGroup(name string) Handler        { return h }

func (_ discardHandler) Sync() error      { return nil }
func (_ discardHandler) Close() error     { return nil }
func (_ discardHandler) Level() Level     { return math.MaxInt }
func (_ discardHandler) Sequence() string { return "" }
