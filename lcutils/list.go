package lcutils

// BoolDiff
// list1-list2
// 布尔-差集,
// 有去重逻辑
// 移除list1 中 所有的 list2 元素
func BoolDiff[T comparable](list1, list2 []T) []T {
	set1 := make(map[T]struct{}, len(list1))
	for _, e := range list1 {
		set1[e] = struct{}{} // 利用 map 自动去重 list1
	}

	// 移除 list2 中存在的元素
	for _, e := range list2 {
		delete(set1, e)
	}

	result := make([]T, 0, len(set1))
	for e := range set1 {
		result = append(result, e)
	}

	return result
}

// BoolEq
// 布尔-相等
// 有去重逻辑
func BoolEq[T comparable](list1, list2 []T) bool {
	set1 := make(map[T]struct{}, len(list1))
	for _, e := range list1 {
		set1[e] = struct{}{}
	}

	set1Len := len(set1)
	list2Len := len(list2)
	if set1Len > list2Len {
		return false
	}

	set2 := make(map[T]struct{}, list2Len)
	for _, e := range list2 {
		_, ok := set1[e]
		if !ok {
			return false
		}
		set2[e] = struct{}{}
	}

	return set1Len == len(set2)
}

// BoolUnion
// list1+list2
// 布尔-并,
// 有去重逻辑
// list1 和 list2 元素，并集
func BoolUnion[T comparable](list1, list2 []T) []T {
	// 创建集合存储list1的元素
	set1 := make(map[T]struct{}, len(list1)+len(list2))
	for _, e := range list1 {
		set1[e] = struct{}{}
	}
	for _, e := range list2 {
		set1[e] = struct{}{}
	}

	result := make([]T, 0, len(set1))
	for e := range set1 {
		result = append(result, e)
	}
	return result
}

// BoolIntersection
// list1，list2相同的元素
// 布尔-交集,
// 有去重逻辑
// list1 和 list2 元素，交集
func BoolIntersection[T comparable](list1, list2 []T) []T {
	// 创建集合存储list1的元素
	set1 := make(map[T]struct{}, len(list1))
	for _, e := range list1 {
		set1[e] = struct{}{}
	}
	result := make([]T, 0, len(set1))
	for _, e := range list2 {
		if _, ok := set1[e]; ok {
			result = append(result, e)
		}
	}
	return result
}

// 去重函数，适用于任何可比较的类型
func RemoveDuplicates[T comparable](slice []T) []T {
	return BoolUnion(slice, []T{})
}

// 集合中是否包含item
func ListHas[T comparable](slice []T, item T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

type ListToolBuilder struct {
	list   []any
	hasAll []any
	hasAny []any
	notAll []any
}

func ListTool(list ...any) *ListToolBuilder {
	return &ListToolBuilder{
		list:   list,
		hasAll: []any{},
		hasAny: []any{},
		notAll: []any{},
	}
}

func (t *ListToolBuilder) HasAll(list ...any) *ListToolBuilder {
	t.hasAll = append(t.hasAll, list...)
	return t
}
func (t *ListToolBuilder) HasAny(list ...any) *ListToolBuilder {
	t.hasAny = append(t.hasAny, list...)
	return t
}
func (t *ListToolBuilder) NotAll(list ...any) *ListToolBuilder {
	t.notAll = append(t.notAll, list...)
	return t
}
func (t *ListToolBuilder) Check() bool {
	if len(t.hasAll) > 0 {
		var c = len(BoolIntersection(t.list, t.hasAll)) == len(t.hasAll)
		if !c {
			return false
		}
	}
	if len(t.hasAny) > 0 {
		var c = len(BoolIntersection(t.list, t.hasAny)) > 0
		if !c {
			return false
		}
	}
	if len(t.notAll) > 0 {
		var c = len(BoolIntersection(t.list, t.notAll)) == 0
		if !c {
			return false
		}

	}
	return true
}
