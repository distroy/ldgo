/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"go.uber.org/zap"
)

const (
	lvlD = LevelDebug
	lvlI = LevelInfo
	lvlW = LevelWarn
	lvlE = LevelError
	lvlF = LevelFatal
	lvlP = LevelPanic
)

func newLogger(log core) *Logger {
	return &Logger{
		core: log,
	}
}

type Logger struct {
	core
}

func (l *Logger) Wrapper() *Wrapper { return (*Wrapper)(l) }

func (l *Logger) clone() *Logger {
	cp := *l
	return &cp
}

func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	log := l.Core().WithOptions(opts...)
	l = l.clone()
	l.log = log
	l.sugar = log.Sugar()
	return l
}

func (l *Logger) With(attrs ...Attr) *Logger {
	if len(attrs) == 0 {
		return l
	}
	log := l.Core().WithLazy(attrs...)
	l = l.clone()
	l.log = log
	l.sugar = log.Sugar()
	return l
}

func (l *Logger) WithSequence(seq string) *Logger {
	if seq == "" || l.Sequence() == seq {
		return l
	}
	l = l.clone()
	l.core = *l.core.withSequence(seq)
	return l
}

func (l *Logger) WithEnabler(p Enabler) *Logger {
	if p == nil {
		p = defaultEnabler{}
	}
	if p == l.enabler {
		return l
	}
	l = l.clone()
	l.enabler = p
	return l
}

func (l *Logger) Debug(msg string, attrs ...Attr) { l.logAttrs(nil, lvlD, 1, msg, attrs...) }
func (l *Logger) Info(msg string, attrs ...Attr)  { l.logAttrs(nil, lvlI, 1, msg, attrs...) }
func (l *Logger) Warn(msg string, attrs ...Attr)  { l.logAttrs(nil, lvlW, 1, msg, attrs...) }
func (l *Logger) Error(msg string, attrs ...Attr) { l.logAttrs(nil, lvlE, 1, msg, attrs...) }
func (l *Logger) Panic(msg string, attrs ...Attr) { l.logAttrs(nil, lvlP, 1, msg, attrs...) }

func (l *Logger) Debugf(fmt string, args ...any) { l.logFmt(nil, lvlD, 1, fmt, args...) }
func (l *Logger) Infof(fmt string, args ...any)  { l.logFmt(nil, lvlI, 1, fmt, args...) }
func (l *Logger) Warnf(fmt string, args ...any)  { l.logFmt(nil, lvlW, 1, fmt, args...) }
func (l *Logger) Errorf(fmt string, args ...any) { l.logFmt(nil, lvlE, 1, fmt, args...) }
func (l *Logger) Panicf(fmt string, args ...any) { l.logFmt(nil, lvlP, 1, fmt, args...) }

func (l *Logger) Debugln(args ...any) { l.log(nil, lvlD, 1, pw(args)) }
func (l *Logger) Infoln(args ...any)  { l.log(nil, lvlD, 1, pw(args)) }
func (l *Logger) Warnln(args ...any)  { l.log(nil, lvlD, 1, pw(args)) }
func (l *Logger) Errorln(args ...any) { l.log(nil, lvlD, 1, pw(args)) }
func (l *Logger) Panicln(args ...any) { l.log(nil, lvlD, 1, pw(args)) }
