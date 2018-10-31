# goa
a golang http router with regexp support, inspired by `httprouter` and `gin`.

[![Build Status](https://travis-ci.org/lovego/goa.svg?branch=master)](https://travis-ci.org/lovego/goa)
[![Coverage Status](https://img.shields.io/coveralls/github/lovego/goa/master.svg)](https://coveralls.io/github/lovego/goa?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/goa?1)](https://goreportcard.com/report/github.com/lovego/goa)
[![GoDoc](https://godoc.org/github.com/lovego/goa?status.svg)](https://godoc.org/github.com/lovego/goa)

## default middlewares
- logging with error alarm support
- processing list of requests in processing
- CORS check

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
	"net/url"
	"os"
	"strings"

	"github.com/lovego/goa"
	"github.com/lovego/goa/middlewares"
	"github.com/lovego/goa/server"
	"github.com/lovego/goa/utilroutes"
	"github.com/lovego/logger"
)

func main() {
	router := goa.New()
	// logger should be the first, to handle panic and log all requests
	router.Use(middlewares.NewLogger(logger.New(os.Stdout)).Record)
	middlewares.SetupProcessingList(router)
	router.Use(middlewares.NewCORS(allowOrigin).Check)

	utilroutes.Setup(router)

	router.Get("/", func(c *goa.Context) {
		c.Data(index())
	})

	router.Group("/users").
		Get("/", func(c *goa.Context) {
			c.Data(userList())
		}).
		// the "X" suffix indicates regular expression
		GetX(`/(\d+)`, func(c *goa.Context) {
			c.Data(userDetail(c.Param(0)))
		})

	server.ListenAndServe(router)
}

func index() (string, error) {
	// do your whatever business logic here
	return "hello, world", nil
}

func userList() (string, error) {
	// do your whatever business logic here
	return "users list", nil
}

func userDetail(userId string) (string, error) {
	// do your whatever business logic here
	return "user: " + userId, nil
}

func allowOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	hostname := u.Hostname()
	return strings.HasSuffix(hostname, ".example.com") ||
		hostname == "example.com" || hostname == "localhost"
}
```
