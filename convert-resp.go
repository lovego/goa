package goa

import (
	"log"
	"reflect"
)

func newRespWriteFunc(typ reflect.Type) (reflect.Type, func(*Context, reflect.Value)) {
	if typ.Kind() != reflect.Ptr {
		log.Panic("resp parameter of handler func must be a struct pointer.")
	}
	typ = typ.Elem()
	if typ.Kind() != reflect.Struct {
		log.Panic("resp parameter of handler func must be a struct pointer.")
	}
	return typ, nil
}
