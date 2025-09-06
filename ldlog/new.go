/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"io"
	"os"

	"github.com/distroy/ldgo/v3/ldlog/internal/_handler"
)

func GetLogger(h Handler) *Logger { return newLogger(newCore(h)) }

type HandlerOptions = _handler.Options

func NewHandler(w io.Writer, opts *HandlerOptions) Handler {
	if w == nil {
		w = os.Stderr
	}
	if opts == nil {
		opts = &HandlerOptions{
			Caller: true,
			Level:  LevelInfo,
		}
	}
	return _handler.NewHandler(w, opts)
}

func New(h Handler, opts ...Option) *Logger {
	core := newCore(h)
	for _, opt := range opts {
		opt(&core)
	}
	return newLogger(core)
}
