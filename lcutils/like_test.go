package lcutils

import "testing"

func TestFindClosestProvince(t *testing.T) {
	provinces := []string{"北京市", "上海", "广东省", "深圳"}
	address := "深圳市南山区科技园"

	result := Like(address, provinces)
	println(result) // 输出：深圳（因为直接匹配到最长的 "深圳"）

	address2 := "沪上海市浦东新区"
	result2 := Like(address2, provinces)
	println(result2) // 输出：上海（编辑距离最小）
}
