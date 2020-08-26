package goa

import (
	"bufio"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"reflect"
)

const respBodyKey = "responseBody"

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
	data := c.data[respBodyKey]
	if data == nil {
		body := append([]byte{}, content...)
		c.data[respBodyKey] = body
	} else if body, ok := data.([]byte); ok {
		body = append(body, content...)
		c.data[respBodyKey] = body
	}
	return c.ResponseWriter.Write(content)
}

func (c *Context) ResponseBody() []byte {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	if data, ok := c.data[respBodyKey]; ok {
		if body, ok := data.([]byte); ok {
			return body
		}
	}
	return nil
}

func (c *Context) Data(data interface{}, err error) {
	statusCode := http.StatusOK
	body := struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}{}
	if err == nil {
		body.Code = `ok`
		body.Message = `success`
	} else {
		if err2, ok := err.(interface {
			Code() string
			Message() string
		}); ok && err2.Code() != "" {
			body.Code, body.Message = err2.Code(), err2.Message()

			if err3, ok := err.(interface {
				GetError() error
			}); ok && err3.GetError() != nil {
				c.SetError(err3.GetError())
			}
		} else {
			statusCode = http.StatusInternalServerError
			body.Code, body.Message = `server-err`, `Server Error.`
			c.SetError(err)
		}
	}

	if data != nil {
		body.Data = data
	} else if err != nil {
		if err2, ok := err.(interface {
			Data() interface{}
		}); ok && err2.Data() != nil {
			body.Data = err2.Data()
		}
	}
	c.JsonWithCode(body, statusCode)
}

func (c *Context) Ok(message string) {
	body := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "ok",
		Message: message,
	}
	c.JsonWithCode(body, http.StatusOK)
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
