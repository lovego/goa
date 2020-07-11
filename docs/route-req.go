package docs

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/lovego/goa/converters"
	"github.com/lovego/jsondoc"
	"github.com/lovego/struct_tag"
	"github.com/lovego/structs"
)

func (r *Route) Title(method, fullPath string) string {
	var title string
	if f, ok := r.req.FieldByName("Title"); ok {
		title = strings.TrimSpace(string(f.Tag))
	}
	return title + " (" + method + " " + fullPath + ")"
}

func (r *Route) Param(buf *bytes.Buffer, fullPath string) {
	f, ok := r.req.FieldByName("Param")
	if !ok {
		return
	}

	buf.WriteString("\n## 路径中正则参数（子表达式）说明\n")
	names := regexp.MustCompile(fullPath).SubexpNames()
	for i := 1; i < len(names); i++ { // names[0] is always "".
		if f, ok := f.Type.FieldByName(converters.UppercaseFirstLetter(names[i])); ok {
			buf.WriteString(fmt.Sprintf("- %s (%v): %s\n", names[i], f.Type, getComment(f.Tag)))
		}
	}
}

func (r *Route) Query(buf *bytes.Buffer) {
	f, ok := r.req.FieldByName("Query")
	if !ok {
		return
	}

	buf.WriteString("\n## Query参数说明\n")
	structs.Traverse(reflect.New(f.Type).Elem(), true, func(_ reflect.Value, f reflect.StructField) bool {
		buf.WriteString(fmt.Sprintf("- %s (%v): %s\n", converters.LowercaseFirstLetter(f.Name), f.Type, getComment(f.Tag)))
		return true
	})
}

func (r *Route) Header(buf *bytes.Buffer) {
	f, ok := r.req.FieldByName("Header")
	if !ok {
		return
	}

	buf.WriteString("\n## 请求头说明\n")
	structs.Traverse(reflect.New(f.Type).Elem(), true, func(_ reflect.Value, f reflect.StructField) bool {
		name, _ := struct_tag.Lookup(string(f.Tag), "header")
		if name == "" {
			name = f.Name
		}
		buf.WriteString(fmt.Sprintf("- %s (%v): %s\n", name, f.Type, getComment(f.Tag)))
		return true
	})
}

func (r *Route) Body(buf *bytes.Buffer) {
	f, ok := r.req.FieldByName("Body")
	if !ok {
		return
	}

	buf.WriteString("\n## 请求体说明（application/json）\n")
	if b, err := jsondoc.MarshalIndent(
		reflect.Zero(f.Type).Interface(), false, "", "  ",
	); err != nil {
		log.Panic(err)
	} else {
		buf.Write(b)
	}
}
