package converters

import (
	"log"
	"reflect"

	"github.com/lovego/structs"
)

func ValidateQueryOrHeader(typ reflect.Type, name string) {
	if typ.Kind() != reflect.Struct {
		log.Panicf("req.%s must be a struct.", name)
	}
	structs.Traverse(reflect.New(typ).Elem(), func(_ reflect.Value, f reflect.StructField) bool {
		if !isSupportedType(f.Type) {
			log.Panicf("req.%s.%s: type must be string, number or bool.", name, f.Name)
		}
		return false
	})
	return
}

func ConvertQueryOrHeader(value reflect.Value, map2strs map[string][]string) (err error) {
	structs.Traverse(value, func(v reflect.Value, f reflect.StructField) bool {
		values := map2strs[lowercaseFirstLetter(f.Name)]
		if len(values) > 0 && values[0] != "" {
			err = Set(v, values[0])
		}
		return err != nil // if err != nil, stop Traverse
	})
	return
}
