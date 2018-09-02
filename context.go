package router

import (
	"net/http"
)

type Context struct {
	*http.Request
	http.ResponseWriter
	handlers handlersChain
	params   []string
	index    int

	data   map[string]interface{}
	errors []error
}

func (c *Context) Next() {
	c.index++
	if c.index >= len(c.handlers) {
		return
	}
	c.handlers[c.index](c)
}
