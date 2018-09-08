package benchmark

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/lovego/goa"
)

type route struct {
	method string
	path   string
}

type goaRouterTestCase struct {
	router *goa.Router
	routes []route
	hits   int
}

type httpRouterTestCase struct {
	router *httprouter.Router
	routes []route
	hits   int
}

func loadGoaRouterTestCase(routes []route) *goaRouterTestCase {
	var router = goa.New()
	var tc = goaRouterTestCase{router: router, routes: routes}
	var paramRegexp = regexp.MustCompile(`:\w+`)
	for _, route := range routes {
		var path = route.path
		if strings.IndexByte(route.path, ':') > 0 {
			path = paramRegexp.ReplaceAllString(path, `(:\w+)`)
		} else {
			path = strings.Replace(path, `.`, `\.`, -1)
		}
		router.Add(route.method, path, func(*goa.Context) {
			tc.hits++
		})
	}
	fmt.Println(tc.router.Group.String())
	return &tc
}

func loadHttpRouterTestCase(routes []route) *httpRouterTestCase {
	var router = httprouter.New()
	var tc = httpRouterTestCase{router: router, routes: routes}
	for _, route := range routes {
		router.Handle(route.method, route.path, func(http.ResponseWriter, *http.Request, httprouter.Params) {
			tc.hits++
		})
	}
	return &tc
}

func runGoaRouterTestCase(b *testing.B, tc *goaRouterTestCase) {
	b.ReportAllocs()
	tc.hits = 0

	request, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		for _, route := range tc.routes {
			request.Method = route.method
			request.URL.Path = route.path
			tc.router.ServeHTTP(nil, request)
		}
	}
	if tc.hits != b.N*len(tc.routes) {
		b.Errorf("hits want: %d, got: %d\n", b.N*len(tc.routes), tc.hits)
	}
}

func runHttpRouterTestCase(b *testing.B, tc *httpRouterTestCase) {
	b.ReportAllocs()
	tc.hits = 0

	request, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		for _, route := range tc.routes {
			request.Method = route.method
			request.URL.Path = route.path
			tc.router.ServeHTTP(nil, request)
		}
	}
	if tc.hits != b.N*len(tc.routes) {
		b.Errorf("hits want: %d, got: %d\n", b.N*len(tc.routes), tc.hits)
	}
}
