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
	var s = ToJsonStr(list)
	fmt.Println(s)
	objList := ToObjList[string](s)
	for i, o := range objList {
		fmt.Println(i, o)
	}
}
