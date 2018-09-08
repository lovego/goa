package goa

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func h0(*Context) {}
func h1(*Context) {}
func h2(*Context) {}
func h3(*Context) {}
func h4(*Context) {}
func h5(*Context) {}

func Example_newNode_static() {
	fmt.Println(newNode("/", true, []handlerFunc{h0}))
	fmt.Println(newNode("/users", true, []handlerFunc{h1}))
	// Output:
	// { static: /, handlers: [github.com/lovego/goa.h0] }
	// { static: /users, handlers: [github.com/lovego/goa.h1] }
}

func Example_newNode_dynamic() {
	fmt.Println(newNode("/[a-z]+", false, []handlerFunc{h0}))
	fmt.Println(newNode("/users/[0-9]+", false, []handlerFunc{h1}))
	fmt.Println(newNode(`/users/\d+`, false, []handlerFunc{h2})) // should not use like this.
	// Output:
	// { dynamic: ^/[a-z]+, handlers: [github.com/lovego/goa.h0] }
	// { dynamic: ^/users/[0-9]+, handlers: [github.com/lovego/goa.h1] }
	// { dynamic: ^/users/\d+, handlers: [github.com/lovego/goa.h2] }
}

func Example_node_addToChildren_static1() {
	n := newNode("/", true, []handlerFunc{h0})
	n.addToChildren("users", true, []handlerFunc{h1})
	fmt.Println(n)
	// Output:
	// { static: /, handlers: [github.com/lovego/goa.h0], children: [
	//   { static: users, handlers: [github.com/lovego/goa.h1] }
	// ] }
}

func Example_node_addToChildren_static2() {
	n := newNode("/u", true, []handlerFunc{h0})
	n.children = []*node{
		{dynamic: regexp.MustCompile("^/")},
	}
	n.addToChildren("sers", true, []handlerFunc{h1})
	fmt.Println(n)
	// Output:
	// { static: /u, handlers: [github.com/lovego/goa.h0], children: [
	//   { static: sers, handlers: [github.com/lovego/goa.h1] }
	//   { dynamic: ^/ }
	// ] }
}

func Example_node_addToChildren_static3() {
	n := newNode("/u", true, []handlerFunc{h0})
	n.children = []*node{
		{static: "nix"},
		{dynamic: regexp.MustCompile("^/1")},
		{dynamic: regexp.MustCompile("^/2")},
	}
	n.addToChildren("sers", true, []handlerFunc{h1})
	fmt.Println(n)
	// Output:
	// { static: /u, handlers: [github.com/lovego/goa.h0], children: [
	//   { static: nix }
	//   { static: sers, handlers: [github.com/lovego/goa.h1] }
	//   { dynamic: ^/1 }
	//   { dynamic: ^/2 }
	// ] }
}

func Example_node_addToChildren_dynamic1() {
	n := newNode("/u", true, []handlerFunc{h0})
	n.children = []*node{
		{static: "sers"},
		{dynamic: regexp.MustCompile("^/")},
	}
	n.addToChildren("[0-9]+", false, []handlerFunc{h1})
	fmt.Println(n)
	// Output:
	// { static: /u, handlers: [github.com/lovego/goa.h0], children: [
	//   { static: sers }
	//   { dynamic: ^/ }
	//   { dynamic: ^[0-9]+, handlers: [github.com/lovego/goa.h1] }
	// ] }
}

func Example_node_split_static1() {
	n := newNode("/users", true, []handlerFunc{h0})
	n.split("/")
	fmt.Println(n)
	// Output:
	// { static: /, children: [
	//   { static: users, handlers: [github.com/lovego/goa.h0] }
	// ] }
}

func Example_node_split_static2() {
	n := newNode("/users/managers", false, []handlerFunc{h0})
	n.split("/users/")
	fmt.Println(n)
	// Output:
	// { static: /users/, children: [
	//   { static: managers, handlers: [github.com/lovego/goa.h0] }
	// ] }
}

func Example_node_split_dynamic1() {
	n := newNode("/[a-z]+", false, []handlerFunc{h0})
	n.split("/")
	fmt.Println(n)
	// Output:
	// { static: /, children: [
	//   { dynamic: ^[a-z]+, handlers: [github.com/lovego/goa.h0] }
	// ] }
}

func Example_node_split_dynamic2() {
	n := newNode(`/users/[0-9]+`, false, []handlerFunc{h0})
	n.split("/u")
	fmt.Println(n)
	// Output:
	// { static: /u, children: [
	//   { dynamic: ^sers/[0-9]+, handlers: [github.com/lovego/goa.h0] }
	// ] }
}

func Example_node_split_dynamic3() {
	n := newNode(`/([a-z]+)/([0-9]+)`, false, []handlerFunc{h0})
	n.split("/([a-z]+)/")
	fmt.Println(n)
	// Output:
	// { dynamic: ^/([a-z]+)/, children: [
	//   { dynamic: ^([0-9]+), handlers: [github.com/lovego/goa.h0] }
	// ] }
}

