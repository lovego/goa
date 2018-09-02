package router

import (
	"net/http"
	"regexp/syntax"
	"strings"
)

type Router struct {
	handlers handlersChain
	base     string
	routes   map[string]*node
}

func New() *Router {
	return &Router{routes: make(map[string]*node)}
}

func (r *Router) Add(method, path string, handler handlerFunc) {
	method = strings.ToUpper(method)
	path = cleanPath(path)
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
	if len(r.base) > 1 && len(r.handlers) > 0 {
		handlers = append(handlers, r.handlers...)
	}
	return append(handlers, handler)
}

func (r *Router) ServeHTTP(req *http.Request, rw http.ResponseWriter) {
}

func (r *Router) Lookup(method, path string) handlersChain {
	return nil
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
