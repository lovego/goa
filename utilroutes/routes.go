package utilroutes

import (
	"fmt"
	httpPprof "net/http/pprof"
	"runtime"
	"runtime/pprof"
	"strconv"

	"github.com/lovego/goa"
)

func Setup(router *goa.Router) {
	router.Get(`/_alive`, func(ctx *goa.Context) {
		ctx.Write([]byte(`ok`))
	})
	router.Use(recordRequests) // ps middleware

	debug := router.Group(`/_debug`)
	debug.Get(`/`, func(ctx *goa.Context) {
		ctx.Write(debugIndex())
	})

	debug.Get(`/reqs`, func(ctx *goa.Context) {
		ctx.Write(requests.ToJson())
	})

	// pprof
	debug.Get(`/cpu`, func(ctx *goa.Context) {
		// ctx.Write([]byte(instanceName + "\n"))
		httpPprof.Profile(ctx.ResponseWriter, ctx.Request)
	})

	debug.Get(`/(\w+)`, func(ctx *goa.Context) {
		name := ctx.Param(0)
		ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		ctx.ResponseWriter.Header().Set("Instance-Name", instanceName)
		profile := pprof.Lookup(name)
		if profile == nil {
			ctx.WriteHeader(404)
			fmt.Fprintf(ctx, "Unknown profile: %s\n", name)
			return
		}
		if name == "heap" && ctx.FormValue("gc") != `` {
			runtime.GC()
		}
		debugLevel, _ := strconv.Atoi(ctx.FormValue("debug"))
		profile.WriteTo(ctx, debugLevel)
	})

	debug.Get(`/trace`, func(ctx *goa.Context) {
		httpPprof.Trace(ctx.ResponseWriter, ctx.Request)
	})
}
