package goa

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func ExampleContext_RequestBody() {
	buf := []byte("hello world!")
	c := &ContextBeforeLookup{Request: &http.Request{Body: ioutil.NopCloser(bytes.NewBuffer(buf))}}
	body, err := c.RequestBody()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(body))
	}
	body, err = c.RequestBody()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(body))
	}
	// Output:
	// hello world!
	// hello world!
}

func ExampleContext_Scheme() {
	req := httptest.NewRequest("GET", "http://example.com/path", nil)
	c := &ContextBeforeLookup{Request: req}
	fmt.Println(c.Scheme())
	req.Header.Set("X-Forwarded-Proto", "https")
	fmt.Println(c.Scheme())
	// Output:
	// http
	// https
}

func ExampleContext_Url() {
	c := &ContextBeforeLookup{Request: httptest.NewRequest("GET", "/path", nil)}
	fmt.Println(c.Url())

	c = &ContextBeforeLookup{Request: httptest.NewRequest("GET", "http://example.com/path", nil)}
	fmt.Println(c.Url())
	// Output:
	// http://example.com/path
	// http://example.com/path
}

func ExampleContext_ClientAddr() {
	req := httptest.NewRequest("GET", "/path", nil)
	c := &ContextBeforeLookup{Request: req}
	fmt.Println(c.ClientAddr())

	req.Header.Set("X-Real-IP", "192.0.2.2")
	fmt.Println(c.ClientAddr())

	req.Header.Set("X-Forwarded-For", "192.0.2.3")
	fmt.Println(c.ClientAddr())

	// Output:
	// 192.0.2.1
	// 192.0.2.2
	// 192.0.2.3
}
