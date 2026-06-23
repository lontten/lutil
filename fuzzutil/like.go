package fuzzutil

// Like 从 list 中找出与 key 最相似的一项并返回其字符串。
// 规则宽松：minMatchLen=1，编辑距离不设上限，始终返回最佳候选（无候选时返回空字符串）。
// 适用于临时、单次、无词表缓存的模糊匹配；固定词表反复提取请用 Vocabulary.ExtractFromText。
//
// 示例：Like("深圳市", []string{"深圳", "广州"}) // "深圳"
func Like(key string, list []string) string {
	term, _, ok := matchBest(key, list, matchRules{
		minMatchLen:     1,
		maxEditDistance: -1,
	})
	if !ok {
		return ""
	}
	return term
}
