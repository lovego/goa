// +build server

package main

import (
	"net/http"

	"github.com/lovego/goa/server"
)

type testHandler struct{}

func (h testHandler) ServeHTTP(http.ResponseWriter, *http.Request) {
}

func main() {
	server.ListenAndServe(testHandler{})
}
