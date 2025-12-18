package lcutils

import (
	"sort"
)

// MatchProvinceCity 返回地址中最匹配的标准省市名称
// provinces: 所有标准省名称列表
// cityToProvince: 市名称到省名称的映射（例如 "深圳" -> "广东省"）
func MatchProvinceCity(address string, provinces []string, cityToProvince map[string]string) string {
	// 提取所有市名称用于匹配
	cities := make([]string, 0, len(cityToProvince))
	for city := range cityToProvince {
		cities = append(cities, city)
	}

	// 1. 优先匹配市
	if matchedCity := fuzzyMatch(address, cities); matchedCity != "" {
		return cityToProvince[matchedCity] + matchedCity // 例如：广东省深圳
	}

	// 2. 若无市匹配，则匹配省
	if matchedProvince := fuzzyMatch(address, provinces); matchedProvince != "" {
		return matchedProvince
	}

	return "" // 无匹配
}

// fuzzyMatch 通用模糊匹配函数（返回最长的包含匹配，或编辑距离最小的候选词）
func fuzzyMatch(address string, candidates []string) string {
	addressRunes := []rune(address)

	// 步骤1：检查直接包含的候选词，优先选最长的
	var containsMatches []string
	for _, c := range candidates {
		cr := []rune(c)
		if containsRunes(addressRunes, cr) {
			containsMatches = append(containsMatches, c)
		}
	}
	if len(containsMatches) > 0 {
		sort.Slice(containsMatches, func(i, j int) bool { // 按长度降序
			return len([]rune(containsMatches[i])) > len([]rune(containsMatches[j]))
		})
		return containsMatches[0]
	}

	// 步骤2：无包含匹配时，计算编辑距离
	maxScore := -1.0
	bestMatch := ""
	addressLen := len(addressRunes)
	for _, c := range candidates {
		cRunes := []rune(c)
		cLen := len(cRunes)
		distance := levenshteinDistance(addressRunes, cRunes)
		score := 1.0 - float64(distance)/float64(max(addressLen, cLen))
		if score > maxScore {
			maxScore = score
			bestMatch = c
		}
	}
	return bestMatch
}
