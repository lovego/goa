package goa

import (
	"log"
	"reflect"
	"regexp"
)

func newReqConverterFunc(typ reflect.Type, path string) func(*Context) (reflect.Value, error) {
	isPtr := false
	if typ.Kind() == reflect.Ptr {
		isPtr = true
		typ = typ.Elem()
	}
	convertFuncs := getReqFieldsConvertFuncs(typ, path)

	return func(ctx *Context) (reflect.Value, error) {
		ptr := reflect.New(typ)
		req := ptr.Elem()

		for _, convertFunc := range convertFuncs {
			if err := convertFunc(req, ctx); err != nil {
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
			funcs = append(funcs, newParamConvertFunc(f, path))
		case "Query":
			funcs = append(funcs, newQueryConvertFunc(f))
		case "Header":
			funcs = append(funcs, newHeaderConvertFunc(f))
		case "Body":
			funcs = append(funcs, newBodyConvertFunc(f))
		case "Session":
			funcs = append(funcs, newSessionConvertFunc(f))
		case "Ctx":
			funcs = append(funcs, newCtxConvertFunc(f))
		case "Title", "Desc": // just for doc, does not care here
		default:
			log.Panicf("Unknow field: %s in req parameter.", f.Name)
		}
	}
	return
}

func newParamConvertFunc(f reflect.StructField, path string) convertFunc {
	if f.Type.Kind() != reflect.Struct {
		log.Panic("Param field of req parameter must be a struct.")
	}
	re := regexp.MustCompile(path)
	names := re.SubexpNames()
	if len(names) == 0 {
		log.Panic("Param field of req parameter error: no named parenthesized subexpression in path.")
	}
	var fields []struct {
		Index      []int
		ParamIndex int
	}

}

func newQueryConvertFunc(f reflect.StructField) convertFunc {
	switch f.Type.Kind() {
	case reflect.Struct, reflect.String:
	default:
		log.Panic("Query field of req parameter must be a struct or string.")
	}
}

func newHeaderConvertFunc(f reflect.StructField) convertFunc {
	switch f.Type.Kind() {
	case reflect.Struct, reflect.Map:
	default:
		log.Panic("Query field of req parameter must be a struct or map.")
	}
}

func newBodyConvertFunc(f reflect.StructField) convertFunc {
	switch f.Type.Kind() {
	case reflect.Struct, reflect.Map:
	default:
		log.Panic("Query field of req parameter must be a struct or map.")
	}
}

func newSessionConvertFunc(f reflect.StructField) convertFunc {
	if f.Type != typeContextPtr {
		log.Panic("Ctx field of req parameter must be of type '*goa.Context'.")
	}
}

var typeContextPtr = reflect.TypeOf((*Context)(nil))

func newCtxConvertFunc(f reflect.StructField) convertFunc {
	if f.Type != typeContextPtr {
		log.Panic("Ctx field of req parameter must be of type '*goa.Context'.")
	}
}
