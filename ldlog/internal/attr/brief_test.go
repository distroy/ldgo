/*
 * Copyright (C) distroy
 */

package attr

import (
	"bytes"
	"io"
	"log"
	"log/slog"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

type testStructEmbedded struct {
	Int64 int `json:"int64,omitempty"`
	Unt64 int `json:"uint64,omitempty"`
}

type testStruct struct {
	*testStructEmbedded

	Bool      bool           `json:"bool,omitempty"`
	String    string         `json:"string,omitempty"`
	Int       int            `json:"int,omitempty"`
	Uint      uint           `json:"uint,omitempty"`
	Float     float64        `json:"float,omitempty"`
	Complex   complex128     `json:"complex,omitempty"`
	Array     []any          `json:"array,omitempty"`
	Map       map[string]any `json:"map,omitempty"`
	Struct    *testStruct    `json:"struct,omitempty"`
	Time      time.Time      `json:"time,omitempty"`
	Duration  time.Duration  `json:"duration,omitempty"`
	Interface any            `json:"interface,omitempty"`
}

func newTestLog(writer io.Writer) *slog.Logger {
	opts := &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	}
	h := slog.NewJSONHandler(writer, opts)
	// h := slog.NewTextHandler(writer, opts)
	return slog.New(h)
}

func testGetLogString(buf *bytes.Buffer, _ *slog.Logger) []byte {
	d := buf.Bytes()
	buf.Reset()
	if l := len(d) - 1; l >= 0 && d[l] == '\n' {
		d = d[:l]
	}

	// x := `{"time":{datetime},`
	// idx := bytes.Index(d, []byte(`"level"`))
	// d = d[idx-len(x):]
	// _ = append(d[:0], x...)
	//
	return d
}

func testRemoveLogPrefix(b []byte) []byte {
	i0 := bytes.Index(b, []byte(`"msg"`))
	b = b[i0:]
	i1 := bytes.Index(b, []byte(`,`))
	b = b[i1:]
	b[0] = '{'
	return b
}

