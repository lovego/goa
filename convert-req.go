package goa

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/lovego/goa/converters"
	"github.com/lovego/structs"
)

func newReqConvertFunc(typ reflect.Type, path string) func(*Context) (reflect.Value, error) {
	isPtr := false
	if typ.Kind() == reflect.Ptr {
		isPtr = true
		typ = typ.Elem()
	}
	param := validateReqFields(typ, path)

	return func(ctx *Context) (reflect.Value, error) {
		ptr := reflect.New(typ)
		req := ptr.Elem()

		var err error
		structs.Traverse(req, true, func(value reflect.Value, f reflect.StructField) bool {
			switch f.Name {
			case "Param":
				err = param.Convert(value, ctx.params)
			case "Query":
				err = converters.ConvertQuery(value, ctx.Request.URL.Query())
			case "Header":
				err = converters.ConvertHeader(value, ctx.Request.Header)
			case "Body":
				err = convertReqBody(value, ctx)
			case "Session":
				if sess := ctx.Get("session"); sess != nil {
					err = converters.ConvertSession(value, reflect.ValueOf(sess))
				}
			case "Ctx":
				value.Set(reflect.ValueOf(ctx))
			}
			return err == nil
		})
		if err != nil {
			return reflect.Value{}, err
		}

		if isPtr {
			return ptr, nil
		} else {
			return req, nil
		}
	}
}

var typeContextPtr = reflect.TypeOf((*Context)(nil))

func validateReqFields(typ reflect.Type, path string) (param converters.ParamConverter) {
	if typ.Kind() != reflect.Struct {
		log.Panic("req parameter of handler func must be a struct or struct pointer.")
	}

	structs.TraverseType(typ, func(f reflect.StructField) {
		switch f.Name {
		case "Param":
			param = converters.ForParam(f.Type, path)
		case "Query":
			converters.ValidateQuery(f.Type)
		case "Header":
			converters.ValidateHeader(f.Type)
		case "Ctx":
			if f.Type != typeContextPtr {
				log.Panic("Ctx field of req parameter must be of type '*goa.Context'.")
			}
		case "Body", "Session": // can be any type, donn't how to validate here.
		case "Title", "Desc": // just for doc, does not care here
		default:
			log.Panicf("Unknow field: req.%s.", f.Name)
		}
	})
	return
}

func convertReqBody(value reflect.Value, ctx *Context) error {
	body, err := ctx.RequestBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(body, value.Addr().Interface())
}
