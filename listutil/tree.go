package listutil

import (
	"fmt"
	"reflect"
	"sort"
)

// ==================== 配置结构体 ====================

// TreeConfig 树形配置
type TreeConfig struct {
	IDField          string // ID字段名
	ParentIDField    string // 父ID字段名
	SortField        string // 排序字段名
	ChildrenField    string // 子节点字段名
	MaxDepth         int    // 最大深度，-1表示不限制
	RootParentID     any    // 根节点的父ID值
	AddRootNode      bool   // 是否添加根节点
	RootNodeID       any    // 根节点ID
	RootNodeName     string // 根节点名称
	OnlyAddWhenMulti bool   // 只有多个根节点时才添加顶级
}

// TreeBuilder 树构建器
type TreeBuilder struct {
	config     *TreeConfig
	fieldCache map[reflect.Type]*fieldInfo
}

// fieldInfo 字段信息缓存
type fieldInfo struct {
	idField       int
	parentIdField int
	sortField     int
	childrenField int
	hasChildren   bool
}

// ==================== 构造函数 ====================

// NewTreeBuilder 创建新的树构建器
func NewTreeBuilder(config *TreeConfig) *TreeBuilder {
	if config == nil {
		config = &TreeConfig{
			IDField:       "ID",
			ParentIDField: "ParentID",
			SortField:     "Sort",
			ChildrenField: "Children",
			MaxDepth:      -1,
		}
	}

	return &TreeBuilder{
		config:     config,
		fieldCache: make(map[reflect.Type]*fieldInfo),
	}
}

// ==================== 核心方法 ====================

// BuildTree 构建树形结构（要求 []*T 类型的指针切片）
func (tb *TreeBuilder) BuildTree(items any) (any, error) {
	if items == nil {
		return nil, fmt.Errorf("input items cannot be nil")
	}

	// 获取反射值
	sliceVal := reflect.ValueOf(items)

	// 验证输入类型
	if sliceVal.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input must be a slice, got %v", sliceVal.Kind())
	}

	// 处理空切片
	if sliceVal.Len() == 0 {
		return items, nil
	}

	// 验证是否为指针切片
	firstElem := sliceVal.Index(0)
	if firstElem.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("slice elements must be pointers, got %v", firstElem.Kind())
	}

	// 获取结构体类型信息
	elemType := firstElem.Type().Elem()
	fields, err := tb.getFieldInfo(elemType)
	if err != nil {
		return nil, err
	}

	// 构建节点映射
	idToPtr := make(map[any]reflect.Value)
	allItems := make([]reflect.Value, sliceVal.Len())

	for i := 0; i < sliceVal.Len(); i++ {
		itemPtr := sliceVal.Index(i)

		// 验证指针有效性
		if itemPtr.Kind() != reflect.Ptr || itemPtr.IsNil() {
			return nil, fmt.Errorf("slice element %d is not a valid pointer", i)
		}

		item := itemPtr.Elem()

		// 获取ID
		id := item.Field(fields.idField).Interface()
		if id == nil {
			return nil, fmt.Errorf("element %d has nil ID", i)
		}

		// 清空现有的子节点
		if fields.hasChildren {
			item.Field(fields.childrenField).Set(reflect.Zero(item.Field(fields.childrenField).Type()))
		}

		// 保存映射
		idToPtr[id] = itemPtr
		allItems[i] = itemPtr
	}

	// 构建树形结构
	roots := tb.buildTreeStructure(allItems, idToPtr, fields)

	// 排序
	tb.sortTree(roots, fields)

	// 构建结果切片
	result := reflect.MakeSlice(sliceVal.Type(), len(roots), len(roots))
	for i, root := range roots {
		result.Index(i).Set(root)
	}

	// 处理根节点包装
	return tb.processRootNode(result.Interface(), len(roots)), nil
}

// ==================== 私有方法 ====================

