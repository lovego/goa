package goa

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleConvertReq() {
	router := New()
	type T struct {
		Type string
		Id   int
		Flag bool
	}

	router.Get(`/(?P<type>\w+)/(?P<id>\d+)/(?P<flag>(true|false))`, func(req struct {
		Param T
		Query struct {
			Page int
			T
		}
		Header struct {
			Cookie string
		}
		Body struct {
			Name string
			T
		}
		Ctx *Context
	}, resp *struct {
		Error  error
		Data   interface{}
		Header struct {
			SetCookie string `header:"Set-Cookie"`
		}
	}) {
		fmt.Printf("req.Param: %+v\n", req.Param)
		fmt.Printf("req.Query: %+v\n", req.Query)
		fmt.Printf("req.Header: %+v\n", req.Header)
		fmt.Printf("req.Body: %+v\n", req.Body)
		fmt.Printf("req.Ctx not nil: %v\n", req.Ctx != nil)

		resp.Data = []int{1, 2, 3}
		resp.Header.SetCookie = "c=d"
	})

	req, err := http.NewRequest(
		"GET",
		"http://localhost/users/123/true?page=3&type=users&Id=123&flag=true",
		strings.NewReader(`{"name":"张三", "type":"users", "id": 123, "flag": true}`),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cookie", "a=b")

	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	resp := rw.Result()
	fmt.Printf("resp.Status: %v\n", resp.Status)
	delete(resp.Header, "Content-Type")
	fmt.Printf("resp.Header: %v\n", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp.Body: %v %v\n", string(body), err)

	// Output:
	// req.Param: {Type:users Id:123 Flag:true}
	// req.Query: {Page:3 T:{Type:users Id:123 Flag:true}}
	// req.Header: {Cookie:a=b}
	// req.Body: {Name:张三 T:{Type:users Id:123 Flag:true}}
	// req.Ctx not nil: true
	// resp.Status: 200 OK
	// resp.Header: map[Set-Cookie:[c=d]]
	// resp.Body: {"code":"ok","message":"success","data":[1,2,3]}
	//  <nil>
}

func ExampleConvertReq_param() {
	router := New()

	router.Get(`/(\d+)`, func(req struct {
		Param int
	}, resp *struct {
	}) {
		fmt.Println(req.Param)
	})
	req, err := http.NewRequest("GET", "http://localhost/123", nil)
	if err != nil {
		panic(err)
	}
	router.ServeHTTP(nil, req)
	// Output:
	// 123
}

func ExampleConvertReq_pointerFields() {
	router := New()
	type T struct {
		Type string
		Id   int
		Flag bool
	}

	router.Get(`/(?P<type>\w+)/(?P<id>\d+)/(?P<flag>(true|false))`, func(req struct {
		Param *T
		Query *struct {
			Page int
			T
		}
		Header *struct {
			Cookie string
		}
		Body *struct {
			Name string
			T
		}
		Ctx *Context
	}, resp *struct {
		Error  error
		Data   interface{}
		Header *struct {
			SetCookie string `header:"Set-Cookie"`
		}
	}) {
		fmt.Printf("req.Param: %+v\n", req.Param)
		fmt.Printf("req.Query: %+v\n", req.Query)
		fmt.Printf("req.Header: %+v\n", req.Header)
		fmt.Printf("req.Body: %+v\n", req.Body)
		fmt.Printf("req.Ctx not nil: %v\n", req.Ctx != nil)

		resp.Data = []int{1, 2, 3}
		resp.Header = &struct {
			SetCookie string `header:"Set-Cookie"`
		}{SetCookie: "c=d"}
	})

	req, err := http.NewRequest(
		"GET",
		"http://localhost/users/123/true?page=3&type=users&Id=123&flag=true",
		strings.NewReader(`{"name":"张三", "type":"users", "id": 123, "flag": true}`),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cookie", "a=b")

	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	resp := rw.Result()
	fmt.Printf("resp.Status: %v\n", resp.Status)
	delete(resp.Header, "Content-Type")
	fmt.Printf("resp.Header: %v\n", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp.Body: %v %v\n", string(body), err)

	// Output:
	// req.Param: &{Type:users Id:123 Flag:true}
	// req.Query: &{Page:3 T:{Type:users Id:123 Flag:true}}
	// req.Header: &{Cookie:a=b}
	// req.Body: &{Name:张三 T:{Type:users Id:123 Flag:true}}
	// req.Ctx not nil: true
	// resp.Status: 200 OK
	// resp.Header: map[Set-Cookie:[c=d]]
	// resp.Body: {"code":"ok","message":"success","data":[1,2,3]}
	//  <nil>
}
