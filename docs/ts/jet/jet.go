package jet

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"github.com/CloudyKit/jet/v6"
)

func RawWriter(w io.Writer, b []byte) {
	w.Write(b)
}

// 生成默认的模板处理函数
func NewFuncMap(f template.FuncMap) template.FuncMap {
	m := template.FuncMap{
		"FistUp":       FistUp,
		"FistLower":    FistLower,
		"ToUpperCamel": ToUpperCamel,
		"ToLowerCamel": ToLowerCamel,
		"ToLine":       ToLine,
		"ToMiddleLine": ToMiddleLine,
		"Lower":        strings.ToLower,
		"Up":           strings.ToUpper,
		"IsZero":       IsZero,
		"Trim":         strings.TrimSpace,
		//"DbType2GoType":     type_util.DbType2GoType,
		//"DbType2GoZeroType": type_util.DbType2GoZeroType,
	}
	if f == nil || len(f) == 0 {
		return m
	}

	for key, value := range f {
		m[key] = value
	}

	return m
}

func Tpl(tpl []byte, data interface{}) (*bytes.Buffer, error) {
	l := jet.NewInMemLoader()

	l.Set("tmp", string(tpl))

	views := jet.NewSet(l, jet.WithSafeWriter(RawWriter), jet.WithDelims("≤", "≥"))
	funcM := NewFuncMap(nil)
	for s, fun := range funcM {
		views.AddGlobal(s, fun)
	}

	t, err := views.GetTemplate("tmp")
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	vars := make(jet.VarMap)

	err = t.Execute(buf, vars, data)
	if err != nil {
		return nil, err
	}
	return buf, nil

}
