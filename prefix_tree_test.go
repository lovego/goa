package router

import (
	"regexp"
	"strings"
	"testing"
)

func ExampleNewNode() {
}

func BenchmarkStringHasPrefix(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		if !strings.HasPrefix("/company-skus/search", "/company-skus") {
			b.Error("not matched")
		}
	}
}

var testRegexp = regexp.MustCompile("^/company-skus")

func BenchmarkRegexpMatch(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		if len(testRegexp.FindStringSubmatch("/company-skus/search")) == 0 {
			b.Error("not matched")
		}
	}
}
