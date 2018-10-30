package goa

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/lovego/regex_tree"
)

type HandlerFunc func(*Context)

func (h HandlerFunc) String() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

type HandlerFuncs []HandlerFunc

func (hs HandlerFuncs) String() string {
	return hs.StringIndent("")
}

func (hs HandlerFuncs) StringIndent(indent string) string {
	if len(hs) == 0 {
		return "[ ]"
	}
	var buf bytes.Buffer
	buf.WriteString("[\n")
	for _, h := range hs {
		buf.WriteString(indent + "  " + fmt.Sprint(h) + "\n")
	}
	buf.WriteString(indent + "]")
	return buf.String()
}

type Router struct {
	RouterGroup
	notFound     HandlerFunc
	fullNotFound HandlerFuncs
}

func New() *Router {
	return &Router{
		RouterGroup:  RouterGroup{routes: make(map[string]*regex_tree.Node)},
		notFound:     defaultNotFound,
		fullNotFound: HandlerFuncs{defaultNotFound},
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handlers, params := r.Lookup(req.Method, req.URL.Path)
	c := &Context{Request: req, ResponseWriter: rw, handlers: handlers, params: params, index: -1}
	if len(handlers) == 0 {
		c.handlers = r.fullNotFound
	}
	c.Next()
}

func (r *Router) Use(handlers ...HandlerFunc) {
	r.handlers = append(r.handlers, handlers...)
	r.fullNotFound = r.concatHandlers(r.notFound)
}

func (r *Router) NotFound(handler HandlerFunc) {
	r.notFound = handler
	r.fullNotFound = r.concatHandlers(r.notFound)
}

func defaultNotFound(c *Context) {
	if c.ResponseWriter != nil {
		c.ResponseWriter.WriteHeader(404)
		c.ResponseWriter.Write([]byte(`{"code":"404","message":"Not Found."}`))
	}
}

func (r *Router) String() string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	buf.WriteString(r.RoutesString())
	if r.notFound != nil {
		buf.WriteString("  notFound: " + r.notFound.String() + "\n")
	}
	buf.WriteString("}\n")
	return buf.String()
}
