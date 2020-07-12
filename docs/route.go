package docs

import (
	"bytes"
	"reflect"
	"regexp"
	"strings"

	"github.com/lovego/struct_tag"
)

type Route struct {
	req, resp reflect.Type
}

func (r *Route) Parse(handler interface{}) bool {
	typ := reflect.TypeOf(handler)
	if typ.NumIn() != 2 {
		return false
	}
	r.req, r.resp = typ.In(0), typ.In(1).Elem()
	return true
}

func (r *Route) Doc(method, fullPath string) []byte {
	buf := bytes.NewBufferString("# " + r.Title(method, fullPath) + "\n")

	r.Param(buf, fullPath)
	r.Query(buf)
	r.Header(buf)
	r.Body(buf)

	r.ResHeader(buf)
	r.ResBody(buf)
	r.ResError(buf)

	return buf.Bytes()
}

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

var whitespaceRegexp = regexp.MustCompile(`\s+`)
