package converters

import (
	"fmt"
	"reflect"
)

func ConvertSession(sess, value reflect.Value) error {
	typ := sess.Type()
	valueTyp := value.Type()
	if valueTyp == typ {
		sess.Set(value)
		return nil
	} else if typ.Kind() == reflect.Ptr {
		if valueTyp.Kind() != reflect.Ptr && valueTyp == typ.Elem() {
			ptr := reflect.New(typ.Elem())
			ptr.Elem().Set(value)
			sess.Set(ptr)
			return nil
		}
	} else {
		if valueTyp.Kind() == reflect.Ptr && valueTyp.Elem() == typ {
			if !value.IsNil() {
				sess.Set(value.Elem())
			}
			return nil
		}
	}
	return fmt.Errorf("req.Session(%v): got unexpected type %v", typ, valueTyp)
}
