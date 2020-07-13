package docs

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/lovego/jsondoc"
	"github.com/lovego/struct_tag"
	"github.com/lovego/structs"
)

func (r *Route) ResHeader(buf *bytes.Buffer) {
	field, ok := r.resp.FieldByName("Header")
	if !ok {
		return
	}

	buf.WriteString("\n## 返回头说明\n")
	if desc := strings.TrimSpace(string(field.Tag)); desc != "" {
		buf.WriteString(desc)
	}
	structs.Traverse(reflect.New(field.Type).Elem(), true, func(_ reflect.Value, f reflect.StructField) bool {
		name, _ := struct_tag.Lookup(string(f.Tag), "header")
		if name == "" {
			name = f.Name
		}
		buf.WriteString(fmt.Sprintf("- %s: %s\n", name, getComment(f.Tag)))
		return true
	})
}

func (r *Route) ResBody(buf *bytes.Buffer) {
	field, ok := r.resp.FieldByName("Data")
	if !ok {
		return
	}

	buf.WriteString("\n## 返回体说明（application/json）\n")
	if desc := strings.TrimSpace(string(field.Tag)); desc != "" {
		buf.WriteString(desc)
	}
	buf.WriteString("```json5\n")
	if b, err := jsondoc.MarshalIndent(
		resBody{Data: reflect.Zero(field.Type).Interface()}, false, "", "  ",
	); err != nil {
		log.Panic(err)
	} else {
		buf.Write(b)
	}
	buf.WriteString("\n```\n")
}

type resBody struct {
	Code    string      `json:"code"    c:"ok表示成功，其他表示错误代码"`
	Message string      `json:"message" c:"与code对应的描述信息"`
	Data    interface{} `json:"data"    c:"返回的数据"`
}

func (r *Route) ResError(buf *bytes.Buffer) {
	field, _ := r.resp.FieldByName("Error")

	if desc := strings.TrimSpace(string(field.Tag)); desc != "" {
		buf.WriteString("\n## 错误码说明（application/json）\n")
		buf.WriteString(desc)
	}
}
