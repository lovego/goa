package docs_test

import (
	"path/filepath"

	"github.com/lovego/fs"
	"github.com/lovego/goa"
)

func ExampleGroup() {
	router := goa.New()
	router.DocDir(filepath.Join(fs.SourceDir(), "testdata"))

	router.Group("/accounts", "账号", "用户、公司、员工、角色、权限")
	router.Group("/goods", "商品")
	router.Group("/bill", "单据", "采购、销售")
	router.Group("/storage", "库存")

	// Output:
}
