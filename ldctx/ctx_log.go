/*
 * Copyright (C) distroy
 */

package ldctx

import (
	"context"
	"fmt"
	_ "unsafe"

	"github.com/distroy/ldgo/v3/ldlog"
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

//go:linkname logFmt github.com/distroy/ldgo/v3/ldlog.logFmt
func logFmt(l *ldlog.Logger, lvl ldlog.Level, skip int, fmt string, args ...any)

//go:linkname logAttrs github.com/distroy/ldgo/v3/ldlog.logAttrs
func logAttrs(l *ldlog.Logger, lvl ldlog.Level, skip int, msg string, args ...ldlog.Attr)

func ctxLogFmt(c Context, lvl ldlog.Level, fmt string, args ...any) {
	format(fmt, args...)
	logFmt(GetLogger(c), lvl, 2, fmt, args...)
}

func ctxLogAttr(c Context, lvl ldlog.Level, msg string, args ...ldlog.Attr) {
	logAttrs(GetLogger(c), lvl, 2, msg, args...)
}

const (
	lvlD = ldlog.LevelDebug
	lvlI = ldlog.LevelInfo
	lvlW = ldlog.LevelWarn
	lvlE = ldlog.LevelError
	lvlP = ldlog.LevelPanic
)

func format(format string, args ...any) {
	if formatFlag {
		_ = fmt.Sprintf(format, args...)
	}
}

func LogD(c Context, msg string, fields ...ldlog.Attr) { ctxLogAttr(c, lvlD, msg, fields...) }
func LogI(c Context, msg string, fields ...ldlog.Attr) { ctxLogAttr(c, lvlI, msg, fields...) }
func LogW(c Context, msg string, fields ...ldlog.Attr) { ctxLogAttr(c, lvlW, msg, fields...) }
func LogE(c Context, msg string, fields ...ldlog.Attr) { ctxLogAttr(c, lvlE, msg, fields...) }
func LogP(c Context, msg string, fields ...ldlog.Attr) { ctxLogAttr(c, lvlP, msg, fields...) }

func LogDf(c Context, fmt string, args ...any) { ctxLogFmt(c, lvlD, fmt, args...) }
func LogIf(c Context, fmt string, args ...any) { ctxLogFmt(c, lvlI, fmt, args...) }
func LogWf(c Context, fmt string, args ...any) { ctxLogFmt(c, lvlW, fmt, args...) }
func LogEf(c Context, fmt string, args ...any) { ctxLogFmt(c, lvlE, fmt, args...) }
func LogPf(c Context, fmt string, args ...any) { ctxLogFmt(c, lvlP, fmt, args...) }
