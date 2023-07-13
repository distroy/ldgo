/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"bytes"
	"io"
	"strconv"

	"github.com/distroy/ldgo/ldconv"
	"github.com/distroy/ldgo/lderr"
	"go.uber.org/zap"
)

var (
	lineEnd = []byte("\r\n")
)

type ChunkedRenderer struct {
	Code        int
	ContentType string
	Headers     map[string]string
	Reader      io.ReadCloser
}

func (r ChunkedRenderer) Render(c Context) {
	reader := r.Reader
	defer func() {
		reader.Close()
	}()

	c.Status(r.Code)
	r.writeHeaders(c)

	// size := 4096
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			return

		} else if err != nil {
			c.LogE("read chunk fail", zap.Error(err))
			e := lderr.WithDetail(lderr.ErrHttpRenderBody, err.Error())
			c.AbortWithError(e)
			return
		}

		chunk := buf[:n]
		if err := r.writeChunk(c, chunk); err != nil {
			// c.LogE("write chunk fail", zap.ByteString("chunk", chunk), zap.Error(err))
			break
		}
	}
}

func (r ChunkedRenderer) writeChunk(c Context, chunk []byte) lderr.Error {
	writer := c.Writer

	size := len(chunk)
	sizeBytes := ldconv.StrToBytesUnsafe(strconv.Itoa(size))

	if _, err := io.Copy(writer, bytes.NewBuffer(sizeBytes)); err != nil {
		c.LogE("write chunk fail", zap.Int("size", size), zap.ByteString("chunk", chunk), zap.Error(err))
		return lderr.WithDetail(lderr.ErrHttpRenderBody, err.Error())
	}

	if _, err := io.Copy(writer, bytes.NewBuffer(lineEnd)); err != nil {
		c.LogE("write chunk fail", zap.Int("size", size), zap.ByteString("chunk", chunk), zap.Error(err))
		return lderr.WithDetail(lderr.ErrHttpRenderBody, err.Error())
	}

	if size > 0 {
		if _, err := io.Copy(writer, bytes.NewBuffer(chunk)); err != nil {
			c.LogE("write chunk fail", zap.Int("size", size), zap.ByteString("chunk", chunk), zap.Error(err))
			return lderr.WithDetail(lderr.ErrHttpRenderBody, err.Error())
		}
	}

	if _, err := io.Copy(writer, bytes.NewBuffer(lineEnd)); err != nil {
		c.LogE("write chunk fail", zap.Int("size", size), zap.ByteString("chunk", chunk), zap.Error(err))
		return lderr.WithDetail(lderr.ErrHttpRenderBody, err.Error())
	}

	return nil
}

func (r ChunkedRenderer) writeHeaders(c Context) {
	header := c.Writer.Header()
	for k, v := range r.Headers {
		if v != "" {
			header.Set(k, v)
		}
	}
	if r.ContentType != "" {
		header.Set("Content-Type", r.ContentType)
	}
	header.Set("Transfer-Encoding", "chunked")
}
