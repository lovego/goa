package api_type

import (
	"errors"
	"reflect"
)

// 自定义类型
func ApiType(typ reflect.Type) (string, bool, error) {
	switch typ.Kind() {
	case reflect.Struct:
		return typ.Name(), true, nil
	case reflect.String, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return typ.Name(), false, nil
	case reflect.Map:
		valueType, isCustom, err := ApiType(typ.Elem())
		if err != nil {
			return "", false, err
		}
		return valueType, isCustom, nil
	case reflect.Array, reflect.Slice:
		valueType, isCustom, err := ApiType(typ.Elem())
		if err != nil {
			return "", false, err
		}
		return valueType, isCustom, nil
	case reflect.Interface:
		return "any", false, nil
	case reflect.Ptr:
		return ApiType(typ.Elem())
	}

	return "", false, errors.New("unsupported type " + typ.Name())
}

// 自定义类型
func ApiTypeCustomList(typeList []reflect.Type) ([]string, error) {
	var list []string
	for _, typ := range typeList {
		typeName, isCustom, err := ApiType(typ)
		if err != nil {
			return nil, err
		}
		if !isCustom {
			continue
		}
		list = append(list, typeName)
	}
	return list, nil
}
