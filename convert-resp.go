package goa

import (
	"log"
	"reflect"

	"github.com/lovego/goa/converters"
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

		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			switch f.Name {
			case "Error":
				if e := resp.FieldByIndex(f.Index).Interface(); e != nil {
					err = e.(error)
				}
			case "Data":
				data = resp.FieldByIndex(f.Index).Interface()
			case "Header":
				converters.WriteRespHeader(resp.FieldByIndex(f.Index), ctx.ResponseWriter.Header())
			}
		}
		ctx.Data(data, err)
	}
}

func validateRespFields(typ reflect.Type) {
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
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
	}
}
