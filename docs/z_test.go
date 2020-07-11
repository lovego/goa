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
	accounts.Group("/users", "用户")
	accounts.Group("/companies", "公司")

	router.Group("/goods", "商品")
	router.Group("/bill", "单据", "采购、销售")
	router.Group("/storage", "库存")

	// Output:
}
