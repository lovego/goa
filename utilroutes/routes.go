package utilroutes

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/lovego/goa"
)

var instanceName = getInstanceName()

// debugs for register custom debug data
func Setup(router *goa.RouterGroup, debugs ...map[string]interface{}) {
	router.Get(`/_alive`, func(ctx *goa.Context) {
		ctx.Write([]byte(`ok`))
	})
	router.Use(recordRequests) // ps middleware

	debug := router.Group(`/_debug`)
	debug.Use(func(ctx *goa.Context) {
		ctx.ResponseWriter.Header().Set("Instance-Name", instanceName)
		ctx.Next()
	})
	debug.Get(`/`, func(ctx *goa.Context) {
		ctx.Write(debugIndex())
	})
	debug.Get(`/reqs`, func(ctx *goa.Context) {
		ctx.Write(requests.ToJson())
	})
	for i := range debugs {
		for path, data := range debugs[i] {
			debug.Get(path, func(ctx *goa.Context) {
				bytes, err := json.Marshal(map[string]interface{}{
					"Instance": instanceName,
					"Data":     data,
				})
				if err != nil {
					ctx.Write([]byte(fmt.Sprint(err)))
				} else {
					ctx.Write(bytes)
				}
			})
		}
	}

	// pprof
	debug.Get(`/cpu`, func(ctx *goa.Context) {
		cpuProfile(ctx.ResponseWriter, ctx.Request)
	})
	debug.Get(`/(\w+)`, func(ctx *goa.Context) {
		getProfile(ctx.Param(0), ctx.ResponseWriter, ctx.Request)
	})

	debug.Get(`/trace`, func(ctx *goa.Context) {
		runTrace(ctx.ResponseWriter, ctx.Request)
	})
}

func getInstanceName() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}
	return fmt.Sprintf(
		"%s (%s) (Listen At %s)", hostname, strings.Join(ipv4Addrs(), ", "), ListenAddr(),
	)
}

func ipv4Addrs() (result []string) {
	ifcs, err := net.Interfaces()
	if err != nil {
		log.Panic(err)
	}
	for _, ifc := range ifcs {
		if ifc.Flags&net.FlagLoopback == 0 {
			result = append(result, ipv4AddrsOfInterface(ifc)...)
		}
	}
	return result
}
func ipv4AddrsOfInterface(ifc net.Interface) (result []string) {
	addrs, err := ifc.Addrs()
	if err != nil {
		log.Panic(err)
	}
	for _, addr := range addrs {
		if str := addr.String(); strings.IndexByte(str, '.') > 0 { // ipv4
			if i := strings.IndexByte(str, '/'); i >= 0 {
				str = str[:i]
			}
			result = append(result, str)
		}
	}
	return result
}

func ListenAddr() string {
	port := os.Getenv(`ProPORT`)
	if port == `` {
		port = `3000`
	}
	return `:` + port
}
