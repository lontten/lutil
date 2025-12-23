package listutil

import (
	"reflect"
	"sort"
)

// ==================== 配置结构体 ====================

// TreeBuilder 树构建器
type TreeBuilder[T any] struct {
	source       map[int]*T
	maxDepth     int // 最大深度，-1表示不限制,root是0
	sortFun      func(a, b T) int
	isParentNode func(a, b T) bool // 判断a是否为b父节点, true为父节点,
	isRootNode   func(a T) bool

	childrenField string
}

// ==================== 构造函数 ====================

func (c *TreeBuilder[T]) ToTree(list []*T) []*T {
	if c.isRootNode == nil {
		panic("isRootNode 函数未设置")
	}
	if c.isParentNode == nil {
		panic("isParentNode 函数未设置")
	}
	if c.childrenField == "" {
		panic("childrenField 未设置")
	}
	c.source = make(map[int]*T)
	var res = make([]*T, 0)

	for i, v := range list {
		b := c.isRootNode(*v)
		if b {
			res = append(res, v)
		} else {
			c.source[i] = v
		}
	}
	for _, re := range res {
		c.list2Tree(re, 1)
	}
	return res
}
func (c *TreeBuilder[T]) list2Tree(target *T, deepth int) {
	if deepth > c.maxDepth && c.maxDepth != -1 {
		return
	}
	var arr = make([]*T, 0)

	v := reflect.ValueOf(target).Elem()
	childrenField := v.FieldByName(c.childrenField)

	var indexs []int
	for k, v := range c.source {
		b := c.isParentNode(*target, *v)
		if b {
			arr = append(arr, v)
			indexs = append(indexs, k)
		}
	}
	for _, index := range indexs {
		delete(c.source, index)
	}
	if c.sortFun != nil {
		sort.SliceStable(arr, func(i, j int) bool {
			return c.sortFun(*arr[i], *arr[j]) < 0
		})
	}
	for _, re := range arr {
		c.list2Tree(re, deepth+1)
	}
	childrenField.Set(reflect.ValueOf(arr))
	return
}

// NewTreeBuilder 创建新的树构建器
func NewTreeBuilder[T any]() *TreeBuilder[T] {
	return &TreeBuilder[T]{
		childrenField: "Children",
		maxDepth:      -1,
	}
}

func (c *TreeBuilder[T]) MaxDepth(maxDepth int) *TreeBuilder[T] {
	c.maxDepth = maxDepth
	return c
}

func (c *TreeBuilder[T]) ChildrenField(childrenField string) *TreeBuilder[T] {
	var t T
	of := reflect.TypeOf(t)
	name, ok := of.FieldByName(childrenField)
	if !ok {
		panic("Children 字段不存在")
	}

	// 检查字段类型是否为 []*T 类型
	if name.Type.Kind() != reflect.Slice {
		panic("Children 字段必须是切片类型")
	}

	// 检查切片元素类型是否为 *T
	elemType := name.Type.Elem()
	if elemType.Kind() != reflect.Ptr || elemType.Elem() != reflect.TypeOf(t) {
		panic("Children 字段必须是 []*T 类型")
	}

	c.childrenField = childrenField
	return c
}

func (c *TreeBuilder[T]) SortFun(sortFun func(a, b T) int) *TreeBuilder[T] {
	c.sortFun = sortFun
	return c
}

func (c *TreeBuilder[T]) IsParentNode(isParentNode func(a, b T) bool) *TreeBuilder[T] {
	c.isParentNode = isParentNode
	return c
}
func (c *TreeBuilder[T]) IsRootNode(isRootNode func(a T) bool) *TreeBuilder[T] {
	c.isRootNode = isRootNode
	return c
}
