// Package listutil 提供切片集合运算与 ListTool 条件检查工具。
package listutil

// BoolDiff 计算 list1 与 list2 的差集（list1 - list2）。
// 有去重逻辑；结果顺序不稳定。
func BoolDiff[T comparable](list1, list2 []T) []T {
	set1 := make(map[T]struct{}, len(list1))
	for _, e := range list1 {
		set1[e] = struct{}{}
	}

	for _, e := range list2 {
		delete(set1, e)
	}

	result := make([]T, 0, len(set1))
	for e := range set1 {
		result = append(result, e)
	}

	return result
}

// BoolEq 判断两个切片是否表示相同集合（有去重逻辑）。
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

// BoolUnion 计算 list1 与 list2 的并集。
// 有去重逻辑；结果顺序不稳定。
func BoolUnion[T comparable](list1, list2 []T) []T {
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

// BoolIntersection 计算 list1 与 list2 的交集。
// 有去重逻辑；结果按 list2 中首次出现的顺序排列。
func BoolIntersection[T comparable](list1, list2 []T) []T {
	set1 := make(map[T]struct{}, len(list1))
	for _, e := range list1 {
		set1[e] = struct{}{}
	}
	result := make([]T, 0)
	seen := make(map[T]struct{})
	for _, e := range list2 {
		if _, ok := set1[e]; !ok {
			continue
		}
		if _, dup := seen[e]; dup {
			continue
		}
		seen[e] = struct{}{}
		result = append(result, e)
	}
	return result
}

// RemoveDuplicates 对切片去重，返回无重复元素的新切片。
// 结果顺序不稳定。
func RemoveDuplicates[T comparable](slice []T) []T {
	return BoolUnion(slice, []T{})
}

// ListHas 判断 slice 中是否包含 item。
func ListHas[T comparable](slice []T, item T) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// ListToolBuilder 链式构建列表条件并执行 Check。
type ListToolBuilder struct {
	list   []any
	hasAll []any
	hasAny []any
	notAll []any
}

// ListTool 创建 ListToolBuilder，list 为待检查的列表元素。
func ListTool(list ...any) *ListToolBuilder {
	return &ListToolBuilder{
		list:   list,
		hasAll: []any{},
		hasAny: []any{},
		notAll: []any{},
	}
}

// HasAll 要求 list 包含所有指定元素（集合语义，重复项不增加要求）。
func (t *ListToolBuilder) HasAll(list ...any) *ListToolBuilder {
	t.hasAll = append(t.hasAll, list...)
	return t
}

// HasAny 要求 list 至少包含一个指定元素。
func (t *ListToolBuilder) HasAny(list ...any) *ListToolBuilder {
	t.hasAny = append(t.hasAny, list...)
	return t
}

// NotAll 要求 list 不包含 notAll 中的全部元素（至少缺一个）。
func (t *ListToolBuilder) NotAll(list ...any) *ListToolBuilder {
	t.notAll = append(t.notAll, list...)
	return t
}

// Check 根据已设置的条件判断 list 是否满足。
func (t *ListToolBuilder) Check() bool {
	if len(t.hasAll) > 0 {
		var c = len(BoolIntersection(t.list, t.hasAll)) == len(RemoveDuplicates(t.hasAll))
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
