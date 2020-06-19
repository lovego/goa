package goa

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
)

type HandlerFunc func(*Context)

func (h HandlerFunc) String() string {
	return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
}

type HandlerFuncs []HandlerFunc

func (hs HandlerFuncs) String() string {
	return hs.StringIndent("")
}

func (hs HandlerFuncs) StringIndent(indent string) string {
	if len(hs) == 0 {
		return "[ ]"
	}
	var buf bytes.Buffer
	buf.WriteString("[\n")
	for _, h := range hs {
		buf.WriteString(indent + "  " + fmt.Sprint(h) + "\n")
	}
	buf.WriteString(indent + "]")
	return buf.String()
}

func convert(h interface{}) (HandlerFunc, error) {
	// reflect.ValueOf(h)
	return nil, nil
}
