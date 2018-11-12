package goa

import (
	"bytes"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/lovego/regex_tree"
)

type RouterGroup struct {
	basePath string
	handlers HandlerFuncs
	routes   map[string]*regex_tree.Node
}

func (g *RouterGroup) Group(path string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		basePath: g.concatPath(regexp.QuoteMeta(path)),
		handlers: g.concatHandlers(handlers...),
		routes:   g.routes,
	}
}

func (g *RouterGroup) GroupX(path string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		basePath: g.concatPath(path),
		handlers: g.concatHandlers(handlers...),
		routes:   g.routes,
	}
}

func (g *RouterGroup) Use(handlers ...HandlerFunc) {
	g.handlers = append(g.handlers, handlers...)
}

func (g *RouterGroup) Get(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("GET", regexp.QuoteMeta(path), handler)
}

func (g *RouterGroup) Post(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("POST", regexp.QuoteMeta(path), handler)
}

func (g *RouterGroup) Put(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("PUT", regexp.QuoteMeta(path), handler)
}

func (g *RouterGroup) Patch(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("PATCH", regexp.QuoteMeta(path), handler)
}

func (g *RouterGroup) Delete(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("DELETE", regexp.QuoteMeta(path), handler)
}

func (g *RouterGroup) GetX(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("GET", path, handler)
}

func (g *RouterGroup) PostX(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("POST", path, handler)
}

func (g *RouterGroup) PutX(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("PUT", path, handler)
}

func (g *RouterGroup) PatchX(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("PATCH", path, handler)
}

func (g *RouterGroup) DeleteX(path string, handler HandlerFunc) *RouterGroup {
	return g.Add("DELETE", path, handler)
}

func (g *RouterGroup) Add(method, path string, handler HandlerFunc) *RouterGroup {
	method = strings.ToUpper(method)
	path = g.concatPath(path)
	// remove trailing slash
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	if handler == nil {
		return g
	}
	handlers := g.concatHandlers(handler)

	rootNode := g.routes[method]
	if rootNode == nil {
		rootNode, err := regex_tree.New(path, handlers)
		if err != nil {
			log.Panic(err)
		}
		g.routes[method] = rootNode
	} else if err := rootNode.Add(path, handlers); err != nil {
		log.Panic(err)
	}
	return g
}

func (g RouterGroup) concatPath(path string) string {
	path = g.basePath + path
	if len(path) == 0 {
		log.Panic(`router path should not be empty.`)
	}
	if path[0] != '/' {
		log.Panic(`router path should begin with "/".`)
	}
	return path
}

func (g RouterGroup) concatHandlers(handlers ...HandlerFunc) HandlerFuncs {
	result := make(HandlerFuncs, len(g.handlers)+len(handlers))
	copy(result, g.handlers)
	copy(result[len(g.handlers):], handlers)
	return result
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
