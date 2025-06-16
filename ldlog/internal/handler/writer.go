/*
 * Copyright (C) distroy
 */

package handler

import (
	"io"
	"os"
)

type logWriter interface {
	io.WriteCloser

	Sync() error
}

type writeSyncer interface {
	io.Writer

	Sync() error
}

func wrapWriter(w io.Writer) logWriter {
	switch ww := w.(type) {
	case *os.File:
		if ww == os.Stdin || ww == os.Stdout || ww == os.Stderr {
			return writeSyncerWrapper{ww}
		}
		return ww

	case logWriter:
		return ww

	case io.WriteCloser:
		return writeCloserWrapper{ww}

	case writeSyncer:
		return writeSyncerWrapper{ww}
	}
	return writerWrapper{w}
}

type writerWrapper struct {
	io.Writer
}

func (w writerWrapper) Unwrap() io.Writer { return w.Writer }
func (w writerWrapper) Sync() error       { return nil }
func (w writerWrapper) Close() error      { return nil }

type writeSyncerWrapper struct {
	writeSyncer
}

func (w writeSyncerWrapper) Unwrap() io.Writer { return w.writeSyncer }
func (w writeSyncerWrapper) Close() error      { return nil }

type writeCloserWrapper struct {
	io.WriteCloser
}

func (w writeCloserWrapper) Unwrap() io.Writer { return w.WriteCloser }
func (w writeCloserWrapper) Sync() error       { return nil }
