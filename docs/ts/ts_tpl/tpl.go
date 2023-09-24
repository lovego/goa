package ts_tpl

import (
	"bytes"
	_ "embed"
	"os"
	"strings"

	"github.com/lovego/goa/docs/ts/jet"
	"github.com/lovego/goa/docs/ts/ts/api_type"
)

//go:embed req.pro.ts
var RequestTpl string

type ApiInfo struct {
	File         string            `yaml:"File" json:"File"`
	Title        string            `yaml:"title" json:"title" c:"接口标题"`
	Desc         string            `yaml:"desc" json:"desc" c:"接口描述"`
	Method       string            `yaml:"method" json:"method"`
	Router       string            `yaml:"router" json:"router"`
	TypeList     []api_type.Object `yaml:"typeList" json:"typeList"`
	Header       *api_type.Object  `yaml:"header" json:"header" c:"请求头说明"`
	Query        *api_type.Object  `yaml:"query" json:"query"`
	Body         *api_type.Object  `yaml:"body" json:"body"`
	Param        *api_type.Object  `yaml:"param" json:"param"`
	Resp         *api_type.Object  `yaml:"resp" json:"resp"`
	FunctionName string            `yaml:"functionName" json:"functionName" c:"函数名"`
	RawRouter    string            `yaml:"rawRouter" json:"rawRouter" c:"原始路径"`
}

func (a ApiInfo) Run() error {
	buf, err := jet.Tpl([]byte(RequestTpl), a)
	if err != nil {
		return err
	}

	//  删除连续空行
	buf = deleteSpaceLine(buf.String())

	err = os.WriteFile(a.File, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func deleteSpaceLine(s string) *bytes.Buffer {
	lines := strings.Split(s, "\n")

	buf := new(bytes.Buffer)

	var isSpaceLine bool

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			if isSpaceLine {
				continue
			}
			isSpaceLine = true
			buf.WriteString("\n")
			continue
		}
		isSpaceLine = false

		buf.WriteString(line)
		buf.WriteString("\n")
	}
	return buf
}
