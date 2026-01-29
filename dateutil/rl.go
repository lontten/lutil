package dateutil

import "github.com/lontten/lcore/v2/types"

// IsNowR 判断时间是否在 当前时间 右侧，之后
func IsNowR(v *types.LocalDate) bool {
	if v == nil {
		return false
	}
	now := types.NowDate()

	return v.After(now)
}

// IsNowL 判断时间是否在 当前时间 左侧，之前
func IsNowL(v *types.LocalDate) bool {
	if v == nil {
		return false
	}
	now := types.NowDate()

	return v.Before(now)
}
