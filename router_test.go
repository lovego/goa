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
		fmt.Printf("update user: %s\n", ctx.Param(0))
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
		{"GET", "/users/101"},
		{"PUT", "/users/101"},
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
	// update user: 101
	// delete user: 101
}
