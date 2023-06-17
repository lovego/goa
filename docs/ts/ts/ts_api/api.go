package ts_api

import (
	"gt/cmds/zero/zero/template"
	"gt/pkg/comment"
	"gt/pkg/parser_zero_api"
	"gt/pkg/stringx"
	"gt/pkg/tpl"
	"strings"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

type Api struct {
	FuncComment string `yaml:"funcComment"`
	FuncName    string `yaml:"funcName"`
	Router      string `yaml:"router"`
	Method      string `yaml:"method"`
	Group       string `yaml:"group"`
	Req         string `yaml:"req"`
	Resp        string `yaml:"resp"`
}

type ApiInfo struct {
	File      string `yaml:"File"`
	NameSpace string `yaml:"nameSpace"`
	Apis      []*Api `yaml:"apis"`
}

func GenApi(outPath, apiFile string) error {

	api, err := parser_zero_api.ParserZeroApi(apiFile)
	if err != nil {
		return err
	}

	var list []*ApiInfo

	for _, group := range api.Service.Groups {
		api, err := apiInfo(group)
		if err != nil {
			return err
		}
		list = append(list, api)
	}

	err = tpl.New(list, template.TsApi).Gen(outPath)
	if err != nil {
		return err
	}

	return nil
}

func apiInfo(group spec.Group) (*ApiInfo, error) {
	var ts ApiInfo

	groupPath := group.GetAnnotation("group")
	ts.NameSpace = getSpace(groupPath)
	ts.File = groupPath + ".ts"

	for _, route := range group.Routes {
		a := &Api{
			FuncComment: comment.CleanComment(strings.Join(route.HandlerComment, " ")),
			FuncName:    route.Handler,
			Router:      route.Path,
			Method:      route.Method,
			Group:       group.GetAnnotation("group"),
			Req:         route.RequestTypeName(),
			Resp:        route.ResponseTypeName(),
		}
		ts.Apis = append(ts.Apis, a)
	}

	return &ts, nil
}

func getSpace(space string) string {
	arr := strings.Split(space, "/")
	space = ""

	for _, s := range arr {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		space = space + stringx.ToUpperCamel(s)
	}
	return space
}
