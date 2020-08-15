package converters

import (
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/lovego/date"
)

func SetArray(v reflect.Value, array []string) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	var length = len(array)
	switch v.Kind() {
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), length, length))
	case reflect.Array:
		if v.Len() < length {
			length = v.Len()
		}
	default:
		if array[0] != "" {
			return Set(v, array[0])
		}
		return nil
	}

	for i := 0; i < length; i++ {
		if err := Set(v.Index(i), array[i]); err != nil {
			return err
		}
	}
	return nil
}

var dateType = reflect.TypeOf(date.Date{})

func Set(v reflect.Value, s string) error {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		v = v.Elem()
	}

	switch v.Type() {
	case dateType:
		if dt, err := date.New(s); err != nil {
			return err
		} else {
			v.Set(reflect.ValueOf(*dt))
			return nil
		}
	}

	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
		return nil
	case reflect.Bool:
		if b, err := strconv.ParseBool(s); err != nil {
			return err
		} else {
			v.SetBool(b)
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return SetInt(v, s)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return SetUint(v, s)
	case reflect.Float32, reflect.Float64:
		return SetFloat(v, s)
	default:
		return json.Unmarshal([]byte(s), v.Addr().Interface())
	}
}

func SetInt(v reflect.Value, s string) error {
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
	return nil
}

func SetUint(v reflect.Value, s string) error {
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
	return nil
}

func SetFloat(v reflect.Value, s string) error {
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
	return nil
}
