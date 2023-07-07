package ts

import (
	"bytes"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/lovego/goa/convert"
	"github.com/lovego/goa/docs/ts/ts/api_type"
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

	if field.Type.Kind() == reflect.Ptr {
		field.Type = field.Type.Elem()
	}

	names := regexp.MustCompile(fullPath).SubexpNames()[1:] // names[0] is always "".
	var arr []string
	for _, name := range names {
		if strings.TrimSpace(name) == "" {
			continue
		}
		arr = append(arr, strings.TrimSpace(name))
	}

	names = arr

	if len(names) == 0 {
		return nil, nil
	}

	for _, name := range names {
		if name == "" {
			continue
		}

		for i := 0; i < field.Type.NumField(); i++ {
			f := field.Type.Field(i)
			if strings.TrimSpace(strings.Split(f.Tag.Get("json"), ",")[0]) == name {
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

// sdfsdf/sdfsdf/($1)/($2)/sdf

func getTsParamPath(fullPath string, names []string) string {
	if len(names) == 0 {
		return fullPath
	}
	var p string

	paths := strings.Split(fullPath, "/(")
	for _, path := range paths {

		pt := path
		n := strings.LastIndex(path, ")") + 1
		if n > 0 {
			path = path[:n]
			pt = pt[n:]
		}

		if strings.HasSuffix(path, ")") {
			path = `${param.` + getNamesValue(names, path) + `}` + pt
		}
		p += "/" + path
	}
	p = "/" + strings.TrimLeft(p, "/")

	//fmt.Println(p)

	return p
}

func getNamesValue(names []string, path string) string {
	for _, s := range names {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if strings.Contains(path, s) {
			return s
		}
	}
	return path
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

	if field.Type.Kind() == reflect.Ptr {
		field.Type = field.Type.Elem()
	}

	name := field.Type.Name()

	ob := api_type.ObjectMap{}

	err := api_type.GetObjectMap(&ob, []reflect.Type{field.Type}, "typescript", api_type.MemberTypeJson)
	if err != nil {
		return nil, nil
	}
	var queryReq api_type.Object
	for s, object := range ob {
		if s == name {
			object.Name = "QueryReq"
			object.Comment = "Query请求参数"
			object.JsonName = field.Tag.Get("json")
			queryReq = object
			//ob[object.Name] = object
			if !ob.IsExistMember(name) {
				delete(ob, name)
			}
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

	if field.Type.Kind() == reflect.Ptr {
		field.Type = field.Type.Elem()
	}

	ob := api_type.ObjectMap{}

	err := api_type.GetObjectMap(&ob, []reflect.Type{field.Type}, "typescript", api_type.MemberTypeJson)
	if err != nil {
		return nil
	}
	var ReqHeader api_type.Object

	name := field.Type.Name()

	for s, object := range ob {
		if s == name {
			object.Name = "Header"
			object.Comment = "Header说明"
			object.JsonName = field.Tag.Get("json")
			//ob[object.Name] = object
			ReqHeader = object
			if !ob.IsExistMember(name) {
				delete(ob, name)
			}
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

	if field.Type.Kind() == reflect.Ptr {
		field.Type = field.Type.Elem()
	}
	ob := api_type.ObjectMap{}

	err := api_type.GetObjectMap(&ob, []reflect.Type{field.Type}, "typescript", api_type.MemberTypeJson)
	if err != nil {
		return nil, nil
	}
	var bodyReq api_type.Object
	name := field.Type.Name()

	for s, object := range ob {
		if s == name {
			object.Name = "BodyReq"
			object.Comment = "Body请求参数"
			object.JsonName = field.Tag.Get("json")
			//ob[object.Name] = object
			bodyReq = object
			if !ob.IsExistMember(name) {
				delete(ob, name)
			}
		}
	}
	return ob.ToList(), &bodyReq
}
