package goa

import (
	"log"
	"reflect"
)

func convert(h interface{}, path string) HandlerFunc {
	if handler, ok := h.(HandlerFunc); ok {
		return handler
	}

	val := reflect.ValueOf(h)

	typ := val.Type()
	if typ.Kind() != reflect.Func {
		log.Panic("handler must be a func.")
	}
	if typ.NumIn() != 2 {
		log.Panic("handler func must have exactly two parameters.")
	}
	if typ.NumOut() != 0 {
		log.Panic("handler func must have no return values.")
	}

	reqConvertFunc := newReqConvertFunc(typ.In(0), path)
	respTyp, respWriteFunc := newRespWriteFunc(typ.In(1))

	return func(ctx *Context) {
		req, err := reqConvertFunc(ctx)
		if err != nil {
			ctx.Data(nil, err)
			return
		}
		resp := reflect.New(respTyp)
		val.Call([]reflect.Value{req, resp})
		if respWriteFunc != nil {
			respWriteFunc(ctx, resp)
		}
	}
}

// handler example
func handlerExample(req *struct {
	Title  string
	Desc   string
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
	}
	Session struct {
		UserId int64
	}
	Ctx *Context
}, resp *struct {
	Error   error
	Headers struct {
		SetCookie string
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
