package goa

import (
	"bytes"
	"net/http"
	"reflect"
	"runtime"

	"github.com/lovego/regex_tree"
)

type Router struct {
	RouterGroup
	notFound []func(*Context)
}

func New() *Router {
	return &Router{
		RouterGroup: RouterGroup{routes: make(map[string]*regex_tree.Node)},
		notFound:    []func(*Context){defaultNotFound},
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handlers, params := r.Lookup(req.Method, req.URL.Path)
	c := &Context{Request: req, ResponseWriter: rw, handlers: handlers, params: params, index: -1}
	if len(handlers) == 0 {
		c.handlers = r.notFound
	}
	c.Next()
}

func (r *Router) Use(handlers ...func(*Context)) {
	r.RouterGroup.Use(handlers...)
	last := len(r.notFound) - 1
	notFound := r.notFound[last]
	r.notFound = append(r.notFound[:last], handlers...)
	r.notFound = append(r.notFound, notFound)
}

func (r *Router) NotFound(handler func(*Context)) {
	last := len(r.notFound) - 1
	r.notFound[last] = handler
}

func defaultNotFound(c *Context) {
	if c.ResponseWriter != nil {
		c.WriteHeader(404)
		c.Write([]byte(`{"code":"404","message":"Not Found."}`))
	}
}

func (r *Router) String() string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	buf.WriteString(r.RoutesString())
	if len(r.notFound) > 0 {
		buf.WriteString("  notFound: " + funcName(r.notFound[len(r.notFound)-1]) + "\n")
	}
	buf.WriteString("}\n")
	return buf.String()
}

func funcName(fun interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
}
