package utilroutes

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var instanceName = getInstanceName()

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
		panic(err)
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
		panic(err)
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
	port := os.Getenv(`GOPORT`)
	if port == `` {
		port = `3000`
	}
	return `:` + port
}
