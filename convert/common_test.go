package convert

import (
	"fmt"
	"log"
	"reflect"
)

func ExampleSetArray() {
	var v = struct {
		Slice []string
		Array [3]int
		Json  []struct {
			Id   int
			Name string
		}
	}{}

	if err := SetArray(reflect.ValueOf(&v.Slice), []string{"a", "bc", "d"}); err != nil {
		log.Panic(err)
	}
	if err := SetArray(reflect.ValueOf(&v.Array), []string{"1", "2"}); err != nil {
		log.Panic(err)
	}
	if err := SetArray(reflect.ValueOf(&v.Json), []string{
		`{"id": 3, "name": "f" }`, "{}",
	}); err != nil {
		log.Panic(err)
	}
	fmt.Println(v.Slice)
	fmt.Println(v.Array)
	fmt.Println(v.Json)
	// Output:
	// [a bc d]
	// [1 2 0]
	// [{3 f} {0 }]
}

type testUnmarshaler struct {
	V string
}

func (u *testUnmarshaler) UnmarshalJSON(b []byte) error {
	u.V = string(b)
	return nil
}

func ExampleSet() {
	var t struct {
		V string
	}
	fmt.Println(Set(reflect.ValueOf(&t), "XX"))
	fmt.Println(t)

	var t2 testUnmarshaler
	fmt.Println(Set(reflect.ValueOf(&t2), "XX"))
	fmt.Println(t2)

	// Output:
	// invalid character 'X' looking for beginning of value
	// {}
	// <nil>
	// {XX}
}
