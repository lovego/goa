package goa

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

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

func (c *Context) ClientAddr() string {
	if addrs := c.Request.Header.Get("X-Forwarded-For"); addrs != `` {
		addr := strings.SplitN(addrs, `, `, 2)[0]
		if addr != `unknown` {
			return addr
		}
	}
	if addr := c.Request.Header.Get("X-Real-IP"); addr != `` && addr != `unknown` {
		return addr
	}
	host, _, _ := net.SplitHostPort(c.RemoteAddr)
	return host
}
