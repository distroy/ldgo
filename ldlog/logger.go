/*
 * Copyright (C) distroy
 */

package ldlog

const (
	lvlD = LevelDebug
	lvlI = LevelInfo
	lvlW = LevelWarn
	lvlE = LevelError
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

func (l *Logger) With(attrs ...Attr) *Logger {
	if len(attrs) == 0 {
		return l
	}
	l = l.clone()
	l.handler = wrapHandler(l.handler.WithAttrs(attrs))
	return l
}

func (l *Logger) WithOptions(opts ...Option) *Logger {
	c := l.core
	for _, opt := range opts {
		opt(&c)
	}
	if c == l.core {
		return l
	}
	l = l.clone()
	l.core = c
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

func (l *Logger) Debugln(args ...any) { l.logln(nil, lvlD, 1, args) }
func (l *Logger) Infoln(args ...any)  { l.logln(nil, lvlD, 1, args) }
func (l *Logger) Warnln(args ...any)  { l.logln(nil, lvlD, 1, args) }
func (l *Logger) Errorln(args ...any) { l.logln(nil, lvlD, 1, args) }
func (l *Logger) Panicln(args ...any) { l.logln(nil, lvlD, 1, args) }
