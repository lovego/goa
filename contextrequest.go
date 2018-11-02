package goa

import (
	"bytes"
	"io/ioutil"
	"net"
	"strings"
)

const reqBodyKey = "requestBody"

func (c *Context) RequestBody() ([]byte, error) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if data, ok := c.data[reqBodyKey]; ok {
		if body, ok := data.([]byte); ok {
			return body, nil
		}
		return nil, nil
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.data[reqBodyKey] = body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

func (c *Context) Scheme() string {
	if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != `` {
		return proto
	}
	return `http`
}

func (c *Context) Url() string {
	return c.Scheme() + `://` + c.Request.Host + c.Request.URL.RequestURI()
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
