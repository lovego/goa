package api_type

import (
	"strings"
)

const (
	formTagKey   = "form"
	pathTagKey   = "path"
	headerTagKey = "header"
)

func cleanComment(str string) string {
	str = strings.TrimSpace(str)
	str = strings.TrimLeft(str, "//")
	str = strings.TrimSpace(str)
	return str
}
