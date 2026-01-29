package logutil

import (
	"fmt"
	"testing"
)

type User struct {
	Name string
	Age  int
}

func TestLog(t *testing.T) {
	var u = User{Name: "john", Age: 22}
	fmt.Printf("%v", &u)
}
