package utilroutes

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestDebugIndex(t *testing.T) {
	if err := ioutil.WriteFile("t.html", debugIndex(), 0660); err != nil {
		log.Panic(err)
	}
}

func TestInstanceName(t *testing.T) {
	fmt.Println(instanceName)
}
