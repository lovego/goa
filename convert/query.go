package convert

import (
	"log"
	"reflect"
	"strings"

	"github.com/lovego/strs"
	"github.com/lovego/structs"
)

func ValidateQuery(typ reflect.Type) {
	if !isStructOrStructPtr(typ) {
		log.Panic("req.Query must be struct or pointer to struct.")
	}
}

func isStructOrStructPtr(typ reflect.Type) bool {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return typ.Kind() == reflect.Struct
}

func Query(value reflect.Value, map2strs map[string][]string) (err error) {
	structs.Traverse(value, true, func(v reflect.Value, f reflect.StructField) bool {
		var values []string
		if tag := f.Tag.Get("json"); tag == "-" {
			return true
		} else if idx := strings.Index(tag, ","); idx > 0 {
			var name = tag[:idx]
			values = map2strs[name]
			if len(values) == 0 {
				switch f.Type.Kind() {
				case reflect.Slice, reflect.Array:
					values = map2strs[name+"[]"]
				}
			}
		} else {
			values = map2strs[f.Name]
			var lowercaseName string
			if len(values) == 0 {
				lowercaseName = strs.FirstLetterToLower(f.Name)
				values = map2strs[lowercaseName]
			}
			if len(values) == 0 {
				switch f.Type.Kind() {
				case reflect.Slice, reflect.Array:
					if values = map2strs[f.Name+"[]"]; len(values) == 0 {
						values = map2strs[lowercaseName+"[]"]
					}
				}
			}
		}
		if len(values) > 0 {
			err = SetArray(v, values)
		}
		return err == nil // if err == nil, go on Traverse
	})
	return
}
