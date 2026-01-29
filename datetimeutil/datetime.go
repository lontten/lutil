package datetimeutil

import "github.com/lontten/lcore/v2/types"

// Max 返回最大值，至少需要一个参数
func Max(list ...types.LocalDateTime) types.LocalDateTime {
	if len(list) == 0 {
		panic("timeutil.Max: at least one argument required")
	}
	m := list[0]
	for _, v := range list[1:] {
		if v.After(m) {
			m = v
		}
	}
	return m
}

// Min 返回最小值，至少需要一个参数
func Min(list ...types.LocalDateTime) types.LocalDateTime {
	if len(list) == 0 {
		panic("timeutil.Min: at least one argument required")
	}
	m := list[0]
	for _, v := range list[1:] {
		if v.Before(m) {
			m = v
		}
	}
	return m
}

// MaxP 返回最大值的指针，忽略nil，如果没有非nil值则返回nil
func MaxP(list ...*types.LocalDateTime) *types.LocalDateTime {
	var m *types.LocalDateTime
	for _, v := range list {
		if v == nil {
			continue
		}
		if m == nil || v.After(*m) {
			m = v
		}
	}
	return m
}

// MinP 返回最小值的指针，忽略nil，如果没有非nil值则返回nil
func MinP(list ...*types.LocalDateTime) *types.LocalDateTime {
	var n *types.LocalDateTime
	for _, v := range list {
		if v == nil {
			continue
		}
		if n == nil || v.Before(*n) {
			n = v
		}
	}
	return n
}

// MaxNow 从指针列表中返回最大值，以当前时间为基准，忽略nil
func MaxNow(list ...*types.LocalDateTime) types.LocalDateTime {
	m := types.NowDateTime()
	for _, v := range list {
		if v == nil {
			continue
		}
		if v.After(m) {
			m = *v
		}
	}
	return m
}

// MinNow 从指针列表中返回最小值，以当前时间为基准，忽略nil
func MinNow(list ...*types.LocalDateTime) types.LocalDateTime {
	n := types.NowDateTime()
	for _, v := range list {
		if v == nil {
			continue
		}
		if v.Before(n) {
			n = *v
		}
	}
	return n
}

// MaxNowP 从指针列表中返回最大值的指针，以当前时间为基准，忽略nil
func MaxNowP(list ...*types.LocalDateTime) *types.LocalDateTime {
	m := types.NowDateTime()
	for _, v := range list {
		if v == nil {
			continue
		}
		if v.After(m) {
			m = *v
		}
	}
	return &m
}

// MinNowP 从指针列表中返回最小值的指针，以当前时间为基准，忽略nil
func MinNowP(list ...*types.LocalDateTime) *types.LocalDateTime {
	n := types.NowDateTime()
	for _, v := range list {
		if v == nil {
			continue
		}
		if v.Before(n) {
			n = *v
		}
	}
	return &n
}
