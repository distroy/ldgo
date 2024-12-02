/*
 * Copyright (C) distroy
 */

package ldstr

import (
	"unicode"
	"unicode/utf8"

	"github.com/distroy/ldgo/v2/ldconv"
	"github.com/distroy/ldgo/v2/ldsort"
)

func ToSnakeCase(s string, sep ...rune) string {
	if s == "" {
		return s
	}
	c := rune('_')
	if len(sep) > 0 && sep[0] != 0 {
		c = sep[0]
	}
	return ldconv.BytesToStrUnsafe(toSnakeCase(ldconv.StrToBytesUnsafe(s), c))
}

func toSnakeCase(s []byte, sep rune) []byte {
	reader := &buffer{buf: s}
	writer := make([]byte, 0, len(s)*2)

	lastSep := true
	writeSep := func() {
		if lastSep {
			return
		}
		writer = utf8.AppendRune(writer, sep)
		lastSep = true
	}
	writeRune := func(c rune) {
		if c == sep {
			writeSep()
			return
		}
		writer = utf8.AppendRune(writer, c)
		lastSep = false
	}

	loopReadRunes(reader, func(i, size int, curr rune) bool {
		if unicode.IsSpace(curr) {
			writeRune(sep)
			return true

		} else if !unicode.IsUpper(curr) {
			writeRune(curr)
			return true

		} else if i > 0 {
			writeRune(sep)
		}

		last := unicode.ToLower(curr)
		loopReadRunes(reader, func(i, _ int, curr rune) bool {
			if unicode.IsSpace(curr) {
				writeRune(last)
				last = sep
				return false

			} else if unicode.IsUpper(curr) {
				writeRune(last)
				last = unicode.ToLower(curr)
				return true

			} else if i > 0 {
				writeRune(sep)
			}

			writeRune(last)
			last = unicode.ToLower(curr)
			return false
		})
		writeRune(last)
		return true
	})

	raw := writer
	if l := len(raw); rune(raw[l-1]) == sep {
		raw = raw[:l-1]
	}
	return raw
}

func ToCamelCase(s string, seps ...[]rune) string {
	if s == "" {
		return s
	}
	sep := []rune{'-', '_'}
	if len(seps) > 0 && len(seps[0]) > 0 {
		// ldsort.SortInt32s(seps[0])
		// return toCamelCase(s, seps[0])
		sep = seps[0]
		ldsort.SortInt32s(sep)
	}
	return ldconv.BytesToStrUnsafe(toCamelCase(ldconv.StrToBytesUnsafe(s), sep))
}

func toCamelCase(s []byte, seps []rune) []byte {
	reader := &buffer{buf: s}
	writer := make([]byte, 0, len(s))

	lastSep := true
	loopReadRunes(reader, func(i, size int, curr rune) bool {
		if unicode.IsSpace(curr) || ldsort.IndexInt32s(seps, curr) >= 0 {
			lastSep = true
			return true

		} else if unicode.IsUpper(curr) {
			lastSep = false
			writer = utf8.AppendRune(writer, curr)
			return true

		} else if lastSep {
			lastSep = false
			curr = unicode.ToUpper(curr)
			writer = utf8.AppendRune(writer, curr)
			return true
		}

		lastSep = false
		writer = utf8.AppendRune(writer, curr)
		return true
	})

	return writer
}

type buffer struct {
	buf []byte
	idx int
}

func loopReadRunes(buf *buffer, fBody func(i int, size int, c rune) bool) {
	l := len(buf.buf)
	for i := 0; buf.idx < l; i++ {
		c, size := utf8.DecodeRune(buf.buf[buf.idx:])
		buf.idx += size
		if ok := fBody(i, size, c); !ok {
			break
		}
	}
}
