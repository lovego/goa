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
