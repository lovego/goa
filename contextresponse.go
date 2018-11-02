package goa

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"reflect"
)

const rspBodyKey = "responseBody"

func (c *Context) Status() int64 {
	status := reflect.ValueOf(c.ResponseWriter).Elem().FieldByName(`status`)
	if status.IsValid() {
		return status.Int()
	} else {
		return 0
	}
}

func (c *Context) ResponseBodySize() int64 {
	s := reflect.ValueOf(c.ResponseWriter).Elem().FieldByName(`written`)
	if s.IsValid() {
		return s.Int()
	} else {
		return 0
	}
}

func (c *Context) Write(content []byte) (int, error) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if data, ok := c.data[rspBodyKey]; ok {
		if body, ok := data.([]byte); ok {
			body = append(body, content...)
			c.data[rspBodyKey] = body
		}
	}
	return c.ResponseWriter.Write(content)
}

func (c *Context) ResponseBody() []byte {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if data, ok := c.data[rspBodyKey]; ok {
		if body, ok := data.([]byte); ok {
			return body
		}
	}
	return nil
}

func (c *Context) Json(data interface{}) {
	bytes, err := json.Marshal(data)
	if err == nil {
		c.ResponseWriter.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
		c.Write(bytes)
	} else {
		log.Panic(err)
	}
}

func (c *Context) Json2(data interface{}, err error) {
	if err != nil {
		c.SetError(err)
	}
	c.Json(data)
}

func (c *Context) Ok(message string) {
	result := make(map[string]interface{})
	result["code"] = "ok"
	result["message"] = message
	c.Json(result)
}

func (c *Context) Data(data interface{}, err error) {
	c.DataWithKey(data, err, `data`)
}

func (c *Context) DataWithKey(data interface{}, err error, key string) {
	result := make(map[string]interface{})
	if err == nil {
		result[`code`] = `ok`
		result[`message`] = `success`
	} else {
		if erro, ok := err.(interface {
			Code() string
			Message() string
		}); ok && erro.Code() != "" {
			result[`code`] = erro.Code()
			result[`message`] = erro.Message()
			if e, ok := err.(interface {
				Err() error
			}); ok && e.Err() != nil {
				c.SetError(err)
			}
		} else {
			c.WriteHeader(500)
			result[`code`] = `server-err`
			result[`message`] = `Server Error.`
			c.SetError(err)
		}
	}

	if data != nil {
		result[key] = data
	} else if err != nil {
		if erro, ok := err.(interface {
			Data() interface{}
		}); ok && erro.Data() != nil {
			result[key] = erro.Data()
		}
	}
	c.Json(result)
}

func (c *Context) Redirect(path string) {
	c.ResponseWriter.Header().Set("Location", path)
	c.ResponseWriter.WriteHeader(302)
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
