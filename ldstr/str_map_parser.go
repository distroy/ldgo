/*
 * Copyright (C) distroy
 */

package ldstr

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/distroy/ldgo/v3/lderr"
)

func mustMarshalJson(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

var (
	strMapFieldsPool = &sync.Pool{
		New: func() interface{} {
			return make([]strMapField, 0, 8)
		},
	}

	strMapParseContextPool = &sync.Pool{
		New: func() interface{} {
			return &strMapParseContext{
				Fields: make([]strMapFieldContext, 0, 8),
			}
		},
	}
)

type strMapField struct {
	Key       string // key
	After     string // after text
	DupIndex  int    // duplicate key index
	DupFirst  bool   // is first duplicate key
	PrevIndex int    // prev index
	NextIndex int    // next index
	Skip      bool
}

type strMapFieldContext struct {
	Field     *strMapField
	Prev      *strMapFieldContext `json:"-"`
	Next      *strMapFieldContext `json:"-"`
	Duplicate *strMapFieldContext `json:"-"`

	Text      string
	Value     string
	LastIndex int
}

type strMapParseContext struct {
	Fields []strMapFieldContext
}

type StrMapParser struct {
	template    string
	left, right string

	before string
	fields []strMapField
}

func (p *StrMapParser) Template() string { return p.template }
func (p *StrMapParser) String() string   { return p.template }

func (p *StrMapParser) Done() {
	p.delFieldsBuf(p.fields)
	p.template = ""
	p.left = ""
	p.right = ""
	p.fields = nil
	p.before = ""
}

func (p *StrMapParser) Init(tmpl string, splits ...string) error {
	l, r := getReplaceSplits(splits)

	p.Done()
	p.fields = p.newFieldsBuf()

	p.template = tmpl
	p.left = l
	p.right = r

	if err := p.initFieldsByTemplate(); err != nil {
		return err
	}

	if err := p.initFieldsDuplicate(); err != nil {
		return err
	}

	if err := p.initFieldsPrevIndex(); err != nil {
		return err
	}
	return nil
}

func (p *StrMapParser) initFieldsByTemplate() error {
	tmpl := p.template
	l, r := p.left, p.right

	for {
		lIdx := strings.Index(tmpl, l)
		if lIdx < 0 {
			break
		}
		rIdx := strings.Index(tmpl[lIdx+1:], r)
		if rIdx < 0 {
			break
		}

		rIdx = lIdx + 1 + rIdx
		lIdx = strings.LastIndex(tmpl[:rIdx], l)

		before := tmpl[:lIdx]
		key := tmpl[lIdx+1 : rIdx]
		if len(p.fields) == 0 {
			p.before = before

		} else {
			p.fields[len(p.fields)-1].After = before
		}

		p.fields = append(p.fields, strMapField{
			Key:       key,
			After:     "",
			PrevIndex: -1,
			NextIndex: -1,
			DupIndex:  -1,
		})

		tmpl = tmpl[rIdx+1:]
	}

	if len(p.fields) == 0 {
		p.before = tmpl
	} else {
		p.fields[len(p.fields)-1].After = tmpl
	}

	return nil
}

func (p *StrMapParser) initFieldsDuplicate() error {
	if len(p.fields) <= 32 {
		return p.initFieldsDuplicateSmall()
	}
	return p.initFieldsDuplicateBig()
}
func (p *StrMapParser) initFieldsDuplicateBig() error {
	fieldKeys := make(map[string]int)
	for i := range p.fields {
		v := &p.fields[i]
		v.DupIndex = -1

		if v.Key == "" {
			continue
		}

		if idx, ok := fieldKeys[v.Key]; ok {
			v.DupIndex = idx
			p.fields[idx].DupFirst = true
			continue
		}

		fieldKeys[v.Key] = i
	}
	return nil
}
func (p *StrMapParser) initFieldsDuplicateSmall() error {
	for i := 1; i < len(p.fields); i++ {
		vi := &p.fields[i]
		vi.DupIndex = -1

		if vi.Key == "" {
			continue
		}

		for j := 0; j < i; j++ {
			vj := &p.fields[j]
			if vi.Key == vj.Key {
				vi.DupIndex = j
				p.fields[j].DupFirst = true
				break
			}
		}
	}

	return nil
}

func (p *StrMapParser) initFieldsPrevIndex() error {
	for i := 1; i < len(p.fields); i++ {
		vi := &p.fields[i]
		prev := &p.fields[i-1]
		vi.PrevIndex = i - 1
		prev.NextIndex = i
	}

	// log.Printf(" === abc %s", mustMarshalJson(p.fields))
	for i := 0; i < len(p.fields); i++ {
		vi := &p.fields[i]
		if vi.After != "" || vi.DupIndex >= 0 || vi.DupFirst {
			continue
		}

		j := p.indexFieldNext(i)
		i = j
	}
	return nil
}

func (p *StrMapParser) indexFieldNext(pos int) int {
	curr := &p.fields[pos]

	i := pos + 1
	for ; i < len(p.fields); i++ {
		vi := &p.fields[i]
		if vi.After == "" && vi.DupIndex < 0 && !vi.DupFirst {
			vi.Skip = true
			vi.PrevIndex = -1
			vi.NextIndex = -1
			continue
		}

		break
	}

	if i >= len(p.fields) {
		return i
	}

	vi := &p.fields[i]
	if vi.DupIndex >= 0 || vi.DupFirst {
		curr.NextIndex = i
		vi.PrevIndex = pos
		return i
	}

	curr.After = vi.After
	vi.After = ""
	vi.Skip = true
	vi.PrevIndex = -1
	vi.NextIndex = -1

	if i+1 < len(p.fields) {
		next := &p.fields[i+1]

		curr.NextIndex = i + 1
		next.PrevIndex = pos
	}

	return i + 1
}

func (p *StrMapParser) Parse(text string) (map[string]string, error) {
	if len(p.fields) == 0 {
		if p.before != text {
			return nil, lderr.ErrInvalidTemplateSyntax
		}

		return make(map[string]string), nil
	}

	if !strings.HasPrefix(text, p.before) {
		return nil, lderr.ErrInvalidTemplateSyntax
	}

	text = text[len(p.before):]

	c := p.newContext()
	defer p.delContext(c)

	head := &c.Fields[0]
	head.LastIndex = -1
	head.Text = text

	for field := head; ; {
		newField, ok := p.parseField(field)
		// log.Printf("=== %#v", c.Fields)
		// log.Printf("=== %s", mustMarshalJson(c.Fields))

		if newField == nil {
			if ok {
				break
			}
			return nil, lderr.ErrInvalidTemplateSyntax
		}

		field = newField
	}

	res := make(map[string]string, len(p.fields))
	for i := range c.Fields {
		v := &c.Fields[i]
		if v.Field.Key == "" {
			continue
		}

		res[v.Field.Key] = v.Value
	}
	return res, nil
}

func (p *StrMapParser) parseField(c *strMapFieldContext) (*strMapFieldContext, bool) {
	if c.Duplicate != nil {
		return p.parseFieldDuplicate(c)
	}

	after := c.Field.After
	if after == "" {
		return p.parseFieldEmptyAfter(c)
	}

	text := c.Text
	if c.LastIndex >= 0 {
		pos := c.LastIndex + len(after) - 1
		text = text[:pos]
	}

	idx := strings.LastIndex(text, after)
	if idx < 0 {
		return p.prevField(c)
	}

	c.Value = text[:idx]
	c.LastIndex = idx

	return p.nextField(c)
}
func (p *StrMapParser) parseFieldDuplicate(c *strMapFieldContext) (*strMapFieldContext, bool) {
	after := c.Field.After
	text := c.Text

	if c.LastIndex >= 0 {
		return p.nextField(c)
	}

	c.Value = c.Duplicate.Value
	idx := len(c.Value)
	if !strings.HasPrefix(text, c.Value) {
		return p.prevField(c)

	} else if !strings.HasPrefix(text[idx:], after) {
		return p.prevField(c)
	}

	c.LastIndex = idx

	return p.nextField(c)
}
func (p *StrMapParser) parseFieldEmptyAfter(c *strMapFieldContext) (*strMapFieldContext, bool) {
	text := c.Text
	idx := c.LastIndex

	if idx < 0 {
		idx = len(text)
	} else {
		idx--
	}

	if idx < 0 {
		return p.prevField(c)
	}

	c.LastIndex = idx
	c.Value = text[:idx]
	return p.nextField(c)
}

func (p *StrMapParser) nextField(c *strMapFieldContext) (*strMapFieldContext, bool) {
	nextIdx := c.LastIndex + len(c.Field.After)
	nextText := c.Text[nextIdx:]

	next := c.Next
	if next != nil {
		next.LastIndex = -1
		next.Text = nextText
		return next, true
	}

	if nextText == "" {
		return nil, true
	}

	return c, false
}

func (p *StrMapParser) prevField(c *strMapFieldContext) (*strMapFieldContext, bool) {
	return c.Prev, false
}

func (p *StrMapParser) newFieldsBuf() []strMapField {
	buf := strMapFieldsPool.Get().([]strMapField)
	buf = buf[:0]
	return buf
}

func (p *StrMapParser) delFieldsBuf(buf []strMapField) {
	if buf != nil {
		strMapFieldsPool.Put(buf)
	}
}

func (p *StrMapParser) newContext() *strMapParseContext {
	c := strMapParseContextPool.Get().(*strMapParseContext)

	c.Fields = c.Fields[:0]

	for i := range p.fields {
		v := &p.fields[i]
		c.Fields = append(c.Fields, strMapFieldContext{
			Field:     v,
			Prev:      nil,
			Next:      nil,
			Duplicate: nil,
		})
	}

	for i := 1; i < len(c.Fields); i++ {
		curr := &c.Fields[i]
		if idx := curr.Field.PrevIndex; idx >= 0 {
			prev := &c.Fields[idx]
			prev.Next = curr
			curr.Prev = prev
		}

		if idx := curr.Field.DupIndex; idx >= 0 {
			curr.Duplicate = &c.Fields[idx]
		}
	}

	return c
}
func (p *StrMapParser) delContext(c *strMapParseContext) {
	c.Fields = c.Fields[:0]
	strMapParseContextPool.Put(c)
}
