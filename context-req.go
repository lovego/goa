package goa

import (
	"bytes"
	"io/ioutil"
	"net"
	"strings"
)

const reqBodyKey = "requestBody"

func (c *ContextBeforeLookup) ParseForm() error {
	c.RequestBody() // record the request body
	return c.Request.ParseForm()
}

func (c *ContextBeforeLookup) RequestBody() ([]byte, error) {
	if c.Request.Body == nil {
		return nil, nil
	}
	if data, ok := c.data[reqBodyKey]; ok {
		if body, ok := data.([]byte); ok {
			return body, nil
		}
		return nil, nil
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.SetError(err)
		return nil, err
	}
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[reqBodyKey] = body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

func (c *ContextBeforeLookup) Scheme() string {
	if proto := c.Request.Header.Get("X-Forwarded-Proto"); proto != `` {
		return proto
	}
	return `http`
}

func (c *ContextBeforeLookup) Origin() string {
	return c.Scheme() + "://" + c.Request.Host
}

func (c *ContextBeforeLookup) Url() string {
	return c.Origin() + c.Request.URL.RequestURI()
}

func (c *ContextBeforeLookup) ClientAddr() string {
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

func (c *ContextBeforeLookup) RequestId() string {
	return c.Request.Header.Get("X-Request-Id")
}
