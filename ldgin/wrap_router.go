/*
 * Copyright (C) distroy
 */

package ldgin

import (
	"net/http"
	"path"
	"strings"

	"github.com/distroy/ldgo/v3/ldctx"
	"go.uber.org/zap"
)

// Router is http router
type Router struct {
	router routerAdapter
}

func wrapRouter(r routerAdapter) Router {
	return Router{
		router: r,
	}
}

func (w Router) Group(relativePath string, midwares ...Midware) Router {
	w.router = w.router.Group(relativePath, midwares...)
	return w
}

func (w Router) Use(midwares ...Midware) Router {
	w.router = w.router.Use(midwares...)
	return w
}

func (w Router) WithAppPath(path string) Router {
	w.router = w.router.WithAppPath(path)
	return w
}

func (w Router) BasePath() string {
	return w.router.BasePath()
}

func (w Router) Handle(method, path string, h Handler, ms ...Midware) {
	w.router.Handle(method, path, h, ms...)
}

func (w Router) GET(p string, h Handler, ms ...Midware)    { w.Handle(http.MethodGet, p, h, ms...) }
func (w Router) DELETE(p string, h Handler, ms ...Midware) { w.Handle(http.MethodDelete, p, h, ms...) }
func (w Router) HEAD(p string, h Handler, ms ...Midware)   { w.Handle(http.MethodHead, p, h, ms...) }
func (w Router) PATCH(p string, h Handler, ms ...Midware)  { w.Handle(http.MethodPatch, p, h, ms...) }
func (w Router) POST(p string, h Handler, ms ...Midware)   { w.Handle(http.MethodPost, p, h, ms...) }
func (w Router) PUT(p string, h Handler, ms ...Midware)    { w.Handle(http.MethodPut, p, h, ms...) }
func (w Router) OPTIONS(p string, h Handler, ms ...Midware) {
	w.Handle(http.MethodOptions, p, h, ms...)
}

func (w Router) Match(methods []string, path string, h Handler, ms ...Midware) {
	for _, method := range methods {
		w.Handle(method, path, h, ms)
	}
}

func (w Router) Any(p string, h Handler, ms ...Midware) {
	anyMethods := []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
	w.Match(anyMethods, p, h, ms...)
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
// router.StaticFile("favicon.ico", "./resources/favicon.ico")
func (w Router) StaticFile(relativePath, filepath string) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	handler := func(c *Context) {
		c.File(filepath)
	}
	w.Handle(http.MethodGet, relativePath, handler)
	w.Handle(http.MethodHead, relativePath, handler)
}

// Static serves files from the given file system root.
// Internally a http.FileServer is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// To use the operating system's file system implementation,
// use :
//
//	router.Static("/static", "/var/www")
func (w Router) Static(relativePath, root string) {
	w.StaticFS(relativePath, Dir(root, false))
}

// StaticFS works just like `Static()` but a custom `http.FileSystem` can be used instead.
// GinUtil by default user: Dir()
func (w Router) StaticFS(relativePath string, fs http.FileSystem) {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := w.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET and HEAD handlers
	w.Handle(http.MethodGet, urlPattern, handler)
	w.Handle(http.MethodHead, urlPattern, handler)
	// return w
}

func (w Router) createStaticHandler(relativePath string, fs http.FileSystem) func(c *Context) {
	fullPath := w.router.calculateFullPath(relativePath)
	fileServer := http.StripPrefix(fullPath, http.FileServer(fs))

	return func(c *Context) {
		file := c.Param("filepath")
		ldctx.LogI(c, "static fs", zap.String("filepath", file))

		if _, noListing := fs.(onlyFilesFS); noListing {
			c.Writer.WriteHeader(http.StatusNotFound)
		}

		// Check if file exists and/or if we have permission to access it
		f, err := fs.Open(file)
		if err != nil {
			ldctx.LogE(c, "filepath is not exists", zap.String("filepath", file))
			c.Writer.WriteHeader(http.StatusNotFound)
			return
		}
		f.Close()

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}
