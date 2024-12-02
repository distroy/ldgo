/*
 * Copyright (C) distroy
 */

package ldflag

import (
	"reflect"
	"strings"
	"testing"

	"github.com/distroy/ldgo/v2/ldptr"
	"github.com/smartystreets/goconvey/convey"
)

const (
	testCommandName = "unittest"
)

type testIncludes []string

func (p *testIncludes) Default() string { return strings.Join(nil, "\n") }
func (p *testIncludes) String() string  { return mustMarshalJson(*p) }
func (p *testIncludes) Set(s string) error {
	*p = append(*p, s)
	return nil
}

type testExcludes []string

func (p *testExcludes) Default() string { return strings.Join([]string{"a", "b"}, "\n") }
func (p *testExcludes) String() string  { return mustMarshalJson(*p) }
func (p *testExcludes) Set(s string) error {
	*p = append(*p, s)
	return nil
}

func newTestFlagSet() *FlagSet {
	s := NewFlagSet()
	s.init()
	s.name = testCommandName
	return s
}

func TestFlagSet_wrigeUsageHeader(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey(`name == "" && no args`, func(c convey.C) {
			s := &FlagSet{
				name: "",
				args: nil,
			}
			b := &strings.Builder{}
			s.writeUsageHeader(b)
			c.So(b.String(), convey.ShouldEqual, "Usage of <command>:\nFlags:\n")
		})
		c.Convey(`name == "abc" && args.meta == ""`, func(c convey.C) {
			s := &FlagSet{
				name: "abc",
				args: &Flag{Meta: ""},
			}
			b := &strings.Builder{}
			s.writeUsageHeader(b)
			c.So(b.String(), convey.ShouldEqual, "Usage: abc [<flags>] [<arg>...]\n\nFlags:\n")
		})
	})
}

func TestFlagSet_wrigeUsage(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		b := &strings.Builder{}
		s := newTestFlagSet()

		c.Convey("normal value", func(c convey.C) {
			type Flags struct {
				Top      int      `ldflag:"name:top; meta:N; usage:show the top <N>"`
				Avg      bool     `ldflag:"usage:show the average complexity"`
				DebugLog bool     `ldflag:"usage:print debug log; bool"`
				Rate     float64  `ldflag:"default:0.65; usage:"`
				Branch   string   `ldflag:"meta:branch; usage:git branch name"`
				Includes []string `ldflag:"name:include; meta:regexp; usage:include file regexps"`
				Excludes []string `ldflag:"name:exclude; meta:regexp; usage:exclude file regexps"`
				Pathes   []string `ldflag:"args; meta:path; default:."`
			}

			flags := &Flags{
				Excludes: []string{
					`(^|/)vendor/`,
					`\.pb\.go$`,
				},
			}
			s.Model(flags)
			s.writeUsage(b)

			c.So(b.String(), convey.ShouldEqual, `Usage: unittest [<flags>] [path...]

Flags:
        -top <N>
                show the top <N>
        -avg <bool>
                show the average complexity
        -debug-log
                print debug log
        -rate <float> (default: 0.65)
        -branch <branch>
                git branch name
        -include <regexp>
                include file regexps
        -exclude <regexp>
                exclude file regexps
                default:
                        "(^|/)vendor/"
                        "\\.pb\\.go$"
`)
		})

		c.Convey("pointer value", func(c convey.C) {
			type Flags struct {
				Top      *int     `ldflag:"name:top; meta:N; usage:show the top <N>"`
				Avg      *bool    `ldflag:"usage:show the average complexity"`
				DebugLog *bool    `ldflag:"usage:print debug log; bool"`
				Rate     *float64 `ldflag:"default:0.65; usage:"`
				Branch   *string  `ldflag:"meta:branch; usage:git branch name"`
				Includes []string `ldflag:"name:include; meta:regexp; usage:include file regexps"`
				Excludes []string `ldflag:"name:exclude; meta:regexp; usage:exclude file regexps"`
				Pathes   []string `ldflag:"args; meta:path; default:."`
			}

			flags := &Flags{
				Excludes: []string{
					`(^|/)vendor/`,
					`\.pb\.go$`,
				},
			}
			s.Model(flags)
			s.writeUsage(b)

			c.So(b.String(), convey.ShouldEqual, `Usage: unittest [<flags>] [path...]

Flags:
        -top <N>
                show the top <N>
        -avg <bool>
                show the average complexity
        -debug-log
                print debug log
        -rate <float> (default: 0.65)
        -branch <branch>
                git branch name
        -include <regexp>
                include file regexps
        -exclude <regexp>
                exclude file regexps
                default:
                        "(^|/)vendor/"
                        "\\.pb\\.go$"
`)
		})
		c.Convey("value with default", func(c convey.C) {
			type Flags struct {
				Includes testIncludes `ldflag:"name:include; meta:regexp; usage:include file regexps"`
				Excludes testExcludes `ldflag:"name:exclude; meta:regexp; usage:exclude file regexps"`
			}

			flags := &Flags{}
			s.Model(flags)
			s.writeUsage(b)

			c.So(b.String(), convey.ShouldEqual, `Usage of unittest:
Flags:
        -include <regexp>
                include file regexps
        -exclude <regexp>
                exclude file regexps
                default:
                        a
                        b
`)
		})

		c.Convey("short flag && no usage", func(c convey.C) {
			type Flags struct {
				Csv bool `ldflag:"bool"`
			}

			flags := &Flags{}

			s.Model(flags)
			s.writeUsage(b)

			c.So(b.String(), convey.ShouldEqual, `Usage of unittest:
Flags:
        -csv
`)
		})
	})
}

