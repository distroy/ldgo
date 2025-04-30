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
	log = log.WithGroup("ga")
	log.Debug("test", slog.String("str", "abc"))
}

func TestHandler(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	log := slog.New(NewHandler(buf, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	testSLogPrint(log)
	fmt.Printf("%s\n", buf.Bytes())
}

func TestSLog(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	log := slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	testSLogPrint(log)
	fmt.Printf("%s\n", buf.Bytes())
}
