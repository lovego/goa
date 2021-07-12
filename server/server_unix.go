// +build linux darwin dragonfly freebsd netbsd openbsd

package server

import (
	"log"
	"syscall"

	"golang.org/x/sys/unix"
)

func init() {
	listenControl = func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			err := unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			if err != nil {
				log.Panic(err)
			}
		})
	}
}
