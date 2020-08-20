package goa

import (
	"encoding/json"
	"log"
	"reflect"

	"github.com/lovego/goa/convert"
	"github.com/lovego/structs"
)

type todoReqFields struct {
	Query  bool
	Header bool
	Body   bool
}

func newReqConvertFunc(typ reflect.Type, path string) func(*Context) (reflect.Value, error) {
	isPtr := false
	if typ.Kind() == reflect.Ptr {
		isPtr = true
		typ = typ.Elem()
	}
	param, todo := validateReqFields(typ, path)

	return func(ctx *Context) (reflect.Value, error) {
		ptr := reflect.New(typ)
		req := ptr.Elem()

		var err error
		structs.Traverse(req, true, func(value reflect.Value, f reflect.StructField) bool {
			switch f.Name {
			case "Param":
				err = param.Convert(value, ctx.params)
			case "Query":
				if todo.Query {
					err = convert.Query(value, ctx.Request.URL.Query())
				}
			case "Header":
				if todo.Header {
					err = convert.Header(value, ctx.Request.Header)
				}
			case "Body":
				if todo.Body {
					err = convertReqBody(value, ctx)
				}
			case "Session":
				if sess := ctx.Get("session"); sess != nil {
					err = convert.Session(value, reflect.ValueOf(sess))
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

func validateReqFields(typ reflect.Type, path string) (
	param convert.ParamConverter, todo todoReqFields,
) {
	if typ.Kind() != reflect.Struct {
		log.Panic("req parameter of handler func must be a struct or struct pointer.")
	}

	structs.TraverseType(typ, func(f reflect.StructField) {
		switch f.Name {
		case "Param":
			param = convert.GetParamConverter(f.Type, path)
		case "Query":
			if !isEmptyStruct(f.Type) {
				convert.ValidateQuery(f.Type)
				todo.Query = true
			}
		case "Header":
			if !isEmptyStruct(f.Type) {
				convert.ValidateHeader(f.Type)
				todo.Header = true
			}
		case "Ctx":
			if f.Type != typeContextPtr {
				log.Panic("Ctx field of req parameter must be of type '*goa.Context'.")
			}
		case "Body":
			if !isEmptyStruct(f.Type) {
				todo.Body = true
			}
		case "Session": // can be any type, don't need to validate here.
		case "Title", "Desc": // just for doc, does not care here
		default:
			log.Panicf("Unknown field: req.%s.", f.Name)
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

func isEmptyStruct(typ reflect.Type) bool {
	return typ.Kind() == reflect.Struct && typ.NumField() == 0
}
