package ts

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/lovego/goa/convert"
	"github.com/lovego/goa/docs/ts/ts/api_type"
	"github.com/lovego/strs"
	"github.com/lovego/struct_tag"
)

func (r *Route) Title() string {
	if f, ok := r.Req.FieldByName("Title"); ok {
		if comment := getComment(f.Tag); comment != "" {
			return comment
		}
		return strings.TrimSpace(string(f.Tag))
	}
	return ""
}

func (r *Route) Desc() string {
	if f, ok := r.Req.FieldByName("Desc"); ok {
		if comment := getComment(f.Tag); comment != "" {
			return comment + "\n\n"
		} else if desc := strings.TrimSpace(string(f.Tag)); desc != "" {
			return desc + "\n\n"
		}
	}
	return ""
}

func (r *Route) Param(fullPath string) (*api_type.Object, []string) {
	field, ok := r.Req.FieldByName("Param")
	if !ok {
		return nil, nil
	}

	paramReq := api_type.Object{
		Name:     "paramReq",
		JsonName: "paramReq",
		Comment:  getComment(field.Tag),
	}

	names := regexp.MustCompile(fullPath).SubexpNames()[1:] // names[0] is always "".
	for _, name := range names {
		if name != "" {
			if f, ok := field.Type.FieldByName(strs.FirstLetterToUpper(name)); ok {

				m := api_type.Member{
					Type:     "string",
					Name:     name,
					Comment:  getComment(f.Tag),
					JsonName: name,
					FormName: name,
					Options:  nil,
					NotMust:  true,
				}

				paramReq.Members = append(paramReq.Members, m)
			}
		}
	}

	return &paramReq, names
}

// 取中间字符串
func Between(str, starting, ending string) string {
	s := strings.Index(str, starting)
	if s < 0 {
		return ""
	}
	s += len(starting)
	e := strings.Index(str[s:], ending)
	if e < 0 {
		return ""
	}
	return str[s : s+e]
}
func getTsParamPath(fullPath string, names []string) string {
	if len(names) == 0 {
		return fullPath
	}
	for {
		if !strings.Contains(fullPath, "(") && !strings.Contains(fullPath, ")") {
			break
		}
		t := Between(fullPath, "(", ")")
		if t == "" {
			continue
		}
		for _, name := range names {
			if !strings.Contains(t, name) {
				continue
			}
			fullPath = strings.ReplaceAll(fullPath, "("+t+")", `${param.`+name+`}`)
		}

	}
	return fullPath
}

func (r *Route) Query() ([]api_type.Object, *api_type.Object) {
	field, ok := r.Req.FieldByName("Query")
	if !ok {
		return nil, nil
	}

	req := api_type.Object{}

	if desc := getComment(field.Tag); desc != "" {
		req.Comment = desc
	}

	ob, err := api_type.GetObjectMap([]reflect.Type{field.Type}, "typescript", api_type.MemberTypeJson)
	if err != nil {
		log.Panic(err)
	}
	var queryReq api_type.Object
	for s, object := range ob {
		if s == "" {
			object.Name = "QueryReq"
			object.Comment = "Query请求参数"
			object.JsonName = field.Tag.Get("json")
			queryReq = object
			//ob[object.Name] = object
			delete(ob, "")
		}
	}

	return ob.ToList(), &queryReq
}

func (r *Route) header(buf *bytes.Buffer) {
	field, ok := r.Req.FieldByName("Header")
	if !ok {
		return
	}

	buf.WriteString("\n## 请求头说明\n")
	if desc := getComment(field.Tag); desc != "" {
		buf.WriteString(desc + "\n\n")
	}
	convert.Traverse(reflect.New(field.Type).Elem(), true,
		func(_ reflect.Value, f reflect.StructField) bool {
			name, _ := struct_tag.Lookup(string(f.Tag), "header")
			if name == "" {
				name = f.Name
			}
			buf.WriteString(fmt.Sprintf("- %s (%v): %s\n", name, f.Type, getComment(f.Tag)))
			return true
		})
}

func (r *Route) Header() *api_type.Object {
	field, ok := r.Req.FieldByName("Header")
	if !ok {
		return nil
	}
	req := api_type.Object{}

	if desc := getComment(field.Tag); desc != "" {
		req.Comment = desc
	}
	ob, err := api_type.GetObjectMap([]reflect.Type{field.Type}, "typescript", api_type.MemberTypeJson)
	if err != nil {
		log.Panic(err)
	}
	var ReqHeader api_type.Object

	for s, object := range ob {
		if s == "" {
			object.Name = "Header"
			object.Comment = "Header说明"
			object.JsonName = field.Tag.Get("json")
			//ob[object.Name] = object
			ReqHeader = object
			delete(ob, "")
		}
	}
	return &ReqHeader
}

func (r *Route) Body() ([]api_type.Object, *api_type.Object) {
	field, ok := r.Req.FieldByName("Body")
	if !ok {
		return nil, nil
	}
	req := api_type.Object{}

	if desc := getComment(field.Tag); desc != "" {
		req.Comment = desc
	}
	ob, err := api_type.GetObjectMap([]reflect.Type{field.Type}, "typescript", api_type.MemberTypeJson)
	if err != nil {
		log.Panic(err)
	}
	var bodyReq api_type.Object

	for s, object := range ob {
		if s == "" {
			object.Name = "BodyReq"
			object.Comment = "Body请求参数"
			object.JsonName = field.Tag.Get("json")
			//ob[object.Name] = object
			bodyReq = object
			delete(ob, "")
		}
	}
	return ob.ToList(), &bodyReq
}
