/*
 * Copyright (C) distroy
 */

package internal

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Options struct {
	Reporter      Reporter
	Retry         int
	RetryInterval time.Duration
	Caller        bool
}

const (
	ldRedisSrcPath = "/ldredis"
	exampleSrcPath = "/ldredis/example"
	goRedisSrcPath = "/github.com/redis/go-redis"
)

func isCallerFilePath(file string) bool {
	if strings.Contains(file, goRedisSrcPath) {
		return false
	}
	if !strings.Contains(file, ldRedisSrcPath) {
		return true
	}
	if strings.HasSuffix(file, "_test.go") {
		return true
	}
	if strings.Contains(file, exampleSrcPath) {
		return true
	}
	return false
}

func GetCallerField(caller bool) zap.Field {
	if !caller {
		return zap.Skip()
	}
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if !isCallerFilePath(file) {
			continue
		}
		return zap.String("caller", fmt.Sprintf("%s:%d", file, line))
	}

	return zap.String("caller", "overflow")
}
