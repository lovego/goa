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

func (c *ContextBeforeLookup) Status() int64 {
	status := reflect.ValueOf(c.ResponseWriter).Elem().FieldByName(`status`)
	if status.IsValid() {
		return status.Int()
	} else {
		return 0
	}
}

func (c *ContextBeforeLookup) ResponseBodySize() int64 {
	s := reflect.ValueOf(c.ResponseWriter).Elem().FieldByName(`written`)
	if s.IsValid() {
		return s.Int()
	} else {
		return 0
	}
}

func (c *ContextBeforeLookup) Write(content []byte) (int, error) {
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

func (c *ContextBeforeLookup) ResponseBody() []byte {
	if data, ok := c.data[respBodyKey]; ok {
		if body, ok := data.([]byte); ok {
			return body
		}
	}
	return nil
}

func (c *ContextBeforeLookup) Data(data interface{}, err error) {
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

	if err != nil {
		if err2, ok := err.(interface {
			Data() interface{}
		}); ok && err2.Data() != nil {
			body.Data = err2.Data()
		} else if data != nil && !isNilValue(data) { // 避免返回"data": null
			body.Data = data
		}
	} else if data != nil && !isNilValue(data) { // 避免返回"data": null
		body.Data = data
	}
	c.StatusJson(statusCode, body)
}

func isNilValue(itfc interface{}) bool {
	v := reflect.ValueOf(itfc)
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface:
		return v.IsNil()
	}
	return false
}

func (c *ContextBeforeLookup) Ok(message string) {
	body := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "ok",
		Message: message,
	}
	c.StatusJson(http.StatusOK, body)
}

func (c *ContextBeforeLookup) Json(data interface{}) {
	c.StatusJson(http.StatusOK, data)
}

func (c *ContextBeforeLookup) Json2(data interface{}, err error) {
	if err != nil {
		c.SetError(err)
	}
	c.StatusJson(http.StatusOK, data)
}

func (c *ContextBeforeLookup) StatusJson(status int, data interface{}) {
	// header should be set before WriteHeader or Write
	c.ResponseWriter.Header().Set(`Content-Type`, `application/json; charset=utf-8`)
	c.WriteHeader(status)

	encoder := json.NewEncoder(c)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(data); err != nil {
		c.SetError(err)
		c.Write([]byte(`{"code":"json-marshal-error","message":"json marshal error"}`))
	}
}

func (c *ContextBeforeLookup) Redirect(url string) {
	c.ResponseWriter.Header().Set("Location", url)
	c.ResponseWriter.WriteHeader(302)
}

func (c *ContextBeforeLookup) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := c.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, errors.New("the ResponseWriter doesn't support hijacking.")
}

func (c *ContextBeforeLookup) Flush() {
	if flusher, ok := c.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
