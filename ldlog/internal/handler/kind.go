/*
 * Copyright (C) distroy
 */

package handler

import "log/slog"

type Kind = slog.Kind

const (
	KindAny Kind = iota
	KindBool
	KindDuration
	KindFloat64
	KindInt64
	KindString
	KindTime
	KindUint64
	KindGroup
	KindLogValuer
)
