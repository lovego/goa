package router

import (
	"regexp"
)

type handleFunc func()

type node struct {
	static   string
	dynamic  *regexp.Regexp
	handler  handleFunc
	children []*node
}

func newNode(path string, handler handleFunc) *node {
	re := regexp.MustCompile(path)
	prefix, complete := re.LiteralPrefix()
	var n = &node{}
	if len(prefix) > 0 {
		n.static = prefix
	}
	if complete {
		n.handler = handler
	} else {
		n.children = []*node{{dynamic: path, handler: handler}}
	}
	return
}

/*
最长公共前缀
*/

func (n *node) add(path string, handler handleFunc) {
	re := regexp.MustCompile(path)
	prefix, complete := re.LiteralPrefix()
	if len(n.static) == 0 && n.dynamic == nil {
		if len(prefix) > 0 {
			n.static = prefix
		}
		n.static = prefix
	}
	return
	// longest common prefix
	if n.static != "" {
	}
}

func (n *node) lookup(path string) handleFunc {
}
