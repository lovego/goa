package api_type

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/lovego/struct_tag"
)

type Member struct {
	Type     string   `yaml:"type"`
	Name     string   `yaml:"name"`
	Comment  string   `yaml:"comment"`
	JsonName string   `yaml:"jsonName"`
	FormName string   `yaml:"formName"`
	Options  []string `yaml:"options"`
	NotMust  bool     `yaml:"notMust"`
}

type Object struct {
	Name     string   `yaml:"name"`
	JsonName string   `yaml:"jsonName"`
	Comment  string   `yaml:"comment"`
	Members  []Member `yaml:"fields"`
}

type ObjectMap map[string]Object

func (o *ObjectMap) Delete(key string) {
	delete(*o, key)
}
func (o *ObjectMap) Deletes(keys []string) {
	for _, key := range keys {
		o.Delete(key)
	}
}

func (o *ObjectMap) Len() int {
	return len(*o)
}
func (o *ObjectMap) ToList() []Object {
	list := make([]Object, 0, o.Len())

	for _, obj := range *o {
		list = append(list, obj)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list
}

func (o *ObjectMap) ToTypeScript() {
	for key, obj := range *o {
		for i, member := range obj.Members {
			(*o)[key].Members[i].Type = TsType(member.Type)
		}
	}
	*o = *o
}

func (o *ObjectMap) NoPrimMembers(lang string) []string {
	var list []string
	for _, object := range *o {
		for _, member := range object.Members {
			_, err := TsPrimitiveType(member.Type)
			if err != nil {
				list = append(list, member.Type)
			}
		}
	}

	return list
}

var whitespaceRegexp = regexp.MustCompile(`\s+`)

// extract comment from struct field tags
func getComment(tag reflect.StructTag) string {
	tagStr := string(tag)
	comment, _ := struct_tag.Lookup(tagStr, `comment`)
	if comment == `` {
		comment, _ = struct_tag.Lookup(tagStr, `c`)
	}
	if comment != `` {
		comment = strings.TrimSpace(comment)
	}
	if comment != `` {
		comment = whitespaceRegexp.ReplaceAllString(comment, " ")
	}
	return comment
}

func GetObjectMap(types []reflect.Type, lang, memberType string) (ObjectMap, error) {
	var list = ObjectMap{}

	for _, typ := range types {
		if typ == nil {
			continue
		}

		obj := Object{
			Name:    typ.Name(),
			Comment: typ.Name(),
			Members: nil,
		}
		obj.Comment = cleanComment(obj.Comment)

		members, specList, err := GetMembers(typ, lang, memberType)
		if err != nil {
			return nil, err
		}
		obj.Members = members
		list[obj.Name] = obj
		if len(specList) == 0 {
			continue
		}
		l, err := GetObjectMap(specList, lang, memberType)
		if err != nil {
			return nil, err
		}
		if len(l) == 0 {
			continue
		}
		for k, v := range l {
			list[k] = v
		}
	}

	return list, nil
}

const (
	MemberTypeJson    = "json"
	MemberTypeNonBody = "!json"
	MemberTypeForm    = "form"
	MemberTypeAll     = ""
)

// memberType json:带有json标签的body字段 !json:除了json标签的其他字段  form:带有form标签的表单字段 为空表示所有字段
func GetMembers(tp reflect.Type, lang, memberType string) ([]Member, []reflect.Type, error) {
	definedType := tp.Kind()

	if definedType != reflect.Struct {

		if definedType == reflect.Pointer {
			return GetMembers(tp.Elem(), lang, memberType)
		}

		return nil, nil, fmt.Errorf("type %s not supported", tp.Name())
	}

	var specTypeList []reflect.Type
	var fields []Member

	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)

		if f.Tag.Get("json") == "-" {
			continue
		}

		if f.Type.Kind() == reflect.Pointer && f.Type.Elem().Kind() == reflect.Struct {
			specTypeList = append(specTypeList, f.Type.Elem())
		}
		if f.Type.Kind() == reflect.Struct {
			if f.Anonymous && f.Tag.Get("json") != "" {
				specTypeList = append(specTypeList, f.Type)
			}
			if !f.Anonymous {
				specTypeList = append(specTypeList, f.Type)
			} else {
				mem, list, err := GetMembers(f.Type, lang, memberType)
				if err != nil {
					return nil, nil, err
				}
				fields = append(fields, mem...)
				specTypeList = append(specTypeList, list...)
				continue
			}
		}

		t, err := GenMemberType(f.Type, lang)
		if err != nil {
			return nil, nil, err
		}

		m := Member{
			Type:     t,
			Name:     f.Name,
			Comment:  getComment(f.Tag),
			JsonName: f.Tag.Get("json"),
			FormName: f.Tag.Get("form"),
		}
		m.Comment = cleanComment(m.Comment)

		if m.JsonName == "" {
			m.JsonName = m.Name
		}
		if m.FormName == "" {
			m.FormName = m.Name
		}
		fields = append(fields, m)

	}

	return fields, specTypeList, nil
}

func GenMemberType(m reflect.Type, lang string) (string, error) {
	switch lang {
	case "typescript":
		t, err := genTsMemberType(m)
		if err != nil {
			return "", err
		}
		return t, nil
	default:
		return m.Name(), nil

	}
}

type Types []reflect.Type

func (a *Types) GetType(typeName string) reflect.Type {
	for _, s := range *a {
		if s.Name() == typeName {
			return s
		}
	}
	return nil
}

func (a *Types) GetTypes(typeNameList []string) Types {
	var list Types
	for _, s := range typeNameList {
		t := a.GetType(s)
		if t == nil {
			continue
		}
		list = append(list, t)
	}
	return list

}
