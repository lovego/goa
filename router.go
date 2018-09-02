package router

import (
	//  "fmt"
	"regexp/syntax"
)

type Router struct {
	middlewares []handlerFunc
	routes      map[string]*node
}

func (r *Router) Add(method, path string, handler handlerFunc) {

}

func cleanPath(path string) string {
	// 所有路径，无论静态还是动态，都必须以"/"开头
	if len(path) == 0 || path[0] != '/' {
		panic(`router path must begin with "/": ` + path)
	}
	re, err := syntax.Parse(path, syntax.Perl)
	if err != nil {
		panic(err)
	}
	return re.String()
}
