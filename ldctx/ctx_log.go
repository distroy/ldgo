/*
 * Copyright (C) distroy
 */

package ldctx

import (
	"context"
	"fmt"
	_ "unsafe"

	"github.com/distroy/ldgo/v2/ldlog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	formatFlag = false
)

type ctxKeyType int

const (
	ctxKeyLogger ctxKeyType = iota
	ctxKeyMap
)

var (
	defaultContext = context.Background()
	consoleContext = WithLogger(context.Background(), ldlog.Console())
	discardContext = WithLogger(context.Background(), ldlog.Discard())
)

func defaultLogger() *ldlog.Logger { return ldlog.Default() }

type stringer interface {
	String() string
}

//go:linkname zCoreByLogger github.com/distroy/ldgo/v2/ldlog.zCoreByLogger
func zCoreByLogger(l *ldlog.Logger, lvl zapcore.Level, skip int) *zap.Logger

//go:linkname zSugarByLogger github.com/distroy/ldgo/v2/ldlog.zSugarByLogger
func zSugarByLogger(l *ldlog.Logger, lvl zapcore.Level, skip int) *zap.SugaredLogger

func zCore(c context.Context, lvl zapcore.Level, skip int) *zap.Logger {
	return zCoreByLogger(GetLogger(c), lvl, skip+1)
}
func zSugar(c context.Context, lvl zapcore.Level, skip int) *zap.SugaredLogger {
	return zSugarByLogger(GetLogger(c), lvl, skip+1)
}

const (
	lvlD = zapcore.DebugLevel
	lvlI = zapcore.InfoLevel
	lvlW = zapcore.WarnLevel
	lvlE = zapcore.ErrorLevel
	lvlP = zapcore.PanicLevel
	lvlF = zapcore.FatalLevel
)

func format(format string, args ...interface{}) {
	if formatFlag {
		_ = fmt.Sprintf(format, args...)
	}
}

func LogD(c Context, msg string, fields ...zap.Field) { zCore(c, lvlD, 1).Debug(msg, fields...) }
func LogI(c Context, msg string, fields ...zap.Field) { zCore(c, lvlI, 1).Info(msg, fields...) }
func LogW(c Context, msg string, fields ...zap.Field) { zCore(c, lvlW, 1).Warn(msg, fields...) }
func LogE(c Context, msg string, fields ...zap.Field) { zCore(c, lvlE, 1).Error(msg, fields...) }
func LogP(c Context, msg string, fields ...zap.Field) { zCore(c, lvlP, 1).Panic(msg, fields...) }
func LogF(c Context, msg string, fields ...zap.Field) { zCore(c, lvlF, 1).Fatal(msg, fields...) }

func LogDf(c Context, fmt string, args ...interface{}) {
	format(fmt, args...)
	zSugar(c, lvlD, 1).Debugf(fmt, args...)
}
func LogIf(c Context, fmt string, args ...interface{}) {
	format(fmt, args...)
	zSugar(c, lvlI, 1).Infof(fmt, args...)
}
func LogWf(c Context, fmt string, args ...interface{}) {
	format(fmt, args...)
	zSugar(c, lvlW, 1).Warnf(fmt, args...)
}
func LogEf(c Context, fmt string, args ...interface{}) {
	format(fmt, args...)
	zSugar(c, lvlE, 1).Errorf(fmt, args...)
}
func LogPf(c Context, fmt string, args ...interface{}) {
	format(fmt, args...)
	zSugar(c, lvlP, 1).Panicf(fmt, args...)
}
func LogFf(c Context, fmt string, args ...interface{}) {
	format(fmt, args...)
	zSugar(c, lvlF, 1).Fatalf(fmt, args...)
}
