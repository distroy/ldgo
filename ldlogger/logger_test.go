/*
 * Copyright (C) distroy
 */

package ldlogger

import (
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {
	l := NewLogger().With(zap.String("abc", "xxx"))
	l.Error("error message")
	l.Warn("warn message")
	l.Info("info message")
	l.Infoln("infoln message", "abc")
}
