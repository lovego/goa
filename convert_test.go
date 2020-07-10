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
	fmt.Printf("resp.Header: %v\n", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp.Body: %v %v\n", string(body), err)

	// Output:
	// req.Param: {Type:users Id:123 Flag:true}
	// req.Query: {Page:3 T:{Type:users Id:123 Flag:true}}
	// req.Header: {Cookie:a=b}
	// req.Body: {Name:张三 T:{Type:users Id:123 Flag:true}}
	// resp.Status: 200 OK
	// resp.Header: map[Content-Type:[application/json; charset=utf-8] Set-Cookie:[c=d]]
	// resp.Body: {"code":"ok","data":[1,2,3],"message":"success"}
	//  <nil>
}
