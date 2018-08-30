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

// 新建节点
func newNode(path string, handlers []handleFunc) *node {
	var n = &node{handlers: handlers}
	_, allStatic := regexp.MustCompile(path).LiteralPrefix()
	if allStatic {
		n.static = path
	} else {
		n.dynamic = regexp.MustCompile(path)
	}
	return n
}

// 添加到节点
func (n *node) add(path string, handlers []handleFunc) bool {
	commonPrefix := n.commonPrefix(path)
	if len(commonPrefix) == 0 {
		return false
	} else
	// 公共前缀比当前节点路径短，则分裂
	if len(commonPrefix) < len(n.static) ||
		n.dynamic != nil && len(commonPrefix) < len(n.dynamic.String()) {
		n.split(commonPrefix)
	}
	childPath := path[len(commonPrefix):]
	// 子节点路径为空
	if len(childPath) == 0 {
		if n.handlers == nil {
			n.handlers = handlers
			return true
		} else {
			panic(`router path conflicts: ` + path)
		}
	}
	n.addToChildren(childPath, handlers)
	return true
}

func (n *node) addToChildren(path string, handlers []handleFunc) {
	for _, child := range n.children {
		if child.add(path, handlers) {
			return
		}
	}
	child := newNode(path, handlers)
	// 静态路径优先匹配，所以将静态子节点放在动态子节点前边
	if l := len(n.children); l > 0 && len(child.static) > 0 && n.children[l-1].dynamic != nil {
		i := 0
		for ; i < l && len(n.children[i].static) > 0; i++ {
		}
		children := append(make([]*node, 0, l+1), n.children[:l]...)
		children = append(children, child)
		n.children = append(children, n.children[l:]...)
	} else {
		n.children = append(n.children, child)
	}
}

func (n *node) commonPrefix(path string) string {
	if len(n.static) > 0 {
		if static, _ := regexp.MustCompile(path).LiteralPrefix(); len(static) > 0 {
			return stringCommonPrefix(n.static, static)
		}
	} else if n.dynamic != nil {
		return regexpCommonPrefix(n.dynamic.String(), path)
	} else {
		panic("both static and dynamic are empty.") // should not happen
	}
	return ""
}

// 分裂为父节点和子节点
func (n *node) split(path string) {
	var child = &node{}
	if len(n.static) > 0 {
		child.static = n.static[len(path):]
		n.static = path
	} else if n.dynamic != nil {
		childPath := regexp.MustCompile(n.dynamic.String()[len(path):])
		if _, allStatic := childPath.LiteralPrefix(); allStatic {
			child.static = childPath.String()
		} else {
			child.dynamic = childPath
		}
		parentPath := regexp.MustCompile(path)
		if _, allStatic := parentPath.LiteralPrefix(); allStatic {
			n.static = parentPath.String()
			n.dynamic = nil
		} else {
			n.dynamic = parentPath
		}
	} else {
		panic("both static and dynamic are empty.") // should not happen
	}

	if n.handlers != nil {
		child.handlers = n.handlers
		n.handlers = nil
	}
	if n.children != nil {
		child.children = n.children
	}
	n.children = []*node{child}
}

func (n *node) lookup(path string) handleFunc {
	return nil
}
