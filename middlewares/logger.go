package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/lovego/goa"
	loggerPkg "github.com/lovego/logger"
	"github.com/lovego/tracer"
)

type Logger struct {
	Logger            *loggerPkg.Logger
	PanicHandler      func(*goa.Context)
	ShouldLogReqBody  func(*goa.Context) bool
	ShouldLogRespBody func(*goa.Context) bool
}

func NewLogger(logger *loggerPkg.Logger) *Logger {
	return &Logger{
		Logger:            logger,
		PanicHandler:      defaultPanicHandler,
		ShouldLogReqBody:  defaultShouldLogBody,
		ShouldLogRespBody: defaultShouldLogBody,
	}
}

func (l *Logger) Record(c *goa.Context) {
	debug := c.URL.Query()["_debug"] != nil
	l.Logger.RecordWithContext(c.Context(), func(tracerCtx context.Context) error {
		if debug {
			tracerCtx = tracer.SetDebug(tracerCtx)
		}
		c.Set("context", tracerCtx)
		c.Next()
		return c.GetError()
	}, func() {
		if l.PanicHandler != nil {
			l.PanicHandler(c)
		}
	}, func(fields *loggerPkg.Fields) {
		l.setFields(fields, c, debug)
	})

}

func (l *Logger) setFields(f *loggerPkg.Fields, c *goa.Context, debug bool) {
	req := c.Request
	f.With("requestId", c.RequestId())
	f.With("host", req.Host)
	f.With("method", req.Method)
	f.With("path", req.URL.Path)
	f.With("rawQuery", req.URL.RawQuery)
	// f.With("query", req.URL.Query())
	f.With("status", c.Status())
	f.With("reqBodySize", req.ContentLength)
	f.With("respBodySize", c.ResponseBodySize())
	f.With("ip", c.ClientAddr())
	f.With("agent", req.UserAgent())
	f.With("refer", req.Referer())
	if sess := c.Get("session"); sess != nil {
		f.With("session", sess)
	}
	if debug || l.ShouldLogReqBody != nil && l.ShouldLogReqBody(c) {
		reqBody, err := c.RequestBody()
		if err != nil {
			if c.GetError() == nil {
				c.SetError(err)
			}
		} else {
			f.With("reqBody", string(reqBody))
		}
	}
	if debug || l.ShouldLogRespBody != nil && l.ShouldLogRespBody(c) {
		f.With("respBody", string(c.ResponseBody()))
	}
}

func defaultPanicHandler(c *goa.Context) {
	if c.ResponseBodySize() <= 0 {
		c.StatusJson(500, map[string]string{"code": "server-err", "message": "Fatal Server Error."})
	} else {
		c.WriteHeader(500)
	}
}

func defaultShouldLogBody(c *goa.Context) bool {
	method := c.Request.Method
	if method == http.MethodPatch ||
		method == http.MethodPut ||
		method == http.MethodDelete {
		return true
	}
	if method == http.MethodPost &&
		!strings.HasSuffix(c.Request.URL.Path, "query") &&
		!strings.HasSuffix(c.Request.URL.Path, "search") {
		return true
	}
	return false
}

func tryUnmarshal(b []byte) interface{} {
	var v map[string]interface{}
	err := json.Unmarshal(b, &v)
	if err == nil {
		return v
	} else {
		return string(b)
	}
}
