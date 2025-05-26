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

type (
	ch = commonHandler
)

func NewHandler(w io.Writer, opts *Options) Handler {
	if opts == nil {
		opts = &Options{}
	}
	return Handler{
		&commonHandler{
			json: false,
			w:    w,
			opts: *opts,
			mu:   &sync.Mutex{},
		},
	}
}

type Handler struct {
	*ch
}

func (h Handler) GetSequence() string  { return h.seqId }
func (h Handler) SetSequence(s string) { h.seqId = s }

func (h Handler) Enabled(c context.Context, lvl slog.Level) bool  { return h.enabled(lvl) }
func (h Handler) Handle(c context.Context, rec slog.Record) error { return h.handle(c, GetRecord(rec)) }
func (h Handler) WithAttrs(as []slog.Attr) slog.Handler           { return Handler{h.withAttrs(GetAttrs(as))} }
func (h Handler) WithGroup(name string) slog.Handler              { return Handler{h.withGroup(name)} }

func (h Handler) Sync() error {
	w := h.ch.w
	switch ww := w.(type) {
	case interface{ Sync() error }:
		return ww.Sync()

	case interface{ Sync() }:
		ww.Sync()
	}
	return nil
}

type commonHandler struct {
	json              bool // true => output JSON; false => output text
	opts              Options
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
		mu:                h.mu, // mutex shared among all clones of this handler
		w:                 h.w,
		seqId:             h.seqId,
	}
}

func (h *commonHandler) handle(_ context.Context, r Record) error {
	seqId := h.seqId
	if seqId == "" {
		seqId = "-"
	}

	state := h.newHandleState(newBuffer(), true, ",")
	defer state.free()

	buf := state.buf

	buf.WriteTime(r.Time, "")
	buf.WriteByte('|')
	buf.WriteString(r.Level.String())
	buf.WriteByte('|')
	buf.WriteString(seqId)
	buf.WriteByte('|')
	if !h.opts.AddSource {
		buf.WriteByte('-')
	} else {
		src := r.Source()
		buf.WriteString(src.Caller())
	}
	buf.WriteByte('|')
	buf.WriteString(r.Message)
	if pfa := h.preformattedAttrs; len(pfa) > 0 {
		buf.WriteByte(',')
		buf.Write(pfa)
	}

	if r.NumAttrs() > 0 {
		buf.WriteByte('|')
	}

	state.sep = ""
	state.appendNonBuiltIns(r)
	buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*state.buf)
	return err
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

func (h *commonHandler) withAttrs(as []Attr) *commonHandler {
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
func (h *commonHandler) attrSep() string { return "," }

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
