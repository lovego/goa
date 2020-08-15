package goa

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
)

func ExampleContext_Param() {
	c := &Context{params: []string{"123", "sdf"}}
	fmt.Println(c.Param(0), "-")
	fmt.Println(c.Param(1), "-")
	fmt.Println(c.Param(2), "-")
	// Output:
	// 123 -
	// sdf -
	//  -
}

func ExampleContext_Next() {
	c := &Context{}
	c.Next()
	c.Next()

	// Output:
}

func ExampleContext_Context() {
	c := &Context{Request: httptest.NewRequest("GET", "/", nil)}
	fmt.Println(c.Context())
	ctx := context.WithValue(context.Background(), "custom", 1)
	c.Set("context", ctx)
	fmt.Println(c.Context() == ctx)
	// Output:
	// context.Background
	// true
}

func ExampleContext_Get() {
	c := &Context{}
	fmt.Println(c.Get("a"))
	c.Set("a", "韩梅梅")
	fmt.Println(c.Get("a"))
	// Output:
	// <nil>
	// 韩梅梅
}

func ExampleContext_Set() {
	c := &Context{}
	c.Set("a", "韩梅梅")
	fmt.Println(c.Get("a"))
	// Output:
	// 韩梅梅
}

func ExampleContext_GetError_SetError() {
	c := &Context{}
	fmt.Println(c.GetError())
	c.SetError(errors.New("the-error"))
	fmt.Println(c.GetError())
	// Output:
	// <nil>
	// the-error
}
