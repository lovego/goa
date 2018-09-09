package goa

import (
	"fmt"
)

func ExampleGroup_Add_conflict() {
}

func ExampleGroup_concatPath_error() {
	defer func() {
		fmt.Println(recover())
	}()
	fmt.Println(Group{}.concatPath("abc"))
	// Output: router path should begin with "/".
}

func ExampleGroup_concatPath_basic() {
	fmt.Println(Group{}.concatPath("/"))
	fmt.Println(Group{}.concatPath("/users/"))
	fmt.Println(Group{basePath: "/admin"}.concatPath(`/users/(\d+)`))
	// Output:
	// /
	// /users/
	// /admin/users/(\d+)
}
