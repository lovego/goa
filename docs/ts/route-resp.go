package ts

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/lovego/goa/convert"
	"github.com/lovego/goa/docs/ts/ts/api_type"
	"github.com/lovego/struct_tag"
)

func (r *Route) RespHeader(buf *bytes.Buffer) {
	field, ok := r.Resp.FieldByName("Header")
	if !ok {
		return
	}

	buf.WriteString("\n## 返回头说明\n")
	if desc := getComment(field.Tag); desc != "" {
		buf.WriteString(desc + "\n\n")
	}
	convert.Traverse(reflect.New(field.Type).Elem(), true,
		func(_ reflect.Value, f reflect.StructField) bool {
			name, _ := struct_tag.Lookup(string(f.Tag), "header")
			if name == "" {
				name = f.Name
			}
			buf.WriteString(fmt.Sprintf("- %s: %s\n", name, getComment(f.Tag)))
			return true
		})
}

func (r *Route) RespBody() ([]api_type.Object, *api_type.Object) {
	field, ok := r.Resp.FieldByName("Data")
	if !ok {
		return nil, nil
	}

	req := api_type.Object{}

	if desc := getComment(field.Tag); desc != "" {
		req.Comment = desc
	}

	if field.Type.Kind() == reflect.Pointer {
		field.Type = field.Type.Elem()
	}
	ob := api_type.ObjectMap{}

	err := api_type.GetObjectMap(&ob,[]reflect.Type{field.Type}, "typescript", api_type.MemberTypeJson)
	if err != nil {
		return nil, nil
	}
	var resp api_type.Object

	name := field.Type.Name()

	for s, object := range ob {
		if s == name {
			object.Name = "Resp"
			object.Comment = "返回内容"
			object.JsonName = field.Tag.Get("json")
			resp = object
			//ob[object.Name] = object
			delete(ob, name)
		}
	}
	return ob.ToList(), &resp
}

type respBody struct {
	Code    string `json:"code"           c:"ok表示成功，其他表示错误代码"`
	Message string `json:"message"        c:"与code对应的描述信息"`
}

type respBodyWithData struct {
	respBody
	Data interface{} `json:"data" c:"返回的数据"`
}

func (r *Route) RespError(buf *bytes.Buffer) {
	field, _ := r.Resp.FieldByName("Error")

	if desc := getComment(field.Tag); desc != "" {
		buf.WriteString("\n## 错误码说明\n")
		buf.WriteString(desc + "\n\n")
	}
}
