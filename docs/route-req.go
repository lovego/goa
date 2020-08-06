package docs

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/lovego/jsondoc"
	"github.com/lovego/strs"
	"github.com/lovego/struct_tag"
	"github.com/lovego/structs"
)

func (r *Route) Title() string {
	if f, ok := r.req.FieldByName("Title"); ok {
		return strings.TrimSpace(string(f.Tag))
	}
	return ""
}

func (r *Route) Desc(buf *bytes.Buffer) {
	if f, ok := r.req.FieldByName("Desc"); ok {
		if desc := strings.TrimSpace(string(f.Tag)); desc != "" {
			buf.WriteString(desc + "\n\n")
		}
	}
}

func (r *Route) Param(buf *bytes.Buffer, fullPath string) {
	field, ok := r.req.FieldByName("Param")
	if !ok {
		return
	}

	buf.WriteString("\n## 路径中正则参数（子表达式）说明\n")
	if desc := strings.TrimSpace(string(field.Tag)); desc != "" {
		buf.WriteString(desc + "\n\n")
	}

	names := regexp.MustCompile(fullPath).SubexpNames()[1:] // names[0] is always "".
	for _, name := range names {
		if name != "" {
			if f, ok := field.Type.FieldByName(strs.FirstLetterToUpper(name)); ok {
				buf.WriteString(fmt.Sprintf("- %s (%v): %s\n", name, f.Type, getComment(f.Tag)))
			}
		}
	}
}

func (r *Route) Query(buf *bytes.Buffer) {
	field, ok := r.req.FieldByName("Query")
	if !ok {
		return
	}

	buf.WriteString("\n## Query参数说明\n")
	if desc := strings.TrimSpace(string(field.Tag)); desc != "" {
		buf.WriteString(desc + "\n\n")
	}
	structs.Traverse(reflect.New(field.Type).Elem(), true,
		func(_ reflect.Value, f reflect.StructField) bool {
			buf.WriteString(fmt.Sprintf(
				"- %s (%v): %s\n", strs.FirstLetterToLower(f.Name), f.Type, getComment(f.Tag),
			))
			return true
		})
}

func (r *Route) Header(buf *bytes.Buffer) {
	field, ok := r.req.FieldByName("Header")
	if !ok {
		return
	}

	buf.WriteString("\n## 请求头说明\n")
	if desc := strings.TrimSpace(string(field.Tag)); desc != "" {
		buf.WriteString(desc + "\n\n")
	}
	structs.Traverse(reflect.New(field.Type).Elem(), true,
		func(_ reflect.Value, f reflect.StructField) bool {
			name, _ := struct_tag.Lookup(string(f.Tag), "header")
			if name == "" {
				name = f.Name
			}
			buf.WriteString(fmt.Sprintf("- %s (%v): %s\n", name, f.Type, getComment(f.Tag)))
			return true
		})
}

func (r *Route) Body(buf *bytes.Buffer) {
	field, ok := r.req.FieldByName("Body")
	if !ok {
		return
	}

	buf.WriteString("\n## 请求体说明（application/json）\n")
	if desc := strings.TrimSpace(string(field.Tag)); desc != "" {
		buf.WriteString(desc + "\n\n")
	}
	buf.WriteString("```json5\n")
	if b, err := jsondoc.MarshalIndent(
		reflect.Zero(field.Type).Interface(), false, "", "  ",
	); err != nil {
		log.Panic(err)
	} else {
		buf.Write(b)
	}
	buf.WriteString("\n```\n")
}
