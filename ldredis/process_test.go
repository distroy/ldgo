/*
 * Copyright (C) distroy
 */

package ldredis

import (
	"context"
	"testing"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/lderr"
	redis "github.com/redis/go-redis/v9"
	"github.com/smartystreets/goconvey/convey"
)

type testStringCmd struct {
	*StringCmd
	err error
}

func (cmd *testStringCmd) Parse(c context.Context) error {
	return cmd.err
}

func Test_hook_process(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := MustNewTestRedis()
		defer rds.Close()

		key := "test-key"
		ctx := ldctx.Default()

		cmd := rds.Set(ctx, key, "test-value", 0)
		c.So(cmd.Err(), convey.ShouldBeNil)
		c.So(cmd.Val(), convey.ShouldNotBeEmpty)

		c.Convey("parse error", func(c convey.C) {
			cmd := &testStringCmd{
				StringCmd: redis.NewStringCmd(ctx, "get", key),
				err:       lderr.ErrUnkown,
			}

			err := rds.Process(ctx, cmd)
			c.So(err, convey.ShouldResemble, lderr.ErrUnkown)
			c.So(cmd.Err(), convey.ShouldResemble, lderr.ErrUnkown)
			c.So(cmd.Val(), convey.ShouldEqual, "test-value")
		})
		c.Convey("parse succ", func(c convey.C) {
			rds.cmdable.AddHook(newHook(rds))
			cmd := &testStringCmd{
				StringCmd: redis.NewStringCmd(ctx, "get", key),
				// err:       lderr.ErrUnkown,
			}

			err := rds.Process(ctx, cmd)
			c.So(err, convey.ShouldBeNil)
			c.So(cmd.Err(), convey.ShouldBeNil)
			c.So(cmd.Val(), convey.ShouldEqual, "test-value")
		})
	})
}

func Test_hook_processPipeline(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		rds := MustNewTestRedis()
		defer rds.Close()

		key := "test-key"
		ctx := ldctx.Default()

		c.Convey("pipeline", func(c convey.C) {
			pipeline := rds.Pipeline()

			pipeline.Set(ctx, key, "test-value", 0)

			cmd := &testStringCmd{
				StringCmd: redis.NewStringCmd(ctx, "get", key),
				err:       lderr.ErrUnkown,
			}
			pipeline.Process(ctx, cmd)

			cmds, err := pipeline.Exec(ctx)
			c.So(err, convey.ShouldResemble, lderr.ErrUnkown)
			c.So(cmds, convey.ShouldHaveLength, 2)
			c.So(cmds[0].Err(), convey.ShouldBeNil)
			c.So(cmds[1].Err(), convey.ShouldResemble, lderr.ErrUnkown)
			c.So(cmds[1].(*testStringCmd).Val(), convey.ShouldEqual, "test-value")
		})
		c.Convey("pipelined", func(c convey.C) {
			rds.cmdable.AddHook(newHook(rds))
			cmds, err := rds.Pipelined(ctx, func(p redis.Pipeliner) error {
				p.Set(ctx, key, "test-value", 0)

				cmd := &testStringCmd{
					StringCmd: redis.NewStringCmd(ctx, "get", key),
					// err:       lderr.ErrUnkown,
				}
				return p.Process(ctx, cmd)
			})

			c.So(err, convey.ShouldBeNil)
			c.So(cmds, convey.ShouldHaveLength, 2)
			c.So(cmds[0].Err(), convey.ShouldBeNil)
			c.So(cmds[1].Err(), convey.ShouldBeNil)
			c.So(cmds[1].(*testStringCmd).Val(), convey.ShouldEqual, "test-value")
		})
	})
}
