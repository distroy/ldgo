/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"io"

	"github.com/distroy/ldgo/v3/ldlog/internal/handler"
)

func GetLogger(h Handler) *Logger { return newLogger(newCore(h)) }

func newHandler(w io.Writer) Handler {
	return handler.NewHandler(w, &handler.Options{
		Caller: true,
		Level:  LevelInfo,
	})
}

func New(w io.Writer, opts ...Option) *Logger {
	core := newCore(newHandler(w))
	for _, opt := range opts {
		opt(&core)
	}

	return newLogger(core)
}