func TestFlagSet_Model(t *testing.T) {
	testClearFlagSet := func(s *FlagSet) {
		for _, v := range s.flagSlice {
			v.tags = nil
			v.val = reflect.Value{}
			v.Value = nil
		}
		if args := s.args; args != nil {
			args.tags = nil
			args.val = reflect.Value{}
			args.Value = nil
		}
	}

	type Flags struct {
		Top      int          `ldflag:"name:top; meta:N; usage:show the top <N>"`
		Avg      bool         `ldflag:"usage:show the average complexity"`
		DebugLog bool         `ldflag:"usage:print debug log; bool"`
		Rate     float64      `ldflag:"default:0.65; usage:"`
		Branch   string       `ldflag:"meta:branch; usage:git branch name"`
		Pathes   []string     `ldflag:"args; meta:path; default:."`
		Includes testExcludes `ldflag:"-"`
	}

	convey.Convey(t.Name(), t, func(c convey.C) {
		c.Convey("", func(c convey.C) {
			flags := &Flags{}

			s := NewFlagSet()
			s.Model(flags)
			testClearFlagSet(s)

			c.So(s.flagSlice, convey.ShouldResemble, []*Flag{
				{
					lvl:     0,
					Name:    "top",
					Meta:    "N",
					Default: "0",
					Usage:   "show the top <N>",
					IsArgs:  false,
					Bool:    false,
				},
				{
					lvl:     0,
					Name:    "avg",
					Meta:    "",
					Default: "false",
					Usage:   "show the average complexity",
					IsArgs:  false,
					Bool:    false,
				},
				{
					lvl:     0,
					Name:    "debug-log",
					Meta:    "",
					Default: "false",
					Usage:   "print debug log",
					IsArgs:  false,
					Bool:    true,
				},
				{
					lvl:     0,
					Name:    "rate",
					Meta:    "",
					Default: "0.65",
					Usage:   "",
					IsArgs:  false,
					Bool:    false,
				},
				{
					lvl:     0,
					Name:    "branch",
					Meta:    "branch",
					Default: "",
					Usage:   "git branch name",
					IsArgs:  false,
					Bool:    false,
				},
			})
			c.So(s.args, convey.ShouldResemble, &Flag{
				lvl:     0,
				Name:    "pathes",
				Meta:    "path",
				Default: ".",
				Usage:   "",
				IsArgs:  true,
				Bool:    false,
			})
		})
	})
}

