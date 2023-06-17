package ts_api_Type

import (
	"gt/pkg/comment"
	"gt/pkg/stringx"
	"gt/pkg/tpl"
	"strings"

	"github.com/lovego/goa/docs/ts/ts/api_type"
)

type ApiTypeInfo struct {
	File           string            `yaml:"File"`
	NameSpace      string            `yaml:"nameSpace"`
	CommonTypeList []api_type.Object `yaml:"commonTypeList"`
}
//type ApiType struct {
//	FuncComment string            `yaml:"funcComment"`
//	FuncName    string            `yaml:"funcName"`
//	Router      string            `yaml:"router"`
//	Method      string            `yaml:"method"`
//	Req         string            `yaml:"req"`
//	ReqFields   []api_type.Member `yaml:"reqFields"`
//	Resp        string            `yaml:"resp"`
//	RespFields  []api_type.Member `yaml:"respFields"`
//}

func GenApiType(outPath, apiFile string) error {
	api, err := parser_zero_api.ParserZeroApi(apiFile)
	if err != nil {
		return err
	}

	var list []*ApiTypeInfo
	for _, group := range api.Service.Groups {
		api, err := apiTypeInfo(group, api.Types)
		if err != nil {
			return err
		}
		list = append(list, api)
	}
	err = tpl.New(list, template.TsApiType).Gen(outPath)
	if err != nil {
		return err
	}

	return nil
}

func apiTypeInfo(group spec.Group, types api_type.Types) (*ApiTypeInfo, error) {

	//fileName := group.GetAnnotation("prefix")
	groupPath := group.GetAnnotation("group")

	typeInfo := &ApiTypeInfo{
		File:           group.GetAnnotation("group") + ".typings.d.ts",
		NameSpace:      getSpace(groupPath),
		CommonTypeList: nil,
	}

	var typeList []spec.Type
	for _, route := range group.Routes {
		typeList = append(typeList, types.GetType(route.RequestTypeName()))
		typeList = append(typeList, types.GetType(route.ResponseTypeName()))
	}

	ls, err := types.GetCustomMemberType(typeList)
	if err != nil {
		return nil, err
	}
	typeList = append(typeList, ls...)

	objectMap, err := api_type.GetObjectMap(typeList, "typescript", api_type.MemberTypeJson)
	if err != nil {
		return nil, err
	}

	typeInfo.CommonTypeList = objectMap.ToList()

	for i, _ := range typeInfo.CommonTypeList {
		for j, _ := range typeInfo.CommonTypeList[i].Members {
			if typeInfo.CommonTypeList[i].Members[j].NotMust {
				typeInfo.CommonTypeList[i].Members[j].JsonName += "?"
				typeInfo.CommonTypeList[i].Members[j].Comment = "选填 " + typeInfo.CommonTypeList[i].Members[j].Comment
				typeInfo.CommonTypeList[i].Members[j].Comment = comment.CleanComment(typeInfo.CommonTypeList[i].Members[j].Comment)
			}
		}
	}

	return typeInfo, nil
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
