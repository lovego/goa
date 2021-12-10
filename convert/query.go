package convert

import (
	"fmt"
	"log"
	"reflect"
	"strings"
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
	if len(map2strs) == 0 {
		return nil
	}
	Traverse(value, true, func(v reflect.Value, f reflect.StructField) bool {
		paramName, arrayParamName := queryParamName(f)
		if paramName == "" {
			return true
		}
		// value is always empty, so Set only when len(values) > 0
		if values := queryParamValues(map2strs, paramName, arrayParamName); len(values) > 0 {
			if err = SetArray(v, values); err != nil {
				err = fmt.Errorf("req.Query.%s: %s", f.Name, err.Error())
			}
			return err == nil // if err == nil, go on Traverse
		}
		return true // go on Traverse
	})
	return
}

func queryParamName(field reflect.StructField) (string, string) {
	tag := field.Tag.Get("json")
	if tag == "-" {
		return "", ""
	}
	name := field.Name
	if tag != "" {
		if idx := strings.Index(tag, ","); idx > 0 {
			name = tag[:idx]
		} else if idx < 0 {
			name = tag
		}
	}
	if kind := field.Type.Kind(); kind == reflect.Slice || kind == reflect.Array {
		return name, name + "[]"
	}
	return name, ""
}

func queryParamValues(map2strs map[string][]string, paramName, arrayParamName string) []string {
	if values, ok := map2strs[paramName]; ok {
		return values
	}
	if arrayParamName != "" {
		if values, ok := map2strs[arrayParamName]; ok {
			return values
		}
	}

	paramName, arrayParamName = strings.ToLower(paramName), strings.ToLower(arrayParamName)
	for key, values := range map2strs {
		key = strings.ToLower(key)
		if key == paramName || arrayParamName != "" && key == arrayParamName {
			return values
		}
	}

	return nil
}