func TestFlagSet_Parse(t *testing.T) {
	convey.Convey(t.Name(), t, func(c convey.C) {
		s := newTestFlagSet()

		c.Convey("normal value", func(c convey.C) {
			type Flags struct {
				Top      int      `ldflag:"name:top; meta:N; usage:show the top <N>"`
				Avg      bool     `ldflag:"usage:show the average complexity"`
				DebugLog bool     `ldflag:"usage:print debug log; bool"`
				Rate     float64  `ldflag:"default:0.65; usage:"`
				Branch   string   `ldflag:"meta:branch; usage:git branch name"`
				Pathes   []string `ldflag:"args; meta:path; default:."`
			}

			c.Convey("no set default", func(c convey.C) {
				flags := &Flags{}
				s.EnableDefault(false)
				s.Model(flags)

				err := s.Parse([]string{
					"-top", "5",
					"-avg", "1",
					"-debug-log",
				})
				c.So(err, convey.ShouldBeNil)
				c.So(flags, convey.ShouldResemble, &Flags{
					Top:      5,
					Avg:      true,
					DebugLog: true,
					Rate:     0,
					Branch:   "",
					Pathes:   []string{},
				})
			})

			c.Convey("set default", func(c convey.C) {
				flags := &Flags{}
				s.EnableDefault(true)
				s.Model(flags)

				err := s.Parse([]string{
					"-top", "5",
					"-avg", "1",
					"-debug-log",
				})
				c.So(err, convey.ShouldBeNil)
				c.So(flags, convey.ShouldResemble, &Flags{
					Top:      5,
					Avg:      true,
					DebugLog: true,
					Rate:     0.65,
					Branch:   "",
					Pathes:   []string{"."},
				})
			})
		})

		c.Convey("ptr value", func(c convey.C) {
			type Flags struct {
				Over     *int     `ldflag:"name:over; meta:N; default:15; usage:show functions with complexity <N>"`
				Top      *int     `ldflag:"name:top; meta:N; usage:show the top <N>"`
				Avg      *bool    `ldflag:"usage:show the average complexity"`
				DebugLog *bool    `ldflag:"usage:print debug log; bool"`
				Rate     *float64 `ldflag:"default:0.65; usage:"`
				Branch   *string  `ldflag:"meta:branch; usage:git branch name"`
				Pathes   []string `ldflag:"args; meta:path; default:."`
			}

			c.Convey("no set default", func(c convey.C) {
				flags := &Flags{}
				s.EnableDefault(false)
				s.Model(flags)

				err := s.Parse([]string{
					"-top", "5",
					"-avg", "1",
					"-debug-log",
				})
				c.So(err, convey.ShouldBeNil)
				c.So(flags, convey.ShouldResemble, &Flags{
					Over:     nil,
					Top:      ldptr.New(5),
					Avg:      ldptr.New(true),
					DebugLog: ldptr.New(true),
					Rate:     nil,
					Branch:   nil,
					Pathes:   []string{},
				})
			})

			c.Convey("set default", func(c convey.C) {
				flags := &Flags{}
				s.EnableDefault(true)
				s.Model(flags)

				err := s.Parse([]string{
					"-top", "5",
					"-avg", "1",
					"-debug-log",
				})
				c.So(err, convey.ShouldBeNil)
				c.So(flags, convey.ShouldResemble, &Flags{
					Over:     ldptr.New(15),
					Top:      ldptr.New(5),
					Avg:      ldptr.New(true),
					DebugLog: ldptr.New(true),
					Rate:     ldptr.New(0.65),
					Branch:   nil,
					Pathes:   []string{"."},
				})
			})
		})

		c.Convey("options", func(c convey.C) {
			type Flags struct {
				Method *string `ldflag:"name:method; options: a, b, c"`
			}

			c.Convey("succ", func(c convey.C) {
				flags := &Flags{}
				s.EnableDefault(false)
				s.Model(flags)

				err := s.Parse([]string{
					"-method", "a",
				})
				c.So(err, convey.ShouldBeNil)
				c.So(flags, convey.ShouldResemble, &Flags{
					Method: ldptr.New("a"),
				})
			})

			c.Convey("fail", func(c convey.C) {
				flags := &Flags{}
				s.EnableDefault(true)
				s.Model(flags)

				err := s.Parse([]string{
					"-method", "d",
				})
				c.So(err, convey.ShouldNotBeNil)
			})
		})

		c.Convey("options with default", func(c convey.C) {
			type Flags struct {
				Method *string `ldflag:"name:method; default: a; options: a, b, c"`
			}

			c.Convey("succ", func(c convey.C) {
				flags := &Flags{}
				s.EnableDefault(true)
				s.Model(flags)

				err := s.Parse([]string{})
				c.So(err, convey.ShouldBeNil)
				c.So(flags, convey.ShouldResemble, &Flags{
					Method: ldptr.New("a"),
				})
			})

			c.Convey("fail", func(c convey.C) {
				flags := &Flags{}
				s.EnableDefault(true)
				s.Model(flags)

				err := s.Parse([]string{
					"-method", "d",
				})
				c.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
