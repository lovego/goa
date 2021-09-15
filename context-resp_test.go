package goa

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/lovego/errs"
)

func ExampleContext_Status() {
	c := &ContextBeforeLookup{ResponseWriter: httptest.NewRecorder()}
	fmt.Println(c.Status())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ContextBeforeLookup{ResponseWriter: w}
		fmt.Println(c.Status())
		c.WriteHeader(http.StatusOK)
		fmt.Println(c.Status())
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.StatusCode)

	// Output:
	// 0
	// 0
	// 200
	// 200
}

func ExampleContext_ResponseBody() {
	c := &ContextBeforeLookup{ResponseWriter: httptest.NewRecorder()}
	fmt.Println("empty" + string(c.ResponseBody()))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ContextBeforeLookup{ResponseWriter: w}
		fmt.Println("empty" + string(c.ResponseBody()))

		c.Write([]byte("1"))
		fmt.Println(string(c.ResponseBody()))

		c.ResponseWriter.Write([]byte("23"))
		fmt.Println(string(c.ResponseBody()))

		c.Write([]byte("456"))
		fmt.Println(string(c.ResponseBody()))
	}))

	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))

	// Output:
	// empty
	// empty
	// 1
	// 1
	// 1456
	// 123456
}

func ExampleContext_ResponseBodySize() {
	c := &ContextBeforeLookup{ResponseWriter: httptest.NewRecorder()}
	fmt.Println(c.ResponseBodySize())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ContextBeforeLookup{ResponseWriter: w}
		fmt.Println(c.ResponseBodySize())

		c.Write([]byte("1"))
		fmt.Println(c.ResponseBodySize())

		c.ResponseWriter.Write([]byte("23"))
		fmt.Println(c.ResponseBodySize())

		c.Write([]byte("456"))
		fmt.Println(c.ResponseBodySize())
	}))
	_, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// 0
	// 0
	// 1
	// 3
	// 6
}

func ExampleContext_Json() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	c.Json(map[string]interface{}{"k": "<value>"})
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)

	// Output:
	// 200
	// map[Content-Type:[application/json; charset=utf-8]]
	// {"k":"<value>"}
	//  <nil>
}

func ExampleContext_Json_error() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	c.Json(map[string]interface{}{"a": "a", "b": func() {}})
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)

	// Output:
	// 200
	// map[Content-Type:[application/json; charset=utf-8]]
	// {"code":"json-marshal-error","message":"json marshal error"} <nil>
}

func ExampleContext_Json2() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	c.Json2(map[string]interface{}{"k": "<value>"}, errors.New("the error"))
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)
	fmt.Println(c.GetError())

	// Output:
	// 200
	// map[Content-Type:[application/json; charset=utf-8]]
	// {"k":"<value>"}
	//  <nil>
	// the error
}

func ExampleContext_Ok() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	c.Ok("success")
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)

	// Output:
	// 200
	// map[Content-Type:[application/json; charset=utf-8]]
	// {"code":"ok","message":"success"}
	//  <nil>
}

func ExampleContext_Data() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	c.Data([]string{"data1", "data2"}, nil)
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)

	// Output:
	// 200
	// map[Content-Type:[application/json; charset=utf-8]]
	// {"code":"ok","message":"success","data":["data1","data2"]}
	//  <nil>
}

func ExampleContext_Data_error() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	c.Data(nil, errors.New("some error"))
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)
	fmt.Println(c.GetError())

	// Output:
	// 500
	// map[Content-Type:[application/json; charset=utf-8]]
	// {"code":"server-err","message":"Server Error."}
	//  <nil>
	// some error
}

func ExampleContext_Data_code_message() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	c.Data(nil, errs.New("the-code", "the message"))
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)
	fmt.Println(c.GetError())

	// Output:
	// 200
	// map[Content-Type:[application/json; charset=utf-8]]
	// {"code":"the-code","message":"the message"}
	//  <nil>
	// <nil>
}

func ExampleContext_Data_code_message_error() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	theErr := errs.New("the-code", "the message")
	theErr.SetError(errors.New("some error"))
	c.Data(nil, theErr)
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)
	fmt.Println(c.GetError())

	// Output:
	// 200
	// map[Content-Type:[application/json; charset=utf-8]]
	// {"code":"the-code","message":"the message"}
	//  <nil>
	// some error
}

func ExampleContext_Data_code_message_data_error() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	theErr := errs.New("the-code", "the message")
	theErr.SetError(errors.New("some error"))
	theErr.SetData([]string{"data1", "data2"})
	c.Data(nil, theErr)
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)
	fmt.Println(c.GetError())

	// Output:
	// 200
	// map[Content-Type:[application/json; charset=utf-8]]
	// {"code":"the-code","message":"the message","data":["data1","data2"]}
	//  <nil>
	// some error
}

func ExampleContext_Redirect() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	c.Redirect("http://example.com/path")
	res := recorder.Result()
	fmt.Println(res.StatusCode)
	fmt.Println(res.Header)
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body), err)
	fmt.Println(c.GetError())

	// Output:
	// 302
	// map[Location:[http://example.com/path]]
	//  <nil>
	// <nil>
}

func ExampleContext_Hijack() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	fmt.Println(c.Hijack())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ContextBeforeLookup{ResponseWriter: w}
		conn, buf, err := c.Hijack()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		buf.Write([]byte("HTTP/1.1 200 OK\r\nHeader-Key: Header-Value\r\n\r\nBody"))
		buf.Flush()
	}))
	defer ts.Close()

	if res, err := http.Get(ts.URL); err != nil {
		log.Fatal(err)
	} else {
		fmt.Println(res.StatusCode)
		fmt.Println(res.Header)
		body, err := ioutil.ReadAll(res.Body)
		fmt.Println(string(body), err)
	}

	// Output:
	// <nil> <nil> the ResponseWriter doesn't support hijacking.
	// 200
	// map[Header-Key:[Header-Value]]
	// Body <nil>
}

func ExampleContext_Flush() {
	recorder := httptest.NewRecorder()
	c := &ContextBeforeLookup{ResponseWriter: recorder}
	c.Flush()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ContextBeforeLookup{ResponseWriter: w}
		c.Flush()
	}))
	defer ts.Close()

	if _, err := http.Get(ts.URL); err != nil {
		log.Fatal(err)
	}

	// Output:
}
