// +build server

package main

import (
	"net/http"
	"time"

	"github.com/lovego/goa/server"
)

type testHandler struct{}

func (h testHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
	time.Sleep(5 * time.Second)
}

func main() {
	server.ListenAndServe(testHandler{})
}
