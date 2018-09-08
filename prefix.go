package goa

import (
	"regexp"
	"regexp/syntax"
	"strings"
)

func (n *node) lookupCommonPrefix(path string) (string, []string) {
	if len(n.static) > 0 {
		if strings.HasPrefix(path, n.static) {
			return n.static, nil
		}
	} else if captures := n.dynamic.FindStringSubmatch(path); len(captures) > 0 {
		return captures[0], captures[1:]
	}
	return "", nil
}

func (n *node) commonPrefix(path string, static bool) string {
	var pathPrefix string
	if !static {
		pathPrefix, static = regexp.MustCompile(path).LiteralPrefix()
	}
	if len(n.static) > 0 {
		if static {
			return stringCommonPrefix(n.static, path)
		} else if len(pathPrefix) > 0 {
			return stringCommonPrefix(n.static, pathPrefix)
		}
	} else {
		if static {
			if prefix, _ := regexp.MustCompile(n.dynamic.String()[1:]).LiteralPrefix(); len(prefix) > 0 {
				return stringCommonPrefix(prefix, path)
			}
		} else {
			return regexpCommonPrefix(n.dynamic.String()[1:], path)
		}
	}
	return ""
}

// 字符串公共前缀
func stringCommonPrefix(a, b string) string {
	if len(a) > len(b) {
		a, b = b, a
	}
	for i := range a {
		if a[i] != b[i] {
			return a[:i]
		}
	}
	return a
}

// 正则表达式公共前缀，a和b任意一个都不能为纯字符串
func regexpCommonPrefix(aStr, bStr string) string {
	a, err := syntax.Parse(aStr, syntax.Perl)
	if err != nil {
		panic(err)
	}

	b, err := syntax.Parse(bStr, syntax.Perl)
	if err != nil {
		panic(err)
	}

	if a.Equal(b) || b.Op == syntax.OpConcat && len(b.Sub) > 0 && b.Sub[0].Equal(a) {
		return a.String()
	}
	if a.Op == syntax.OpConcat && len(a.Sub) > 0 && a.Sub[0].Equal(b) {
		return b.String()
	}

	if a.Op == syntax.OpConcat && b.Op == syntax.OpConcat {
		return concatRegexpCommonPrefix(a, b)
	}
	return ""
}

func concatRegexpCommonPrefix(a, b *syntax.Regexp) string {
	if len(a.Sub) > len(b.Sub) {
		a, b = b, a
	}

	var common string
	for i, sub := range a.Sub {
		if sub.Equal(b.Sub[i]) {
			common += sub.String()
		} else if sub.Op == syntax.OpLiteral && b.Sub[i].Op == syntax.OpLiteral {
			return common + stringCommonPrefix(sub.String(), b.Sub[i].String())
		} else {
			return common
		}
	}
	return common
}
