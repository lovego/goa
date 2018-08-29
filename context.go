package router

import (
	"net/http"
)

type Context struct {
	*http.Request
	http.ResponseWriter
	data map[string]interface{}
}
