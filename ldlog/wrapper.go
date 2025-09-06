/*
 * Copyright (C) distroy
 */

package ldlog

type Wrapper struct {
	core
}

func (l *Wrapper) Logger() *Logger { return (*Logger)(l) }

func (l *Wrapper) Debugf(fmt string, args ...any)   { l.logFmt(nil, lvlD, 1, fmt, args...) }
func (l *Wrapper) Debug(args ...any)                { l.logln(nil, lvlD, 1, args) }
func (l *Wrapper) Debugln(args ...any)              { l.logln(nil, lvlD, 1, args) }
func (l *Wrapper) Debugz(msg string, attrs ...Attr) { l.logAttrs(nil, lvlD, 1, msg, attrs...) }

func (l *Wrapper) Infof(fmt string, args ...any)   { l.logFmt(nil, lvlI, 1, fmt, args...) }
func (l *Wrapper) Info(args ...any)                { l.logln(nil, lvlI, 1, args) }
func (l *Wrapper) Infoln(args ...any)              { l.logln(nil, lvlI, 1, args) }
func (l *Wrapper) Infoz(msg string, attrs ...Attr) { l.logAttrs(nil, lvlI, 1, msg, attrs...) }

func (l *Wrapper) Printf(fmt string, args ...any)   { l.logFmt(nil, lvlI, 1, fmt, args...) }
func (l *Wrapper) Print(args ...any)                { l.logln(nil, lvlI, 1, args) }
func (l *Wrapper) Println(args ...any)              { l.logln(nil, lvlI, 1, args) }
func (l *Wrapper) Printz(msg string, attrs ...Attr) { l.logAttrs(nil, lvlI, 1, msg, attrs...) }

func (l *Wrapper) Logf(fmt string, args ...any)   { l.logFmt(nil, lvlI, 1, fmt, args...) }
func (l *Wrapper) Log(args ...any)                { l.logln(nil, lvlI, 1, args) }
func (l *Wrapper) Logln(args ...any)              { l.logln(nil, lvlI, 1, args) }
func (l *Wrapper) Logz(msg string, attrs ...Attr) { l.logAttrs(nil, lvlI, 1, msg, attrs...) }

func (l *Wrapper) Warnf(fmt string, args ...any)   { l.logFmt(nil, lvlW, 1, fmt, args...) }
func (l *Wrapper) Warn(args ...any)                { l.logln(nil, lvlW, 1, args) }
func (l *Wrapper) Warnln(args ...any)              { l.logln(nil, lvlW, 1, args) }
func (l *Wrapper) Warnz(msg string, attrs ...Attr) { l.logAttrs(nil, lvlW, 1, msg, attrs...) }

func (l *Wrapper) Warningf(fmt string, args ...any)   { l.logFmt(nil, lvlW, 1, fmt, args...) }
func (l *Wrapper) Warning(args ...any)                { l.logln(nil, lvlW, 1, args) }
func (l *Wrapper) Warningln(args ...any)              { l.logln(nil, lvlW, 1, args) }
func (l *Wrapper) Warningz(msg string, attrs ...Attr) { l.logAttrs(nil, lvlW, 1, msg, attrs...) }

func (l *Wrapper) Errorf(fmt string, args ...any)   { l.logFmt(nil, lvlE, 1, fmt, args...) }
func (l *Wrapper) Error(args ...any)                { l.logln(nil, lvlE, 1, args) }
func (l *Wrapper) Errorln(args ...any)              { l.logln(nil, lvlE, 1, args) }
func (l *Wrapper) Errorz(msg string, attrs ...Attr) { l.logAttrs(nil, lvlE, 1, msg, attrs...) }

func (l *Wrapper) Panicf(fmt string, args ...any)   { l.logFmt(nil, lvlP, 1, fmt, args...) }
func (l *Wrapper) Panic(args ...any)                { l.logln(nil, lvlP, 1, args) }
func (l *Wrapper) Panicln(args ...any)              { l.logln(nil, lvlP, 1, args) }
func (l *Wrapper) Panicz(msg string, attrs ...Attr) { l.logAttrs(nil, lvlP, 1, msg, attrs...) }

func (l *Wrapper) V(v int) bool {
	if v <= 0 {
		return !l.Enabled(nil, LevelDebug)
	}
	return true
}
