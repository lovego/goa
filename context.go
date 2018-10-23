package goa

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
)

type Context struct {
	*http.Request
	http.ResponseWriter
	handlers []HandlerFunc
	params   []string
	index    int

	data map[string]interface{}
	err  error
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

func (c *Context) Get(key string) interface{} {
	value, ok := c.data[key]
	if ok {
		return value
	}
	panic(fmt.Sprintf("context get :%s not exists", key))
}

func (c *Context) Status() int64 {
	status := reflect.ValueOf(c.ResponseWriter).Elem().FieldByName(`status`)
	if status.IsValid() {
		return status.Int()
	} else {
		return 0
	}
}

func (c *Context) RequestBody() []byte {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if data, ok := c.data["requestBody"]; ok {
		if body, ok := data.([]byte); ok {
			return body
		}
		return nil
	}
	body, bodyReader := readAndClone(c.Request.Body)
	c.data["requestBody"] = body
	c.Request.Body = bodyReader

	return body
}

func (c *Context) Write(content []byte) (int, error) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if data, ok := c.data["responseBody"]; ok {
		if body, ok := data.([]byte); ok {
			body = append(body, content...)
			c.data["responseBody"] = body
		}
	}
	return c.ResponseWriter.Write(content)
}

func (c *Context) ResponseBody() []byte {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if data, ok := c.data["responseBody"]; ok {
		if body, ok := data.([]byte); ok {
			return body
		}
	}
	return nil
}

func (c *Context) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := c.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
}

func (c *Context) Flush() {
	if flusher, ok := c.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func readAndClone(reader io.ReadCloser) ([]byte, io.ReadCloser) {
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("read http request body error: %v", err)
		return nil, nil
	}
	return body, ioutil.NopCloser(bytes.NewBuffer(body))
}

func (c *Context) Scheme() string {
	if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != `` {
		return proto
	}
	return `http`
}

func (c *Context) Url() string {
	return c.Scheme() + `://` + c.Request.Host + c.Request.RequestURI
}

func (c *Context) Redirect(path string) {
	c.ResponseWriter.Header().Set("Location", path)
	c.ResponseWriter.WriteHeader(302)
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
