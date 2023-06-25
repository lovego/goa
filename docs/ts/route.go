package ts

import (
	"fmt"
	"html"
	"reflect"
	"regexp"
	"strings"

	"github.com/lovego/goa/docs/ts/ts_tpl"
	"github.com/lovego/struct_tag"
)

type Route struct {
	Req, Resp reflect.Type
}

func (r *Route) Parse(handler interface{}) bool {
	typ := reflect.TypeOf(handler)
	if typ.NumIn() != 2 {
		return false
	}
	r.Req, r.Resp = typ.In(0), typ.In(1).Elem()
	return true
}

func (r *Route) TypeScriptSdk(method, fullPath, tsFile string) error {

	param, names := r.Param(fullPath)
	if param != nil {
		return nil
	}
	commQuery, reqQuery := r.Query()
	fmt.Println(reqQuery)

	reqHeader := r.Header()

	commBody, reqBody := r.Body()
	fmt.Println(reqBody)

	//r.RespHeader(buf)
	commResp, resp := r.RespBody()
	fmt.Println(resp)
	//r.RespError(buf)

	comm := append(commQuery, commBody...)
	comm = append(comm, commResp...)

	fmt.Println(fullPath)
	fmt.Println("tsfile:", tsFile)

	api := ts_tpl.ApiInfo{
		File:         tsFile,
		Title:        strings.TrimSpace(r.Title()),
		Desc:         strings.TrimSpace(r.Desc()),
		Method:       method,
		Router:       getTsParamPath(fullPath, names),
		RawRouter:    r.MethodPath(method, fullPath),
		TypeList:     comm,
		Header:       reqHeader,
		Query:        reqQuery,
		Body:         reqBody,
		Param:        param,
		Resp:         resp,
		FunctionName: "test",
	}

	err := api.Run()
	if err != nil {
		return err
	}
	return nil
}

func (r *Route) MethodPath(method, fullPath string) string {
	return method + " " + html.EscapeString(fullPath)
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
