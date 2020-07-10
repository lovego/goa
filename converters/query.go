package converters

import (
	"log"
	"reflect"

	"github.com/lovego/structs"
)

func ValidateQuery(typ reflect.Type) {
	if typ.Kind() != reflect.Struct {
		log.Panic("req.Query must be a struct.")
	}
	structs.Traverse(reflect.New(typ).Elem(), true, func(_ reflect.Value, f reflect.StructField) bool {
		if !isSupportedType(f.Type) {
			log.Panicf("req.Query.%s: type must be string, number or bool.", f.Name)
		}
		return false
	})
	return
}

func ConvertQuery(value reflect.Value, map2strs map[string][]string) (err error) {
	structs.Traverse(value, true, func(v reflect.Value, f reflect.StructField) bool {
		values := map2strs[f.Name]
		if len(values) == 0 {
			values = map2strs[lowercaseFirstLetter(f.Name)]
		}
		if len(values) > 0 && values[0] != "" {
			err = Set(v, values[0])
		}
		return err == nil // if err == nil, go on Traverse
	})
	return
}
