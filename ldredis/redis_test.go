package ldredis

import (
	"testing"
	"time"

	"github.com/distroy/ldgo/v3/ldctx"
	"github.com/distroy/ldgo/v3/ldhook"
	redis "github.com/redis/go-redis/v9"
	"github.com/smartystreets/goconvey/convey"
)

func Test_New(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := ldhook.NewPatches()
		defer patches.Reset()

		convey.Convey("redis cluster", func() {
			patches.Applys([]ldhook.Hook{
				ldhook.FuncHook{
					Target: redis.NewClusterClient,
					Double: ldhook.Values{&redis.ClusterClient{}},
				},
				ldhook.FuncHook{
					Target: (*redis.ClusterClient).AddHook,
					Double: ldhook.Values{},
				},
			})
			cfg := &Config{
				Cluster:  true,
				Addrs:    []string{"1.1.1.1"},
				Password: "now password",
			}
			convey.So(NewByConfig(cfg), convey.ShouldNotBeNil)
		})

		convey.Convey("redis not cluster", func() {
			patches.Applys([]ldhook.Hook{
				ldhook.FuncHook{
					Target: redis.NewClient,
					Double: ldhook.Values{&redis.Client{}},
				},
				ldhook.FuncHook{
					Target: (*redis.Client).AddHook,
					Double: ldhook.Values{},
				},
			})
			cfg := &Config{
				Cluster:  false,
				Addrs:    []string{"1.1.1.1"},
				Addr:     "1.1.1.1",
				Password: "now password",
			}
			convey.So(NewByConfig(cfg), convey.ShouldNotBeNil)
		})
	})
}

type testReporter struct {
	commands  [][]interface{}
	pipelines [][][]interface{}
}

func (r *testReporter) Report(cmd Cmder, d time.Duration) {
	r.commands = append(r.commands, cmd.Args())
}

func (r *testReporter) ReportPipeline(cmds []Cmder, d time.Duration) {
	argss := make([][]interface{}, len(cmds))
	for _, cmd := range cmds {
		argss = append(argss, cmd.Args())
	}
	r.pipelines = append(r.pipelines, argss)
}

func TestRedis_WithReport(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		ctx := ldctx.Discard()

		rdb := MustNewTestRedis()
		defer rdb.Close()

		reporter := &testReporter{}
		rdb1 := rdb.WithReport(reporter)
		rdb1.wrapper.AddHook(newHook(rdb1))

		const (
			key = "test-key"
		)

		cmd0 := rdb.Set(ctx, key, "123", 0)
		c.So(cmd0.Err(), convey.ShouldBeNil)
		c.So(cmd0.Val(), convey.ShouldNotBeEmpty)

		c.So(reporter, convey.ShouldResemble, &testReporter{})

		cmd1 := rdb.Get(ctx, key)
		c.So(cmd1.Err(), convey.ShouldBeNil)
		c.So(cmd1.Val(), convey.ShouldNotBeEmpty)
		c.So(reporter, convey.ShouldResemble, &testReporter{})

		cmd2 := rdb1.Get(ctx, key)
		c.So(cmd2.Err(), convey.ShouldBeNil)
		c.So(cmd2.Val(), convey.ShouldNotBeEmpty)
		c.So(reporter, convey.ShouldResemble, &testReporter{
			commands: [][]interface{}{
				{"get", key},
			},
		})
	})
}
