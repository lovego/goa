package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/lovego/goa/utilroutes"
)

func ListenAndServe(handler http.Handler) {
	server := &http.Server{Handler: handler}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, os.Interrupt)

	go func() {
		if err := server.Serve(getListener()); err != nil && err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()

	<-ch
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	var wait = time.Minute
	if s := os.Getenv("ProShutdownWait"); s != "" {
		if w, err := time.ParseDuration(s); err == nil {
			wait = w
		}
	}
	beginAt := time.Now()
	c, cancel := context.WithDeadline(context.Background(), beginAt.Add(wait))
	defer cancel()
	if err := server.Shutdown(c); err == nil {
		log.Printf("shutdown. (waited %v)\n", time.Since(beginAt))
	} else {
		log.Println("shutdown error: ", err)
	}
}

var listenControl func(network, address string, c syscall.RawConn) error

func getListener() net.Listener {
	addr := utilroutes.ListenAddr()
	listenConfig := net.ListenConfig{Control: listenControl}
	listener, err := listenConfig.Listen(context.Background(), `tcp`, addr)
	if err != nil {
		log.Panic(err)
	}
	log.Println(color.GreenString(`backend started.(` + addr + `)`))
	return tcpKeepAliveListener{listener.(*net.TCPListener)}
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
