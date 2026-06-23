// Package numutil 提供数值相关的工具函数。
package numutil

import "math"

// Round45 对浮点数进行四舍五入，返回 int64。
// num 为待四舍五入的浮点数。
func Round45(num float64) int64 {
	if num >= 0 {
		// 正数：加 0.5 后向下取整
		return int64(math.Floor(num + 0.5))
	}
	// 负数：减 0.5 后向上取整
	return int64(math.Ceil(num - 0.5))
}

// Round45i 对浮点数进行四舍五入，返回 int。
// num 为待四舍五入的浮点数。
func Round45i(num float64) int {
	if num >= 0 {
		// 正数：加 0.5 后向下取整
		return int(math.Floor(num + 0.5))
	}
	// 负数：减 0.5 后向上取整
	return int(math.Ceil(num - 0.5))
}
