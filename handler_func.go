package goa

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"runtime"
)

type HandlerFunc func(*Context)

func (h HandlerFunc) String() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

type HandlerFuncs []HandlerFunc

func (hs HandlerFuncs) String() string {
	return hs.StringIndent("")
}

func (hs HandlerFuncs) StringIndent(indent string) string {
	if len(hs) == 0 {
		return "[ ]"
	}
	var buf bytes.Buffer
	buf.WriteString("[\n")
	for _, h := range hs {
		buf.WriteString(indent + "  " + fmt.Sprint(h) + "\n")
	}
	buf.WriteString(indent + "]")
	return buf.String()
}

func convert(h interface{}) HandlerFunc {
	if handler, ok := h.(HandlerFunc); ok {
		return handler
	}
	val := reflect.ValueOf(h)
	typ := val.Type()

	if typ.Kind() != reflect.Func {
		log.Panic("handler must be a func.")
	}
	if typ.NumIn() != 1 {
		log.Panic("handler func must have exactly one parameter.")
	}
	paramType := typ.In(0)
	if paramType.Kind() != reflect.Struct {
		log.Panic("handler parameter must be a struct.")
	}

	return func(c *Context) {
		val.Call([]reflect.Value{})
	}
}

// handler example
func handlerExample(req *struct {
	Title   string
	Desc    string
	Ctx     *Context
	Session struct {
		UserId int64
	}
	Params struct {
		Id int64
	}
	Query struct {
		Id   int64
		Page int64
	}
	Headers struct {
		Cookie string
	}
	Body struct {
		Id   int64
		Name string
	} `当XXX时的请求体格式`
	Body2 struct {
		Id   int64
		Name string
	} `当YYY时的请求体格式`
}, resp *struct {
	Headers struct {
		Cookie string
	}
	Body struct {
		Id   int64
		Name string
	} `当XXX时的返回体格式`
	Body2 struct {
		Id   int64
		Name string
	} `当YYY时的返回体格式`
}) {
	// resp.Body, resp.Error = users.Get(req.Params.Id)
}
