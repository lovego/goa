package goa

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/lovego/regex_tree"
    "github.com/lovego/tracer"
    "time"
)

type HandlerFunc func(*Context)

func (h HandlerFunc) String() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

type Router struct {
	Group
	notFound     HandlerFunc
}

func New() *Router {
	return &Router{
		Group:        Group{routes: make(map[string]*regex_tree.Node)},
		notFound:     defaultNotFound,
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handlers, params := r.Lookup(req.Method, req.URL.Path)
	ctx := &Context{ResponseWriter: rw, handlers: handlers, params: params, index: -1}
    span := &tracer.Span{At: time.Now()}
	tracerCtx := tracer.Context(req.Context(), span)
	ctx.Request = req.WithContext(tracerCtx)
	if len(handlers) == 0 {
	    r.notFound(ctx)
	    return
	}
	ctx.Next()
}

func (r *Router) Use(handlers ...HandlerFunc) {
	r.handlers = append(r.handlers, handlers...)
}

func (r *Router) NotFound(handler HandlerFunc) {
	r.notFound = handler
}

func defaultNotFound(ctx *Context) {
	if ctx.ResponseWriter != nil {
		ctx.ResponseWriter.WriteHeader(404)
		ctx.ResponseWriter.Write([]byte(`{"code":"404","message":"Not Found."}`))
	}
}
