package api_type

import (
	"errors"
	"fmt"
	"reflect"
)

func TsPrimitiveType(typeName string) (string, error) {
	switch typeName {
	case "string":
		return "string", nil
	case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32", "int64", "uint64":
		return "number", nil
	case "float", "float32", "float64":
		return "number", nil
	case "bool":
		return "boolean", nil
	case "[]byte":
		return "Blob", nil
	case "interface{}":
		return "any", nil
	}
	return "", errors.New("unsupported primitive type " + typeName)
}
func TsRawType(typeName string) (string, error) {
	t, err := TsPrimitiveType(typeName)
	if err == nil {
		return "", errors.New(typeName + " is not raw type")
	}

	return t, nil
}
func TsType(typeName string) string {
	t, err := TsPrimitiveType(typeName)
	if err == nil {
		return t
	}
	return typeName
}
func genTsMemberType(m reflect.Type) (typ string, err error) {
	typ, err = ToTypeScriptType(m)
	return
}
func ToTypeScriptType(typ reflect.Type) (string, error) {
	switch typ.Kind() {
	case reflect.Struct:
		if typ.Name() == "Time" ||
			typ.Name() == "Decimal" ||
			typ.Name() == "Date" {
			return "string", nil
		}
		return typ.Name(), nil
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		r, err := TsPrimitiveType(typ.Kind().String())
		if err != nil {
			return "", err
		}

		return r, nil
	case reflect.Map:
		valueType, err := ToTypeScriptType(typ.Elem())
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("{ [key: string]: %s }", valueType), nil
	case reflect.Array, reflect.Slice:
		if typ.Name() == "[]byte" {
			return "Blob", nil
		}

		valueType, err := ToTypeScriptType(typ.Elem())
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("Array<%s>", valueType), nil
	case reflect.Interface:
		return "any", nil
	case reflect.Pointer:
		return ToTypeScriptType(typ.Elem())
	}

	return "", errors.New("unsupported type " + typ.Name())
}
