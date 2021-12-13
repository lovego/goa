# goa
A golang http router with regexp and document generation support.

[![Build Status](https://github.com/lovego/goa/actions/workflows/go.yml/badge.svg)](https://github.com/lovego/goa/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/lovego/goa/badge.svg?branch=master&1)](https://coveralls.io/github/lovego/goa)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovego/goa)](https://goreportcard.com/report/github.com/lovego/goa)
[![Documentation](https://pkg.go.dev/badge/github.com/lovego/goa)](https://pkg.go.dev/github.com/lovego/goa@v0.3.2)


## Usage
### demo `main` package
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
	utilroutes.Setup(&router.RouterGroup)

	if os.Getenv("GOA_DOC") != "" {
		router.DocDir(filepath.Join(fs.SourceDir(), "docs", "apis"))
	}

	// If don't need document, use this simple style.
	router.Get("/", func(c *goa.Context) {
		c.Data("index", nil)
	})

	// If need document, use this style for automated routes document generation.
	router.Group("/users", "用户", "用户相关的接口").
		Get("/", func(req struct {
			Title   string        `用户列表`
			Desc    string        `根据搜索条件获取用户列表`
			Query   users.ListReq
			Session users.Session
		}, resp *struct {
			Data  users.ListResp
			Error error
		}) {
			resp.Data, resp.Error = req.Query.Run(&req.Session)
		}).
		Get(`/(\d+)`, func(req struct {
			Title string       `用户详情`
			Desc  string       `根据用户ID获取用户详情`
			Param int64        `用户ID`
			Ctx   *goa.Context 
		}, resp *struct {
			Data  users.DetailResp
			Error error
		}) {
			resp.Data, resp.Error = users.Detail(req.Param)
		})

	if os.Getenv("GOA_DOC") != "" {
		return
	}

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

### demo `users` package
```go
package users

import "time"

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

type Session struct {
	UserId  int64
	LoginAt time.Time
}

func (l *ListReq) Run(sess *Session) (ListResp, error) {
	return ListResp{}, nil
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

func Detail(userId int64) (DetailResp, error) {
	return DetailResp{}, nil
}
```


## Handler func

### The `req` (request) parameter 
The `req` parameter must be a struct, and it can have the following 8 fields.
The `Title` and `Desc` fields are just used for document generation, so can be of any type, and their values are untouched;
The other fields's values are set in advance from the `http.Request`, so can be used directly in the handler.
Except `Session` and `Ctx`, other fields's full tag is used as description for the corresponding object in the document.
For `Param`, `Query`, `Header` and `Body`, if it's a struct, the struct fields's `comment` or `c` tag is used as the description for the fields in the document.

1. `Title`: It's full tag is used as title of the route in document. 
2. `Desc`:  It's full tag is used as description the route in document.
3. `Param`: Subexpression parameters in regular expression path. If there is only one subexpression in the path and it's not named, the whole `Param` is set to the match properly. Otherwise, the `Param` must be a struct, and it's fields are set properly to the corresponding named subexpression. The first letter of the subexpression name is changed to uppercase to find the corresponding field. 
4. `Query`: Query parameters in the the request, `Query` must be a struct, and it's fields are set properly to the corresponding query paramter. The first letter of the query parameter name is changed to uppercase to find the corresponding field in the struct.
5. `Header`: Headers in the request. `Header` must be a struct, and it's fields are set properly to the corresponding header. The field's `header` tag(if present) or it's name is used as the corresponding header name. 
6. `Body`: The request body is set to `Body` using `json.Unmarshal`.
7. `Session`: `Session` is set to `goa.Context.Get("session")`, so the type must be exactly the same. 
8. `Ctx`: `Ctx` must be of type `*goa.Context`.

### The `resp` (response) parameter.
The `resp` parameter must be a struct pointer, and it can have the following 3 fields.
The fields's full tag is used as description for the corresponding object in the document.
For `Data` and `Header`, if it's a struct, the struct fields's `comment` or `c` tag is used as the description for the fields in the document.

1. `Data`: `Data` is writed in response body as a `data` field using `json.Marshal`.
2. `Error`: `Error` is writed in response body as the `code` and `message` fields.
3. `Header`: `Header` is writed in response headers.

see [full examples](docs/z_test.go) and [the generated documents](docs/testdata/README.md).

## Default middlewares
- logging with error alarm support
- list of requests in processing
- CORS check

## Attentions
- static route is always matched before regexp route.
- call `c.Next()` in middleware to pass control to the next midlleware or route,
  if you don't call `c.Next()` no remaining midlleware or route will be executed.
- generally don't use midlleware after routes, because generally the routes don't call `c.Next()`.

