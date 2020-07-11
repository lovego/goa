package docs

import (
	"bytes"
	"fmt"
	"log"
	"reflect"

	"github.com/lovego/jsondoc"
	"github.com/lovego/struct_tag"
	"github.com/lovego/structs"
)

func (r *Route) ResHeader(buf *bytes.Buffer) {
	f, ok := r.resp.FieldByName("Header")
	if !ok {
		return
	}

	buf.WriteString("\n## 返回头说明\n")
	structs.Traverse(reflect.New(f.Type).Elem(), true, func(_ reflect.Value, f reflect.StructField) bool {
		name, _ := struct_tag.Lookup(string(f.Tag), "header")
		if name == "" {
			name = f.Name
		}
		buf.WriteString(fmt.Sprintf("- %s: %s\n", name, getComment(f.Tag)))
		return true
	})
}

func (r *Route) ResBody(buf *bytes.Buffer) {
	f, ok := r.resp.FieldByName("Body")
	if !ok {
		return
	}

	buf.WriteString("\n## 返回体说明（application/json）\n")
	if b, err := jsondoc.MarshalIndent(
		resBody{Data: reflect.Zero(f.Type).Interface()}, false, "", "  ",
	); err != nil {
		log.Panic(err)
	} else {
		buf.Write(b)
	}
}

type resBody struct {
	Code    string      `json:"code"    c:"ok表示成功，其他表示错误代码"`
	Message string      `json:"message" c:"与code对应的描述信息"`
	Data    interface{} `json:"data"    c:"返回的数据"`
}

func (r *Route) ResError(buf *bytes.Buffer) {
	f, _ := r.resp.FieldByName("Error")

	if len(f.Tag) == 0 {
		return
	}

	buf.WriteString("\n## 错误码说明（application/json）\n")
	buf.WriteString(string(f.Tag))
}