func Example_node_split_dynamic4() {
	n := newNode("/users/[0-9]+", false, []handlerFunc{h0})
	n.split("/users/")
	fmt.Println(n)
	// Output:
	// { static: /users/, children: [
	//   { dynamic: ^[0-9]+, handlers: [github.com/lovego/goa.h0] }
	// ] }
}

func Example_node_add_1() {
	root := newNode("/", false, []handlerFunc{h0})
	root.add("/users", false, []handlerFunc{h1})
	root.add(`/users/(\d+)`, false, []handlerFunc{h2})
	fmt.Println(root)
	// Output:
	// { static: /, handlers: [github.com/lovego/goa.h0], children: [
	//   { static: users, handlers: [github.com/lovego/goa.h1], children: [
	//     { dynamic: ^/(\d+), handlers: [github.com/lovego/goa.h2] }
	//   ] }
	// ] }
}

func Example_node_add_2() {
	root := newNode("/", true, []handlerFunc{h0})
	root.add("/users", true, []handlerFunc{h1})
	root.add("/users/([0-9]+)", false, []handlerFunc{h2})
	root.add("/unix/([a-z]+)", false, []handlerFunc{h3})
	root.add("/users/root", false, []handlerFunc{h4})
	root.add("/([0-9]+)", false, []handlerFunc{h5})
	fmt.Println(root)
	// Output:
	// { static: /, handlers: [github.com/lovego/goa.h0], children: [
	//   { static: u, children: [
	//     { static: sers, handlers: [github.com/lovego/goa.h1], children: [
	//       { static: /, children: [
	//         { static: root, handlers: [github.com/lovego/goa.h4] }
	//         { dynamic: ^([0-9]+), handlers: [github.com/lovego/goa.h2] }
	//       ] }
	//     ] }
	//     { dynamic: ^nix/([a-z]+), handlers: [github.com/lovego/goa.h3] }
	//   ] }
	//   { dynamic: ^([0-9]+), handlers: [github.com/lovego/goa.h5] }
	// ] }
}

func Example_node_add_conflict1() {
	root := newNode("/", true, []handlerFunc{h0})
	fmt.Println(root.add("/", true, []handlerFunc{h1}))
	// Output: 2
}

func Example_node_add_conflict2() {
	root := newNode("/", true, []handlerFunc{h0})
	root = newNode("/users", true, []handlerFunc{h0})
	fmt.Println(root.add("/users", true, []handlerFunc{h1}))
	// Output: 2
}

func Example_node_add_conflict3() {
	root := newNode("/users", false, []handlerFunc{h0})
	root.add("/", false, []handlerFunc{h1})
	fmt.Println(root.add("/users", false, []handlerFunc{h2}))
	// Output: 2
}

func Example_node_add_conflict4() {
	root := newNode("/users/active", false, []handlerFunc{h0})
	root.add("/", false, []handlerFunc{h1})
	root.add("/users", false, []handlerFunc{h2})
	fmt.Println(root.add("/users/active", false, []handlerFunc{h3}))
	// Output: 2
}

func Example_node_add_conflict5() {
	root := newNode("/users/([0-9]+)", false, []handlerFunc{h0})
	root.add("/", false, []handlerFunc{h1})
	root.add("/users", false, []handlerFunc{h2})
	fmt.Println(root.add("/users/([0-9]+)", false, []handlerFunc{h3}))
	// Output: 2
}

func Example_node_lookup_1() {
	root := newNode("/", false, []handlerFunc{h0})
	root.add("/users", false, []handlerFunc{h1})
	root.add("/users/([0-9]+)", false, []handlerFunc{h2})
	root.add("/unix/([a-z]+)/([0-9.]+)", false, []handlerFunc{h3})
	root.add("/users/root", false, []handlerFunc{h4})
	root.add("/([0-9]+)", false, []handlerFunc{h5})
	fmt.Println(root.lookup("/"))
	fmt.Println(root.lookup("/users"))
	fmt.Println(root.lookup("/users/123"))
	fmt.Println(root.lookup("/unix/linux/4.4.0"))
	fmt.Println(root.lookup("/users/root"))
	fmt.Println(root.lookup("/987"))
	fmt.Println(root.lookup("404"))
	fmt.Println(root.lookup("/users404"))
	fmt.Println(root.lookup("/unix/linux/4.4.0a"))

	// Output:
	// true [github.com/lovego/goa.h0] []
	// true [github.com/lovego/goa.h1] []
	// true [github.com/lovego/goa.h2] [123]
	// true [github.com/lovego/goa.h3] [linux 4.4.0]
	// true [github.com/lovego/goa.h4] []
	// true [github.com/lovego/goa.h5] [987]
	// false [] []
	// true [] []
	// true [] []
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

func BenchmarkRegexpFindSubmatch(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i <= b.N; i++ {
		if len(testRegexp.FindStringSubmatch("/company-skus/search")) == 0 {
			b.Error("not matched")
		}
	}
}
