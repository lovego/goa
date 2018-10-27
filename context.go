package goa

import (
	"context"
	"net/http"
)

type Context struct {
	*http.Request
	http.ResponseWriter
	handlers HandlerFuncs
	params   []string
	index    int

	data map[string]interface{}
	err  error
}

func (c *Context) Param(i int) string {
	if i < len(c.params) {
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

func (c *Context) Context() context.Context {
	if data := c.Get("context"); data != nil {
		if ctx, ok := data.(context.Context); ok {
			return ctx
		}
	}
	return c.Request.Context()
}

func (c *Context) Get(key string) interface{} {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if value, ok := c.data[key]; ok {
		return value
	}
	return nil
}

func (c *Context) Set(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
}

func (ctx *Context) SetError(err error) {
	ctx.err = err
}

func (ctx *Context) GetError() error {
	return ctx.err
}
