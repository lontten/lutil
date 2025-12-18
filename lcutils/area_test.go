package lcutils

import "testing"

func TestFindStandardAddress(t *testing.T) {
	provinces := []string{"北京市", "上海", "广东省", "江苏省", "四川省"}
	cityToProvince := map[string]string{
		"深圳":     "广东省",
		"南京市":   "江苏省",
		"浦东新区": "上海",
		"成都":     "四川省",
	}

	address1 := "深圳市南山区科技园"
	result1 := MatchProvinceCity(address1, provinces, cityToProvince)
	println(result1) // 输出：广东省深圳（假设映射中市名为"深圳"）

	address2 := "沪上海市浦东新区"
	result2 := MatchProvinceCity(address2, provinces, cityToProvince)
	println(result2) // 输出：上海浦东新区（假设"浦东新区"映射到"上海"）

	address3 := "江苏"
	result3 := MatchProvinceCity(address3, provinces, cityToProvince)
	println(result3) // 输出：江苏省

	address4 := "成都市金牛区金府路111号金府西部五金机电城19幢3号"
	result4 := MatchProvinceCity(address4, provinces, cityToProvince)
	println(result4) // 输出：江苏省
}
