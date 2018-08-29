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
	var commonPrefix string
	if len(n.static) > 0 {
		if static, _ := regexp.MustCompile(path).LiteralPrefix(); len(static) > 0 {
			commonPrefix = n.addToStaticNode(static, path, handlers)
		}
	} else if n.dynamic != nil {
		if static, _ := regexp.MustCompile(path).LiteralPrefix(); len(static) == 0 {
			commonPrefix = n.addToDynamicNode(path, handlers)
		}
	} else {
		panic("both static and dynamic are empty.") // should not happen
	}

	if len(commonPrefix) == 0 {
		return false
	}
	path = path[len(commonPrefix):]
	if len(path) == 0 {
		return true
	}
	// add to children
	for i, child := range n.children {
		if child.add() {
			return true
		}
	}
	n.addToChildren(static, dynamic, handlers)
}

// 尝试加入静态节点，如果加入成功则返回共同前缀
func (n *node) addToStaticNode(static, path string, handlers []handleFunc) string {
	commonPrefix := longestCommonPrefix(n.static, static)
	if len(commonPrefix) == 0 {
		return ""
	}
	// 公共前缀比当前节点静态路径短，则分裂
	if len(commonPrefix) < len(n.static) {
		n.split(commonPrefix)
	}
	// 公共前缀就是完整的路径
	if len(commonPrefix) == len(path) {
		if n.handlers == nil {
			n.handlers = handlers
		} else {
			panic(`router path conflicts: ` + path)
		}
	}
	return commonPrefix
}

func (n *node) addToDynamicNode(path string, handlers []handleFunc) bool {
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
	n.children = []*node{child}
}

// 添加到子节点
func (n *node) addToChildren(static, dynamic string, handlers []handleFunc) {
}

// 最长公共前缀
func longestCommonPrefix(a, b string) string {
	if len(a) > len(b) {
		a, b = b, a
	}
	for i, char := range a {
		if a[i] != b[i] {
			break
		}
	}
	return a[:i]
}

func (n *node) lookup(path string) handleFunc {
}
