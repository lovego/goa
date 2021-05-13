package docs_test

import (
	"path/filepath"
	"time"

	"github.com/lovego/fs"
	"github.com/lovego/goa"
)

func ExampleGroup() {
	router := goa.New()
	router.DocDir(filepath.Join(fs.SourceDir(), "testdata"))

	accounts := router.Group("/", "账号", "用户、公司、员工、角色、权限")
	accounts.Group("/users", "用户").
		Get(`/`, testHandler).
		Get(`/(?P<type>\w+)/(?P<id>\d+)`, testHandler2)

	accounts.Group("/companies", "公司")

	router.Group("/goods", "商品")
	router.Group("/bill", "单据", "采购、销售")
	router.Group("/storage", "库存")

	// Output:
}

type T struct {
	Type string `c:"类型"`
	Id   *int    `c:"ID"`
	Flag bool   `json:"-" c:"标志"`
}

func testHandler(req struct {
	Title string `用户列表`
	Desc  string `列出所有的用户`
	Query struct {
		Page int `c:"页码"`
		T
	}
	Header struct {
		Cookie string `c:"Cookie中包含会话信息"`
	}
	Body *struct {
		Name *string `c:"名称"`
		T
	}
	Session struct {
		UserId  int64
		LoginAt time.Time
	}
	Ctx *goa.Context
}, resp *struct {
	Data []*struct {
		Id   *int    `c:"ID"`
		Name *string `c:"名称"`
	}
	Error error
}) {
}

func testHandler2(req struct {
	Title string `用户详情`
	Desc  string `获取用户的详细信息`
	Param T
}, resp *struct {
	Header struct {
		SetCookie string `header:"Set-Cookie" c:"返回登录会话"`
	}
	Data struct {
		Id   int    `c:"ID"`
		Name string `c:"名称"`
	}
	Error error
}) {
}
