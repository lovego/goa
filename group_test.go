package goa

import (
	"fmt"
)

func Example_cleanPath_1() {
	defer func() {
		fmt.Println(recover())
	}()
	fmt.Println(cleanPath(""))
	// Output: router path must begin with "/":
}

func Example_cleanPath_2() {
	defer func() {
		fmt.Println(recover())
	}()
	fmt.Println(cleanPath("abc"))
	// Output: router path must begin with "/": abc
}

func Example_cleanPath_3() {
	defer func() {
		fmt.Println(recover())
	}()
	fmt.Println(cleanPath("/("))
	// Output: error parsing regexp: missing closing ): `/(`
}

func Example_cleanPath_4() {
	fmt.Println(cleanPath("/"))
	fmt.Println(cleanPath("/users/"))
	fmt.Println(cleanPath(`/users/(\d+)`))
	// Output:
	// /
	// /users
	// /users/([0-9]+)
}
