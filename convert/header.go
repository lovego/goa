package convert

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/lovego/struct_tag"
	"github.com/lovego/structs"
)

func ValidateHeader(typ reflect.Type) {
	if !isStructOrStructPtr(typ) {
		log.Panic("req.Header must be struct or pointer to struct.")
	}
}

func Header(value reflect.Value, map2strs map[string][]string) (err error) {
	structs.Traverse(value, true, func(v reflect.Value, f reflect.StructField) bool {
		if f.Tag.Get("json") == "-" {
			return true
		}

		key, _ := struct_tag.Lookup(string(f.Tag), "header")
		if key == "" {
			key = f.Name
		}
		values := map2strs[key]
		if len(values) > 0 && values[0] != "" {
			if err = Set(v, values[0]); err != nil {
				err = fmt.Errorf("req.Header.%s: %s", f.Name, err.Error())
			}
		}
		return err == nil // if err == nil, go on Traverse
	})
	return
}

func ValidateRespHeader(typ reflect.Type) {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		log.Panic("resp.Header must be struct or pointer to struct.")
	}
	structs.Traverse(reflect.New(typ).Elem(), false, func(_ reflect.Value, f reflect.StructField) bool {
		if f.Type.Kind() != reflect.String {
			log.Panicf("resp.Header.%s: type must be string.", f.Name)
		}
		return true
	})
	return
}

func WriteRespHeader(value reflect.Value, header http.Header) {
	structs.Traverse(value, false, func(v reflect.Value, f reflect.StructField) bool {
		if value := v.String(); value != "" {
			key, _ := struct_tag.Lookup(string(f.Tag), "header")
			if key == "" {
				key = f.Name
			}
			header.Set(key, value)
		}
		return false
	})
	return
}