func TestBrief(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		buf := bytes.NewBuffer(nil)
		log := newTestLog(buf)

		SetBriefStringLen(10)
		c.So(briefStringLen, convey.ShouldEqual, 10)
		SetBriefArrayLen(1)
		c.So(briefArrayLen, convey.ShouldEqual, 1)

		getLogString := func() string {
			b := testGetLogString(buf, log)
			b = testRemoveLogPrefix(b)
			return b2s(b)
		}

		c.Convey("no field", func(c convey.C) {
			log.Debug("test", String("key", "value"))
			s := getLogString()
			// c.So(s, convey.ShouldEqual, `2024-06-13T10:50:01.011+0800|debug|test|{"key": "value"}`)
			c.So(s, convey.ShouldEqual, `{"key":"value"}`)
		})
		c.Convey("brief string", func(c convey.C) {
			c.Convey("brief string len 15", func(c convey.C) {
				SetBriefStringLen(15)

				c.Convey("string len 10", func(c convey.C) {
					log.Info("test", BriefString("key", "0123456789"))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `2024-06-13T10:50:01.011+0800|info|test|{"key": "0123456789"}`)
					c.So(s, convey.ShouldEqual, `{"key":"0123456789"}`)
				})
				c.Convey("string len 15", func(c convey.C) {
					log.Warn("test", BriefString("key", "012345678901234"))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `2024-06-13T10:50:01.011+0800|warn|test|{"key": "012345678901234"}`)
					c.So(s, convey.ShouldEqual, `{"key":"012345678901234"}`)
				})
				c.Convey("string len 16", func(c convey.C) {
					log.Error("test", BriefString("key", "0123456789012345"))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `2024-06-13T10:50:01.011+0800|error|test|{"key": {"<len>": 16, "<type>": "string", "<brief>": "012345678901234"}}`)
					c.So(s, convey.ShouldEqual, `{"key":{"<len>":16,"<type>":"string","<brief>":"012345678901234"}}`)
				})
			})
			c.Convey("brief string len 10", func(c convey.C) {
				SetBriefStringLen(10)

				c.Convey("string len 10", func(c convey.C) {
					log.Info("test", BriefString("key", "0123456789"))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": "0123456789"}`)
					c.So(s, convey.ShouldEqual, `{"key":"0123456789"}`)
				})
				c.Convey("string len 15", func(c convey.C) {
					log.Info("test", BriefString("key", "012345678901234"))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": {"<len>": 15, "<type>": "string", "<brief>": "0123456789"}}`)
					c.So(s, convey.ShouldEqual, `{"key":{"<len>":15,"<type>":"string","<brief>":"0123456789"}}`)
				})
			})
		})
		c.Convey("brief byte string", func(c convey.C) {
			c.Convey("byte string len 10", func(c convey.C) {
				log.Info("test", BriefByteString("key", []byte("0123456789")))
				s := getLogString()
				// c.So(s, convey.ShouldEqual, `{"key": "0123456789"}`)
				c.So(s, convey.ShouldEqual, `{"key":"0123456789"}`)
			})
			c.Convey("byte string len 15", func(c convey.C) {
				log.Info("test", BriefByteString("key", []byte("012345678901234")))
				s := getLogString()
				// c.So(s, convey.ShouldEqual, `{"key": {"<len>": 15, "<type>": "string", "<brief>": "0123456789"}}`)
				c.So(s, convey.ShouldEqual, `{"key":{"<len>":15,"<type>":"string","<brief>":"0123456789"}}`)
			})
		})
		c.Convey("brief strings", func(c convey.C) {
			c.Convey("brief array len 2", func(c convey.C) {
				SetBriefArrayLen(2)

				c.Convey("strings len 2", func(c convey.C) {
					log.Info("test", BriefStrings("key", []string{"0123456789", "012345678901234"}))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": ["0123456789", {"<len>": 15, "<type>": "string", "<brief>": "0123456789"}]}`)
					c.So(s, convey.ShouldEqual, `{"key":["0123456789",{"<len>":15,"<type>":"string","<brief>":"0123456789"}]}`)
				})
				c.Convey("strings len 3", func(c convey.C) {
					log.Info("test", BriefStrings("key", []string{"0123456789", "012345678901234", "abcdefghijklmnopqrstuvwxyz"}))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": {"<len>": 3, "<type>": "array", "<brief>": ["0123456789", {"<len>": 15, "<type>": "string", "<brief>": "0123456789"}]}}`)
					c.So(s, convey.ShouldEqual, `{"key":{"<len>":3,"<type>":"array","<brief>":["0123456789",{"<len>":15,"<type>":"string","<brief>":"0123456789"}]}}`)
				})
			})
			c.Convey("brief array len 1", func(c convey.C) {
				SetBriefArrayLen(1)

				c.Convey("strings len 1", func(c convey.C) {
					log.Info("test", BriefStrings("key", []string{"0123456789"}))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": ["0123456789"]}`)
					c.So(s, convey.ShouldEqual, `{"key":["0123456789"]}`)
				})
				c.Convey("strings len 3", func(c convey.C) {
					log.Info("test", BriefStrings("key", []string{"0123456789", "012345678901234", "abcdefghijklmnopqrstuvwxyz"}))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": {"<len>": 3, "<type>": "array", "<brief>": ["0123456789"]}}`)
					c.So(s, convey.ShouldEqual, `{"key":{"<len>":3,"<type>":"array","<brief>":["0123456789"]}}`)
				})
			})
		})
		c.Convey("brief reflect", func(c convey.C) {
			c.Convey("array", func(c convey.C) {
				c.Convey("array len 2", func(c convey.C) {
					SetBriefArrayLen(2)
					c.Convey("array len 2", func(c convey.C) {
						log.Info("test", BriefReflect("key", [...]any{123456789, []byte("012345678901234")}))
						s := getLogString()
						// c.So(s, convey.ShouldEqual, `{"key": [123456789, {"<len>": 15, "<type>": "string", "<brief>": "0123456789"}]}`)
						c.So(s, convey.ShouldEqual, `{"key":[123456789,{"<len>":15,"<type>":"string","<brief>":"0123456789"}]}`)
					})
					c.Convey("array len 3", func(c convey.C) {
						log.Info("test", BriefReflect("key", []any{123456789, "012345678901234", "abcdefghijklmnopqrstuvwxyz"}))
						s := getLogString()
						// c.So(s, convey.ShouldEqual, `{"key": {"<len>": 3, "<type>": "array", "<brief>": [123456789, {"<len>": 15, "<type>": "string", "<brief>": "0123456789"}]}}`)
						c.So(s, convey.ShouldEqual, `{"key":{"<len>":3,"<type>":"array","<brief>":[123456789,{"<len>":15,"<type>":"string","<brief>":"0123456789"}]}}`)
					})
				})
				c.Convey("array len 1", func(c convey.C) {
					SetBriefArrayLen(1)
					c.Convey("array len 1", func(c convey.C) {
						log.Info("test", BriefReflect("key", []any{123456789}))
						s := getLogString()
						// c.So(s, convey.ShouldEqual, `{"key": [123456789]}`)
						c.So(s, convey.ShouldEqual, `{"key":[123456789]}`)
					})
					c.Convey("array len 3", func(c convey.C) {
						log.Info("test", BriefReflect("key", []any{123456789, "012345678901234", "abcdefghijklmnopqrstuvwxyz"}))
						s := getLogString()
						// c.So(s, convey.ShouldEqual, `{"key": {"<len>": 3, "<type>": "array", "<brief>": [123456789]}}`)
						c.So(s, convey.ShouldEqual, `{"key":{"<len>":3,"<type>":"array","<brief>":[123456789]}}`)
					})
				})
			})
			c.Convey("map", func(c convey.C) {
				c.Convey("map len 10", func(c convey.C) {
					SetBriefMapLen(10)
					c.Convey("without brief string & array", func(c convey.C) {
						log.Info("test", BriefReflect("key", map[string]any{
							"bool":    true,
							"string":  `0123456789`,
							"int":     1001,
							"uint":    uint(1002),
							"float":   float64(1003),
							"complex": complex(1004, 1005),
							"array":   []any{},
							"map":     map[string]any{},
							"struct":  &testStruct{},
						}))
						s := getLogString()
						// c.So(s, convey.ShouldEqual, `{"key": {"array": [], "bool": true, "complex": "1004+1005i", "float": 1003, "int": 1001, "map": {}, "string": "0123456789", "struct": {}, "uint": 1002}}`)
						c.So(s, convey.ShouldEqual, `{"key":{"array":[],"bool":true,"complex":"1004+1005i","float":1003,"int":1001,"map":{},"string":"0123456789","struct":{},"uint":1002}}`)
					})
					c.Convey("with brief string & array", func(c convey.C) {
						log.Info("test", BriefReflect("key", map[string]any{
							"bool":    true,
							"string":  `012345678901234`,
							"int":     1001,
							"uint":    uint(1002),
							"float":   float64(1003),
							"complex": complex(1004, 1005),
							"array":   []any{123456789, 123456789, 123456789},
							"map":     map[string]any{},
							"struct":  &testStruct{},
						}))
						s := getLogString()
						// c.So(s, convey.ShouldEqual, `{"key": {"array": {"<len>": 3, "<type>": "array", "<brief>": [123456789]}, "bool": true, "complex": "1004+1005i", "float": 1003, "int": 1001, "map": {}, "string": {"<len>": 15, "<type>": "string", "<brief>": "0123456789"}, "struct": {}, "uint": 1002}}`)
						c.So(s, convey.ShouldEqual, `{"key":{"array":{"<len>":3,"<type>":"array","<brief>":[123456789]},"bool":true,"complex":"1004+1005i","float":1003,"int":1001,"map":{},"string":{"<len>":15,"<type>":"string","<brief>":"0123456789"},"struct":{},"uint":1002}}`)
					})
				})
				c.Convey("map len 5", func(c convey.C) {
					SetBriefMapLen(5)
					c.Convey("without brief string & array", func(c convey.C) {
						log.Info("test", BriefReflect("key", map[string]any{
							"bool":    true,
							"string":  `0123456789`,
							"int":     1001,
							"uint":    uint(1002),
							"float":   float64(1003),
							"complex": complex(1004, 1005),
							"array":   []any{},
							"map":     map[string]any{},
							"struct":  &testStruct{},
						}))
						s := getLogString()
						// c.So(s, convey.ShouldEqual, `{"key": {"<len>": 9, "<type>": "map", "<brief>": {"array": [], "bool": true, "complex": "1004+1005i", "float": 1003, "int": 1001}}}`)
						c.So(s, convey.ShouldEqual, `{"key":{"<len>":9,"<type>":"map","<brief>":{"array":[],"bool":true,"complex":"1004+1005i","float":1003,"int":1001}}}`)
					})
					c.Convey("with brief string & array", func(c convey.C) {
						log.Info("test", BriefReflect("key", map[string]any{
							"bool":    true,
							"string":  `012345678901234`,
							"int":     1001,
							"uint":    uint(1002),
							"float":   float64(1003),
							"complex": complex(1004, 1005),
							"array":   []any{123456789, 123456789, 123456789},
							"map":     map[string]any{},
							"struct":  &testStruct{},
						}))
						s := getLogString()
						// c.So(s, convey.ShouldEqual, `{"key": {"<len>": 9, "<type>": "map", "<brief>": {"array": {"<len>": 3, "<type>": "array", "<brief>": [123456789]}, "bool": true, "complex": "1004+1005i", "float": 1003, "int": 1001}}}`)
						c.So(s, convey.ShouldEqual, `{"key":{"<len>":9,"<type>":"map","<brief>":{"array":{"<len>":3,"<type>":"array","<brief>":[123456789]},"bool":true,"complex":"1004+1005i","float":1003,"int":1001}}}`)
					})
				})
			})
			c.Convey("struct", func(c convey.C) {
				c.Convey("without brief string & array", func(c convey.C) {
					log.Info("test", BriefReflect("key", &testStruct{
						Bool:    true,
						String:  `0123456789`,
						Int:     1001,
						Uint:    uint(1002),
						Float:   float64(1003),
						Complex: complex(1004, 1005),
						Array:   []any{},
						Map:     map[string]any{},
						Struct:  &testStruct{},
					}))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": {"bool": true, "string": "0123456789", "int": 1001, "uint": 1002, "float": 1003, "complex": "1004+1005i", "array": [], "map": {}, "struct": {}}}`)
					c.So(s, convey.ShouldEqual, `{"key":{"bool":true,"string":"0123456789","int":1001,"uint":1002,"float":1003,"complex":"1004+1005i","array":[],"map":{},"struct":{}}}`)
				})
				c.Convey("with brief string & array", func(c convey.C) {
					log.Info("test", BriefReflect("key", &testStruct{
						Bool:    true,
						String:  `012345678901234`,
						Int:     1001,
						Uint:    uint(1002),
						Float:   float64(1003),
						Complex: complex(1004, 1005),
						Array:   []any{123456789, 123456789, 123456789},
						Map:     map[string]any{},
						Struct:  &testStruct{},
					}))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": {"bool": true, "string": {"<len>": 15, "<type>": "string", "<brief>": "0123456789"}, "int": 1001, "uint": 1002, "float": 1003, "complex": "1004+1005i", "array": {"<len>": 3, "<type>": "array", "<brief>": [123456789]}, "map": {}, "struct": {}}}`)
					c.So(s, convey.ShouldEqual, `{"key":{"bool":true,"string":{"<len>":15,"<type>":"string","<brief>":"0123456789"},"int":1001,"uint":1002,"float":1003,"complex":"1004+1005i","array":{"<len>":3,"<type>":"array","<brief>":[123456789]},"map":{},"struct":{}}}`)
				})
				c.Convey("embedded struct", func(c convey.C) {
					log.Info("test", BriefReflect("key", &testStruct{
						testStructEmbedded: &testStructEmbedded{
							Int64: 2001,
							Unt64: 2002,
						},
					}))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": {"int64": 2001, "uint64": 2002}}`)
					c.So(s, convey.ShouldEqual, `{"key":{"int64":2001,"uint64":2002}}`)
				})
				c.Convey("time & duration & fmt.Stringer", func(c convey.C) {
					log.Info("test", BriefReflect("key", &testStruct{
						Time:      time.Unix(1718247001, 11*1000*1000),
						Duration:  time.Second,
						Interface: string_t("012345678901234"),
					}))
					s := getLogString()
					// c.So(s, convey.ShouldEqual, `{"key": {"time": "2024-06-13T10:50:01.011+0800", "duration": 1, "interface": {"<len>": 15, "<type>": "string", "<brief>": "0123456789"}}}`)
					c.So(s, convey.ShouldEqual, `{"key":{"time":"2024-06-13T10:50:01.011+0800","duration":"1s","interface":{"<len>":15,"<type>":"string","<brief>":"0123456789"}}}`)
				})
			})
		})
	})
}
