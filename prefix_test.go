package router

import (
	"fmt"
)

func Example_node_commonPrefix_static() {
	fmt.Println(newNode("/", nil).commonPrefix("/"))
	fmt.Println(newNode("users", nil).commonPrefix("members"))
	fmt.Println(newNode("/users", nil).commonPrefix("/"))
	fmt.Println(newNode("users", nil).commonPrefix("([0-9]+)"))
	fmt.Println(newNode("/", nil).commonPrefix("/users"))
	fmt.Println(newNode("users/managers", nil).commonPrefix("users/([0-9]+)"))
	// Output:
	// /
	//
	// /
	//
	// /
	// users/
}

func Example_node_commonPrefix_dynamic() {
	fmt.Println(newNode("/([a-z]+)", nil).commonPrefix("/"))
	fmt.Println(newNode("users/([0-9]+)", nil).commonPrefix("members"))
	fmt.Println(newNode("/([a-z]+)", nil).commonPrefix("/([a-z]+)/members"))
	fmt.Println(newNode("users/([0-9]+)", nil).commonPrefix("users/([a-z]+)"))
	fmt.Println(newNode("users/([0-9]+)", nil).commonPrefix("users/managers"))
	// Output:
	// /
	//
	// /([a-z]+)
	// users/
	// users/
}

func Example_stringCommonPrefix_empty1() {
	fmt.Println(stringCommonPrefix("users", "managers"))
	// Output:
}
func Example_stringCommonPrefix_same1() {
	fmt.Println(stringCommonPrefix("/", "/"))
	// Output: /
}

func Example_stringCommonPrefix_same2() {
	fmt.Println(stringCommonPrefix("/users", "/users"))
	// Output: /users
}

func Example_stringCommonPrefix_leftLonger1() {
	fmt.Println(stringCommonPrefix("/users", "/"))
	// Output: /
}

func Example_stringCommonPrefix_leftLonger2() {
	fmt.Println(stringCommonPrefix("/users/root", "/users"))
	// Output: /users
}

func Example_stringCommonPrefix_rightLonger1() {
	fmt.Println(stringCommonPrefix("/", "/users"))
	// Output: /
}

func Example_stringCommonPrefix_rightLonger2() {
	fmt.Println(stringCommonPrefix("/users", "/users/root"))
	// Output: /users
}

func Example_stringCommonPrefix_differentSuffix() {
	fmt.Println(stringCommonPrefix("/users/list", "/users/root"))
	// Output: /users/
}

func Example_regexpCommonPrefix_empty1() {
	fmt.Println(regexpCommonPrefix(`user_(\d+)/xyz`, `manager_(\d+)/def`))
	// Output:
}

func Example_regexpCommonPrefix_empty2() {
	fmt.Println(regexpCommonPrefix(`(\d+)/xyz`, `(\w+)/def`))
	// Output:
}

func Example_regexpCommonPrefix_empty3() {
	fmt.Println(regexpCommonPrefix(`(\d+)`, `(\w+)`))
	// Output:
}

func Example_regexpCommonPrefix_same1() {
	fmt.Println(regexpCommonPrefix(`/(\d+)`, `/(\d+)`))
	// Output: /([0-9]+)
}

func Example_regexpCommonPrefix_same2() {
	fmt.Println(regexpCommonPrefix(`/user_(\d+)`, `/user_(\d+)`))
	// Output: /user_([0-9]+)
}

func Example_regexpCommonPrefix_leftLonger1() {
	fmt.Println(regexpCommonPrefix(`/(\d+)`, `/`))
	// Output: /
}

func Example_regexpCommonPrefix_leftLonger2() {
	fmt.Println(regexpCommonPrefix(`/users/(\d+)`, `/users/`))
	// Output: /users/
}

func Example_regexpCommonPrefix_leftLonger3() {
	fmt.Println(regexpCommonPrefix(`/users/(\w+)/(\d+)`, `/users/(\w+)`))
	// Output: /users/([0-9A-Z_a-z]+)
}

func Example_regexpCommonPrefix_rightLonger1() {
	fmt.Println(regexpCommonPrefix(`/`, `/(\d+)`))
	// Output: /
}

func Example_regexpCommonPrefix_rightLonger2() {
	fmt.Println(regexpCommonPrefix(`/users/`, `/users/(\d+)`))
	// Output: /users/
}

func Example_regexpCommonPrefix_differentSuffix1() {
	fmt.Println(regexpCommonPrefix(`/user_([0-9]+)/xyz`, `/user_([0-9]+)/def`))
	// Output: /user_([0-9]+)/
}

func Example_regexpCommonPrefix_differentSuffix2() {
	// should not use like so
	fmt.Println(regexpCommonPrefix("users/([0-9]+)", "users/managers"))
	// Output:
}

func Example_regexpCommonPrefix_differentSuffix3() {
	fmt.Println(regexpCommonPrefix("/([a-z]+)/members/([0-9]+)", "/([a-z]+)/managers/([0-9]+)"))
	// Output: /([a-z]+)/m
}

func Example_regexpCommonPrefix_panic1() {
	defer func() {
		fmt.Println(recover())
	}()
	fmt.Println(regexpCommonPrefix("(", "/"))
	// Output: error parsing regexp: missing closing ): `(`
}

func Example_regexpCommonPrefix_panic2() {
	defer func() {
		fmt.Println(recover())
	}()
	fmt.Println(regexpCommonPrefix("/", ")"))
	// Output: error parsing regexp: unexpected ): `)`
}
