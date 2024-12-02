/*
 * Copyright (C) distroy
 */

package ldflag

import (
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/distroy/ldgo/v2/ldtagmap"
)

type Flag struct {
	lvl  int
	val  reflect.Value
	tags ldtagmap.Tags

	Name    string
	Value   Value
	Meta    string
	Default string
	Usage   string
	IsArgs  bool
	Bool    bool
	Options []string
}

func (f *Flag) inOptions(s string) bool {
	if f.Default != "" && s == f.Default {
		return true
	}
	for _, v := range f.Options {
		if v == s {
			return true
		}
	}
	return false
}

type FlagSet struct {
	command   *flag.FlagSet
	name      string
	flagSlice []*Flag
	flagMap   map[string]*Flag
	args      *Flag
	model     reflect.Value
	noDefault bool
}

func NewFlagSet() *FlagSet {
	s := &FlagSet{}
	s.init()
	return s
}

func (s *FlagSet) init() {
	if s.flagMap != nil {
		return
	}

	name := os.Args[0]
	s.command = flag.NewFlagSet(name, flag.ExitOnError)
	s.name = name
	s.flagMap = make(map[string]*Flag)
	s.command.Usage = s.printUsage
	s.noDefault = false
}

func (s *FlagSet) Args() []string {
	s.init()
	return s.command.Args()
}

// Deprecated
func (s *FlagSet) EnableDefault(on bool) {
	s.init()
	s.noDefault = !on
}

func (s *FlagSet) SetOutput(w io.Writer) {
	s.init()
	s.command.SetOutput(w)
}

func (s *FlagSet) PrintUsage() {
	s.printUsage()
}

func (s *FlagSet) printUsage() {
	w := s.command.Output()
	s.writeUsage(w)
}

func (s *FlagSet) WriteUsage(w io.Writer) {
	s.writeUsage(w)
}

func (s *FlagSet) writeUsage(w io.Writer) {
	// log.Printf("flags: %s", jsoncore.MustMarshalToString(s.model.Interface()))

	s.writeUsageHeader(w)

	flags := s.flagSlice
	// flags := s.sortedFlags()
	for _, f := range flags {
		s.writeFlagUsage(w, f)
	}
}

func (s *FlagSet) writeUsageHeader(w io.Writer) {
	name := s.name
	if name == "" {
		name = "<command>"
	}

	if s.args != nil {
		meta := s.args.Meta
		if meta == "" {
			meta = "<arg>"
		}

		fmt.Fprintf(w, "Usage: %s [<flags>] [%s...]\n", name, meta)
		fmt.Fprint(w, "\n")
		fmt.Fprintf(w, "Flags:\n")
		return
	}

	fmt.Fprintf(w, "Usage of %s:\n", name)
	fmt.Fprintf(w, "Flags:\n")
}

func (s *FlagSet) writeFlagUsage(w io.Writer, f *Flag) {
	const (
		tab           = "        "
		nameSize      = len(tab) * 2
		namePrefix    = tab
		usagePrefix   = "\n" + namePrefix + tab
		defaultPrefix = usagePrefix + tab
	)

	b := &strings.Builder{}

	fmt.Fprintf(b, "%s-%s", namePrefix, f.Name) // Two spaces before -; see next two comments.
	meta, usage := unquoteUsage(f)
	if len(meta) > 0 {
		fmt.Fprintf(b, " %s", meta)
	}

	// Boolean flags of one ASCII letter are so common we
	// treat them specially, putting their usage on the same line.
	if usage != "" {
		if b.Len() <= nameSize-2 { // space, space, '-', 'x'.
			for i := b.Len(); i < nameSize; i++ {
				fmt.Fprint(b, " ")
			}
		}

		// Four spaces before the tab triggers good alignment
		// for both 4- and 8-space tab stops.
		fmt.Fprint(b, usagePrefix)
	}

	fmt.Fprint(b, strings.ReplaceAll(usage, "\n", usagePrefix))

	if isFlagDefaultZero(f) {
		fmt.Fprint(w, b.String(), "\n")
		return
	}

	switch v := f.Value.(type) {
	default:
		if strings.Index(f.Default, "\n") > 0 {
			fmt.Fprint(b, usagePrefix, "default:")
			fmt.Fprint(b, defaultPrefix, strings.ReplaceAll(f.Default, "\n", defaultPrefix))
		} else {
			fmt.Fprintf(b, " (default: %v)", f.Default)
		}

	case *stringValue, *stringPtrValue:
		fmt.Fprintf(b, " (default: %q)", f.Default)

	case *stringsValue:
		fmt.Fprint(b, usagePrefix, "default:")
		for _, s := range *v {
			fmt.Fprintf(b, "%s%q", defaultPrefix, s)
		}
	}

	fmt.Fprint(w, b.String(), "\n")
}

func (s *FlagSet) sortedFlags() []*Flag {
	res := make([]*Flag, len(s.flagSlice))
	copy(res, s.flagSlice)

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})
	return res
}

func (s *FlagSet) MustParse(args ...[]string) {
	s.init()

	a := os.Args[1:]
	if len(args) > 0 {
		a = args[0]
	}

	err := s.parse(a)
	if err != nil {
		panic(fmt.Errorf("parse flag set fail. args:%v, err:%v", a, err))
	}
}

func (s *FlagSet) Parse(args ...[]string) error {
	s.init()

	if len(args) > 0 {
		return s.parse(args[0])
	}

	return s.parse(os.Args[1:])
}

