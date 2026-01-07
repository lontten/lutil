package numutil

import "math"

// Round45 对浮点数进行四舍五入，返回整数
// 参数：num 待四舍五入的浮点数
// 返回值：四舍五入后的整数
func Round45(num float64) int64 {
	if num >= 0 {
		// 正数：加 0.5 后向下取整
		return int64(math.Floor(num + 0.5))
	}
	// 负数：减 0.5 后向上取整
	return int64(math.Ceil(num - 0.5))
}
func Round45i(num float64) int {
	if num >= 0 {
		// 正数：加 0.5 后向下取整
		return int(math.Floor(num + 0.5))
	}
	// 负数：减 0.5 后向上取整
	return int(math.Ceil(num - 0.5))
}
