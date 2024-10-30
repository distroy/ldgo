/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/distroy/ldgo/v3/lderr"
	"github.com/distroy/ldgo/v3/ldhook"
	"github.com/distroy/ldgo/v3/ldlog"
	"github.com/gin-gonic/gin"
	"github.com/smartystreets/goconvey/convey"
)

func Test_shoudBind(t *testing.T) {
	convey.Convey(t.Name(), t, func() {
		patches := ldhook.NewPatches()
		defer patches.Reset()
		patches.Applys([]ldhook.Hook{
			ldhook.FuncHook{
				Target: ldlog.Default,
				Double: ldhook.Values{ldlog.Discard()},
			},
		})

		gin.SetMode(gin.TestMode)
		g, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx := newContext(g)

		convey.Convey("finally succ - ShouldBindJSON", func() {
			type Request struct {
				ProjectId int64  `uri:"project_id"`
				ChannelId int64  `uri:"channel_id"`
				Page      int    `form:"page"`
				Where     string `json:"where"`
			}

			body := `{
				"Where": "abc",
				"channel_id": 123
			}`

			g.Params = append(g.Params, gin.Param{Key: "project_id", Value: "101"})
			g.Params = append(g.Params, gin.Param{Key: "channel_id", Value: "201"})
			g.Request = httptest.NewRequest("GET", "http://github.com/?page=301", strings.NewReader(body))

			req := &Request{}
			convey.So(shouldBind(ctx, req), convey.ShouldBeNil)
			convey.So(req, convey.ShouldResemble, &Request{
				ProjectId: 101,
				ChannelId: 201,
				Page:      301,
				Where:     "abc",
			})
		})
		convey.Convey("bind json but no body", func() {
			type Request struct {
				ProjectId int64  `uri:"project_id"`
				ChannelId int64  `uri:"channel_id"`
				Page      int    `form:"page"`
				Where     string `json:"where"`
			}

			g.Params = append(g.Params, gin.Param{Key: "project_id", Value: "101"})
			g.Params = append(g.Params, gin.Param{Key: "channel_id", Value: "201"})
			g.Request = httptest.NewRequest("GET", "http://github.com/?page=301", nil)

			req := &Request{}
			convey.So(shouldBind(ctx, req), convey.ShouldResemble, lderr.ErrParseRequest)
		})
	})
}
