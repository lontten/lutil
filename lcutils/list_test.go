package lcutils

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTool_HasAll(t *testing.T) {
	as := assert.New(t)
	check := ListTool(1, 2, 3).
		HasAll(1, 2).
		Check()
	as.True(check)

	check = ListTool(1, 2, 3).
		HasAll(1, 2, 3).
		Check()
	as.True(check)

	check = ListTool(1, 2, 3).
		HasAll(1, 2, 3, 4).
		Check()
	as.False(check)
}

func TestListTool_HasAny(t *testing.T) {
	as := assert.New(t)
	check := ListTool(1, 2, 3).
		HasAny(1, 2, 4).
		Check()
	as.True(check)

	check = ListTool(1, 2, 3).
		HasAny(1, 2, 3, 4).
		Check()
	as.True(check)

	check = ListTool(1, 2, 3).
		HasAny(3, 4).
		Check()
	as.True(check)

	check = ListTool(1, 2, 3).
		HasAny(4).
		Check()
	as.False(check)

}

func TestListTool_NotAll(t *testing.T) {
	as := assert.New(t)
	check := ListTool(1, 2, 3).
		NotAll(1, 2, 4).
		Check()
	as.False(check)

	check = ListTool(1, 2, 3).
		NotAll(1, 2, 3, 4).
		Check()
	as.False(check)

	check = ListTool(1, 2, 3).
		NotAll(3, 4).
		Check()
	as.False(check)

	check = ListTool(1, 2, 3).
		NotAll(4).
		Check()
	as.True(check)
}

type CheckUser struct {
	Name int
}

// 当一个结构体的所有字段都是可比较类型（comparable）时，这个结构体本身就是可比较的（comparable）
func TestListTool_other(t *testing.T) {
	as := assert.New(t)

	var u1 = CheckUser{Name: 1}
	var u2 = CheckUser{Name: 2}

	check := ListTool(u1, u2).
		HasAll(u1).
		Check()
	as.True(check)
}

func TestBoolEq(t *testing.T) {
	as := assert.New(t)
	as.True(BoolEq([]int{1, 2, 3}, []int{1, 2, 3}))
	as.False(BoolEq([]int{1, 2, 3}, []int{1, 2, 4}))
	as.False(BoolEq([]int{1, 2, 3}, []int{1, 2}))
	as.False(BoolEq([]int{1, 2, 3}, []int{1, 2, 3, 4}))
	as.True(BoolEq([]int{1, 2, 3}, []int{1, 2, 3, 3}))
	as.True(BoolEq([]int{1, 2, 3}, []int{3, 1, 2}))
}

func TestBoolDiff(t *testing.T) {
	as := assert.New(t)
	as.Equal([]int{3}, BoolDiff([]int{1, 2, 3}, []int{1, 2}))
	as.Equal([]int{1}, BoolDiff([]int{1, 2, 3}, []int{2, 3}))

	list := BoolDiff([]int{1, 2, 3}, []int{3})
	sort.Ints(list)
	as.Equal([]int{1, 2}, list)

	as.Equal([]int{}, BoolDiff([]int{}, []int{1, 2, 3}))
	as.Equal([]int{}, BoolDiff([]int{}, []int{}))
	as.Equal([]int{}, BoolDiff([]int{1, 2, 3}, []int{1, 2, 3}))
}
