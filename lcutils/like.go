package lcutils

import (
	"sort"
)

// Like 模糊匹配，返回最匹配的
func Like(key string, list []string) string {
	keyRunes := []rune(key)

	// 步骤1：检查所有包含的，并选择最长匹配
	var candidates []string
	for _, p := range list {
		pr := []rune(p)
		if containsRunes(keyRunes, pr) {
			candidates = append(candidates, p)
		}
	}
	if len(candidates) > 0 {
		// 按字符长度降序排序，选择最长的
		sort.Slice(candidates, func(i, j int) bool {
			return len([]rune(candidates[i])) > len([]rune(candidates[j]))
		})
		return candidates[0]
	}

	// 步骤2：无包含匹配，使用编辑距离计算相似度
	maxScore := -1.0
	best := ""
	keyLen := len(keyRunes)
	for _, p := range list {
		pRunes := []rune(p)
		pLen := len(pRunes)
		distance := levenshteinDistance(keyRunes, pRunes)
		maxLen := max(keyLen, pLen)
		if maxLen == 0 {
			continue
		}
		score := 1.0 - float64(distance)/float64(maxLen)
		if score > maxScore {
			maxScore = score
			best = p
		}
	}

	return best
}

// containsRunes 检查 rune 数组 `s` 是否包含 `t` 作为子序列
func containsRunes(s, t []rune) bool {
	if len(t) == 0 {
		return true
	}
	if len(s) < len(t) {
		return false
	}
	for i := 0; i <= len(s)-len(t); i++ {
		match := true
		for j := 0; j < len(t); j++ {
			if s[i+j] != t[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// levenshteinDistance 计算两个 rune 数组的编辑距离
func levenshteinDistance(a, b []rune) int {
	alen, blen := len(a), len(b)
	if alen == 0 {
		return blen
	}
	if blen == 0 {
		return alen
	}

	// 初始化二维矩阵
	matrix := make([][]int, alen+1)
	for i := range matrix {
		matrix[i] = make([]int, blen+1)
		matrix[i][0] = i
	}
	for j := 1; j <= blen; j++ {
		matrix[0][j] = j
	}

	// 计算编辑距离
	for i := 1; i <= alen; i++ {
		for j := 1; j <= blen; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			matrix[i][j] = min(
				matrix[i-1][j]+1,      // 删除
				matrix[i][j-1]+1,      // 插入
				matrix[i-1][j-1]+cost, // 替换
			)
		}
	}

	return matrix[alen][blen]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
