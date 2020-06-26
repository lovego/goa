package converters

import (
	"reflect"
	"strconv"
)

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

func Set(v reflect.Value, s string) error {
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

func capitalize(s string) string {
	if len(s) > 0 && s[0] >= 'a' && s[0] <= 'z' {
		b := []byte(s)
		b[0] -= ('a' - 'A')
		return string(b)
	}
	return s
}
