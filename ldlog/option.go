/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"github.com/distroy/ldgo/v3/ldlog/internal/_handler"
)

const (
	defaultLogLevel        = "INFO"
	defaultLogEnableCaller = true
)

func GetLevelKey() string  { return _handler.LevelKey }
func GetCallerKey() string { return _handler.CallerKey }

func SetSequenceKey(key string) { _handler.SequenceKey = key }
func GetSequenceKey() string    { return _handler.SequenceKey }

type Option func(l *core)

func SetLevel(lvl Level) Option   { return func(l *core) { l.withAttrs(Any(GetLevelKey(), lvl)) } }
func SetEnabler(e Enabler) Option { return func(l *core) { l.enabler = e } }
func SetSequence(s string) Option { return func(l *core) { l.withAttrs(String(GetSequenceKey(), s)) } }

func EnableCaller(e bool) Option    { return func(l *core) { l.withAttrs(Bool(GetCallerKey(), e)) } }
func AddStackSkip(delta int) Option { return func(l *core) { l.stackSkip += delta } }
