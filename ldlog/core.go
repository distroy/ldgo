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
		handler: h,
		enabler: defaultEnabler{},
	}
}

type core struct {
	handler   Handler
	enabler   Enabler
	stackSkip int
}

func (l *core) Enabler() Enabler { return l.enabler }
func (l *core) Handler() Handler { return l.handler }

func (l *core) Sync() error {
	switch h := l.handler.(type) {
	case interface{ Sync() error }:
		return h.Sync()
	}
	return nil
}

func (l *core) close() error {
	switch h := l.handler.(type) {
	case interface{ Close() error }:
		return h.Close()
	}
	return nil
}

func (l *core) Sequence() string {
	if h, _ := l.handler.(interface{ Sequence() string }); h != nil {
		return h.Sequence()
	}
	return ""
}

func (l *core) withSequence(seq string) *core { return l.withAttr(String(GetSequenceKey(), seq)) }
func (l *core) withLevel(lvl Level) *core     { return l.withAttr(Int(GetLevelKey(), int(lvl))) }
func (l *core) withEnabler(e Enabler) *core   { return l.clone(func(l *core) { l.enabler = e }) }
func (l *core) withAddStackSkip(delta int) *core {
	return l.clone(func(l *core) { l.stackSkip += delta })
}

func (l *core) withAttr(attrs ...Attr) *core {
	return l.clone(func(l *core) { l.handler = l.handler.WithAttrs(attrs) })
}

func (l *core) clone(fs ...func(l *core)) *core {
	cp := *l
	for _, f := range fs {
		f(&cp)
	}
	return &cp
}

func (l *core) ctx(c context.Context) context.Context {
	if c == nil {
		c = context.Background()
	}
	return c
}

func (l *core) Enabled(c context.Context, lvl Level) bool { return l.enabled(l.ctx(c), lvl, 1) }
func (l *core) enabled(c context.Context, lvl Level, skip int) bool {
	if !l.handler.Enabled(c, lvl) {
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
	runtime.Callers(skip+1+l.stackSkip, pcs[:])
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
	if l == nil || !l.enabled(c, lvl, skip+1) {
		return
	}
	pc := l.getCaller(skip + 1)
	r := slog.NewRecord(time.Now(), lvl, msg, pc)
	r.Add(args...)
	l.writeRecord(c, lvl, &r)
}

func (l *core) logFmt(c context.Context, lvl Level, skip int, format string, args ...any) {
	c = l.ctx(c)
	if l == nil || !l.enabled(c, lvl, skip+1) {
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
	if l == nil || !l.enabled(c, lvl, skip+1) {
		return
	}
	pc := l.getCaller(skip + 1)
	r := slog.NewRecord(time.Now(), lvl, msg, pc)
	r.AddAttrs(attrs...)
	l.writeRecord(c, lvl, &r)
}
