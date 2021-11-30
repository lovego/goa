package middlewares

import (
	"context"
	"encoding/json"
	"mime"
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
		ShouldLogReqBody:  shouldLogReqBody,
		ShouldLogRespBody: shouldLogRespBody,
	}
}

func (l *Logger) Record(c *goa.Context) {
	debug := c.URL.Query()["_debug"] != nil
	l.Logger.RecordWithContext(c.Context(), func(tracerCtx context.Context) error {
		if debug {
			tracerCtx = tracer.SetDebug(tracerCtx)
		}
		if debug || l.ShouldLogReqBody != nil && l.ShouldLogReqBody(c) {
			c.RequestBody()
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

func shouldLogReqBody(c *goa.Context) bool {
	if v := c.Request.Header.Get("Content-Type"); v != "" {
		mediaType, _, _ := mime.ParseMediaType(v)
		// multipart is often use by file uploading.
		if strings.HasPrefix(mediaType, "multipart/") {
			return false
		}
	}
	return true
}

func shouldLogRespBody(c *goa.Context) bool {
	method := c.Request.Method
	if method == http.MethodGet || method == http.MethodPost &&
		strings.Contains(c.Request.URL.Path, "query") &&
		strings.Contains(c.Request.URL.Path, "search") {
		return false
	}
	return true
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
