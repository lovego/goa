package router

import (
	"reflect"
	"runtime"
)

type handlerFunc func(*Context)

func (h handlerFunc) String() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

type handlersChain []handlerFunc

func (hs handlersChain) Handle(ctx *Context) {

}
