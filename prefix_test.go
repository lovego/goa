package router

import (
	"fmt"
)

func ExampleStringCommonPrefix_Empty1() {
	fmt.Println(stringCommonPrefix("users", "managers"))
	// Output:
}
func ExampleStringCommonPrefix_Same1() {
	fmt.Println(stringCommonPrefix("/", "/"))
	// Output: /
}

func ExampleStringCommonPrefix_Same2() {
	fmt.Println(stringCommonPrefix("/users", "/users"))
	// Output: /users
}

func ExampleStringCommonPrefix_LeftLonger1() {
	fmt.Println(stringCommonPrefix("/users", "/"))
	// Output: /
}

func ExampleStringCommonPrefix_LeftLonger2() {
	fmt.Println(stringCommonPrefix("/users/root", "/users"))
	// Output: /users
}

func ExampleStringCommonPrefix_RightLonger1() {
	fmt.Println(stringCommonPrefix("/", "/users"))
	// Output: /
}

func ExampleStringCommonPrefix_RightLonger2() {
	fmt.Println(stringCommonPrefix("/users", "/users/root"))
	// Output: /users
}

func ExampleStringCommonPrefix_DifferentSuffix() {
	fmt.Println(stringCommonPrefix("/users/list", "/users/root"))
	// Output: /users/
}

func ExampleRegexpCommonPrefix_Empty1() {
	fmt.Println(regexpCommonPrefix(`user_(\d+)/xyz`, `manager_(\d+)/def`))
	// Output:
}

func ExampleRegexpCommonPrefix_Empty2() {
	fmt.Println(regexpCommonPrefix(`(\d+)/xyz`, `(\w+)/def`))
	// Output:
}

func ExampleRegexpCommonPrefix_Empty3() {
	fmt.Println(regexpCommonPrefix(`(\d+)`, `(\w+)`))
	// Output:
}

func ExampleRegexpCommonPrefix_Same1() {
	fmt.Println(regexpCommonPrefix(`/(\d+)`, `/(\d+)`))
	// Output: /([0-9]+)
}

func ExampleRegexpCommonPrefix_Same2() {
	fmt.Println(regexpCommonPrefix(`/user_(\d+)`, `/user_(\d+)`))
	// Output: /user_([0-9]+)
}

func ExampleRegexpCommonPrefix_LeftLonger1() {
	fmt.Println(regexpCommonPrefix(`/(\d+)`, `/`))
	// Output: /
}

func ExampleRegexpCommonPrefix_LeftLonger2() {
	fmt.Println(regexpCommonPrefix(`/users/(\d+)`, `/users/`))
	// Output: /users/
}

func ExampleRegexpCommonPrefix_LeftLonger3() {
	fmt.Println(regexpCommonPrefix(`/users/(\w+)/(\d+)`, `/users/(\w+)`))
	// Output: /users/([0-9A-Z_a-z]+)
}

func ExampleRegexpCommonPrefix_RightLonger1() {
	fmt.Println(regexpCommonPrefix(`/`, `/(\d+)`))
	// Output: /
}

func ExampleRegexpCommonPrefix_RightLonger2() {
	fmt.Println(regexpCommonPrefix(`/users/`, `/users/(\d+)`))
	// Output: /users/
}

func ExampleRegexpCommonPrefix_DifferentSuffix() {
	fmt.Println(regexpCommonPrefix(`/user_([0-9]+)/xyz`, `/user_([0-9]+)/def`))
	// Output: /user_([0-9]+)
}
