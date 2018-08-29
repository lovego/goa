package router

import (
	"regexp/syntax"
)

// 最长公共前缀
func longestCommonPrefix(a, b string) string {
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

// 正则非字面量前缀
func regexpNonLiteralPrefix(expr string) string {
	re, err := syntax.Parse(expr, syntax.Perl)
	if err != nil {
		panic(err)
	}
	re = re.Simplify()
	if re.Op != syntax.OpConcat {
		return expr
	}
	var prefix string
	for _, sub := range re.Sub {
		if sub.Op == syntax.OpLiteral {
			break
		}
		prefix += sub.String()
	}
	if len(prefix) > 0 {
		return prefix
	}
	return expr
}
