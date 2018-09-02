package goa

import (
	"net/http"
	"reflect"
	"runtime"
)

type handlerFunc func(*Context)

func (h handlerFunc) String() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

type Router struct {
	Group
	notFound     handlerFunc
	fullNotFound []handlerFunc
}

func New() *Router {
	return &Router{
		Group:        Group{routes: make(map[string]*node)},
		notFound:     defaultNotFound,
		fullNotFound: []handlerFunc{defaultNotFound},
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handlers, params := r.Lookup(req.Method, req.URL.Path)
	ctx := &Context{Request: req, ResponseWriter: rw, handlers: handlers, params: params, index: -1}
	if len(handlers) == 0 {
		ctx.handlers = r.fullNotFound
	}
	ctx.Next()
}

func (r *Router) Use(handlers ...handlerFunc) {
	r.handlers = append(r.handlers, handlers...)
	r.fullNotFound = r.combineHandlers(r.notFound)
}

func (r *Router) NotFound(handler handlerFunc) {
	r.notFound = handler
	r.fullNotFound = r.combineHandlers(r.notFound)
}

func defaultNotFound(ctx *Context) {
	if ctx.ResponseWriter != nil {
		ctx.ResponseWriter.WriteHeader(404)
		ctx.ResponseWriter.Write([]byte(`{"code":"404","message":"Not Found."}`))
	}
}
