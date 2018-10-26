package middlewares

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/lovego/goa"
	loggerPkg "github.com/lovego/logger"
	"github.com/lovego/tracer"
)

type Logger struct {
	Logger        *loggerPkg.Logger
	PanicHandler  func(*goa.Context)
	ShouldLogBody func(*goa.Context) bool
}

func NewLogger(logger *loggerPkg.Logger) *Logger {
	return &Logger{
		Logger:        logger,
		PanicHandler:  defaultPanicHandler,
		ShouldLogBody: defaultShouldLogBody,
	}
}

func (l *Logger) Middleware(ctx *goa.Context) {
	debug := ctx.URL.Query()["_debug"] != nil
	l.Logger.RecordWithContext(ctx.Context(), func(tracerCtx context.Context) error {
		if debug {
			tracer.GetSpan(tracerCtx).SetDebug(true)
		}
		ctx.Set("context", tracerCtx)
		ctx.Next()
		return ctx.GetError()
	}, func() {
		if l.PanicHandler != nil {
			l.PanicHandler(ctx)
		}
	}, func(fields *loggerPkg.Fields) {
		l.setFields(fields, ctx, debug)
	})

}

func (l *Logger) setFields(f *loggerPkg.Fields, ctx *goa.Context, debug bool) {
	req := ctx.Request
	f.With("host", req.Host)
	f.With("method", req.Method)
	f.With("path", req.URL.Path)
	f.With("rawQuery", req.URL.RawQuery)
	f.With("query", req.URL.Query())
	f.With("status", ctx.Status())
	f.With("reqBodySize", req.ContentLength)
	f.With("resBodySize", ctx.ResponseBodySize())
	f.With("ip", ctx.ClientAddr())
	f.With("agent", req.UserAgent())
	f.With("refer", req.Referer())
	if sess := ctx.Get("session"); sess != nil {
		f.With("session", sess)
	}
	if debug || l.ShouldLogBody != nil && l.ShouldLogBody(ctx) {
		f.With("reqBody", tryUnmarshal(ctx.RequestBody()))
		f.With("resBody", tryUnmarshal(ctx.ResponseBody()))
	}
}

func defaultPanicHandler(ctx *goa.Context) {
	ctx.WriteHeader(500)
	if ctx.ResponseBodySize() <= 0 {
		ctx.Json(map[string]string{"code": "server-err", "message": "Fatal Server Error."})
	}
}

func defaultShouldLogBody(ctx *goa.Context) bool {
	method := ctx.Request.Method
	return method == http.MethodPost ||
		method == http.MethodDelete ||
		method == http.MethodPut
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
