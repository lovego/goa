package goa

import (
	"log"
	pathPkg "path"
	"regexp"
	"strings"

	"github.com/lovego/regex_tree"
)

func (g *RouterGroup) DocDir(dir string) *RouterGroup {
	g.docGroup.SetDir(dir)
	return g
}

func (g *RouterGroup) Group(path string, descs ...string) *RouterGroup {
	// make a copy to prevent handlers of children groups overwritten.
	handlersCopy := make([]interface{}, len(g.handlers))
	copy(handlersCopy, g.handlers)
	newGroup := &RouterGroup{
		basePath: g.concatPath(path),
		handlers: handlersCopy,
		routes:   g.routes,
	}
	if g.docGroup.Dir != "" {
		newGroup.docGroup = g.docGroup.Child(path, newGroup.basePath, descs)
	}
	return newGroup
}

// Use adds middlewares to the group, which will be executed for all routes in this group.
func (g *RouterGroup) Use(handlers ...func(*Context)) *RouterGroup {
	for _, handler := range handlers {
		g.handlers = append(g.handlers, handler)
	}
	return g
}

// Watch watchs every route in the group,
// and optionally return a middleware to be executed only for this route.
func (g *RouterGroup) Watch(
	watchers ...func(method, fullPath string, args []interface{}) func(*Context),
) *RouterGroup {
	for _, watcher := range watchers {
		g.handlers = append(g.handlers, watcher)
	}
	return g
}

func (g *RouterGroup) Add(method, path string, handler interface{}, args ...interface{}) *RouterGroup {
	if handler == nil {
		return g
	}
	method = strings.ToUpper(method)
	fullPath := g.makeFullPath(path)
	handlers := g.makeHandlers(method, fullPath, args, handler)

	if tree := g.routes[method]; tree == nil {
		if tree, err := regex_tree.New(fullPath, handlers); err != nil {
			log.Panic(err)
		} else {
			g.routes[method] = tree
		}
	} else if err := tree.Add(fullPath, handlers); err != nil {
		log.Panic(err)
	}

	if g.docGroup.Dir != "" {
		g.docGroup.Route(method, path, fullPath, handler)
	}
	return g
}

func (g RouterGroup) makeHandlers(
	method, fullPath string, args []interface{}, routeHandler interface{},
) HandlerFuncs {
	handlers := make([]func(*Context), 0, len(g.handlers)+1)
	for _, v := range g.handlers {
		switch h := v.(type) {
		case func(*Context):
			handlers = append(handlers, h)
		case func(method, fullPath string, args []interface{}) func(*Context):
			if handler := h(method, fullPath, args); handler != nil {
				handlers = append(handlers, handler)
			}
		default:
			log.Panicf("Unknown handler: %v\n", h)
		}
	}
	return append(handlers, convertHandler(routeHandler, fullPath))
}

func (g RouterGroup) makeFullPath(path string) string {
	full := g.concatPath(path)
	// remove trailing slash
	if len(full) > 1 && full[len(full)-1] == '/' {
		full = full[:len(full)-1]
	}
	return full
}

func (g RouterGroup) concatPath(path string) string {
	// if path contains "(" or ")" it should be a regular expression
	if !strings.ContainsAny(path, "()") {
		path = regexp.QuoteMeta(path)
	}

	path = pathPkg.Join(g.basePath, path)
	if len(path) == 0 {
		log.Panic(`router path must not be empty.`)
	}
	if path[0] != '/' {
		log.Panic(`router path must begin with "/".`)
	}
	return path
}
