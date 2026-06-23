// Package fuzzutil 提供字符串模糊匹配（Like）与关系链词表提取（Vocabulary）。
//
// 临时单次匹配：
//
//	fuzzutil.Like("深圳市", []string{"深圳", "广州"})
//
// 关系链词表提取（对每条链逐节点加权计分）：
//
//	vocab := fuzzutil.NewVocabulary(nodes)
//	result := vocab.ExtractFromText("深圳市南山区")
package fuzzutil

import (
	"sort"
)

// MatchKind 表示一次匹配的命中方式。
type MatchKind int

const (
	// MatchNone 未匹配。
	MatchNone MatchKind = iota
	// MatchContain 通过子串包含命中（text 中包含候选词）。
	MatchContain
	// MatchFuzzy 通过编辑距离命中（无子串包含时，距离在阈值内）。
	MatchFuzzy
)

// matchRules 控制 matchBest 的匹配规则。
type matchRules struct {
	// minMatchLen 候选词至少需要的 rune 数，低于此长度的候选不参与匹配。
	minMatchLen int
	// maxEditDistance 允许的最大编辑距离；0 表示禁用编辑距离阶段（仅子串包含）；
	// -1 表示不限制距离（Like 使用，始终返回编辑距离最小的候选）。
	maxEditDistance int
}

// matchBest 从 candidates 中找出与 text 最匹配的一项。
// candidates 建议已按 rune 长度降序排列，同阶段多命中时仍会在结果中取最长者。
//
// 两阶段策略：
//  1. 子串包含：text 包含 candidate，且 len(candidate) >= minMatchLen
//  2. 编辑距离：仅当 maxEditDistance != 0 且无包含命中时执行
func matchBest(text string, candidates []string, rules matchRules) (term string, kind MatchKind, ok bool) {
	if len(candidates) == 0 {
		return "", MatchNone, false
	}

	textRunes := []rune(text)

	// minMatchLen>=2 时，text 本身也须达到该长度（避免单字「京」模糊命中「南京」）
	if rules.minMatchLen >= 2 && len(textRunes) < rules.minMatchLen {
		return "", MatchNone, false
	}

	// 阶段 1：子串包含，取最长命中
	var containHits []string
	for _, c := range candidates {
		cRunes := []rune(c)
		if len(cRunes) < rules.minMatchLen {
			continue
		}
		if containsRunes(textRunes, cRunes) {
			containHits = append(containHits, c)
		}
	}
	if len(containHits) > 0 {
		best := containHits[0]
		bestLen := len([]rune(best))
		for _, c := range containHits[1:] {
			if l := len([]rune(c)); l > bestLen {
				best = c
				bestLen = l
			}
		}
		return best, MatchContain, true
	}

	// 阶段 2：编辑距离
	if rules.maxEditDistance == 0 {
		return "", MatchNone, false
	}

	bestTerm := ""
	bestDist := -1
	bestLen := -1

	for _, c := range candidates {
		cRunes := []rune(c)
		if len(cRunes) < rules.minMatchLen {
			continue
		}
		dist := levenshteinDistance(textRunes, cRunes)
		if rules.maxEditDistance > 0 && dist > rules.maxEditDistance {
			continue
		}
		cLen := len(cRunes)
		if bestDist < 0 || dist < bestDist || (dist == bestDist && cLen > bestLen) {
			bestDist = dist
			bestTerm = c
			bestLen = cLen
		}
	}

	if bestTerm != "" {
		return bestTerm, MatchFuzzy, true
	}
	return "", MatchNone, false
}

// containsRunes 检查 s 是否包含连续子串 t（非子序列）。
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

// levenshteinDistance 计算两个 rune 切片的 Levenshtein 编辑距离。
func levenshteinDistance(a, b []rune) int {
	alen, blen := len(a), len(b)
	if alen == 0 {
		return blen
	}
	if blen == 0 {
		return alen
	}

	matrix := make([][]int, alen+1)
	for i := range matrix {
		matrix[i] = make([]int, blen+1)
		matrix[i][0] = i
	}
	for j := 1; j <= blen; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= alen; i++ {
		for j := 1; j <= blen; j++ {
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			matrix[i][j] = minInt(
				matrix[i-1][j]+1,
				matrix[i][j-1]+1,
				matrix[i-1][j-1]+cost,
			)
		}
	}
	return matrix[alen][blen]
}

func minInt(a, b, c int) int {
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

// sortNamesByRuneDesc 按 rune 长度降序排序，保证「南京市」优先于「南京」。
func sortNamesByRuneDesc(names []string) {
	sort.Slice(names, func(i, j int) bool {
		return len([]rune(names[i])) > len([]rune(names[j]))
	})
}
