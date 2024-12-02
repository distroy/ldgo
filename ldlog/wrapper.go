/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Wrapper struct {
	core
}

func (l *Wrapper) Logger() *Logger { return (*Logger)(l) }

func (l *Wrapper) Debugf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlD, 1).Debugf(fmt, args...)
}
func (l *Wrapper) Debug(args ...interface{})              { l.zSugar(lvlD, 1).Debug(pw(args)) }
func (l *Wrapper) Debugln(args ...interface{})            { l.zSugar(lvlD, 1).Debug(pw(args)) }
func (l *Wrapper) Debugz(fmt string, fields ...zap.Field) { l.zCore(lvlD, 1).Debug(fmt, fields...) }

func (l *Wrapper) Infof(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlI, 1).Infof(fmt, args...)
}
func (l *Wrapper) Info(args ...interface{})              { l.zSugar(lvlI, 1).Info(pw(args)) }
func (l *Wrapper) Infoln(args ...interface{})            { l.zSugar(lvlI, 1).Info(pw(args)) }
func (l *Wrapper) Infoz(fmt string, fields ...zap.Field) { l.zCore(lvlI, 1).Info(fmt, fields...) }

func (l *Wrapper) Printf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlI, 1).Infof(fmt, args...)
}
func (l *Wrapper) Print(args ...interface{})              { l.zSugar(lvlI, 1).Info(pw(args)) }
func (l *Wrapper) Println(args ...interface{})            { l.zSugar(lvlI, 1).Info(pw(args)) }
func (l *Wrapper) Printz(fmt string, fields ...zap.Field) { l.zCore(lvlI, 1).Info(fmt, fields...) }

func (l *Wrapper) Logf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlI, 1).Infof(fmt, args...)
}
func (l *Wrapper) Log(args ...interface{})              { l.zSugar(lvlI, 1).Info(pw(args)) }
func (l *Wrapper) Logln(args ...interface{})            { l.zSugar(lvlI, 1).Info(pw(args)) }
func (l *Wrapper) Logz(fmt string, fields ...zap.Field) { l.zCore(lvlI, 1).Info(fmt, fields...) }

func (l *Wrapper) Warnf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlW, 1).Warnf(fmt, args...)
}
func (l *Wrapper) Warn(args ...interface{})              { l.zSugar(lvlW, 1).Warn(pw(args)) }
func (l *Wrapper) Warnln(args ...interface{})            { l.zSugar(lvlW, 1).Warn(pw(args)) }
func (l *Wrapper) Warnz(fmt string, fields ...zap.Field) { l.zCore(lvlW, 1).Warn(fmt, fields...) }

func (l *Wrapper) Warningf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlW, 1).Warnf(fmt, args...)
}
func (l *Wrapper) Warning(args ...interface{})              { l.zSugar(lvlW, 1).Warn(pw(args)) }
func (l *Wrapper) Warningln(args ...interface{})            { l.zSugar(lvlW, 1).Warn(pw(args)) }
func (l *Wrapper) Warningz(fmt string, fields ...zap.Field) { l.zCore(lvlW, 1).Warn(fmt, fields...) }

func (l *Wrapper) Errorf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlE, 1).Errorf(fmt, args...)
}
func (l *Wrapper) Error(args ...interface{})              { l.zSugar(lvlE, 1).Error(pw(args)) }
func (l *Wrapper) Errorln(args ...interface{})            { l.zSugar(lvlE, 1).Error(pw(args)) }
func (l *Wrapper) Errorz(fmt string, fields ...zap.Field) { l.zCore(lvlE, 1).Error(fmt, fields...) }

func (l *Wrapper) Panicf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlP, 1).Panicf(fmt, args...)
}
func (l *Wrapper) Panic(args ...interface{})              { l.zSugar(lvlP, 1).Panic(pw(args)) }
func (l *Wrapper) Panicln(args ...interface{})            { l.zSugar(lvlP, 1).Panic(pw(args)) }
func (l *Wrapper) Panicz(fmt string, fields ...zap.Field) { l.zCore(lvlP, 1).Panic(fmt, fields...) }

func (l *Wrapper) Fatalf(fmt string, args ...interface{}) {
	l.format(fmt, args...)
	l.zSugar(lvlF, 1).Fatalf(fmt, args...)
}
func (l *Wrapper) Fatal(args ...interface{})              { l.zSugar(lvlF, 1).Fatal(pw(args)) }
func (l *Wrapper) Fatalln(args ...interface{})            { l.zSugar(lvlF, 1).Fatal(pw(args)) }
func (l *Wrapper) Fatalz(fmt string, fields ...zap.Field) { l.zCore(lvlF, 1).Fatal(fmt, fields...) }

func (l *Wrapper) V(v int) bool {
	if v <= 0 {
		return !l.Enabled(zapcore.DebugLevel)
	}
	return true
}
