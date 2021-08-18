// +build windows

package server

import (
	"log"
	"syscall"

	"golang.org/x/sys/windows"
)

func init() {
	listenControl = func(network, address string, c syscall.RawConn) error {
		return c.Control(func(fd uintptr) {
			err := windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_REUSEADDR, 1)
			if err != nil {
				log.Panic(err)
			}
		})
	}
}
