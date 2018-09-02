package goa

import (
	"net/http"
	"reflect"
	"regexp/syntax"
	"runtime"
	"strings"
)

type handlerFunc func(*Context)

func (h handlerFunc) String() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

type Router struct {
	basePath string
	handlers []handlerFunc
	routes   map[string]*node
	notFound []handlerFunc
}

func New() *Router {
	return &Router{routes: make(map[string]*node)}
}

func (r *Router) Add(method, path string, handler handlerFunc) *Router {
	method = strings.ToUpper(method)
	path = cleanPath(r.basePath + path)
	handlers := r.getHandlers(handler)
	rootNode := r.routes[method]
	if rootNode == nil {
		r.routes[method] = newNode(path, handlers)
	} else if rootNode.add(path, handlers) == addResultConflict {
		panic("router conflict: " + method + " " + path)
	}
	return r
}

func (r *Router) getHandlers(handler handlerFunc) []handlerFunc {
	if handler == nil {
		panic("handler func should not be nil")
	}
	var handlers []handlerFunc
	if len(r.handlers) > 0 {
		handlers = append(handlers, r.handlers...)
	}
	return append(handlers, handler)
}

func (r *Router) ServeHTTP(req *http.Request, rw http.ResponseWriter) {
	handlers, params := r.Lookup(req.Method, req.URL.Path)
	ctx := &Context{Request: req, ResponseWriter: rw, handlers: handlers, params: params, index: -1}
	if len(handlers) == 0 {
		ctx.handlers = r.notFound
	}
	ctx.Next()
}

func (r *Router) Lookup(method, path string) ([]handlerFunc, []string) {
	if method == `HEAD` {
		method = `GET`
	}
	rootNode := r.routes[method]
	if rootNode == nil {
		return nil, nil
	}
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	_, handlers, params := rootNode.lookup(path)
	return handlers, params
}

func cleanPath(path string) string {
	// 所有路径，无论静态还是动态，都必须以"/"开头
	if len(path) == 0 || path[0] != '/' {
		panic(`router path must begin with "/": ` + path)
	}
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	re, err := syntax.Parse(path, syntax.Perl)
	if err != nil {
		panic(err)
	}
	return re.String()
}

func (r *Router) Get(path string, handler handlerFunc) *Router {
	return r.Add("GET", path, handler)
}

func (r *Router) Post(path string, handler handlerFunc) *Router {
	return r.Add("POST", path, handler)
}

func (r *Router) Put(path string, handler handlerFunc) *Router {
	return r.Add("PUT", path, handler)
}

func (r *Router) Patch(path string, handler handlerFunc) *Router {
	return r.Add("PATCH", path, handler)
}

func (r *Router) Delete(path string, handler handlerFunc) *Router {
	return r.Add("DELETE", path, handler)
}
