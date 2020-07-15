package goa

import (
	"bytes"
	"log"
	pathPkg "path"
	"regexp"
	"sort"
	"strings"

	"github.com/lovego/goa/docs"
	"github.com/lovego/regex_tree"
)

type RouterGroup struct {
	basePath string
	handlers HandlerFuncs
	routes   map[string]*regex_tree.Node
	docGroup docs.Group
}

func (g *RouterGroup) DocDir(dir string) *RouterGroup {
	g.docGroup.SetDir(dir)
	return g
}

func (g *RouterGroup) Group(path string, descs ...string) *RouterGroup {
	newGroup := &RouterGroup{
		basePath: g.concatPath(quotePath(path)),
		handlers: g.concatHandlers(),
		routes:   g.routes,
	}
	if g.docGroup.Dir != "" {
		newGroup.docGroup = g.docGroup.Child(path, newGroup.basePath, descs)
	}
	return newGroup
}

func (g *RouterGroup) Use(handlers ...HandlerFunc) *RouterGroup {
	g.handlers = append(g.handlers, handlers...)
	return g
}

func (g *RouterGroup) Add(method, path string, handler interface{}) *RouterGroup {
	method = strings.ToUpper(method)
	fullPath := g.concatPath(quotePath(path))
	// remove trailing slash
	if len(fullPath) > 1 && fullPath[len(fullPath)-1] == '/' {
		fullPath = fullPath[:len(fullPath)-1]
	}
	if handler == nil {
		return g
	}
	handlerFunc := convertHandler(handler, fullPath)
	handlers := g.concatHandlers(handlerFunc)

	rootNode := g.routes[method]
	if rootNode == nil {
		rootNode, err := regex_tree.New(fullPath, handlers)
		if err != nil {
			log.Panic(err)
		}
		g.routes[method] = rootNode
	} else if err := rootNode.Add(fullPath, handlers); err != nil {
		log.Panic(err)
	}

	if g.docGroup.Dir != "" {
		g.docGroup.Route(method, path, fullPath, handler)
	}
	return g
}

func (g *RouterGroup) Get(path string, handler interface{}) *RouterGroup {
	return g.Add("GET", path, handler)
}

func (g *RouterGroup) Post(path string, handler interface{}) *RouterGroup {
	return g.Add("POST", path, handler)
}

func (g *RouterGroup) GetPost(path string, handler interface{}) *RouterGroup {
	g.Add("GET", path, handler)
	return g.Add("POST", path, handler)
}

func (g *RouterGroup) Put(path string, handler interface{}) *RouterGroup {
	return g.Add("PUT", path, handler)
}

func (g *RouterGroup) Patch(path string, handler interface{}) *RouterGroup {
	return g.Add("PATCH", path, handler)
}

func (g *RouterGroup) Delete(path string, handler interface{}) *RouterGroup {
	return g.Add("DELETE", path, handler)
}

func (g *RouterGroup) Lookup(method, path string) (HandlerFuncs, []string) {
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
		return handlers.(HandlerFuncs), params
	}
	return nil, nil
}

func (g *RouterGroup) String() string {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	if g.basePath != "" {
		buf.WriteString("  basePath: " + g.basePath + "\n")
	}
	buf.WriteString(g.RoutesString())
	buf.WriteString("}\n")
	return buf.String()
}

func (g *RouterGroup) RoutesString() string {
	var buf bytes.Buffer
	if len(g.handlers) > 0 {
		buf.WriteString("  handlers: " + g.handlers.StringIndent("  ") + "\n")
	}
	if len(g.routes) > 0 {
		buf.WriteString("  routes: {\n")
		methods := make([]string, 0, len(g.routes))
		for method := range g.routes {
			methods = append(methods, method)
		}
		sort.Strings(methods)
		for _, method := range methods {
			buf.WriteString("    " + method + ":\n" + g.routes[method].StringIndent("    ") + "\n")
		}
		buf.WriteString("  }\n")
	}
	return buf.String()
}

func (g RouterGroup) concatPath(path string) string {
	path = pathPkg.Join(g.basePath, path)
	if len(path) == 0 {
		log.Panic(`router path must not be empty.`)
	}
	if path[0] != '/' {
		log.Panic(`router path must begin with "/".`)
	}
	return path
}

func (g RouterGroup) concatHandlers(handlers ...HandlerFunc) HandlerFuncs {
	result := make(HandlerFuncs, len(g.handlers)+len(handlers))
	copy(result, g.handlers)
	copy(result[len(g.handlers):], handlers)
	return result
}

func quotePath(path string) string {
	// if path contains "(" or ")" it should be a regular expression
	if strings.ContainsAny(path, "()") {
		return path
	}
	return regexp.QuoteMeta(path)
}
