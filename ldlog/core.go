/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"
)

const (
	formatFlag = false
)

func newCore(h Handler) core {
	return core{
		handler: wrapHandler(h),
		enabler: defaultEnabler{},
	}
}

type core struct {
	handler   logHandler
	enabler   Enabler
	stackSkip int
}

func (l *core) Enabler() Enabler { return l.enabler }
func (l *core) Handler() Handler { return l.handler }

func (l *core) Sync() error  { return l.handler.Sync() }
func (l *core) close() error { return l.handler.Close() }

func (l *core) Level() Level     { return l.handler.Level() }
func (l *core) Sequence() string { return l.handler.Sequence() }

func (l *core) setAttr(attr ...Attr) { l.handler = wrapHandler(l.handler.WithAttrs(attr)) }

func (l *core) ctx(c context.Context) context.Context {
	if c == nil {
		c = context.Background()
	}
	return c
}

func (l *core) Enabled(c context.Context, lvl Level) bool { return l.enabled(l.ctx(c), lvl, 1) }
func (l *core) enabled(c context.Context, lvl Level, skip int) bool {
	if l == nil || l.handler == nil || !l.handler.Enabled(c, lvl) {
		return false
	}
	return l.enabler.Enable(lvl, skip+1)
}

func (l *core) format(format string, args ...any) {
	if formatFlag {
		_ = fmt.Sprintf(format, args...)
	}
}

func (l *core) getCaller(skip int) uintptr {
	var pcs [1]uintptr
	// skip [runtime.Callers, this function, this function's caller]
	runtime.Callers(skip+2+l.stackSkip, pcs[:])
	return pcs[0]
}

func (l *core) writeRecord(c context.Context, lvl Level, r *Record) {
	c = l.ctx(c)
	_ = l.Handler().Handle(c, *r)
	// if lvl < LevelPanic {
	// 	return
	// }
	// l.Sync()
	// panic(rec2err(r))
}

func (l *core) log(c context.Context, lvl Level, skip int, msg string, args ...any) {
	c = l.ctx(c)
	if !l.enabled(c, lvl, skip+1) {
		return
	}
	pc := l.getCaller(skip + 1)
	r := slog.NewRecord(time.Now(), lvl, msg, pc)
	r.Add(args...)
	l.writeRecord(c, lvl, &r)
}

func (l *core) logFmt(c context.Context, lvl Level, skip int, format string, args ...any) {
	c = l.ctx(c)
	if !l.enabled(c, lvl, skip+1) {
		return
	}
	pc := l.getCaller(skip + 1)
	msg := fmt.Sprintf(format, args...)
	r := slog.NewRecord(time.Now(), lvl, msg, pc)
	// r.Add(args...)
	l.writeRecord(c, lvl, &r)
}

func (l *core) logAttrs(c context.Context, lvl Level, skip int, msg string, attrs ...Attr) {
	c = l.ctx(c)
	if !l.enabled(c, lvl, skip+1) {
		return
	}
	pc := l.getCaller(skip + 1)
	r := slog.NewRecord(time.Now(), lvl, msg, pc)
	r.AddAttrs(attrs...)
	l.writeRecord(c, lvl, &r)
}

func (l *core) logln(c context.Context, lvl Level, skip int, args ...any) {
	c = l.ctx(c)
	if !l.enabled(c, lvl, skip+1) {
		return
	}
	pc := l.getCaller(skip + 1)
	msg := sprintln(args)
	r := slog.NewRecord(time.Now(), lvl, msg, pc)
	// r.Add(args...)
	l.writeRecord(c, lvl, &r)
}
