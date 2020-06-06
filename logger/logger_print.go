/*
 * Copyright (C) distroy
 */

package logger

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

const (
	_DEBUG  = zapcore.DebugLevel
	_INFO   = zapcore.InfoLevel
	_WARN   = zapcore.WarnLevel
	_ERROR  = zapcore.ErrorLevel
	_DPANIC = zapcore.DPanicLevel
	_PANIC  = zapcore.PanicLevel
	_FATAL  = zapcore.FatalLevel
)

func sprintln(args []interface{}) string {
	if len(args) == 0 {
		return ""
	}

	text := fmt.Sprintln(args...)
	size := len(text)
	if size == 0 {
		return ""
	}
	if text[size-1] == '\n' {
		return text[:size-1]
	}
	return text
}
