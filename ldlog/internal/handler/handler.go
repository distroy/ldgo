/*
 * Copyright (C) distroy
 */

package handler

import (
	"context"
	"io"
	"log/slog"
	"slices"
	"sync"
)

var (
	SequenceKey = "request_id"
)

var (
	_ slog.Handler = (*Handler)(nil)
)

type Handler struct {
	ch *commonHandler
}

func (h Handler) Enabled(c context.Context, lvl slog.Level) bool  { return h.enabled(c, lvl) }
func (h Handler) Handle(c context.Context, rec slog.Record) error { return h.handle(c, rec) }
func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler        { return h.withAttrs(attrs) }
func (h Handler) WithGroup(name string) slog.Handler              { return h.withGroup(name) }

func (h Handler) withAttrs(attrs []slog.Attr) Handler            { return Handler{h.ch.withAttrs(attrs)} }
func (h Handler) withGroup(name string) Handler                  { return Handler{h.ch.withGroup(name)} }
func (h Handler) enabled(c context.Context, lvl slog.Level) bool { return h.ch.enabled(lvl) }

func (h Handler) handle(c context.Context, r slog.Record) error {
	seqId := h.ch.seqId
	if seqId == "" {
		seqId = "-"
	}

	buf := newBuffer()
	buf.WriteTime(r.Time, "")
	buf.WriteByte('|')
	buf.WriteString(r.Level.String())
	buf.WriteByte('|')
	buf.WriteString(seqId)
	buf.WriteByte('|')
	buf.WriteString(r.Message)
	if pfa := h.ch.preformattedAttrs; len(pfa) > 0 {
		buf.WriteByte(',')
		buf.Write(pfa)
	}

	buf.WriteByte('|')

	state := h.ch.newHandleState(buf, true, ",")
	defer state.free()
	state.appendNonBuiltIns(r)
	return h.ch.handle(r)
}

type commonHandler struct {
	json              bool // true => output JSON; false => output text
	opts              slog.HandlerOptions
	preformattedAttrs []byte
	// groupPrefix is for the text handler only.
	// It holds the prefix for groups that were already pre-formatted.
	// A group will appear here when a call to WithGroup is followed by
	// a call to WithAttrs.
	groupPrefix string
	groups      []string // all groups started from WithGroup
	nOpenGroups int      // the number of groups opened in preformattedAttrs
	mu          *sync.Mutex
	w           io.Writer
	seqId       string
}

func (h *commonHandler) clone() *commonHandler {
	// We can't use assignment because we can't copy the mutex.
	return &commonHandler{
		json:              h.json,
		opts:              h.opts,
		preformattedAttrs: slices.Clip(h.preformattedAttrs),
		groupPrefix:       h.groupPrefix,
		groups:            slices.Clip(h.groups),
		nOpenGroups:       h.nOpenGroups,
		w:                 h.w,
		mu:                h.mu, // mutex shared among all clones of this handler
	}
}

// enabled reports whether l is greater than or equal to the
// minimum level.
func (h *commonHandler) enabled(l slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return l >= minLevel
}

func (h *commonHandler) withAttrs(as []slog.Attr) *commonHandler {
	// We are going to ignore empty groups, so if the entire slice consists of
	// them, there is nothing to do.
	if countEmptyGroups(as) == len(as) {
		return h
	}
	h2 := h.clone()
	// Pre-format the attributes as an optimization.
	state := h2.newHandleState((*Buffer)(&h2.preformattedAttrs), false, "")
	defer state.free()
	state.prefix.WriteString(h.groupPrefix)
	if pfa := h2.preformattedAttrs; len(pfa) > 0 {
		state.sep = h.attrSep()
		if h2.json && pfa[len(pfa)-1] == '{' {
			state.sep = ""
		}
	}
	// Remember the position in the buffer, in case all attrs are empty.
	pos := state.buf.Len()
	state.openGroups()
	if !state.appendAttrs(as) {
		state.buf.SetLen(pos)
	} else {
		// Remember the new prefix for later keys.
		h2.groupPrefix = state.prefix.String()
		// Remember how many opened groups are in preformattedAttrs,
		// so we don't open them again when we handle a Record.
		h2.nOpenGroups = len(h2.groups)
	}
	return h2
}

func (h *commonHandler) withGroup(name string) *commonHandler {
	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	return h2
}

// attrSep returns the separator between attributes.
func (h *commonHandler) attrSep() string {
	return ","
}

func (h *commonHandler) newHandleState(buf *Buffer, freeBuf bool, sep string) handleState {
	s := handleState{
		h:       h,
		buf:     buf,
		freeBuf: freeBuf,
		sep:     sep,
		prefix:  newBuffer(),
	}
	if h.opts.ReplaceAttr != nil {
		gs := groupPool.Get()
		s.groups = &gs
		*s.groups = append(*s.groups, h.groups[:h.nOpenGroups]...)
	}
	return s
}

// handle is the internal implementation of Handler.Handle
// used by TextHandler and JSONHandler.
func (h *commonHandler) handle(r slog.Record) error {
	state := h.newHandleState(newBuffer(), true, "")
	defer state.free()
	if h.json {
		state.buf.WriteByte('{')
	}
	// Built-in attributes. They are not in a group.
	stateGroups := state.groups
	state.groups = nil // So ReplaceAttrs sees no groups instead of the pre groups.
	rep := h.opts.ReplaceAttr
	// time
	if !r.Time.IsZero() {
		key := slog.TimeKey
		val := r.Time.Round(0) // strip monotonic to match Attr behavior
		if rep == nil {
			state.appendKey(key)
			state.appendTime(val)
		} else {
			state.appendAttr(slog.Time(key, val))
		}
	}
	// level
	key := slog.LevelKey
	val := r.Level
	if rep == nil {
		state.appendKey(key)
		state.appendString(val.String())
	} else {
		state.appendAttr(slog.Any(key, val))
	}
	// source
	if h.opts.AddSource {
		state.appendAttr(slog.Any(slog.SourceKey, getRecord(&r).source()))
	}
	key = slog.MessageKey
	msg := r.Message
	if rep == nil {
		state.appendKey(key)
		state.appendString(msg)
	} else {
		state.appendAttr(slog.String(key, msg))
	}
	state.groups = stateGroups // Restore groups passed to ReplaceAttrs.
	state.appendNonBuiltIns(r)
	state.buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*state.buf)
	return err
}
