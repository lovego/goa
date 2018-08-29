package router

import (
	"regexp"
)

type handleFunc func(*Context)

type node struct {
	static   string
	dynamic  *regexp.Regexp
	handlers []handleFunc
	children []*node
}

// 新建根节点
func newRootNode(path string, handlers []handleFunc) *node {
	prefix, complete := regexp.MustCompile(path).LiteralPrefix()
	var n = &node{static: prefix}
	if complete {
		n.handlers = handlers
	} else {
		n.children = []*node{{
			dynamic:  regexp.MustCompile(path[len(prefix):]),
			handlers: handlers,
		}}
	}
	return n
}

// 添加到节点
func (n *node) add(path string, handlers []handleFunc) bool {
	commonPrefix := n.commonPrefix(path)
	if len(commonPrefix) == 0 {
		return false
	}
	// 公共前缀就是完整的路径
	if len(commonPrefix) == len(path) {
		if n.handlers == nil {
			n.handlers = handlers
			return true
		} else {
			panic(`router path conflicts: ` + path)
		}
	}
	path = path[len(commonPrefix):]
	for _, child := range n.children {
		if child.add(path, handlers) {
			return true
		}
	}
	n.children = append(n.children, &node{})
	return true
}

func (n *node) commonPrefix(path string) (common string) {
	if len(n.static) > 0 {
		if static, _ := regexp.MustCompile(path).LiteralPrefix(); len(static) > 0 {
			common = longestCommonPrefix(n.static, static)
			// 公共前缀比当前节点静态路径短，则分裂
			if len(common) > 0 && len(common) < len(n.static) {
				n.split(common)
			}
		}
	} else if n.dynamic != nil {
		if static, _ := regexp.MustCompile(path).LiteralPrefix(); len(static) == 0 {
			if regexpNonLiteralPrefix(path) == n.dynamic.String() {
				common = n.dynamic.String()
			}
		}
	} else {
		panic("both static and dynamic are empty.") // should not happen
	}
	return
}

// 根据静态路径分裂为父节点和子节点
func (n *node) split(static string) {
	child := node{static: n.static[len(static):]}
	n.static = static
	if n.handlers != nil {
		child.handlers = n.handlers
		n.handlers = nil
	}
	if n.children != nil {
		child.children = n.children
	}
	n.children = []*node{&child}
}

func (n *node) lookup(path string) handleFunc {
	return nil
}
