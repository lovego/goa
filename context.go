package goa

import (
	"net/http"
)

type Context struct {
	*http.Request
	http.ResponseWriter
	handlers []handlerFunc
	params   []string
	index    int

	data   map[string]interface{}
	errors []error
}

func (c *Context) Param(i int) string {
	if i <= len(c.params) {
		return c.params[i]
	}
	return ""
}

func (c *Context) Next() {
	c.index++
	if c.index >= len(c.handlers) {
		return
	}
	c.handlers[c.index](c)
}
