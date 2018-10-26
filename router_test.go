package goa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func ExampleRouter() {
	router := New()

	router.Get("/", func(ctx *Context) {
		fmt.Println("root")
	})
	users := router.Group("/users")

	users.Get("/", func(ctx *Context) {
		fmt.Println("list users")
	})
	users.GetX(`/(\d+)`, func(ctx *Context) {
		fmt.Printf("show user: %s\n", ctx.Param(0))
	})

	users.Post(`/`, func(ctx *Context) {
		fmt.Println("create a user")
	})
	users.PostX(`/postx`, func(ctx *Context) {
	})

	users = users.GroupX(`/(\d+)`)

	users.Put(`/`, func(ctx *Context) {
		fmt.Printf("fully update user: %s\n", ctx.Param(0))
	})
	users.PutX(`/putx`, func(ctx *Context) {
	})

	users.Patch(`/`, func(ctx *Context) {
		fmt.Printf("partially update user: %s\n", ctx.Param(0))
	})
	users.PatchX(`/patchx`, func(ctx *Context) {
	})

	users.Delete(`/`, func(ctx *Context) {
		fmt.Printf("delete user: %s\n", ctx.Param(0))
	})
	users.DeleteX(`/deletex`, func(ctx *Context) {
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

	// Output:
	// root
	// list users
	// create a user
	// show user: 101
	// fully update user: 101
	// partially update user: 101
	// delete user: 101
}

func ExampleRouter_Use() {
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
	router.Get("/", func(ctx *Context) {
		fmt.Println("root")
	})

	request, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}
	router.ServeHTTP(nil, request)
	// Output:
	// middleware 1 pre
	// middleware 2 pre
	// root
	// middleware 2 post
	// middleware 1 post
}

func ExampleRouter_NotFound() {
	router := New()
	router.Use(func(ctx *Context) {
		fmt.Println("middleware")
		ctx.Next()
	})
	router.Get("/", func(ctx *Context) {
		fmt.Println("root")
	})

	request, err := http.NewRequest("GET", "http://localhost/404", nil)
	if err != nil {
		panic(err)
	}
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, request)

	response := rw.Result()
	if body, err := ioutil.ReadAll(response.Body); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode, string(body))
	}

	router.NotFound(func(ctx *Context) {
		fmt.Println("404 not found")
	})
	router.ServeHTTP(nil, request)

	request.URL.Path = "/"
	router.ServeHTTP(nil, request)

	// Output:
	// middleware
	// 404 {"code":"404","message":"Not Found."}
	// middleware
	// 404 not found
	// middleware
	// root
}
