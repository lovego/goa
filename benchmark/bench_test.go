package benchmark

import (
	"testing"
)

var goaRouterStatic = loadGoaRouterTestCase(staticRoutes)
var httpRouterStatic = loadHttpRouterTestCase(staticRoutes)

var goaRouterGithub = loadGoaRouterTestCase(githubAPI)
var httpRouterGithub = loadHttpRouterTestCase(githubAPI)

var goaRouterGooglePlus = loadGoaRouterTestCase(googlePlusAPI)
var httpRouterGooglePlus = loadHttpRouterTestCase(googlePlusAPI)

var goaRouterParseCom = loadGoaRouterTestCase(parseComAPI)
var httpRouterParseCom = loadHttpRouterTestCase(parseComAPI)

func BenchmarkGoaRouter_Static157(b *testing.B) {
	runGoaRouterTestCase(b, goaRouterStatic)
}

func BenchmarkHttpRouter_Static157(b *testing.B) {
	runHttpRouterTestCase(b, httpRouterStatic)
}

func BenchmarkGoaRouter_Github203(b *testing.B) {
	runGoaRouterTestCase(b, goaRouterGithub)
}

func BenchmarkHttpRouter_Github203(b *testing.B) {
	runHttpRouterTestCase(b, httpRouterGithub)
}

func BenchmarkGoaRouter_GooglePlus13(b *testing.B) {
	runGoaRouterTestCase(b, goaRouterGooglePlus)
}

func BenchmarkHttpRouter_GooglePlus13(b *testing.B) {
	runHttpRouterTestCase(b, httpRouterGooglePlus)
}

func BenchmarkGoaRouter_ParseCom26(b *testing.B) {
	runGoaRouterTestCase(b, goaRouterParseCom)
}

func BenchmarkHttpRouter_ParseCom26(b *testing.B) {
	runHttpRouterTestCase(b, httpRouterParseCom)
}
