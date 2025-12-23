package listutil

import (
	"encoding/json"
	"fmt"
	"time"
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

	fmt.Println("\n=== 示例4: 性能测试 ===")
	examplePerformance()
}

func exampleMenu() {
	// 示例数据
	menus := []*Menu{
		{ID: 1, ParentID: 0, Name: "系统管理", Order: 1},
		{ID: 2, ParentID: 1, Name: "用户管理", Order: 2},
		{ID: 3, ParentID: 1, Name: "角色管理", Order: 1},
		{ID: 4, ParentID: 0, Name: "业务管理", Order: 2},
		{ID: 5, ParentID: 4, Name: "订单管理", Order: 1},
		{ID: 6, ParentID: 5, Name: "订单列表", Order: 1},
		{ID: 7, ParentID: 5, Name: "订单详情", Order: 2},
	}

	// 使用链式构造器
	builder := NewBuilder().
		WithIDField("ID").
		WithParentIDField("ParentID").
		WithSortField("Order").
		WithChildrenField("Children").
		WithMaxDepth(3). // 限制3层
		WithRootParentID(0). // 父ID为0的是根节点
		Build()

	// 构建树
	tree, err := builder.BuildTree(menus)
	if err != nil {
		fmt.Printf("构建失败: %v\n", err)
		return
	}

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
	builder := NewBuilder().
		WithMaxDepth(2). // 只要前2层
		WithRootParentID(""). // 父ID为空的是根节点
		Build()

	// 构建树
	tree, err := builder.BuildTree(depts)
	if err != nil {
		fmt.Printf("构建失败: %v\n", err)
		return
	}

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
	builder := NewBuilder().
		WithIDField("ID").
		WithParentIDField("ParentID").
		WithSortField("Order").
		WithChildrenField("Children").
		WithRootParentID(0).
		AddRootWhenMulti(0, "所有菜单"). // 多个根节点时才添加顶级
		Build()

	// 构建树
	tree, err := builder.BuildTree(menus)
	if err != nil {
		fmt.Printf("构建失败: %v\n", err)
		return
	}

	// 输出结果
	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Println(string(jsonData))
}

func examplePerformance() {
	// 性能测试数据结构
	type BenchNode struct {
		ID       int          `tree:"id"`
		ParentID int          `tree:"pid"`
		Name     string       `tree:"name"`
		Sort     int          `tree:"sort"`
		Children []*BenchNode `tree:"children"`
	}

	// 生成测试数据
	generateData := func(count int) []*BenchNode {
		nodes := make([]*BenchNode, count)
		for i := 0; i < count; i++ {
			nodes[i] = &BenchNode{
				ID:       i,
				ParentID: i / 10, // 创建树形结构
				Name:     fmt.Sprintf("Node-%d", i),
				Sort:     i % 10,
			}
		}
		return nodes
	}

	// 测试不同数据量
	testCases := []struct {
		name  string
		count int
	}{
		{"1,000 nodes", 1000},
		{"10,000 nodes", 10000},
		{"50,000 nodes", 50000},
	}

	builder := NewBuilder().
		WithIDField("ID").
		WithParentIDField("ParentID").
		WithSortField("Sort").
		WithChildrenField("Children").
		WithRootParentID(-1).
		Build()

	for _, tc := range testCases {
		fmt.Printf("\n测试 %s:\n", tc.name)
		data := generateData(tc.count)

		start := time.Now()
		_, err := builder.BuildTree(data)
		elapsed := time.Since(start)

		if err != nil {
			fmt.Printf("  错误: %v\n", err)
		} else {
			fmt.Printf("  耗时: %v\n", elapsed)
			fmt.Printf("  平均: %v/节点\n", elapsed/time.Duration(tc.count))
		}
	}
}
