package router

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
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
	re := regexp.MustCompile(path)
	if _, allStatic := re.LiteralPrefix(); allStatic {
		n.static = path
	} else {
		n.dynamic = re
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

// 分裂为父节点和子节点
func (n *node) split(path string) {
	var childPath string
	if len(n.static) > 0 {
		childPath = n.static[len(path):]
	} else if n.dynamic != nil {
		childPath = n.dynamic.String()[len(path):]
	} else {
		panic("both static and dynamic are empty.") // should not happen
	}
	child := newNode(childPath, n.handlers)
	child.children = n.children

	re := regexp.MustCompile(path)
	if _, allStatic := re.LiteralPrefix(); allStatic {
		n.static = re.String()
		n.dynamic = nil
	} else {
		n.static = ""
		n.dynamic = re
	}
	n.handlers = nil
	n.children = []*node{child}
}

func (n *node) lookup(path string) handleFunc {
	return nil
}

func (n *node) String() string {
	return n.string("")
}

func (n *node) string(indent string) string {
	var fields []string
	if n.static != "" {
		fields = append(fields, "static: "+n.static)
	}
	if n.dynamic != nil {
		fields = append(fields, "dynamic: "+n.dynamic.String())
	}
	if n.handlers != nil {
		fields = append(fields, fmt.Sprintf("handlers: %v", n.handlers))
	}
	if len(n.children) > 0 {
		var children bytes.Buffer
		for _, child := range n.children {
			children.WriteString(child.string(indent+"  ") + "\n")
		}
		fields = append(fields, fmt.Sprintf("children: [\n%s%s]", children.String(), indent))
	}

	return indent + "{ " + strings.Join(fields, ", ") + " }"
}
