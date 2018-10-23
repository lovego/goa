package goa

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
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
