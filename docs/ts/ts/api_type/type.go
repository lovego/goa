package api_type

import (
	"encoding/json"
	"errors"
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

func (o *ObjectMap) Get(key string) (Object, bool) {
	s, ok := (*o)[key]
	return s, ok
}
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

func (o *ObjectMap) IsExistMember(name string) bool {
	for _, object := range *o {
		//if object.Name == name {
		//	return true
		//}
		for _, mem := range object.Members {
			if mem.Type == name {
				return true
			}
		}
	}

	return false
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

func GetTyp(typ reflect.Type) reflect.Type {
	switch typ.Kind() {
	case reflect.Struct:
		return typ
	case reflect.Ptr, reflect.Slice, reflect.Array, reflect.Map:
		return GetTyp(typ.Elem())
	default:
		return typ
	}
}

func GetObjectMap(ob *ObjectMap, types []reflect.Type, lang, memberType string) error {
	if ob == nil {
		return errors.New("obj必填")
	}

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

		members, specList, err := GetMembers(ob, typ, lang, memberType)
		if err != nil {
			return err
		}
		obj.Members = members
		(*ob)[obj.Name] = obj
		if len(specList) == 0 {
			continue
		}
		err = GetObjectMap(ob, specList, lang, memberType)
		if err != nil {
			return err
		}
		//if len(l) == 0 {
		//	continue
		//}
		//for k, v := range l {
		//	list[k] = v
		//}
	}

	return nil
}

const (
	MemberTypeJson    = "json"
	MemberTypeNonBody = "!json"
	MemberTypeForm    = "form"
	MemberTypeAll     = ""
)

// memberType json:带有json标签的body字段 !json:除了json标签的其他字段  form:带有form标签的表单字段 为空表示所有字段
func GetMembers(list *ObjectMap, tp reflect.Type, lang, memberType string) ([]Member, []reflect.Type, error) {
	definedType := tp.Kind()

	if definedType != reflect.Struct {

		if definedType == reflect.Ptr {
			return GetMembers(list, tp.Elem(), lang, memberType)
		}

		return nil, nil, fmt.Errorf("type %s not supported", tp.Name())
	}

	var specTypeList []reflect.Type
	var fields []Member

	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)

		if f.Tag.Get("json") == "-" || (f.Name[0] >= 'a' && f.Name[0] <= 'z') {
			continue
		}
		v := reflect.New(f.Type).Interface()

		var name string

		_, unmar := v.(json.Unmarshaler)
		_, mar := v.(json.Marshaler)
		if unmar || mar {
			name = "any"
		} else {

			t := GetTyp(f.Type)

			if t.Kind() == reflect.Struct && !(f.Anonymous && f.Tag.Get("json") == "") {
				name := t.Name()

				_, ok := list.Get(name)
				if !ok {
					//n := strings.ToLower(t.Name())

					if !(unmar || mar) && (!f.Anonymous || (f.Anonymous && f.Tag.Get(
						"json") != "")) {
						specTypeList = append(specTypeList, t)
					}
				}

			}
			if f.Anonymous && f.Tag.Get("json") == "" && t.Name() != "" {
				if f.Type.Kind() == reflect.Ptr {
					f.Type = f.Type.Elem()
				}
				mem, list, err := GetMembers(list, f.Type, lang, memberType)
				if err != nil {
					return nil, nil, err
				}
				fields = append(fields, mem...)
				specTypeList = append(specTypeList, list...)
				continue
			}

			tName, err := GenMemberType(f.Type, lang)
			if err != nil {
				return nil, nil, err
			}
			name = tName
		}

		if name == "" {
			name = f.Name
		}

		m := Member{
			Type:     name,
			Name:     f.Name,
			Comment:  getComment(f.Tag),
			JsonName: strings.ReplaceAll(f.Tag.Get("json"), ",omitempty", ""),
			FormName: strings.ReplaceAll(f.Tag.Get("form"), ",omitempty", ""),
		}
		m.Comment = cleanComment(m.Comment)

		if m.JsonName == "" {
			m.JsonName = m.Name
		}
		if m.FormName == "" {
			m.FormName = m.Name
		}
		// 如果member中有m中已经存在的,则将field中的删掉
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
