# goa
A golang http router with regexp support and document generation.

[![Build Status](https://travis-ci.org/lovego/goa.svg?branch=master)](https://travis-ci.org/lovego/goa)
[![Coverage Status](https://img.shields.io/coveralls/github/lovego/goa/master.svg)](https://coveralls.io/github/lovego/goa?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/goa?1)](https://goreportcard.com/report/github.com/lovego/goa)
[![GoDoc](https://godoc.org/github.com/lovego/goa?status.svg)](https://godoc.org/github.com/lovego/goa)


## Usage
### The `main` package
```go
package main

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/lovego/fs"
	"github.com/lovego/goa"
	"github.com/lovego/goa/benchmark/example/users"
	"github.com/lovego/goa/middlewares"
	"github.com/lovego/goa/server"
	"github.com/lovego/goa/utilroutes"
	"github.com/lovego/logger"
)

func main() {
	router := goa.New()
	// logger should comes first, to handle panic and log all requests
	router.Use(middlewares.NewLogger(logger.New(os.Stdout)).Record)
	router.Use(middlewares.NewCORS(allowOrigin).Check)
	utilroutes.Setup(router)

	if os.Getenv("GOA_DOC") != "" {
		router.DocDir(filepath.Join(fs.SourceDir(), "docs", "apis"))
	}

	// If donn't need documentation, use this simple style.
	router.Get("/", func(c *goa.Context) {
		c.Data("index", nil)
	})

	// If need documentation, use this style for automated routes documentation generation.
	router.Group("/users", "用户", "用户相关的接口").
		Get("/", func(req struct {
			Title string `用户列表`
			Desc  string `根据搜索条件获取用户列表`
			Query users.ListReq
		}, resp *struct {
			Error error
			Data  users.ListResp
		}) {
			resp.Data, resp.Error = req.Query.Run()
		}).
		Get(`/(?P<userId>\d+)`, func(req struct {
			Title string `用户详情`
			Desc  string `根据用户ID获取用户详情`
			Param users.DetailReq
		}, resp *struct {
			Error error
			Data  users.DetailResp
		}) {
			resp.Data, resp.Error = req.Param.Run()
		})

	server.ListenAndServe(router)
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

### The `users` package
```go
package users

type ListReq struct {
	Name     string `c:"用户名称"`
	Type     string `c:"用户类型"`
	Page     int    `c:"页码"`
	PageSize int    `c:"每页数据条数"`
}

type ListResp struct {
	TotalSize int `c:"总数据条数"`
	TotalPage int `c:"总页数"`
	Rows      []struct {
		Id    int    `c:"ID"`
		Name  string `c:"名称"`
		Phone string `c:"电话号码"`
	}
}

func (l *ListReq) Run() (ListResp, error) {
	return ListResp{}, nil
}

type DetailReq struct {
	UserId int64 `c:"用户ID"`
}

type DetailResp struct {
	TotalSize int `c:"总数据条数"`
	TotalPage int `c:"总页数"`
	Rows      []struct {
		Id    int    `c:"ID"`
		Name  string `c:"名称"`
		Phone string `c:"电话号码"`
	}
}

func (l *DetailReq) Run() (DetailResp, error) {
	return DetailResp{}, nil
}
```

## Document generation
see [full examples](docs/z_test.go) and [the generated documents](docs/testdata/README.md).

## Default middlewares
- logging with error alarm support
- list of requests in processing
- CORS check

## Attentions
- static route is always matched before regexp route.
- goa use regular expression trees. when match a request, the request is matched from the root node down to leaf node, until a whole match is found. for better performance, once a node is matched, only the children nodes will be traversed for match, no sibbling nodes will be checked even if no match is found in the children nodes. so sibbling routes should always match different pathes and not have overlaps to avoid this flaw.
- call `c.Next()` in middleware to pass control to the next midlleware or route,
  if you don't call `c.Next()` no remaining midlleware or route will be executed.
- generally don't use midlleware after routes, because generally the routes don't call `c.Next()`.

