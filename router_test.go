package goa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func ExampleRouter() {
	router := New()

	router.Get("/", func(c *Context) {
		fmt.Println("root")
	})
	users := router.Group("/users")

	users.Get("/", func(c *Context) {
		fmt.Println("list users")
	})
	users.GetX(`/(\d+)`, func(c *Context) {
		fmt.Printf("show user: %s\n", c.Param(0))
	})

	users.Post(`/`, func(c *Context) {
		fmt.Println("create a user")
	})
	users.PostX(`/postx`, func(c *Context) {
	})

	users = users.GroupX(`/(\d+)`)

	users.Put(`/`, func(c *Context) {
		fmt.Printf("fully update user: %s\n", c.Param(0))
	})
	users.PutX(`/putx`, func(c *Context) {
	})

	users.Patch(`/`, func(c *Context) {
		fmt.Printf("partially update user: %s\n", c.Param(0))
	})
	users.PatchX(`/patchx`, func(c *Context) {
	})

	users.Delete(`/`, func(c *Context) {
		fmt.Printf("delete user: %s\n", c.Param(0))
	})
	users.DeleteX(`/deletex`, func(c *Context) {
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
func ExampleRouter_GetPost() {
	router := New()
	router.GetPost("/GetPost", func(c *Context) {
		fmt.Println("GetPost")
	})
	router.GetPostX("/GetPostX", func(c *Context) {
		fmt.Println("GetPostX")
	})

	request, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}
	for _, route := range [][2]string{
		{"GET", "/GetPost"},
		{"POST", "/GetPost"},
		{"GET", "/GetPostX"},
		{"POST", "/GetPostX"},
	} {
		request.Method = route[0]
		request.URL.Path = route[1]
		router.ServeHTTP(nil, request)
	}

	// Output:
	// GetPost
	// GetPost
	// GetPostX
	// GetPostX
}

func ExampleRouter_Use() {
	router := New()
	router.Use(func(c *Context) {
		fmt.Println("middleware 1 pre")
		c.Next()
		fmt.Println("middleware 1 post")
	})
	router.Use(func(c *Context) {
		fmt.Println("middleware 2 pre")
		c.Next()
		fmt.Println("middleware 2 post")
	})
	router.Get("/", func(c *Context) {
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
	router.Use(func(c *Context) {
		fmt.Println("middleware")
		c.Next()
	})
	router.Get("/", func(c *Context) {
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

	router.NotFound(func(c *Context) {
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

func ExampleRouter_String() {
	router := New()
	router.Use(func(c *Context) {
		c.Next()
	})

	router.Get("/", func(c *Context) {
		fmt.Println("root")
	})
	users := router.Group("/users")

	users.Get("/", func(c *Context) {
		fmt.Println("list users")
	})
	users.GetX(`/(\d+)`, func(c *Context) {
		fmt.Printf("show user: %s\n", c.Param(0))
	})

	users.Post(`/`, func(c *Context) {
		fmt.Println("create a user")
	})
	users.PostX(`/postx`, func(c *Context) {
	})

	fmt.Println(router)

	// Output:
	// {
	//   handlers: [
	//     github.com/lovego/goa.ExampleRouter_String.func1
	//   ]
	//   routes: {
	//     GET:
	//     { static: /, data: [
	//       github.com/lovego/goa.ExampleRouter_String.func1
	//       github.com/lovego/goa.ExampleRouter_String.func2
	//     ], children: [
	//       { static: users, data: [
	//         github.com/lovego/goa.ExampleRouter_String.func1
	//         github.com/lovego/goa.ExampleRouter_String.func3
	//       ], children: [
	//         { dynamic: ^/([0-9]+), data: [
	//           github.com/lovego/goa.ExampleRouter_String.func1
	//           github.com/lovego/goa.ExampleRouter_String.func4
	//         ] }
	//       ] }
	//     ] }
	//     POST:
	//     { static: /users, data: [
	//       github.com/lovego/goa.ExampleRouter_String.func1
	//       github.com/lovego/goa.ExampleRouter_String.func5
	//     ], children: [
	//       { static: /postx, data: [
	//         github.com/lovego/goa.ExampleRouter_String.func1
	//         github.com/lovego/goa.ExampleRouter_String.func6
	//       ] }
	//     ] }
	//   }
	//   notFound: github.com/lovego/goa.defaultNotFound
	// }
}

func ExampleHandlerFuncs_String() {
	var hs HandlerFuncs
	fmt.Println(hs)
	hs = HandlerFuncs{
		func(*Context) {},
	}
	fmt.Println(hs)
	// Output:
	// [ ]
	// [
	//   github.com/lovego/goa.ExampleHandlerFuncs_String.func1
	// ]
}
