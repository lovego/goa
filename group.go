package goa

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/lovego/regex_tree"
)

type Group struct {
	basePath string
	handlers []HandlerFunc
	routes   map[string]*regex_tree.Node
}

func (g *Group) Group(path string, handlers ...HandlerFunc) *Group {
	return &Group{
		basePath: g.concatPath(regexp.QuoteMeta(path)),
		handlers: g.concatHandlers(handlers...),
		routes:   g.routes,
	}
}

func (g *Group) GroupX(path string, handlers ...HandlerFunc) *Group {
	return &Group{
		basePath: g.concatPath(path),
		handlers: g.concatHandlers(handlers...),
		routes:   g.routes,
	}
}

func (g *Group) Use(handlers ...HandlerFunc) {
	g.handlers = append(g.handlers, handlers...)
}

func (g *Group) Get(path string, handler HandlerFunc) {
	g.Add("GET", regexp.QuoteMeta(path), handler)
}

func (g *Group) Post(path string, handler HandlerFunc) {
	g.Add("POST", regexp.QuoteMeta(path), handler)
}

func (g *Group) Put(path string, handler HandlerFunc) {
	g.Add("PUT", regexp.QuoteMeta(path), handler)
}

func (g *Group) Patch(path string, handler HandlerFunc) {
	g.Add("PATCH", regexp.QuoteMeta(path), handler)
}

func (g *Group) Delete(path string, handler HandlerFunc) {
	g.Add("DELETE", regexp.QuoteMeta(path), handler)
}

func (g *Group) GetX(path string, handler HandlerFunc) {
	g.Add("GET", path, handler)
}

func (g *Group) PostX(path string, handler HandlerFunc) {
	g.Add("POST", path, handler)
}

func (g *Group) PutX(path string, handler HandlerFunc) {
	g.Add("PUT", path, handler)
}

func (g *Group) PatchX(path string, handler HandlerFunc) {
	g.Add("PATCH", path, handler)
}

func (g *Group) DeleteX(path string, handler HandlerFunc) {
	g.Add("DELETE", path, handler)
}

func (g *Group) Add(method, path string, handler HandlerFunc) {
	method = strings.ToUpper(method)
	path = g.concatPath(path)
	// remove trailing slash
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	if handler == nil {
		return
	}
	handlers := g.concatHandlers(handler)

	rootNode := g.routes[method]
	if rootNode == nil {
		rootNode, err := regex_tree.New(path, handlers)
		if err != nil {
			panic(err)
		}
		g.routes[method] = rootNode
	} else if err := rootNode.Add(path, handlers); err != nil {
		panic(err)
	}
}

func (g Group) concatPath(path string) string {
	path = g.basePath + path
	if len(path) == 0 {
		panic(`router path should not be empty.`)
	}
	if path[0] != '/' {
		panic(`router path should begin with "/".`)
	}
	return path
}

func (g Group) concatHandlers(handlers ...HandlerFunc) []HandlerFunc {
	result := make([]HandlerFunc, len(g.handlers)+len(handlers))
	copy(result, g.handlers)
	copy(result[len(g.handlers):], handlers)
	return result
}

func (g *Group) Lookup(method, path string) ([]HandlerFunc, []string) {
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
	handlers, params := rootNode.Lookup(path)
	if handlers != nil {
		return handlers.([]HandlerFunc), params
	}
	return nil, nil
}

func (g *Group) String() string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	if g.basePath != "" {
		buf.WriteString("  basePath: " + g.basePath + "\n")
	}
	if len(g.handlers) > 0 {
		buf.WriteString("  handlers: " + fmt.Sprint(g.handlers) + "\n")
	}
	if len(g.routes) > 0 {
		buf.WriteString("  routes: {\n")
		for method, routes := range g.routes {
			buf.WriteString("    " + method + ":\n" + routes.StringIndent("    ") + "\n")
		}
		buf.WriteString("  }\n")
	}
	buf.WriteString("}\n")

	return buf.String()
}
