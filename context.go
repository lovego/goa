package goa

import (
	"net/http"
    "fmt"
    "strings"
    "net"
)

type Context struct {
	*http.Request
	http.ResponseWriter
	handlers []HandlerFunc
	params   []string
	index    int

	data   map[string]interface{}
	err error
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

func (c *Context) Set(key string, value interface{}) {
    c.data[key] = value
}

func (c *Context) Get(key string)(interface{}){
    value, ok := c.data[key]
    if ok{
        return value
    }
    panic(fmt.Sprintf("context get :%s not exists", key))
}

func (ctx *Context) ClientAddr() string {
    if addrs := ctx.Request.Header.Get("X-Forwarded-For"); addrs != `` {
        addr := strings.SplitN(addrs, `, `, 2)[0]
        if addr != `unknown` {
            return addr
        }
    }
    if addr := ctx.Request.Header.Get("X-Real-IP"); addr != `` && addr != `unknown` {
        return addr
    }
    host, _, _ := net.SplitHostPort(ctx.RemoteAddr)
    return host
}
