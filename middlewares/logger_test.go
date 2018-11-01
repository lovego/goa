package middlewares

import (
    "fmt"
    "net/http"
    "github.com/lovego/goa"
    "github.com/lovego/logger"
    "github.com/lovego/goa/server"
    "time"
    "bytes"
    "testing"
)

func TestRecord(t *testing.T) {
    router := goa.New()
    buf := make([]byte, 10000)
    buff := bytes.NewBuffer(buf)
    l := NewLogger(logger.New(buff))
    router.Use(func(c *goa.Context) {
        fmt.Println("middleware 1 pre")
        c.Next()
        fmt.Println("middleware 1 post")
    })
    router.Use(l.Record)
    router.Get("/", func(c *goa.Context) {
        fmt.Println("ok")
        c.Ok("ok")
    })
    s := server.Server{&http.Server{Addr:"localhost:9999", Handler:router}}
    go func() {
        err := s.Server.ListenAndServe()
        if err != nil{
            fmt.Println(err.Error())
        }
    }()
    time.Sleep(2*time.Second)
    _, err := http.Get("http://localhost:9999?query=123")
    if err != nil{
        fmt.Println(err.Error())
        return
    }
    fmt.Println(string(buff.String()))
}
