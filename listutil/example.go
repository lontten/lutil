package listutil

import (
	"encoding/json"
	"fmt"
)

// 示例结构体
type Menu struct {
	ID       int     `json:"id" tree:"id"`
	ParentID int     `json:"parent_id" tree:"pid"`
	Name     string  `json:"name"`
	Order    float64 `json:"order" tree:"sort"`
	Children []*Menu `json:"children,omitempty" tree:"children"`
}

// 使用示例
func Example() {
	fmt.Println("=== 示例1: 基本菜单树 ===")
	exampleMenu()

	fmt.Println("\n=== 示例2: 部门树（限制层级）===")
	exampleDepartment()

	fmt.Println("\n=== 示例3: 带顶级节点的树 ===")
	exampleWithRoot()

}

func exampleMenu() {
	// 示例数据
	menus := []*Menu{
		{ID: 2, ParentID: 1, Name: "用户管理", Order: 2},
		{ID: 3, ParentID: 1, Name: "角色管理", Order: 1},
		{ID: 4, ParentID: 0, Name: "业务管理", Order: 2},
		{ID: 5, ParentID: 4, Name: "订单管理", Order: 1},
		{ID: 6, ParentID: 5, Name: "订单列表", Order: 1},
		{ID: 7, ParentID: 5, Name: "订单详情", Order: 2},
		{ID: 1, ParentID: 0, Name: "系统管理", Order: 1},
	}

	// 使用链式构造器
	treeTool := NewTreeBuilder[Menu]().
		IsRootNode(func(a Menu) bool {
			return a.ParentID == 0
		}).
		IsParentNode(func(a, b Menu) bool {
			return a.ID == b.ParentID
		}).
		SortFun(func(a, b Menu) int {
			return int(a.Order - b.Order)
		})
	tree := treeTool.ToTree(menus)

	// 输出结果
	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(jsonData))
}

func exampleDepartment() {
	// Department 结构体
	type Department struct {
		DeptID   string        `json:"dept_id" tree:"id"`
		ParentID string        `json:"parent_id" tree:"pid"`
		Name     string        `json:"name"`
		Sort     int           `json:"sort" tree:"sort"`
		Children []*Department `json:"children,omitempty" tree:"children"`
	}

	// 示例数据
	depts := []*Department{
		{DeptID: "1001", ParentID: "", Name: "总公司", Sort: 1},
		{DeptID: "1002", ParentID: "1001", Name: "技术部", Sort: 2},
		{DeptID: "1003", ParentID: "1001", Name: "市场部", Sort: 1},
		{DeptID: "1004", ParentID: "1002", Name: "前端组", Sort: 2},
		{DeptID: "1005", ParentID: "1002", Name: "后端组", Sort: 1},
		{DeptID: "1006", ParentID: "1003", Name: "市场一组", Sort: 1},
		{DeptID: "1007", ParentID: "1006", Name: "市场一组子部门", Sort: 1},
	}

	// 使用链式构造器
	treeTool := NewTreeBuilder[Department]().
		MaxDepth(2).
		IsRootNode(func(a Department) bool {
			return a.ParentID == ""
		}).
		IsParentNode(func(a, b Department) bool {
			return a.DeptID == b.ParentID
		})
	tree := treeTool.ToTree(depts)

	// 输出结果
	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(jsonData))
}

func exampleWithRoot() {
	// 示例数据
	menus := []*Menu{
		{ID: 1, ParentID: 0, Name: "系统管理", Order: 1},
		{ID: 2, ParentID: 0, Name: "业务管理", Order: 2},
		{ID: 3, ParentID: 0, Name: "用户中心", Order: 3},
	}

	// 使用链式构造器（多个根节点时添加顶级）
	// 构建树
	treeTool := NewTreeBuilder[Menu]().
		MaxDepth(2).
		IsRootNode(func(a Menu) bool {
			return a.ParentID == 0
		}).
		IsParentNode(func(a, b Menu) bool {
			return a.ID == b.ParentID
		}).
		SortFun(func(a, b Menu) int {
			return int(a.Order - b.Order)
		})
	tree := treeTool.ToTree(menus)
	// 输出结果
	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(jsonData))
}
