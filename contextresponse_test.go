package goa

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
)

func ExampleContext_Status() {
	c := &Context{ResponseWriter: httptest.NewRecorder()}
	fmt.Println(c.Status())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &Context{ResponseWriter: w}
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
	c := &Context{ResponseWriter: httptest.NewRecorder()}
	fmt.Println("empty" + string(c.ResponseBody()))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &Context{ResponseWriter: w}
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
	c := &Context{ResponseWriter: httptest.NewRecorder()}
	fmt.Println(c.ResponseBodySize())

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &Context{ResponseWriter: w}
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
