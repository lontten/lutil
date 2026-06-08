package jsonutil

import (
	"fmt"
	"testing"
)

func TestToObjList(t *testing.T) {
	var list = make([]string, 0)
	list = append(list, "a")
	list = append(list, "b")
	list = append(list, "c")
	var s = ToJsonStrPanic("a")
	fmt.Println(s)
	objList := ToObjDefault[[]string](s)
	for i, o := range objList {
		fmt.Println(i, o)
	}
	fmt.Println(objList)
}
