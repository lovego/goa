package router

import "fmt"

type Router struct {
	middlewares []handleFunc
	routes      map[string]*node
}

func (r *Router) Add(method, path string, handler handleFunc) {
	// 所有路径，无论静态还是动态，都必须以"/"开头
	if len(path) == 0 || path[0] != '/' {
		panic(`router path must begin with "/": ` + path)
	}
}
