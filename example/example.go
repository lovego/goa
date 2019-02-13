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
