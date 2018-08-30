package router

import (
	"regexp/syntax"
)

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

// 正则表达式公共前缀
func regexpCommonPrefix(aStr, bStr string) string {
	a, err := syntax.Parse(aStr, syntax.Perl)
	if err != nil {
		panic(err)
	}
	a = a.Simplify()

	b, err := syntax.Parse(bStr, syntax.Perl)
	if err != nil {
		panic(err)
	}
	b = b.Simplify()

	if a.Equal(b) {
		return aStr
	}

	if a.Op == syntax.OpConcat && b.Op == syntax.OpConcat {
		if len(a.Sub) > len(b.Sub) {
			a, b = b, a
		}

		var common string
		for i, sub := range a.Sub {
			if !sub.Equal(b.Sub[i]) {
				return common
			}
			common += sub.String()
		}
		return common
	}
	return ""
}
