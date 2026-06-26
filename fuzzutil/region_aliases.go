package fuzzutil

import "sync"

// adminSuffixes 行政区划后缀，按长度从长到短排列。
var adminSuffixes = []string{
	"维吾尔自治区",
	"壮族自治区",
	"回族自治区",
	"特别行政区",
	"自治区",
	"自治州",
	"达斡尔族区",
	"回族区",
	"自治县",
	"自治旗",
	"地区",
	"林区",
	"特区",
	"省",
	"市",
	"盟",
	"旗",
	"县",
	"区",
}

var (
	defaultRegionAliasesOnce sync.Once
	defaultRegionAliases     map[string][]string
)

// DefaultRegionNameAliases 返回内建中国行政区划别名表（固定简称）。
// 返回值为浅拷贝，调用方请勿修改；需自定义请用 MatchOpts().NameAliases。
func DefaultRegionNameAliases() map[string][]string {
	defaultRegionAliasesOnce.Do(func() {
		defaultRegionAliases = map[string][]string{
			"新疆维吾尔自治区": {"新疆"},
			"西藏自治区":    {"西藏"},
			"内蒙古自治区":   {"内蒙古"},
			"广西壮族自治区":  {"广西"},
			"宁夏回族自治区":  {"宁夏"},
			"香港特别行政区":  {"香港"},
			"澳门特别行政区":  {"澳门"},
		}
	})
	out := make(map[string][]string, len(defaultRegionAliases))
	for k, v := range defaultRegionAliases {
		copied := make([]string, len(v))
		copy(copied, v)
		out[k] = copied
	}
	return out
}

// adminSuffixAliases 从节点名剥离行政区划后缀，生成额外匹配候选。
func adminSuffixAliases(name string) []string {
	var result []string
	stripAdminSuffixesRecursive(name, &result)
	return dedupeStrings(result)
}

func stripAdminSuffixesRecursive(name string, result *[]string) {
	runes := []rune(name)
	for _, suffix := range adminSuffixes {
		suffixRunes := []rune(suffix)
		if len(runes) <= len(suffixRunes) {
			continue
		}
		if string(runes[len(runes)-len(suffixRunes):]) != suffix {
			continue
		}
		stripped := string(runes[:len(runes)-len(suffixRunes)])
		if len([]rune(stripped)) < 2 {
			return
		}
		*result = append(*result, stripped)
		stripAdminSuffixesRecursive(stripped, result)
		return
	}
}

func dedupeStrings(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, s := range items {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

func mergeAliasMap(dst map[string][]string, src map[string][]string) map[string][]string {
	if len(src) == 0 {
		return dst
	}
	if dst == nil {
		dst = make(map[string][]string, len(src))
	}
	for k, v := range src {
		merged := append(append([]string{}, dst[k]...), v...)
		dst[k] = dedupeStrings(merged)
	}
	return dst
}
