package fuzzutil

// Like 从 list 中找出与 key 最相似的一项并返回其字符串。
// 规则：minMatchLen=1，minOverlap=1，编辑距离不设上限；至少 1 个 rune 相同才返回（无候选或均不满足时返回空字符串）。
// 适用于临时、单次、无词表缓存的模糊匹配；固定词表反复提取请用 Vocabulary.ExtractFromText。
//
// 示例：Like("深圳市", []string{"深圳", "广州"}) // "深圳"
func Like(key string, list []string) string {
	term, _, ok := matchBest(key, list, matchRules{
		minMatchLen:     1,
		minOverlap:      1,
		maxEditDistance: -1,
	})
	if !ok {
		return ""
	}
	return term
}
