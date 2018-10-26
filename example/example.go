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
	router.Use(middlewares.NewLogger(logger.New(os.Stdout)).Middleware)
	router.Use(middlewares.Ps)

	router.Get("/", func(ctx *goa.Context) {
		ctx.Data("hello, world", nil)
	})

	server.ListenAndServe(router)
}
