package router

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func ExampleNewNode_Static() {
	printNode(newNode("/", nil))
	printNode(newNode("/users", nil))
	// Output:
	// {static:/ dynamic:<nil> handlers:[] children:[]}
	// {static:/users dynamic:<nil> handlers:[] children:[]}
}

func ExampleNewNode_Dynamic() {
	printNode(newNode("/[a-z]+", nil))
	printNode(newNode("/users/[0-9]+", nil))
	printNode(newNode(`/users/\d+`, nil)) // should not use like this.
	// Output:
	// {static: dynamic:/[a-z]+ handlers:[] children:[]}
	// {static: dynamic:/users/[0-9]+ handlers:[] children:[]}
	// {static: dynamic:/users/\d+ handlers:[] children:[]}
}

func printNode(n *node) {
	fmt.Printf(
		"{static:%s dynamic:%v handlers:%v children:%v}\n",
		n.static, n.dynamic, n.handlers, n.children,
	)
}

func BenchmarkStringHasPrefix(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i <= b.N; i++ {
		if !strings.HasPrefix("/company-skus/search", "/company-skus") {
			b.Error("not matched")
		}
	}
}

var testRegexp = regexp.MustCompile("^/company-skus")

func BenchmarkRegexpMatch(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i <= b.N; i++ {
		if len(testRegexp.FindStringSubmatch("/company-skus/search")) == 0 {
			b.Error("not matched")
		}
	}
}
