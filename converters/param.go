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
	names := regexp.MustCompile(path).SubexpNames()[1:] // names[0] is always "".
	if len(names) == 0 {
		log.Panic("req.Param: no parenthesized subexpression in path.")
	}
	if len(names) == 1 && names[0] == "" {
		return ParamConverter{}
	} else if typ.Kind() != reflect.Struct {
		log.Panic("req.Param must be a struct.")
	}

	var fields []ParamField
	for i, name := range names {
		if name != "" {
			if f, ok := typ.FieldByName(UppercaseFirstLetter(name)); ok {
				if isSupportedType(f.Type) {
					fields = append(fields, ParamField{ParamIndex: i, StructField: f})
				} else {
					log.Panicf("req.Param.%s: type must be string, number or bool.", f.Name)
				}
			}
		}
	}
	if len(fields) == 0 {
		log.Panic("req.Param: no matched named parenthesized subexpression in path.")
	}
	return ParamConverter{fields}
}

func (pc ParamConverter) Convert(param reflect.Value, paramsSlice []string) error {
	if len(pc.fields) == 0 {
		return Set(param, paramsSlice[0])
	}
	for _, f := range pc.fields {
		if err := Set(param.FieldByIndex(f.Index), paramsSlice[f.ParamIndex]); err != nil {
			return fmt.Errorf("req.Param.%s: %s", f.Name, err.Error())
		}
	}
	return nil
}

func UppercaseFirstLetter(s string) string {
	if len(s) > 0 && s[0] >= 'a' && s[0] <= 'z' {
		b := []byte(s)
		b[0] -= ('a' - 'A')
		return string(b)
	}
	return s
}

func LowercaseFirstLetter(s string) string {
	if len(s) > 0 && s[0] >= 'A' && s[0] <= 'Z' {
		b := []byte(s)
		b[0] += ('a' - 'A')
		return string(b)
	}
	return s
}
