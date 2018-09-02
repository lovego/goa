package goa

import (
	"fmt"
	"net/http"
)

func Example_Router() {
	router := New()

	router.Get("/", func(ctx *Context) {
		fmt.Println("root")
	})
	router.Get("/users", func(ctx *Context) {
		fmt.Println("list users")
	})
	router.Get(`/users/(\d+)`, func(ctx *Context) {
		fmt.Printf("show user: %s\n", ctx.Param(0))
	})

	router.Post(`/users`, func(ctx *Context) {
		fmt.Println("create a user")
	})
	router.Put(`/users/(\d+)`, func(ctx *Context) {
		fmt.Printf("fully update user: %s\n", ctx.Param(0))
	})
	router.Patch(`/users/(\d+)`, func(ctx *Context) {
		fmt.Printf("partially update user: %s\n", ctx.Param(0))
	})
	router.Delete(`/users/(\d+)`, func(ctx *Context) {
		fmt.Printf("delete user: %s\n", ctx.Param(0))
	})

	request, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}
	for _, route := range [][2]string{
		{"GET", "/"},
		{"GET", "/users"},
		{"POST", "/users"},
		{"GET", "/users/101/"}, // with a trailing slash
		{"PUT", "/users/101"},
		{"PATCH", "/users/101"},
		{"DELETE", "/users/101"},
	} {
		request.Method = route[0]
		request.URL.Path = route[1]
		router.ServeHTTP(request, nil)
	}

	// Output:
	// root
	// list users
	// create a user
	// show user: 101
	// fully update user: 101
	// partially update user: 101
	// delete user: 101
}

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
