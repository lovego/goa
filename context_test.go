package goa

import (
    "testing"
    "context"
    "fmt"
    "bytes"
    "net/http"
    "io/ioutil"
)

func TestContext(t *testing.T){
    c := &Context{index: -1}
    bgd := context.Background()
    c.Set("context", context.WithValue(bgd, "custom", 1))
    ctx := c.Context()
    value := ctx.Value("custom")
    if value == nil{
        t.Fail()
    }
    if v, ok:=  value.(int);!ok{
        t.Fail()
    }else if v != 1{
        t.Fail()
    }
}

func ExampleContext_Param() {
    c := &Context{params:[]string{"123", "sdf"}}
    fmt.Println(c.Param(0), "-")
    fmt.Println(c.Param(1), "-")
    fmt.Println(c.Param(2), "-")
    // Output:
    // 123 -
    // sdf -
    //  -
}

func ExampleContext_RequestBody() {
    buf := []byte("hello world!")
    c := &Context{Request: &http.Request{Body:ioutil.NopCloser(bytes.NewBuffer(buf))}}
    body := c.RequestBody()
    fmt.Println(string(body))
    fmt.Println(string(c.RequestBody()))
    fmt.Println(string(c.data[reqBodyKey].([]byte)))
    // Output:
    // hello world!
    // hello world!
    // hello world!
}