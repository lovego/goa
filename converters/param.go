package converters

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
)

type ParamConverter struct {
	fields []ParamField
}

type ParamField struct {
	ParamIndex int
	reflect.StructField
}

func ForParam(typ reflect.Type, path string) ParamConverter {
	if typ.Kind() != reflect.Struct {
		log.Panic("req.Param must be a struct.")
	}
	names := regexp.MustCompile(path).SubexpNames()
	// names[0] is always "".
	if len(names) <= 1 {
		log.Panic("req.Param: no named parenthesized subexpression in path.")
	}

	var fields []ParamField
	for i := 1; i < len(names); i++ {
		if f, ok := typ.FieldByName(capitalize(names[i])); ok {
			if isSupportedType(f.Type) {
				fields = append(fields, ParamField{ParamIndex: i - 1, StructField: f})
			} else {
				log.Panic("req.Param.%s: type must be string, number or bool.", f.Name)
			}
		}
	}
	if len(fields) == 0 {
		log.Panic("req.Param: no matched named parenthesized subexpression in path.")
	}
	return ParamConverter{fields}
}

func (pc ParamConverter) Convert(param reflect.Value, paramsSlice []string) error {
	for _, field := range pc.fields {
		if err := Set(param.FieldByIndex(field.Index), paramsSlice[field.ParamIndex]); err != nil {
			return fmt.Errorf("req.Param.%s: %s", field.Name, err.Error())
		}
	}
	return nil
}
