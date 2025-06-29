/*
 * Copyright (C) distroy
 */

package main

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/ldlog"
	"github.com/distroy/ldgo/v3/ldredis"
	"github.com/distroy/ldgo/v3/ldredis/ldrediscodec"
)

func init() {
	log := ldlog.New(nil, ldlog.SetLevel(ldlog.LevelDebug))
	ldlog.SetDefault(log)
}

func newRedis(ctx ldctx.Context) *ldredis.Redis {
	return ldredis.MustNewTestRedis()
}

func pipeline(ctx ldctx.Context) {
	rds := newRedis(ctx)
	rds = rds.WithRetry(3)
	// rds = rds.WithCaller(false)
	defer rds.Close()

	keys := []string{
		"test:pipeline:0",
		"test:pipeline:1",
		"test:pipeline:2",
	}

	rds.MSet(ctx, keys[0], "111", keys[1], "aaa", keys[2], "xxx")

	p := rds.Pipeline()
	p.Get(ctx, keys[0])
	p.Get(ctx, keys[1])
	p.Get(ctx, keys[2])
	cmds, err := p.Exec(ctx)
	ldctx.LogI(ctx, "pipeline return", ldlog.Error(err))
	for _, v := range cmds {
		cmd, _ := v.(*ldredis.StringCmd)
		ldctx.LogI(ctx, "pipeline res", ldlog.Reflect("cmd", cmd.Args()), ldlog.String("val", cmd.Val()), ldlog.Error(cmd.Err()))
	}
}

func slice(ctx ldctx.Context) {
	rds := newRedis(ctx)
	key := "test:hash:hmget"
	defer rds.Close()

	rds1 := rds.WithReport(LogReporter{})

	rds1.HSet(ctx, key, "1", 1)
	rds.HSet(ctx, key, "2", "abc")
	rds.HSet(ctx, key, "3", 128.1)

	cmd := rds.HMGet(ctx, key, "1", "2", "3", "4")
	ldctx.LogI(ctx, "", ldlog.Stringer("type", reflect.TypeOf(cmd.Val())), ldlog.Reflect("value", cmd.Val()))
	for i, v := range cmd.Val() {
		ldctx.LogIf(ctx, "idx:%d, type:%T, value:%v", i, v, v)
	}
}

func codecStruct(ctx ldctx.Context) {
	type codecStruct struct {
		Str1 string `json:"str1"`
		Str2 string `json:"str2"`
		Int1 int64  `json:"int1"`
		Int2 int64  `json:"int2"`
	}

	rds := newRedis(ctx)
	defer rds.Close()
	key := "test:codec:struct"

	sCmd := ldrediscodec.New(rds, ldrediscodec.Json[any]()).Set(ctx, key, &codecStruct{
		Str1: "aaa",
		Str2: "bbb",
		Int1: 111,
		Int2: 222,
	}, time.Minute)
	ldctx.LogI(ctx, "cmd", ldlog.Reflect("cmd", sCmd.Args()))

	gCmd0 := ldrediscodec.New(rds, ldrediscodec.Json[*codecStruct]()).Get(ctx, key)
	ldctx.LogIf(ctx, "type:%T, value:%v", gCmd0.Val(), gCmd0.Val())

	gCmd1 := ldrediscodec.New(rds, ldrediscodec.Json[any]()).Get(ctx, key)
	ldctx.LogIf(ctx, "type:%T, value:%v", gCmd1.Val(), gCmd1.Val())
}

func codecBaseType(ctx ldctx.Context) {
	rds := newRedis(ctx)
	defer rds.Close()

	key := "test:codec:basetype"
	cli := ldrediscodec.New(rds, ldrediscodec.Json[any]())
	cli.HMSetMap(ctx, key, map[string]interface{}{
		"i1": 1234,
		"s1": "abc",
		"s2": "134",
		"m1": map[int]interface{}{
			1: 1,
			2: "a",
		},
	})

	cmd := cli.HGetAll(ctx, key)
	ldctx.LogI(ctx, "", ldlog.Reflect("cmd", cmd.Args()), ldlog.Stringer("type", reflect.TypeOf(cmd.Val())), ldlog.Reflect("val", cmd.Val()))
}

func main() {
	ctx := ldctx.Default()
	// pipeline(ctx)
	// fmt.Fprintln(os.Stderr)
	slice(ctx)
	fmt.Fprintln(os.Stderr)
	codecStruct(ctx)
	fmt.Fprintln(os.Stderr)
	codecBaseType(ctx)
	fmt.Fprintln(os.Stderr)
}

type LogReporter struct{}

func (_ LogReporter) Report(cmd ldredis.Cmder, d time.Duration) {
	ldctx.LogI(nil, "report redis cmd", ldlog.Reflect("cmd", cmd.Args()))
}

func (_ LogReporter) ReportPipeline(cmds []ldredis.Cmder, d time.Duration) {
	ldctx.LogI(nil, "report redis pipline cmd")
	for i, cmd := range cmds {
		ldctx.LogI(nil, "report redis pipline cmd", ldlog.Int("idx", i), ldlog.Reflect("cmd", cmd.Args()))
	}
}
