package goa

import (
	"log"
	"reflect"

	"github.com/lovego/goa/converters"
	"github.com/lovego/structs"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func newRespWriteFunc(typ reflect.Type) (reflect.Type, func(*Context, reflect.Value)) {
	if typ.Kind() != reflect.Ptr {
		log.Panic("resp parameter of handler func must be a struct pointer.")
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		log.Panic("resp parameter of handler func must be a struct pointer.")
	}
	validateRespFields(typ)
	return typ, func(ctx *Context, resp reflect.Value) {
		var data interface{}
		var err error

		structs.Traverse(resp, false, func(v reflect.Value, f reflect.StructField) bool {
			switch f.Name {
			case "Error":
				if e := v.Interface(); e != nil {
					err = e.(error)
				}
			case "Data":
				data = v.Interface()
			case "Header":
				converters.WriteRespHeader(v, ctx.ResponseWriter.Header())
			}
			return true
		})
		ctx.Data(data, err)
	}
}

func validateRespFields(typ reflect.Type) {
	structs.TraverseType(typ, func(f reflect.StructField) {
		switch f.Name {
		case "Data":
			// data can be of any type
		case "Error":
			if !f.Type.Implements(errorType) {
				log.Panicf(`resp.Error must be of "error" type.`)
			}
		case "Header":
			converters.ValidateRespHeader(f.Type)
		default:
			log.Panicf("Unknow field: resp.%s.", f.Name)
		}
	})

}
