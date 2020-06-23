package goa

import (
	"fmt"

	"github.com/lovego/regex_tree"
)

func ExampleRouterGroup() {
	g := &RouterGroup{basePath: "/users", routes: make(map[string]*regex_tree.Node)}
	g.Use(func(*Context) {})
	g.Get("/nil", nil)
	g.Get("/", func(*Context) {})

	fmt.Println(g)
	fmt.Println(g.Lookup("HEAD", "/users"))
	fmt.Println(g.Lookup("POST", "/users"))
	// Output:
	// {
	//   basePath: /users
	//   handlers: [
	//     github.com/lovego/goa.ExampleRouterGroup.func1
	//   ]
	//   routes: {
	//     GET:
	//     { static: /users, data: [
	//       github.com/lovego/goa.ExampleRouterGroup.func1
	//       github.com/lovego/goa.ExampleRouterGroup.func2
	//     ] }
	//   }
	// }
	//
	// [
	//   github.com/lovego/goa.ExampleRouterGroup.func1
	//   github.com/lovego/goa.ExampleRouterGroup.func2
	// ] []
	// [ ] []
}

func ExampleRouterGroup_Add_error1() {
	defer func() {
		fmt.Println(recover())
	}()
	g := &RouterGroup{routes: make(map[string]*regex_tree.Node)}
	g.Add("GET", "/(", func(*Context) {})
	// Output: error parsing regexp: missing closing ): `/(`
}

func ExampleRouterGroup_Add_error2() {
	defer func() {
		fmt.Println(recover())
	}()
	g := &RouterGroup{routes: make(map[string]*regex_tree.Node)}
	g.Add("GET", "/", func(*Context) {})
	g.Add("GET", "/", func(*Context) {})
	// Output: path already exists
}

func ExampleRouterGroup_concatPath_basic() {
	fmt.Println(RouterGroup{}.concatPath("/"))
	fmt.Println(RouterGroup{}.concatPath("/users/"))
	fmt.Println(RouterGroup{basePath: "/admin"}.concatPath(`/users/(\d+)`))
	// Output:
	// /
	// /users/
	// /admin/users/(\d+)
}

func ExampleRouterGroup_concatPath_error1() {
	defer func() {
		fmt.Println(recover())
	}()
	fmt.Println(RouterGroup{}.concatPath(""))
	// Output: router path must not be empty.
}

func ExampleRouterGroup_concatPath_error2() {
	defer func() {
		fmt.Println(recover())
	}()
	fmt.Println(RouterGroup{}.concatPath("abc"))
	// Output: router path must begin with "/".
}
