/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/distroy/ldgo/v2/lderr"
	"github.com/gin-gonic/gin"
	"github.com/smartystreets/goconvey/convey"
)

type testReaderForChunked struct {
	Error error
	Lines []string

	index int
}

func (r *testReaderForChunked) Read(buff []byte) (int, error) {
	if r.index < len(r.Lines) {
		data := r.Lines[r.index]
		r.index++

		return copy(buff, data), nil
	}

	err := r.Error
	if err != nil {
		return 0, err
	}
	return 0, io.EOF
}

func TestReaderRenderer_Render(t *testing.T) {
	buildResultString := func(strs []string) string {
		return strings.Join(strs, "")
	}

	convey.Convey(t.Name(), t, func(c convey.C) {
		w := httptest.NewRecorder()
		g, _ := gin.CreateTestContext(w)

		g.Params = append(g.Params, gin.Param{Key: "project_id", Value: "101"})
		g.Params = append(g.Params, gin.Param{Key: "channel_id", Value: "201"})
		g.Request = httptest.NewRequest("GET", "http://github.com/?page=301", nil)

		c.Convey("chunked", func(c convey.C) {
			c.Convey("render succ", func(c convey.C) {
				c.Convey("only reader", func(c convey.C) {
					handler := WrapHandler(func(c *Context) (ReaderRenderer, lderr.Error) {
						renderer := ReaderRenderer{
							Reader: &testReaderForChunked{
								Lines: []string{
									strings.Repeat("abc", 5),
									"",
									strings.Repeat("xyz", 6),
								},
							},
						}
						return renderer, nil
					})

					handler(g)
					c.So(w.Code, convey.ShouldEqual, http.StatusOK)
					c.So(GetError(g), convey.ShouldBeNil)

					c.SkipSo(w.Header().Get(headerContentType), convey.ShouldEqual, "")
					c.So(w.Header().Get(chunkedHeaderKey), convey.ShouldEqual, chunkedHeaderValue)

					c.So(GetRequest(g), convey.ShouldBeNil)
					c.So(GetResponse(g), convey.ShouldBeNil)
					c.So(GetRenderer(g), convey.ShouldNotBeNil)
					c.So(GetContext(g), convey.ShouldNotBeNil)

					c.So(w.Body.String(), convey.ShouldEqual, buildResultString([]string{
						strings.Repeat("abc", 5),
						strings.Repeat("xyz", 6),
					}))
				})

				c.Convey("other options", func(c convey.C) {
					handler := WrapHandler(func(c *Context) (ReaderRenderer, lderr.Error) {
						renderer := ReaderRenderer{
							Code:        220,
							ContentType: "text/html; charset=utf-8",
							Headers: map[string]string{
								"x-abc": "aaaa",
							},
							Reader: &testReaderForChunked{
								Lines: []string{
									strings.Repeat("abc", 5),
									"",
									strings.Repeat("xyz", 6),
								},
							},
						}
						return renderer, nil
					})

					handler(g)
					c.So(w.Code, convey.ShouldEqual, 220)
					c.So(GetError(g), convey.ShouldBeNil)

					c.So(w.Header().Get(headerContentType), convey.ShouldEqual, "text/html; charset=utf-8")
					c.So(w.Header().Get(chunkedHeaderKey), convey.ShouldEqual, chunkedHeaderValue)
					c.So(w.Header().Get("x-abc"), convey.ShouldEqual, "aaaa")

					c.So(w.Body.String(), convey.ShouldEqual, buildResultString([]string{
						strings.Repeat("abc", 5),
						strings.Repeat("xyz", 6),
					}))
				})

				c.Convey("chunked with content length", func(c convey.C) {
					handler := WrapHandler(func(c *Context) (ReaderRenderer, lderr.Error) {
						renderer := ReaderRenderer{
							Code:          220,
							ContentType:   "text/html; charset=utf-8",
							ContentLength: 100,
							Headers: map[string]string{
								"x-abc":          "aaaa",
								chunkedHeaderKey: chunkedHeaderValue,
							},
							Reader: &testReaderForChunked{
								Lines: []string{
									strings.Repeat("abc", 5),
									"",
									strings.Repeat("xyz", 6),
								},
							},
						}
						return renderer, nil
					})

					handler(g)
					c.So(w.Code, convey.ShouldEqual, 220)
					c.So(GetError(g), convey.ShouldBeNil)

					c.So(w.Header().Get(headerContentType), convey.ShouldEqual, "text/html; charset=utf-8")
					c.So(w.Header().Get(chunkedHeaderKey), convey.ShouldEqual, chunkedHeaderValue)
					c.So(w.Header().Get("x-abc"), convey.ShouldEqual, "aaaa")

					c.So(w.Body.String(), convey.ShouldEqual, buildResultString([]string{
						strings.Repeat("abc", 5),
						strings.Repeat("xyz", 6),
					}))
				})
			})

			c.Convey("read fail", func(c convey.C) {
				handler := WrapHandler(func(c *Context) (ReaderRenderer, lderr.Error) {
					renderer := ReaderRenderer{
						Reader: &testReaderForChunked{
							Error: fmt.Errorf("test error"),
							Lines: []string{
								strings.Repeat("abc", 3),
							},
						},
					}
					return renderer, nil
				})

				handler(g)
				c.So(w.Code, convey.ShouldEqual, http.StatusOK)
				c.So(GetError(g), convey.ShouldNotBeNil)
				c.So(lderr.GetCode(GetError(g)), convey.ShouldEqual, lderr.ErrHttpRenderBody.Code())

				c.So(w.Header().Get(chunkedHeaderKey), convey.ShouldEqual, chunkedHeaderValue)

				c.So(GetRenderer(g), convey.ShouldNotBeNil)
				c.So(w.Body.String(), convey.ShouldEqual, buildResultString([]string{
					strings.Repeat("abc", 3),
					crlf, crlf,
					"server happened some errors", crlf,
					fmt.Sprintf("code: %d", lderr.ErrHttpRenderBody.Code()), crlf,
					fmt.Sprintf("message: %s", lderr.ErrHttpRenderBody.Error()), crlf,
					"details:", crlf,
					"\ttest error", crlf,
				}))
			})
		})
	})
}
