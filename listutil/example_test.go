package listutil

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TreeDemo struct {
	Id       int
	Pid      int
	Name     string
	Sort     int
	Children []*TreeDemo
}

func TestExample(t *testing.T) {
	list := make([]*TreeDemo, 0)
	list = append(list, &TreeDemo{
		Id:   1,
		Pid:  0,
		Name: "1",
		Sort: 30,
	})
	list = append(list, &TreeDemo{
		Id:   2,
		Pid:  0,
		Name: "2",
		Sort: 20,
	})
	list = append(list, &TreeDemo{
		Id:   3,
		Pid:  0,
		Name: "3",
		Sort: 10,
	})
	list = append(list, &TreeDemo{
		Id:   4,
		Pid:  1,
		Name: "1-4",
		Sort: 1,
	})

	treeTool := NewTreeBuilder[TreeDemo]().
		IsRootNode(func(a TreeDemo) bool {
			return a.Pid == 0
		}).
		IsParentNode(func(a, b TreeDemo) bool {
			return a.Id == b.Pid
		}).
		SortFun(func(a, b TreeDemo) int {
			return a.Sort - b.Sort
		})
	tree := treeTool.ToTree(list)

	// 输出结果
	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(jsonData))
}

func TestExample2(t *testing.T) {
	list := make([]*TreeDemo, 0)
	list = append(list, &TreeDemo{
		Id:   1,
		Pid:  5,
		Name: "1",
		Sort: 30,
	})
	list = append(list, &TreeDemo{
		Id:   2,
		Pid:  5,
		Name: "2",
		Sort: 20,
	})
	list = append(list, &TreeDemo{
		Id:   3,
		Pid:  5,
		Name: "3",
		Sort: 10,
	})
	list = append(list, &TreeDemo{
		Id:   4,
		Pid:  1,
		Name: "1-4",
		Sort: 1,
	})

	treeTool := NewTreeBuilder[TreeDemo]().
		IsParentNode(func(a, b TreeDemo) bool {
			return a.Id == b.Pid
		}).
		SortFun(func(a, b TreeDemo) int {
			return a.Sort - b.Sort
		})
	tree := treeTool.ToTree(list)

	// 输出结果
	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(jsonData))
}

func TestExample3(t *testing.T) {
	list := make([]*TreeDemo, 0)
	list = append(list, &TreeDemo{
		Id:   1,
		Pid:  5,
		Name: "1",
		Sort: 30,
	})
	list = append(list, &TreeDemo{
		Id:   2,
		Pid:  5,
		Name: "2",
		Sort: 20,
	})
	list = append(list, &TreeDemo{
		Id:   3,
		Pid:  5,
		Name: "3",
		Sort: 10,
	})
	list = append(list, &TreeDemo{
		Id:   4,
		Pid:  1,
		Name: "1-4",
		Sort: 1,
	})

	treeTool := NewTreeBuilder[TreeDemo]().
		IsParentNode(func(a, b TreeDemo) bool {
			return a.Id == b.Pid
		}).
		SortFun(func(a, b TreeDemo) int {
			return a.Sort - b.Sort
		}).
		SetRootNode(TreeDemo{
			Id:       100,
			Pid:      -1,
			Name:     "",
			Sort:     100,
			Children: nil,
		})
	tree := treeTool.ToTree(list)

	// 输出结果
	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(jsonData))
}
