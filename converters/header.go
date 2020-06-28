package converters

import (
	"log"
	"net/http"
	"reflect"
)

func ValidateHeader(typ reflect.Type) {
	if typ.Kind() != reflect.Struct {
		log.Panic("req.Header must be a struct.")
	}
}

func ConvertHeader(header reflect.Value, headerMap http.Header) error {
	return nil
}