func (s *FlagSet) parse(args []string) error {
	err := s.command.Parse(args)
	if err != nil {
		return err
	}

	for _, f := range s.flagSlice {
		if len(f.Options) == 0 {
			continue
		}
		value := f.Value.String()
		if f.inOptions(value) {
			continue
		}
		// msg := fmt.Sprintf("invalid value %q for flag -%s", value, f.Name)
		msg := fmt.Sprintf("the value of flag -%s should be %v", f.Name, f.Options)
		fmt.Fprintln(s.command.Output(), msg)
		s.printUsage()
		return fmt.Errorf("%s", msg)
	}

	// log.Printf(" === %#v", s.args)
	if s.args != nil {
		args := s.command.Args()
		if !s.noDefault && len(args) == 0 && s.args.Default != "" {
			args = []string{s.args.Default}
		}
		// log.Printf(" === %v", args)
		s.args.val.Set(reflect.ValueOf(args))
	}
	return nil
}

func (s *FlagSet) Model(v interface{}) {
	s.init()

	val := reflect.ValueOf(v)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr && typ.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("input flags must be pointer to struct. %s", typ.String()))
	}
	val = val.Elem()
	s.model = val

	s.parseStruct(0, val)
}

func (s *FlagSet) addFlag(f *Flag) {
	if f == nil {
		return
	}

	if f.IsArgs {
		// log.Printf(" === 111 %#v", f)
		if s.args == nil || s.args.lvl > f.lvl {
			s.args = f
		}
		// log.Printf(" === 222 %#v", s.args)
		return
	}

	v, val := s.getFlagValue(f)
	if v == nil {
		return
	}
	f.Value = v

	if !s.noDefault && f.Default != "" && (val.Kind() != reflect.Slice || val.Len() == 0) {
		v.Set(f.Default)
	}

	if f.Default == "" {
		f.Default = v.String()
		if vv, _ := v.(valueWithDefault); vv != nil {
			f.Default = vv.Default()
		}
	}

	s.command.Var(v, f.Name, f.Usage)
	s.flagSlice = append(s.flagSlice, f)
	s.flagMap[f.Name] = f
	// log.Printf(" === %s: %v", typ.String(), val.Interface())
}

func (s *FlagSet) getFlagValue(f *Flag) (flagVal Value, refVal reflect.Value) {
	val := f.val
	if v, ok := val.Interface().(Value); ok {
		return v, val
	}

	if v := getAddrValue(val); v != nil {
		return v, val
	}

	// if val.Kind() == reflect.Ptr {
	// 	val = val.Elem()
	// }

	typ := val.Type()
	// log.Printf(" === 2222 %#v", m)

	fn := fillFlagFuncMap[typ]
	if fn == nil {
		return nil, val
	}

	v := fn(val)
	if f.Bool {
		switch vv := v.(type) {
		case *boolValue:
			return newBoolFlag(vv), val

		case boolPtrValue:
			return newBoolPtrFlag(vv), val
		}
	}

	return v, val
}

func (s *FlagSet) parseStruct(lvl int, val reflect.Value) {
	typ := val.Type()
	// log.Printf(" === %s", typ.String())

	for i, l := 0, typ.NumField(); i < l; i++ {
		field := typ.Field(i)
		if field.Anonymous || field.Tag.Get(tagName) == "-" {
			continue
		}

		fVal := val.Field(i)

		if _, ok := fVal.Interface().(Value); ok {
			if fVal.Kind() == reflect.Ptr && fVal.IsNil() {
				fVal.Set(reflect.New(fVal.Type().Elem()))
			}
			s.parseFieldFlag(lvl, fVal, field)
			continue
		}

		s.parseStructField(lvl, fVal, field)
	}
}

func (s *FlagSet) parseFieldFlag(lvl int, val reflect.Value, field reflect.StructField) {
	tag, ok := field.Tag.Lookup(tagName)
	if !ok || len(tag) == 0 {
		f := &Flag{
			lvl:  lvl,
			val:  val,
			tags: ldtagmap.New(),
			Name: parseFlagName(field),
		}
		s.addFlag(f)
		return
	}

	if tag == "-" {
		return
	}

	tags := ldtagmap.Parse(tag)
	// if tags.Has("-") {
	// 	return
	// }
	// log.Printf(" === %s %#v", tag, tags)

	f := &Flag{
		lvl:     lvl,
		val:     val,
		tags:    tags,
		Name:    tags.Get("name"),
		Meta:    tags.Get("meta"),
		Usage:   tags.Get("usage"),
		Default: strings.TrimSpace(tags.Get("default")),
		IsArgs:  tags.Has("args"),
		Bool:    tags.Has("bool"),
		Options: s.parseOptions(tags.Get("options")),
	}

	if f.Default == "" && len(f.Options) > 0 {
		f.Default = f.Options[0]
	}

	if len(f.Name) == 0 {
		f.Name = parseFlagName(field)
	}
	// log.Printf(" === 1111 %#v", m)

	s.addFlag(f)
	return
}

func (s *FlagSet) parseOptions(str string) []string {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}
	slice := strings.Split(str, ",")
	opts := make([]string, 0, len(slice))
	for _, v := range slice {
		v = strings.TrimSpace(v)
		if v != "" {
			opts = append(opts, v)
		}
	}
	return opts
}

func (s *FlagSet) parseStructField(lvl int, fVal reflect.Value, field reflect.StructField) {
	val := fVal
	typ := field.Type
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
		typ = typ.Elem()
		if val.IsNil() {
			val.Set(reflect.New(typ))
		}
		val = val.Elem()
	}

	if typ.Kind() == reflect.Struct {
		s.parseStruct(lvl+1, val)
		return
	}

	s.parseFieldFlag(lvl, val, field)
}
