/*
 * Copyright (C) distroy
 */

package handler

import (
	"bytes"
	"fmt"
	"log/slog"
	"testing"
)

func testSLogPrint(log *slog.Logger) {
	log = log.With(slog.Int("int", 123))
	log = log.WithGroup("g")
	log.Debug("test", slog.String("str", "abc"), slog.Group("g1", slog.Int("int1", 234)),
		slog.Group("", slog.String("str1", "xyz")))
}

func TestHandler(t *testing.T) {
	opts := &Options{AddSource: true, Level: slog.LevelDebug}
	buf := bytes.NewBuffer(nil)
	log := slog.New(NewHandler(buf, opts))
	testSLogPrint(log)
	fmt.Printf("%s\n", buf.Bytes())
}

func TestSLog(t *testing.T) {
	opts := &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}
	buf := bytes.NewBuffer(nil)
	log := slog.New(slog.NewTextHandler(buf, opts))
	testSLogPrint(log)
	fmt.Printf("%s\n", buf.Bytes())
}
