package jet

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

// 大写 65-90
// 小写 97-122
// 是否是大写字母
func IsUpChar(c uint8) bool {
	return c >= 'A' && c <= 'Z'
}

// 是否是小写字母
func IsLowChar(c uint8) bool {
	return c >= 'a' && c <= 'z'
}

// 是否是小写字母
func IsNumChar(c uint8) bool {
	return c >= '0' && c <= '9'
}

func UpChar(c uint8) uint8 {
	if IsLowChar(c) {
		return c - 'a' + 'A'
	}
	return c
}
func LowChar(c uint8) uint8 {
	if IsUpChar(c) {
		return c - 'A' + 'a'
	}
	return c
}

// ToLowerCamel 转为小驼峰格式
func ToLowerCamel(src string) string {
	return strcase.ToLowerCamel(src)
}

// ToUpperCamel 转为大驼峰格式
func ToUpperCamel(src string) string {
	return strcase.ToCamel(src)
}

// 大写字母转为下划线格式
func ToLine(src string) string {

	var ss []byte
	str := []byte(src)

	for i, _ := range str {
		if IsUpChar(str[i]) {
			//第一个字母不加下划线,并且上一个字母不是大写则往前添加下划线,(解决连续大写的问题)
			if i > 0 && IsLowChar(src[i-1]) {
				ss = append(ss, '_')
			}
			str[i] = LowChar(str[i])
		}
		ss = append(ss, str[i])
	}
	return string(ss)
}

// 大写字母转为下划线格式 中间
func ToMiddleLine(src string) string {

	var ss []byte
	str := []byte(src)

	for i, _ := range str {
		if IsUpChar(str[i]) {
			//第一个字母不加下划线,并且上一个字母不是大写则往前添加下划线,(解决连续大写的问题)
			if i > 0 && IsLowChar(src[i-1]) {
				ss = append(ss, '_')
			}
			str[i] = LowChar(str[i])
		}
		ss = append(ss, str[i])
	}

	return strings.ReplaceAll(string(ss), "_", "-")
}

// 大写 65-90
// 小写 97-122
func FistUp(str string) string {
	ss := []byte(str)
	ss[0] = UpChar(ss[0])
	return string(ss)
}
func FistLower(str string) string {
	ss := []byte(str)
	ss[0] = LowChar(ss[0])
	return string(ss)
}

// go 1.13支持
func IsZero(obj interface{}) bool {
	if obj == nil {
		return true
	}

	object := reflect.ValueOf(obj)
	if object.IsZero() {
		return true
	}
	return false
}

func Trim(str string) string {
	return strings.TrimSpace(str)
}

func LongestCommonPrefix(strs []string) string {
	switch len(strs) {
	case 0:
		return ""
	case 1:
		return strs[0]
	}
	isFinish := true
	i := 0
	for ; isFinish; i++ {
		for j := 1; j < len(strs); j++ {
			if i > len(strs[j])-1 || i > len(strs[j-1])-1 || strs[j-1][i] != strs[j][i] {
				isFinish = false
			}
		}
	}
	return strs[0][:i-1]
}

func IsNumber(num string) bool {
	for _, i2 := range num {
		if i2 == '.' {
			continue
		}
		if i2 < '0' || i2 > '9' {
			return false
		}
	}

	return true
}

// 取中间字符串
func Between(str, starting, ending string) string {
	s := strings.Index(str, starting)
	if s < 0 {
		return ""
	}
	s += len(starting)
	e := strings.Index(str[s:], ending)
	if e < 0 {
		return ""
	}
	return str[s : s+e]
}

func DeleteSpace(str string) string {
	var dest []int32
	var i int32
	for _, value := range str {
		if unicode.IsSpace(value) == true {
			continue
		}

		dest = append(dest, value)
		i++
	}

	return string(dest)
}
