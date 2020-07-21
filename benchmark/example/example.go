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
