package docs_test

import (
	"path/filepath"

	"github.com/lovego/fs"
	"github.com/lovego/goa"
)

func ExampleGroup() {
	router := goa.New()
	router.DocDir(filepath.Join(fs.SourceDir(), "testdata"))

	accounts := router.Group("/", "账号", "用户、公司、员工、角色、权限")
	accounts.Group("/users", "用户").
		Get(`/`, testHandler).
		Get(`/(?P<type>\w+)/(?P<id>\d+)`, testHandler)

	accounts.Group("/companies", "公司")

	router.Group("/goods", "商品")
	router.Group("/bill", "单据", "采购、销售")
	router.Group("/storage", "库存")

	// Output:
}

type T struct {
	Type string
	Id   int
	Flag bool
}

func testHandler(req struct {
	Title string `用户列表`
	Query struct {
		Page int
		T
	}
	Header struct {
		Cookie string
	}
	Body struct {
		Name string
		T
	}
}, resp *struct {
	Error  error
	Data   interface{}
	Header struct {
		SetCookie string `header:"Set-Cookie"`
	}
}) {
}

func testHandler2(req struct {
	Title string `用户详情`
	Param T
	Query struct {
		Page int
		T
	}
	Header struct {
		Cookie string
	}
	Body struct {
		Name string
		T
	}
}, resp *struct {
	Error  error
	Data   interface{}
	Header struct {
		SetCookie string `header:"Set-Cookie"`
	}
}) {
}
