package goa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func Example_basic() {
	router := New()

	router.Get("/", func(ctx *Context) {
		fmt.Println("root")
	})
	router.Get("/users", func(ctx *Context) {
		fmt.Println("list users")
	})
	router.GetX(`/users/(\d+)`, func(ctx *Context) {
		fmt.Printf("show user: %s\n", ctx.Param(0))
	})

	router.Post(`/users`, func(ctx *Context) {
		fmt.Println("create a user")
	})
	router.PutX(`/users/(\d+)`, func(ctx *Context) {
		fmt.Printf("fully update user: %s\n", ctx.Param(0))
	})
	router.PatchX(`/users/(\d+)`, func(ctx *Context) {
		fmt.Printf("partially update user: %s\n", ctx.Param(0))
	})
	router.DeleteX(`/users/(\d+)`, func(ctx *Context) {
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
		router.ServeHTTP(nil, request)
	}
	rw := httptest.NewRecorder()
	request.URL.Path = "/404"
	router.ServeHTTP(rw, request)
	response := rw.Result()
	if body, err := ioutil.ReadAll(response.Body); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode, string(body))
	}

	// Output:
	// root
	// list users
	// create a user
	// show user: 101
	// fully update user: 101
	// partially update user: 101
	// delete user: 101
	// 404 {"code":"404","message":"Not Found."}
}

func ExampleRouter_NotFound() {
	router := New()
	router.Use(func(ctx *Context) {
		fmt.Println("middleware 1 pre")
		ctx.Next()
		fmt.Println("middleware 1 post")
	})
	router.Use(func(ctx *Context) {
		fmt.Println("middleware 2 pre")
		ctx.Next()
		fmt.Println("middleware 2 post")
	})
	router.Get("/", func(ctx *Context){
	    fmt.Println("you got it")
    })
	request, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}
	router.ServeHTTP(nil, request)
	// Output:
	// middleware 1 pre
	// middleware 2 pre
	// you got it
	// middleware 2 post
	// middleware 1 post
}
