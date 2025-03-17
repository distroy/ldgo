/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	formatFlag = false
)

func newCore(log *zap.Logger, sugar *zap.SugaredLogger) core {
	return core{
		log:     log,
		sugar:   sugar,
		enabler: defaultEnabler{},
	}
}

type core struct {
	log      *zap.Logger
	sugar    *zap.SugaredLogger
	sequence string
	enabler  Enabler
}

func (l *core) Sync() { l.log.Sync() }

func (l *core) Enabled(lvl zapcore.Level) bool { return l.enable(lvl, 1) }
func (l *core) Enabler() Enabler               { return l.enabler }

func (l *core) Core() *zap.Logger         { return l.log }
func (l *core) Sugar() *zap.SugaredLogger { return l.sugar }

func (l *core) zCore(lvl zapcore.Level, skip int) *zap.Logger {
	if !l.enable(lvl, skip+1) {
		return Discard().Core()
	}
	return l.log
}
func (l *core) zSugar(lvl zapcore.Level, skip int) *zap.SugaredLogger {
	if !l.enable(lvl, skip+1) {
		return Discard().Sugar()
	}
	return l.sugar
}

func (l *core) enable(lvl zapcore.Level, skip int) bool {
	if !l.log.Core().Enabled(lvl) {
		return false
	}

	return l.enabler.Enable(lvl, skip+1)
}

func (l *core) format(format string, args ...any) {
	if formatFlag {
		_ = fmt.Sprintf(format, args...)
	}
}
