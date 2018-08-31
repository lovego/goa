package router

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func Example_newNode_static() {
	fmt.Println(newNode("/", nil))
	fmt.Println(newNode("/users", nil))
	// Output:
	// { static: / }
	// { static: /users }
}

func Example_newNode_dynamic() {
	fmt.Println(newNode("/[a-z]+", nil))
	fmt.Println(newNode("/users/[0-9]+", nil))
	fmt.Println(newNode(`/users/\d+`, nil)) // should not use like this.
	// Output:
	// { dynamic: /[a-z]+ }
	// { dynamic: /users/[0-9]+ }
	// { dynamic: /users/\d+ }
}

func Example_node_addToChildren_static1() {
	n := newNode("/", nil)
	n.addToChildren("users", nil)
	fmt.Println(n)
	// Output:
	// { static: /, children: [
	//   { static: users }
	// ] }
}

func Example_node_addToChildren_static2() {
	n := newNode("/u", nil)
	n.children = []*node{
		{dynamic: regexp.MustCompile("/")},
	}
	n.addToChildren("sers", nil)
	fmt.Println(n)
	// Output:
	// { static: /u, children: [
	//   { static: sers }
	//   { dynamic: / }
	// ] }
}

func Example_node_addToChildren_static3() {
	n := newNode("/u", nil)
	n.children = []*node{
		{static: "nix"},
		{dynamic: regexp.MustCompile("/1")},
		{dynamic: regexp.MustCompile("/2")},
	}
	n.addToChildren("sers", nil)
	fmt.Println(n)
	// Output:
	// { static: /u, children: [
	//   { static: nix }
	//   { static: sers }
	//   { dynamic: /1 }
	//   { dynamic: /2 }
	// ] }
}

func Example_node_addToChildren_dynamic1() {
	n := newNode("/u", nil)
	n.children = []*node{
		{dynamic: regexp.MustCompile("/")},
	}
	n.addToChildren("[0-9]+", nil)
	fmt.Println(n)
	// Output:
	// { static: /u, children: [
	//   { static: sers }
	//   { dynamic: / }
	//   { dynamic: [0-9]+ }
	// ] }
}

func Example_node_split_static1() {
	n := newNode("/users", nil)
	n.split("/")
	fmt.Println(n)
	// Output:
	// { static: /, children: [
	//   { static: users }
	// ] }
}

func Example_node_split_static2() {
	n := newNode("/users/managers", nil)
	n.split("/users/")
	fmt.Println(n)
	// Output:
	// { static: /users/, children: [
	//   { static: managers }
	// ] }
}

func Example_node_split_dynamic1() {
	n := newNode("/[a-z]+", nil)
	n.split("/")
	fmt.Println(n)
	// Output:
	// { static: /, children: [
	//   { dynamic: [a-z]+ }
	// ] }
}

func Example_node_split_dynamic2() {
	n := newNode(`/users/[0-9]+`, nil)
	n.split("/u")
	fmt.Println(n)
	// Output:
	// { static: /u, children: [
	//   { dynamic: sers/[0-9]+ }
	// ] }
}

func Example_node_split_dynamic3() {
	n := newNode(`/([a-z]+)/([0-9]+)`, nil)
	n.split("/([a-z]+)/")
	fmt.Println(n)
	// Output:
	// { dynamic: /([a-z]+)/, children: [
	//   { dynamic: ([0-9]+) }
	// ] }
}

func Example_node_split_dynamic4() {
	n := newNode("/users/[0-9]+", nil)
	n.split("/users/")
	fmt.Println(n)
	// Output:
	// { static: /users/, children: [
	//   { dynamic: [0-9]+ }
	// ] }
}

func h0(*Context) {}
func h1(*Context) {}
func h2(*Context) {}

func Example_node_add_1() {
	root := newNode("/", []handleFunc{h0})
	root.add("/users", []handleFunc{h1})
	root.add("/users/([0-9]+)", []handleFunc{h2})
	fmt.Println(root)
	// Output:
	// { static: /, handlers: [ github.com/lovego/router.h0 ], children: [
	//   { static: users, handlers: [ github.com/lovego/router.h1 ], children: [
	//     { dynamic: /([0-9]+), handlers: [ github.com/lovego/router.h2 ] }
	//   ] }
	// ] }
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
		if !testRegexp.MatchString("/company-skus/search") {
			b.Error("not matched")
		}
	}
}
