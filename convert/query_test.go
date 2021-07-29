package convert

import (
	"fmt"
	"reflect"
)

func ExampleQuery() {
	var v struct {
		Id    int
		Name  string `json:"userName"`
		Desc  string `json:",omitempty"`
		Array []int  `json:"array,omitempty"`
		X     int
		Y     string `json:"-"`
	}
	fmt.Println(Query(reflect.ValueOf(&v), map[string][]string{}))
	fmt.Printf("%+v\n", v)

	fmt.Println(Query(reflect.ValueOf(&v), map[string][]string{
		"id":       {"9"},
		"userName": {"LiLei"},
		"Desc":     {"description"},
		"array[]":  {"1", "2", "3"},
	}))
	fmt.Printf("%+v\n", v)

	fmt.Println(Query(reflect.ValueOf(&v), map[string][]string{
		"x": {"a"},
	}))
	fmt.Println(v.X)

	// Output:
	// <nil>
	// {Id:0 Name: Desc: Array:[] X:0 Y:}
	// <nil>
	// {Id:9 Name:LiLei Desc:description Array:[1 2 3] X:0 Y:}
	// req.Query.X: strconv.ParseInt: parsing "a": invalid syntax
	// 0
}
