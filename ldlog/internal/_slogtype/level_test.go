/*
 * Copyright (C) distroy
 */

package _slogtype

import (
	"fmt"
	"testing"
)

func TestLevel_Values(t *testing.T) {
	ifPanic := func(ok bool) {
		if ok {
			return
		}
		panic(ok)
	}

	ifPanic(LevelTrace < LevelDebug)
	ifPanic(LevelDebug < LevelInfo)
	ifPanic(LevelInfo < LevelWarn)
	ifPanic(LevelWarn < LevelError)
	ifPanic(LevelError < LevelPanic)
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		// name string
		l    Level
		want string
	}{
		{LevelTrace - 100, "TRACE"},
		{LevelTrace, "TRACE"},
		{LevelTrace + 3, "TRACE+3"},
		{LevelDebug, "DEBUG"},
		{LevelDebug + 3, "DEBUG+3"},
		{LevelInfo, "INFO"},
		{LevelInfo + 3, "INFO+3"},
		{LevelError, "ERROR"},
		{LevelError + 3, "ERROR+3"},
		{LevelPanic, "PANIC"},
		{LevelPanic + 3, "PANIC+3"},
	}
	for i, tt := range tests {
		name := fmt.Sprintf("%d:%s", i, tt.want)
		t.Run(name, func(t *testing.T) {
			if got := tt.l.String(); got != tt.want {
				t.Errorf("Level.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
