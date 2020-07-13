package docs

import (
	"bytes"
	"html"
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
	buf := bytes.NewBufferString(
		"# " + r.Title() + "<br>" + r.MethodPath(method, fullPath) + "\n",
	)
	r.Desc(buf)

	r.Param(buf, fullPath)
	r.Query(buf)
	r.Header(buf)
	r.Body(buf)

	r.RespHeader(buf)
	r.RespBody(buf)
	r.RespError(buf)

	return buf.Bytes()
}

func (r *Route) MethodPath(method, fullPath string) string {
	return method + " " + html.EscapeString(fullPath)
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
