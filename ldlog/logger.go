/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	lvlD = zapcore.DebugLevel
	lvlI = zapcore.InfoLevel
	lvlW = zapcore.WarnLevel
	lvlE = zapcore.ErrorLevel
	lvlF = zapcore.FatalLevel
	lvlP = zapcore.PanicLevel
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

func (l *Logger) With(fields ...zap.Field) *Logger {
	if len(fields) == 0 {
		return l
	}
	log := l.Core().WithLazy(fields...)
	l = l.clone()
	l.log = log
	l.sugar = log.Sugar()
	return l
}

func (l *Logger) GetSequence() string { return l.sequence }
func (l *Logger) WithSequence(seq string) *Logger {
	if seq == "" || l.sequence == seq {
		return l
	}
	log := l.Core().With(zap.String(sequenceKey, seq))
	l = l.clone()
	l.log = log
	l.sugar = log.Sugar()
	l.sequence = seq
	return l
}

// Log based on probability(rate). rate should be in [0, 1.0]
//
// Deprecated: use `WithEnabler` instead.
func (l *Logger) WithRate(rate float64) *Logger { return l.WithEnabler(RateEnabler(rate)) }

// Log based on time interval.
//
// Deprecated: use `WithEnabler` instead.
func (l *Logger) WithInterval(d time.Duration) *Logger { return l.WithEnabler(IntervalEnabler(d)) }

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

func (l *Logger) Debug(msg string, fields ...zap.Field) { l.zCore(lvlD, 1).Debug(msg, fields...) }
func (l *Logger) Info(msg string, fields ...zap.Field)  { l.zCore(lvlI, 1).Info(msg, fields...) }
func (l *Logger) Warn(msg string, fields ...zap.Field)  { l.zCore(lvlW, 1).Warn(msg, fields...) }
func (l *Logger) Error(msg string, fields ...zap.Field) { l.zCore(lvlE, 1).Error(msg, fields...) }
func (l *Logger) Panic(msg string, fields ...zap.Field) { l.zCore(lvlP, 1).Panic(msg, fields...) }
func (l *Logger) Fatal(msg string, fields ...zap.Field) { l.zCore(lvlF, 1).Fatal(msg, fields...) }

func (l *Logger) Debugf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlD, 1).Debugf(fmt, args...)
}
func (l *Logger) Infof(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlI, 1).Infof(fmt, args...)
}
func (l *Logger) Warnf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlW, 1).Warnf(fmt, args...)
}
func (l *Logger) Errorf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlE, 1).Errorf(fmt, args...)
}
func (l *Logger) Panicf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlP, 1).Panicf(fmt, args...)
}
func (l *Logger) Fatalf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlF, 1).Fatalf(fmt, args...)
}

func (l *Logger) Debugln(args ...interface{}) { l.zSugar(lvlD, 1).Debug(pw(args)) }
func (l *Logger) Infoln(args ...interface{})  { l.zSugar(lvlI, 1).Info(pw(args)) }
func (l *Logger) Warnln(args ...interface{})  { l.zSugar(lvlW, 1).Warn(pw(args)) }
func (l *Logger) Errorln(args ...interface{}) { l.zSugar(lvlE, 1).Error(pw(args)) }
func (l *Logger) Panicln(args ...interface{}) { l.zSugar(lvlP, 1).Panic(pw(args)) }
func (l *Logger) Fatalln(args ...interface{}) { l.zSugar(lvlF, 1).Fatal(pw(args)) }
