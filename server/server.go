package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/fatih/color"
)

type Server struct {
	*http.Server
}

func ListenAndServe(handler http.Handler) {
	s := Server{&http.Server{}}
	s.ListenAndServe(handler)
}

func (s Server) ListenAndServe(handler http.Handler) {
	s.Server.Handler = handler

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, os.Interrupt)

	go func() {
		if err := s.Server.Serve(getListener()); err != nil && err != http.ErrServerClosed {
			log.Panic(err)
		}
	}()

	<-ch
	s.gracefulShutdown()
}

func (s Server) gracefulShutdown() {
	if runtime.GOOS != "linux" {
		return
	}
	c, cancel := context.WithDeadline(context.Background(), time.Now().Add(7*time.Second))
	defer cancel()
	if err := s.Server.Shutdown(c); err == nil {
		log.Println(`shutdown`)
	} else {
		log.Println("shutdown error: ", err)
	}
}

func getListener() net.Listener {
	port := os.Getenv(`GOPORT`)
	if port == `` {
		port = `3000`
	}
	addr := `:` + port
	ln, err := net.Listen(`tcp`, addr)
	if err != nil {
		log.Panic(err)
	}
	log.Println(color.GreenString(`started.(` + addr + `)`))
	return tcpKeepAliveListener{ln.(*net.TCPListener)}
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
