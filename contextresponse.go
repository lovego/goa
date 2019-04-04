package goa

import (
	"bufio"
	"encoding/json"
	"errors"
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
	data := c.data[rspBodyKey]
	if data == nil {
		body := append([]byte{}, content...)
		c.data[rspBodyKey] = body
	} else if body, ok := data.([]byte); ok {
		body = append(body, content...)
		c.data[rspBodyKey] = body
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

func (c *Context) JsonWithCode(data interface{}, code int) {
	c.ResponseWriter.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
	c.WriteHeader(code)
	encoder := json.NewEncoder(c)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		c.SetError(err)
		c.Write([]byte(`{"code":"json-marshal-error","message":"json marshal error"}`))
	}
}

func (c *Context) Json(data interface{}) {
	c.JsonWithCode(data, http.StatusOK)
}

func (c *Context) Json2(data interface{}, err error) {
	if err != nil {
		c.SetError(err)
	}
	c.JsonWithCode(data, http.StatusOK)
}

func (c *Context) Ok(message string) {
	result := make(map[string]interface{})
	result["code"] = "ok"
	result["message"] = message
	c.JsonWithCode(result, http.StatusOK)
}

func (c *Context) Data(data interface{}, err error) {
	c.DataWithKey(data, err, `data`)
}

func (c *Context) DataWithKey(data interface{}, err error, key string) {
	statusCode := http.StatusOK
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
				GetError() error
			}); ok && e.GetError() != nil {
				c.SetError(e.GetError())
			}
		} else {
			statusCode = http.StatusInternalServerError
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
	c.JsonWithCode(result, statusCode)
}

func (c *Context) Redirect(url string) {
	c.ResponseWriter.Header().Set("Location", url)
	c.ResponseWriter.WriteHeader(302)
}

func (c *Context) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := c.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, errors.New("the ResponseWriter doesn't support hijacking.")
}

func (c *Context) Flush() {
	if flusher, ok := c.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
