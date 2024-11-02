/*
 * Copyright (C) distroy
 */

package ldref

import (
	"fmt"
	"strings"
	"sync"

	"github.com/distroy/ldgo/v3/lderr"
)

var (
	contextPool = &sync.Pool{}
)

func getContext() *context {
	c, _ := contextPool.Get().(*context)
	if c == nil {
		c = &context{
			fields: make([]string, 0, 16),
			errors: make([]string, 0, 16),
		}
	}

	return c
}

func putContext(c *context) {
	c.Reset()
	contextPool.Put(c)
}

type context struct {
	fields []string
	errors []string
}

func (c *context) Reset() {
	c.fields = c.fields[:0]
	c.errors = c.errors[:0]
}

func (c *context) Error() error {
	if len(c.errors) == 0 {
		return nil
	}

	err := lderr.ErrReflectError
	err = lderr.New(err.Status(), err.Code(), c.errors[0])
	return lderr.WithDetails(err, c.errors)
}

// AddErrorf formats according to a format specifier.
func (c *context) AddErrorf(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	c.AddError(text)
}

func (c *context) AddError(text string) {
	fields := strings.Join(c.fields, ".")
	if fields != "" {
		text = fmt.Sprintf("%s: %s", fields, text)
	}

	c.errors = append(c.errors, text)
}

func (c *context) PushField(field string) {
	c.fields = append(c.fields, field)
}

func (c *context) PopField() string {
	length := len(c.fields)
	lastIdx := length - 1

	field := c.fields[lastIdx]
	c.fields = c.fields[:lastIdx]
	return field
}
