package goa

import (
	"log"
	"reflect"

	"github.com/lovego/goa/converters"
)

func newReqConvertFunc(typ reflect.Type, path string) func(*Context) (reflect.Value, error) {
	isPtr := false
	if typ.Kind() == reflect.Ptr {
		isPtr = true
		typ = typ.Elem()
	}
	convertFuncs := getReqFieldsConvertFuncs(typ, path)

	return func(ctx *Context) (reflect.Value, error) {
		ptr := reflect.New(typ)
		req := ptr.Elem()

		for _, convertFn := range convertFuncs {
			if err := convertFn(req, ctx); err != nil {
				return reflect.Value{}, err
			}
		}
		if isPtr {
			return ptr, nil
		} else {
			return req, nil
		}
	}
}

type convertFunc func(reflect.Value, *Context) error

func getReqFieldsConvertFuncs(typ reflect.Type, path string) (funcs []convertFunc) {
	if typ.Kind() != reflect.Struct {
		log.Panic("req parameter of handler func must be a struct or struct pointer.")
	}

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		switch f.Name {
		case "Param":
			funcs = append(funcs, newParamConvertFunc(f.Type, f.Index, path))
		case "Query":
			funcs = append(funcs, newQueryConvertFunc(f.Type, f.Index))
		case "Header":
			funcs = append(funcs, newHeaderConvertFunc(f.Type, f.Index))
		case "Body":
			funcs = append(funcs, newBodyConvertFunc(f.Type, f.Index))
		case "Session":
			funcs = append(funcs, newSessionConvertFunc(f.Type, f.Index))
		case "Ctx":
			funcs = append(funcs, newCtxConvertFunc(f.Type, f.Index))
		case "Title", "Desc": // just for doc, does not care here
		default:
			log.Panicf("Unknow field: req.%s.", f.Name)
		}
	}
	return
}

func newParamConvertFunc(typ reflect.Type, index []int, path string) convertFunc {
	converter := converters.ForParam(typ, path)
	return func(req reflect.Value, ctx *Context) error {
		return converter.Convert(req.FieldByIndex(index), ctx.params)
	}
}

func newQueryConvertFunc(typ reflect.Type, index []int) convertFunc {
	switch typ.Kind() {
	case reflect.String:
	case reflect.Struct:
	default:
		log.Panic("Query field of req parameter must be a struct or string.")
	}
}

func newHeaderConvertFunc(typ reflect.Type, index []int) convertFunc {
	switch typ.Kind() {
	case reflect.Struct, reflect.Map:
	default:
		log.Panic("Query field of req parameter must be a struct or map.")
	}
}

func newBodyConvertFunc(typ reflect.Type, index []int) convertFunc {
	switch typ.Kind() {
	case reflect.Struct, reflect.Map:
	default:
		log.Panic("Query field of req parameter must be a struct or map.")
	}
}

func newSessionConvertFunc(typ reflect.Type, index []int) convertFunc {
	return func(req reflect.Value, ctx *Context) error {
		sess := ctx.Get("session")
		if sess == nil {
			return nil
		}
		return converters.ConvertSession(req.FieldByIndex(index), reflect.ValueOf(sess))
	}
}

var typeContextPtr = reflect.TypeOf((*Context)(nil))

func newCtxConvertFunc(typ reflect.Type, index []int) convertFunc {
	if typ != typeContextPtr {
		log.Panic("Ctx field of req parameter must be of type '*goa.Context'.")
	}
	return func(req reflect.Value, ctx *Context) error {
		req.FieldByIndex(index).Set(reflect.ValueOf(ctx))
		return nil
	}
}
