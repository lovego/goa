package goa

import (
	"fmt"
)

func ExampleRouterGroup_Add_conflict() {
}

func ExampleRouterGroup_concatPath_error() {
	defer func() {
		fmt.Println(recover())
	}()
	fmt.Println(RouterGroup{}.concatPath("abc"))
	// Output: router path should begin with "/".
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
