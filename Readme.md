# GOA web framework
一个功能强大的web框架

### 特性
- 正则路由
- 路由分组
- 支持中间件

### 默认中间件
- 基于cookie的session会话
- 记录pending_request
- 强大的日志记录与报警功能

### 注意事项
当静态路由与动态路由冲突时，静态路由会优先被查找。

当路由Use中间件时，确保所有中间件handler在业务handler前面。

编写中间件时，如果不调用ctx.Next()，则表示不执行接下来的所有handlers。

preRequest部分代码请放在ctx.Next()前，afterResponse部分代码放在ctx.Next()后面。

### 代码示例

```
	router := New()
	router.Use(func(ctx *Context) {
		fmt.Println("middleware 1 pre")
		ctx.Next()
		fmt.Println("middleware 1 post")
	})
	router.Use(func(ctx *Context) {
		fmt.Println("middleware 2 pre")
		ctx.Next()
		fmt.Println("middleware 2 post")
	})
	router.Get("/", func(ctx *Context){
	    fmt.Println("you got it")
    })
	request, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		panic(err)
	}
	router.ServeHTTP(nil, request)

```