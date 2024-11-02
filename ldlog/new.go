/*
 * Copyright (C) distroy
 */

package ldlog

import (
	"github.com/distroy/ldgo/v3/ldconv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(log *zap.Logger, fields ...zap.Field) *Logger {
	if len(fields) > 0 {
		log = log.With(fields...)
	}
	return newLogger(newCore(log, log.Sugar()))
}

func New(opts ...Option) *Logger {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	var level zapcore.Level

	if err := level.UnmarshalText(ldconv.StrToBytesUnsafe(options.level)); err != nil {
		level.UnmarshalText(ldconv.StrToBytesUnsafe("info")) // nolint
	}

	encoder := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	handle := zapcore.Lock(options.writer)
	zapCores := []zapcore.Core{
		// zapcore.NewCore(zapcore.NewConsoleEncoder(encoder), handle, level),
		zapcore.NewCore(options.encoderBuilder(encoder), handle, level),
	}

	// create options with priority for our opts
	defaultOptions := []zap.Option{}
	defaultOptions = append(defaultOptions, zap.AddCallerSkip(1))
	if options.enableCaller {
		// defaultOptions = append(defaultOptions, zap.AddStacktrace(level))
		defaultOptions = append(defaultOptions, zap.AddCaller())
	}

	core := zapcore.NewTee(zapCores...)
	zlog := zap.New(core, defaultOptions...)

	return newLogger(newCore(zlog, zlog.Sugar()))
}
