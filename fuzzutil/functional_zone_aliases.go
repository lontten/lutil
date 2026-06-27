package fuzzutil

//go:generate go run gen_functional_zones.go .

import "sync"

var (
	defaultFunctionalZonesOnce sync.Once
	defaultFunctionalZones     map[string]string
)

// DefaultFunctionalZoneAliases 返回内建国家级功能区/开发区 → 规范区县名映射（浅拷贝）。
// 数据来自 data/functional_zones_national.json，由 gen_functional_zones.go 生成 nationalFunctionalZones。
// 业务扩展请用 MatchOpts().FunctionalZones。
func DefaultFunctionalZoneAliases() map[string]string {
	defaultFunctionalZonesOnce.Do(func() {
		defaultFunctionalZones = make(map[string]string, len(nationalFunctionalZones))
		for k, v := range nationalFunctionalZones {
			defaultFunctionalZones[k] = v
		}
	})
	out := make(map[string]string, len(defaultFunctionalZones))
	for k, v := range defaultFunctionalZones {
		out[k] = v
	}
	return out
}

func mergeFunctionalZoneMap(dst, src map[string]string) map[string]string {
	if len(src) == 0 {
		return dst
	}
	if dst == nil {
		dst = make(map[string]string, len(src))
	}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// functionalZonesToNameAliases 将 功能区名→区县名 反转为 nameAliases 格式（区县名→[]功能区名）。
func functionalZonesToNameAliases(zones map[string]string) map[string][]string {
	if len(zones) == 0 {
		return nil
	}
	out := make(map[string][]string)
	for zoneName, district := range zones {
		if district == "" || zoneName == "" {
			continue
		}
		out[district] = append(out[district], zoneName)
	}
	for k := range out {
		out[k] = dedupeStrings(out[k])
		sortNamesByRuneDesc(out[k])
	}
	return out
}

// NationalFunctionalZoneAliasCount 返回内建国家级功能区 alias 条目数（测试用）。
func NationalFunctionalZoneAliasCount() int {
	return len(nationalFunctionalZones)
}
