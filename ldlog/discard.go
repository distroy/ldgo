/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newDiscard() *Logger {
	core := discardCore{}
	zlog := zap.New(core)
	return newLogger(newCore(zlog, zlog.Sugar()))
}

var _ zapcore.Core = (*discardCore)(nil)

type discardCore struct{}

func (_ discardCore) Enabled(l zapcore.Level) bool                   { return false }
func (_ discardCore) With(f []zapcore.Field) zapcore.Core            { return discardCore{} }
func (_ discardCore) Write(e zapcore.Entry, f []zapcore.Field) error { return nil }
func (_ discardCore) Sync() error                                    { return nil }

func (_ discardCore) Check(e zapcore.Entry, c *zapcore.CheckedEntry) *zapcore.CheckedEntry { return c }
