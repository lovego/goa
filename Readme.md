# goa
a golang http router with regexp support, inspired by `httprouter` and `gin`.

[![Build Status](https://travis-ci.org/lovego/goa.svg?branch=master)](https://travis-ci.org/lovego/goa)
[![Coverage Status](https://img.shields.io/coveralls/github/lovego/goa/master.svg)](https://coveralls.io/github/lovego/goa?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/goa?1)](https://goreportcard.com/report/github.com/lovego/goa)
[![GoDoc](https://godoc.org/github.com/lovego/goa?status.svg)](https://godoc.org/github.com/lovego/goa)

## default middlewares
- logging with error alarm
- the processing requests

## attentions
- static route is always matched before regexp route.
- call `c.Next()` in middleware to pass control to the next midlleware or route,
  if you don't call `c.Next()` no remaining midlleware or route will be executed.
- generally don't use midlleware after routes,
  because generally the routes don't call `c.Next()`.

## usage
```go
package main

import (
	"os"

	"github.com/lovego/goa"
	"github.com/lovego/goa/middlewares"
	"github.com/lovego/goa/server"
	"github.com/lovego/logger"
)

func main() {
	router := goa.New()
	router.Use(middlewares.NewLogger(logger.New(os.Stdout)).Record)
	router.Use(middlewares.Ps)

	router.Get("/", func(c *goa.Context) {
		c.Data("hello, world", nil)
	})

	router.Get("/users", func(c *goa.Context) {
		c.Data("users list", nil)
	})

	// the "X" suffix indicates regular expression
	router.GetX(`/users/(\d+)`, func(c *goa.Context) {
		c.Data("user: "+c.Param(0), nil)
	})

	server.ListenAndServe(router)
}
```