// getFieldInfo 获取字段信息并缓存
func (tb *TreeBuilder) getFieldInfo(typ reflect.Type) (*fieldInfo, error) {
	// 检查缓存
	if info, exists := tb.fieldCache[typ]; exists {
		return info, nil
	}

	info := &fieldInfo{
		idField:       -1,
		parentIdField: -1,
		sortField:     -1,
		childrenField: -1,
		hasChildren:   false,
	}

	// 遍历字段
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// 检查tag (优先级更高)
		if tag := field.Tag.Get("tree"); tag != "" {
			switch tag {
			case "id":
				info.idField = i
			case "parent_id", "pid":
				info.parentIdField = i
			case "sort", "order":
				info.sortField = i
			case "children":
				info.childrenField = i
				info.hasChildren = true
			}
		}

		// 检查字段名
		if fieldName == tb.config.IDField {
			info.idField = i
		} else if fieldName == tb.config.ParentIDField {
			info.parentIdField = i
		} else if fieldName == tb.config.SortField {
			info.sortField = i
		} else if fieldName == tb.config.ChildrenField {
			info.childrenField = i
			info.hasChildren = true
		}
	}

	// 验证必要字段
	if info.idField == -1 {
		return nil, fmt.Errorf("ID field '%s' not found in struct", tb.config.IDField)
	}
	if info.parentIdField == -1 {
		return nil, fmt.Errorf("ParentID field '%s' not found in struct", tb.config.ParentIDField)
	}

	// 缓存结果
	tb.fieldCache[typ] = info
	return info, nil
}

// buildTreeStructure 构建树形结构
func (tb *TreeBuilder) buildTreeStructure(items []reflect.Value, idToPtr map[any]reflect.Value, fields *fieldInfo) []reflect.Value {
	roots := make([]reflect.Value, 0)
	depthMap := make(map[any]int)

	for _, itemPtr := range items {
		item := itemPtr.Elem()
		parentID := item.Field(fields.parentIdField).Interface()

		// 判断是否为根节点
		if tb.isRootNode(parentID) {
			depthMap[item.Field(fields.idField).Interface()] = 1
			roots = append(roots, itemPtr)
			continue
		}

		// 查找父节点
		if parentPtr, exists := idToPtr[parentID]; exists {
			currentDepth := depthMap[parentID] + 1

			// 检查深度限制
			if tb.config.MaxDepth <= 0 || currentDepth <= tb.config.MaxDepth {
				// 添加子节点到父节点
				parent := parentPtr.Elem()
				children := parent.Field(fields.childrenField)

				if children.Kind() == reflect.Slice {
					children.Set(reflect.Append(children, itemPtr))
					depthMap[item.Field(fields.idField).Interface()] = currentDepth
				}
			}
		} else {
			// 父节点不存在，作为根节点
			depthMap[item.Field(fields.idField).Interface()] = 1
			roots = append(roots, itemPtr)
		}
	}

	return roots
}

// isRootNode 判断是否为根节点
func (tb *TreeBuilder) isRootNode(parentID any) bool {
	// 如果配置了根节点的父ID值
	if tb.config.RootParentID != nil {
		return reflect.DeepEqual(parentID, tb.config.RootParentID)
	}

	// 检查常见零值
	return isZeroValue(parentID)
}

// sortTree 排序树
func (tb *TreeBuilder) sortTree(nodes []reflect.Value, fields *fieldInfo) {
	if len(nodes) == 0 {
		return
	}

	// 排序当前层
	sort.Slice(nodes, func(i, j int) bool {
		valI := nodes[i].Elem().Field(fields.sortField).Interface()
		valJ := nodes[j].Elem().Field(fields.sortField).Interface()
		return tb.compareSortValue(valI, valJ)
	})

	// 递归排序子节点
	for _, node := range nodes {
		elem := node.Elem()
		if fields.hasChildren {
			children := elem.Field(fields.childrenField)
			if children.Len() > 0 {
				childNodes := make([]reflect.Value, children.Len())
				for i := 0; i < children.Len(); i++ {
					childNodes[i] = children.Index(i)
				}
				tb.sortTree(childNodes, fields)

				// 重建排序后的子节点切片
				newChildren := reflect.MakeSlice(children.Type(), len(childNodes), len(childNodes))
				for i, child := range childNodes {
					newChildren.Index(i).Set(child)
				}
				elem.Field(fields.childrenField).Set(newChildren)
			}
		}
	}
}

