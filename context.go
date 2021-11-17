package goa

import (
	"context"
	"net/http"

	"github.com/lovego/errs"
)

type ContextBeforeLookup struct {
	*http.Request
	http.ResponseWriter

	data map[string]interface{}
	err  error
}

type Context struct {
	ContextBeforeLookup

	handlers []func(c *Context)
	params   []string
	index    int
}

// Param returns captured subpatterns.
// index begin at 0, so pass 0 to get the first captured subpatterns.
func (c *Context) Param(i int) string {
	if i < len(c.params) {
		return c.params[i]
	}
	return ""
}

// run the next midllware or route handler.
func (c *Context) Next() {
	c.index++
	if c.index >= len(c.handlers) {
		return
	}
	c.handlers[c.index](c)
}

func (c *ContextBeforeLookup) Context() context.Context {
	if data := c.Get("context"); data != nil {
		if c, ok := data.(context.Context); ok {
			return c
		}
	}
	return c.Request.Context()
}

func (c *ContextBeforeLookup) Get(key string) interface{} {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if value, ok := c.data[key]; ok {
		return value
	}
	return nil
}

func (c *ContextBeforeLookup) Set(key string, value interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = value
}

func (c *ContextBeforeLookup) SetError(err error) {
	c.err = errs.Trace(err)
}

func (c *ContextBeforeLookup) GetError() error {
	return c.err
}
