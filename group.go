package goa

import (
	"regexp/syntax"
	"strings"
)

type Group struct {
	basePath string
	handlers []handlerFunc
	routes   map[string]*node
}

func (g *Group) Use(handlers ...handlerFunc) {
	g.handlers = append(g.handlers, handlers...)
}

func (g *Group) Add(method, path string, handler handlerFunc) {
	method = strings.ToUpper(method)
	path = cleanPath(g.basePath + path)
	if handler == nil {
		return
	}
	handlers := g.combineHandlers(handler)

	rootNode := g.routes[method]
	if rootNode == nil {
		g.routes[method] = newNode(path, handlers)
	} else if rootNode.add(path, handlers) == addResultConflict {
		panic("router conflict: " + method + " " + path)
	}
}

func (g *Group) combineHandlers(handler handlerFunc) []handlerFunc {
	size := len(g.handlers)
	if handler != nil {
		size++
	}
	result := make([]handlerFunc, size)
	copy(result, g.handlers)
	if handler != nil {
		result[size-1] = handler
	}
	return result
}

func (g *Group) Lookup(method, path string) ([]handlerFunc, []string) {
	if method == `HEAD` {
		method = `GET`
	}
	rootNode := g.routes[method]
	if rootNode == nil {
		return nil, nil
	}
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	_, handlers, params := rootNode.lookup(path)
	return handlers, params
}

func (g *Group) Get(path string, handler handlerFunc) {
	g.Add("GET", path, handler)
}

func (g *Group) Post(path string, handler handlerFunc) {
	g.Add("POST", path, handler)
}

func (g *Group) Put(path string, handler handlerFunc) {
	g.Add("PUT", path, handler)
}

func (g *Group) Patch(path string, handler handlerFunc) {
	g.Add("PATCH", path, handler)
}

func (g *Group) Delete(path string, handler handlerFunc) {
	g.Add("DELETE", path, handler)
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
