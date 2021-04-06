package utilroutes

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestDebugIndex(t *testing.T) {
	if err := ioutil.WriteFile("t.html", debugIndex(), 0660); err != nil {
		panic(err)
	}
}

func TestInstanceName(t *testing.T) {
	fmt.Println(instanceName)
}
