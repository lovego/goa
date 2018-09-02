package router

import (
	"net/http"
	"regexp/syntax"
	"strings"
)

type Router struct {
	basePath string
	handlers handlersChain
	routes   map[string]*node
	notFound handlersChain
}

func New() *Router {
	return &Router{routes: make(map[string]*node)}
}

func (r *Router) Add(method, path string, handler handlerFunc) {
	method = strings.ToUpper(method)
	path = cleanPath(r.basePath + path)
	handlers := r.getHandlers(handler)
	rootNode := r.routes[method]
	if rootNode == nil {
		r.routes[method] = newNode(path, handlers)
	} else if rootNode.add(path, handlers) == addResultConflict {
		panic("router conflict: " + method + " " + path)
	}
}

func (r *Router) getHandlers(handler handlerFunc) handlersChain {
	if handler == nil {
		panic("handler func should not be nil")
	}
	var handlers handlersChain
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

func (r *Router) Lookup(method, path string) (handlersChain, []string) {
	if method == `HEAD` {
		method = `GET`
	}
	rootNode := r.routes[method]
	if rootNode == nil {
		return nil, nil
	}
	_, handlers, params := rootNode.lookup(strings.TrimRight(path, "/"))
	return handlers, params
}

func cleanPath(path string) string {
	// 所有路径，无论静态还是动态，都必须以"/"开头
	if len(path) == 0 || path[0] != '/' {
		panic(`router path must begin with "/": ` + path)
	}
	path = strings.TrimRight(path, "/")
	re, err := syntax.Parse(path, syntax.Perl)
	if err != nil {
		panic(err)
	}
	return re.String()
}
