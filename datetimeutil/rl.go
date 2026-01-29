package datetimeutil

import "github.com/lontten/lcore/v2/types"

// IsNowR 判断时间是否在 当前时间 右侧，之后
func IsNowR(v *types.LocalDateTime) bool {
	if v == nil {
		return false
	}
	now := types.NowDateTime()

	return v.After(now)
}

// IsNowL 判断时间是否在 当前时间 左侧，之前
func IsNowL(v *types.LocalDateTime) bool {
	if v == nil {
		return false
	}
	now := types.NowDateTime()

	return v.Before(now)
}
