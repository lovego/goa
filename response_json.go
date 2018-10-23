package goa

import (
	"encoding/json"
	"log"
)

func (ctx *Context) Json(data interface{}) {
	bytes, err := json.Marshal(data)
	if err == nil {
		ctx.Response.Header.Set(`Content-Type`, `application/json; charset=utf-8`)
		ctx.ResponseWriter.Write(bytes)
	} else {
		log.Panic(err)
	}
}

func (ctx *Context) Json2(data interface{}, err error) {
	if err != nil {
		ctx.LogError(err)
	}
	if bytes, err := json.Marshal(data); err == nil {
		ctx.Response.Header.Set(`Content-Type`, `application/json; charset=utf-8`)
		ctx.ResponseWriter.Write(bytes)
	} else {
		panic(err)
	}
}

func (ctx *Context) Ok(message string) {
	ctxult := make(map[string]interface{})
	ctxult["code"] = "ok"
	ctxult["message"] = message
	ctx.Json(ctxult)
}

func (ctx *Context) Data(data interface{}, err error) {
	ctx.DataWithKey(data, err, `data`)
}

func (ctx *Context) Result(data interface{}, err error) {
	ctx.DataWithKey(data, err, `ctxult`)
}

func (ctx *Context) DataWithKey(data interface{}, err error, key string) {
	ctxult := make(map[string]interface{})
	if err == nil {
		ctxult[`code`] = `ok`
		ctxult[`message`] = `success`
	} else {
		if erro, ok := err.(interface {
			Code() string
			Message() string
		}); ok && erro.Code() != "" {
			ctxult[`code`] = erro.Code()
			ctxult[`message`] = erro.Message()
			if e, ok := err.(interface {
				Err() error
			}); ok && e.Err() != nil {
				ctx.LogError(err)
			}
		} else {
			ctx.WriteHeader(500)
			ctxult[`code`] = `server-err`
			ctxult[`message`] = `Server Error.`
			ctx.LogError(err)
		}
	}

	if data != nil {
		ctxult[key] = data
	} else if err != nil {
		if erro, ok := err.(interface {
			Data() interface{}
		}); ok && erro.Data() != nil {
			ctxult[key] = erro.Data()
		}
	}
	ctx.Json(ctxult)
}

func (ctx *Context) LogError(err error) {
	ctx.err = err
}

func (ctx *Context) GetError() error {
	return ctx.err
}
