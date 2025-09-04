/*
 * Copyright (C) distroy
 */

package lderr

import (
	"context"
	"io"
	"testing"
)

type testUnwrapError struct {
	error
}

func (e testUnwrapError) Unwrap() error { return e.error }

func TestIs(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   bool
	}{
		{
			name:   "nil & nil",
			err:    nil,
			target: nil,
			want:   true,
		},
		{
			name:   "nil & succ",
			err:    nil,
			target: ErrSuccess,
			want:   true,
		},
		{
			name:   "succ & nil",
			err:    ErrSuccess,
			target: nil,
			want:   true,
		},
		{
			name:   "succ & succ",
			err:    ErrSuccess,
			target: ErrSuccess,
			want:   true,
		},
		{
			name:   "nil & unknown",
			err:    nil,
			target: ErrUnkown,
			want:   false,
		},
		{
			name:   "succ & unknown",
			err:    ErrSuccess,
			target: ErrUnkown,
			want:   false,
		},
		{
			name:   "context canceled",
			err:    ErrCtxCanceled,
			target: context.Canceled,
			want:   true,
		},
		{
			name:   "context deadline exceeded",
			err:    ErrCtxDeadlineExceeded,
			target: context.DeadlineExceeded,
			want:   true,
		},
		{
			name:   "iof & wrap iof",
			err:    io.EOF,
			target: Wrap(io.EOF),
			want:   false,
		},
		{
			name:   "wrap iof & iof",
			err:    Wrap(io.EOF),
			target: io.EOF,
			want:   true,
		},
		{
			name:   "wrap iof & unknown",
			err:    Wrap(io.EOF),
			target: io.EOF,
			want:   true,
		},
		{
			name:   "wrap iof & succ",
			err:    Wrap(io.EOF),
			target: ErrSuccess,
			want:   false,
		},
		{
			name:   "unknown & overwrite unknown",
			err:    ErrUnkown,
			target: Override(ErrUnkown, "abc"),
			want:   true,
		},
		{
			name:   "unknown & string error",
			err:    ErrUnkown,
			target: strError(errMessageUnkonw),
			want:   true,
		},
		{
			name:   "unknown & panic",
			err:    testUnwrapError{ErrUnkown},
			target: ErrServicePanic,
			want:   false,
		},
		{
			name:   "unknown & nil",
			err:    testUnwrapError{ErrUnkown},
			target: nil,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.err, tt.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
