package goa

import (
	"log"
	"reflect"

	"github.com/lovego/goa/convert"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func newRespWriteFunc(typ reflect.Type, hasCtx bool) (reflect.Type, func(*Context, reflect.Value)) {
	if typ.Kind() != reflect.Ptr {
		log.Panic("resp parameter of handler func must be a struct pointer.")
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		log.Panic("resp parameter of handler func must be a struct pointer.")
	}
	if validateRespFields(typ) {
		return typ, nil
	}
	return typ, func(ctx *Context, resp reflect.Value) {
		if hasCtx && ctx.ResponseBodySize() > 0 {
			return
		}

		var data interface{}
		var err error

		convert.Traverse(resp, false, func(v reflect.Value, f reflect.StructField) bool {
			switch f.Name {
			case "Error":
				if e := v.Interface(); e != nil {
					err = e.(error)
				}
			case "Data":
				data = v.Interface()
			case "Header":
				convert.WriteRespHeader(v, ctx.ResponseWriter.Header())
			}
			return true
		})
		ctx.Data(data, err)
	}
}

func validateRespFields(typ reflect.Type) bool {
	empty := true
	convert.TraverseType(typ, func(f reflect.StructField) {
		switch f.Name {
		case "Data":
			// data can be of any type
		case "Error":
			if !f.Type.Implements(errorType) {
				log.Panicf(`resp.Error must be of "error" type.`)
			}
		case "Header":
			convert.ValidateRespHeader(f.Type)
		default:
			log.Panicf("Unknown field: resp.%s.", f.Name)
		}
		empty = false
	})
	return empty
}
