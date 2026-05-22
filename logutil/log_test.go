package logutil

import (
	"testing"
)

type User struct {
	Name string
	Age  int
}

func TestLog(t *testing.T) {
	var u = User{Name: "john", Age: 22}
	Log(&u)
}
