package goa

import (
	"bytes"
	"sort"

	"github.com/lovego/goa/docs"
	"github.com/lovego/regex_tree"
)

type RouterGroup struct {
	basePath string
	handlers []interface{}
	routes   map[string]*regex_tree.Node
	docGroup docs.Group
}

func (g *RouterGroup) Lookup(method, path string) (HandlerFuncs, []string) {
	if method == `HEAD` {
		method = `GET`
	}
	tree := g.routes[method]
	if tree == nil {
		return nil, nil
	}
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	handlers, params := tree.Lookup(path)
	if handlers != nil {
		return handlers.(HandlerFuncs), params
	}
	return nil, nil
}

func (g *RouterGroup) Get(path string, handler interface{}, args ...interface{}) *RouterGroup {
	return g.Add("GET", path, handler, args...)
}

func (g *RouterGroup) Post(path string, handler interface{}, args ...interface{}) *RouterGroup {
	return g.Add("POST", path, handler, args...)
}

func (g *RouterGroup) GetPost(path string, handler interface{}, args ...interface{}) *RouterGroup {
	g.Add("GET", path, handler, args...)
	return g.Add("POST", path, handler, args...)
}

func (g *RouterGroup) Put(path string, handler interface{}, args ...interface{}) *RouterGroup {
	return g.Add("PUT", path, handler, args...)
}

func (g *RouterGroup) Patch(path string, handler interface{}, args ...interface{}) *RouterGroup {
	return g.Add("PATCH", path, handler, args...)
}

func (g *RouterGroup) Delete(path string, handler interface{}, args ...interface{}) *RouterGroup {
	return g.Add("DELETE", path, handler, args...)
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
		buf.WriteString("  handlers: " + handlersStringIndent(g.handlers, "  ") + "\n")
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

func handlersStringIndent(handlers []interface{}, indent string) string {
	if len(handlers) == 0 {
		return "[ ]"
	}
	var buf bytes.Buffer
	buf.WriteString("[\n")
	for _, h := range handlers {
		buf.WriteString(indent + "  " + funcName(h) + "\n")
	}
	buf.WriteString(indent + "]")
	return buf.String()
}

type HandlerFuncs []func(*Context)

// to print regex_tree.Node in unit tests.
func (handlers HandlerFuncs) String() string {
	return handlers.StringIndent("")
}

func (handlers HandlerFuncs) StringIndent(indent string) string {
	if len(handlers) == 0 {
		return "[ ]"
	}
	var buf bytes.Buffer
	buf.WriteString("[\n")
	for _, h := range handlers {
		buf.WriteString(indent + "  " + funcName(h) + "\n")
	}
	buf.WriteString(indent + "]")
	return buf.String()
}
