package goa

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
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
	if typ.Kind() != reflect.Struct {
		log.Panic("req.Param must be a struct.")
	}
	names := regexp.MustCompile(path).SubexpNames()
	// names[0] is always "".
	if len(names) <= 1 {
		log.Panic("req.Param: no named parenthesized subexpression in path.")
	}
	type field struct {
		paramIndex int
		reflect.StructField
	}

	var fields []field
	for i := 1; i < len(names); i++ {
		if f, ok := typ.FieldByName(capitalize(names[i])); ok {
			if isSupportedType(f.Type) {
				fields = append(fields, field{paramIndex: i - 1, StructField: f})
			} else {
				log.Panic("req.Param.%s: type must be string, number or bool.", f.Name)
			}
		}
	}
	if len(fields) == 0 {
		log.Panic("req.Param: no matched named parenthesized subexpression in path.")
	}
	return func(req reflect.Value, ctx *Context) error {
		param := req.FieldByIndex(index)
		for _, field := range fields {
			if err := set(param.FieldByIndex(field.Index), ctx.Param(field.paramIndex)); err != nil {
				return fmt.Errorf("req.Param.%s: %s", field.Name, err.Error())
			}
		}
		return nil
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
	if typ != typeContextPtr {
		log.Panic("Ctx field of req parameter must be of type '*goa.Context'.")
	}
}

var typeContextPtr = reflect.TypeOf((*Context)(nil))

func newCtxConvertFunc(typ reflect.Type, index []int) convertFunc {
	if typ != typeContextPtr {
		log.Panic("Ctx field of req parameter must be of type '*goa.Context'.")
	}
}

func set(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		if b, err := strconv.ParseBool(s); err != nil {
			return err
		} else {
			v.SetBool(b)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var bits int
		switch v.Kind() {
		case reflect.Int:
			bits = 0
		case reflect.Int8:
			bits = 8
		case reflect.Int16:
			bits = 16
		case reflect.Int32:
			bits = 32
		case reflect.Int64:
			bits = 64
		}
		if i, err := strconv.ParseInt(s, 10, bits); err != nil {
			return err
		} else {
			v.SetInt(i)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var bits int
		switch v.Kind() {
		case reflect.Uint:
			bits = 0
		case reflect.Uint8:
			bits = 8
		case reflect.Uint16:
			bits = 16
		case reflect.Uint32:
			bits = 32
		case reflect.Uint64:
			bits = 64
		}
		if u, err := strconv.ParseUint(s, 10, bits); err != nil {
			return err
		} else {
			v.SetUint(u)
		}
	case reflect.Float32, reflect.Float64:
		var bits int
		switch v.Kind() {
		case reflect.Float32:
			bits = 32
		case reflect.Float64:
			bits = 64
		}
		if f, err := strconv.ParseFloat(s, bits); err != nil {
			return err
		} else {
			v.SetFloat(f)
		}
	}
	return nil
}

func isSupportedType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.String, reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func capitalize(s string) string {
	if len(s) > 0 && s[0] >= 'a' && s[0] <= 'z' {
		b := []byte(s)
		b[0] -= ('a' - 'A')
		return string(b)
	}
	return s
}