// compareSortValue 比较排序值
func (tb *TreeBuilder) compareSortValue(a, b any) bool {
	if a == nil && b == nil {
		return false
	}
	if a == nil {
		return true
	}
	if b == nil {
		return false
	}

	// 类型断言比较
	switch av := a.(type) {
	case int:
		if bv, ok := b.(int); ok {
			return av < bv
		}
	case int64:
		if bv, ok := b.(int64); ok {
			return av < bv
		}
	case float64:
		if bv, ok := b.(float64); ok {
			return av < bv
		}
	case string:
		if bv, ok := b.(string); ok {
			return av < bv
		}
	}

	// 尝试数值比较
	fa, ok1 := toFloat64(a)
	fb, ok2 := toFloat64(b)
	if ok1 && ok2 {
		return fa < fb
	}

	// 默认字符串比较
	return toString(a) < toString(b)
}

// processRootNode 处理根节点包装
func (tb *TreeBuilder) processRootNode(tree any, rootCount int) any {
	if !tb.config.AddRootNode {
		return tree
	}

	// 检查是否只在多个根节点时添加顶级
	if tb.config.OnlyAddWhenMulti && rootCount <= 1 {
		return tree
	}

	// 创建根节点
	root := map[string]any{
		tb.config.IDField:       tb.config.RootNodeID,
		"Name":                  tb.config.RootNodeName,
		tb.config.ChildrenField: tree,
	}

	return []map[string]any{root}
}

// ==================== 链式构造器 ====================

// Builder 链式构造器
type Builder struct {
	config *TreeConfig
}

// NewBuilder 创建新的构造器
func NewBuilder() *Builder {
	return &Builder{
		config: &TreeConfig{
			IDField:       "ID",
			ParentIDField: "ParentID",
			SortField:     "Sort",
			ChildrenField: "Children",
			MaxDepth:      -1,
		},
	}
}

// WithIDField 设置ID字段
func (b *Builder) WithIDField(field string) *Builder {
	b.config.IDField = field
	return b
}

// WithParentIDField 设置父ID字段
func (b *Builder) WithParentIDField(field string) *Builder {
	b.config.ParentIDField = field
	return b
}

// WithSortField 设置排序字段
func (b *Builder) WithSortField(field string) *Builder {
	b.config.SortField = field
	return b
}

// WithChildrenField 设置子节点字段
func (b *Builder) WithChildrenField(field string) *Builder {
	b.config.ChildrenField = field
	return b
}

// WithMaxDepth 设置最大深度
func (b *Builder) WithMaxDepth(depth int) *Builder {
	b.config.MaxDepth = depth
	return b
}

// WithRootParentID 设置根节点的父ID值
func (b *Builder) WithRootParentID(id any) *Builder {
	b.config.RootParentID = id
	return b
}

// AddRootNode 添加根节点
func (b *Builder) AddRootNode(id any, name string) *Builder {
	b.config.AddRootNode = true
	b.config.RootNodeID = id
	b.config.RootNodeName = name
	return b
}

// AddRootWhenMulti 多个根节点时才添加顶级
func (b *Builder) AddRootWhenMulti(id any, name string) *Builder {
	b.config.AddRootNode = true
	b.config.OnlyAddWhenMulti = true
	b.config.RootNodeID = id
	b.config.RootNodeName = name
	return b
}

// Build 构建TreeBuilder
func (b *Builder) Build() *TreeBuilder {
	return NewTreeBuilder(b.config)
}

// ==================== 辅助函数 ====================

// isZeroValue 判断是否为零值
func isZeroValue(v any) bool {
	if v == nil {
		return true
	}

	switch val := v.(type) {
	case int:
		return val == 0
	case int64:
		return val == 0
	case float64:
		return val == 0
	case string:
		return val == "" || val == "0" || val == "null" || val == "nil"
	case bool:
		return !val
	default:
		// 尝试反射判断
		return reflect.ValueOf(v).IsZero()
	}
}

// toFloat64 尝试转换为float64
func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

// toString 转换为字符串
func toString(v any) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", v)
	}
}
