/*
 * Copyright (C) distroy
 */

package ldlog

import "log/slog"

type logHandler interface {
	Handler

	Sync() error
	Close() error

	Level() Level
	Sequence() string
}

func wrapHandler(h slog.Handler) logHandler {
	if h == nil {
		return nil
	}
	if hh, _ := h.(logHandler); hh != nil {
		return hh
	}
	return handlerWrapper{h}
}

type handlerWrapper struct {
	Handler
}

func (h handlerWrapper) Sync() error  { return nil }
func (h handlerWrapper) Close() error { return nil }

func (h handlerWrapper) Level() Level     { return 0 }
func (h handlerWrapper) Sequence() string { return "" }
