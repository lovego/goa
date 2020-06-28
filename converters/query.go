package converters

import (
	"log"
	"net/url"
	"reflect"
)

func ValidateQuery(typ reflect.Type) {
	if typ.Kind() != reflect.Struct {
		log.Panic("req.Query must be a struct.")
	}
}

func ConvertQuery(header reflect.Value, queryMap url.Values) error {
	return nil
}
