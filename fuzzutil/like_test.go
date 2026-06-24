package fuzzutil

import "testing"

func TestLike(t *testing.T) {
	tests := []struct {
		name string
		key  string
		list []string
		want string
	}{
		{
			name: "子串包含",
			key:  "深圳市",
			list: []string{"深圳", "广州"},
			want: "深圳",
		},
		{
			name: "最长优先",
			key:  "南京市",
			list: []string{"南京", "南京市"},
			want: "南京市",
		},
		{
			name: "编辑距离",
			key:  "江苏",
			list: []string{"江苏省", "浙江省"},
			want: "江苏省",
		},
		{
			name: "无候选",
			key:  "abc",
			list: nil,
			want: "",
		},
		{
			name: "空 key 返回编辑距离最小候选",
			key:  "",
			list: []string{"深圳"},
			want: "深圳",
		},
		{
			name: "沪上海市浦东新区匹配上海",
			key:  "沪上海市浦东新区",
			list: []string{"北京市", "上海", "广东省", "深圳"},
			want: "上海",
		},
		{
			name: "无相同字不命中",
			key:  "abc",
			list: []string{"xyz"},
			want: "",
		},
		{
			name: "至少一字相同可命中",
			key:  "abc",
			list: []string{"axc"},
			want: "axc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Like(tt.key, tt.list)
			if got != tt.want {
				t.Errorf("Like(%q, %v) = %q, want %q", tt.key, tt.list, got, tt.want)
			}
		})
	}
}

func TestMatchBest_Rules(t *testing.T) {
	strict := matchRules{minMatchLen: 2, maxEditDistance: 1}
	_, _, ok := matchBest("京", []string{"北京", "南京"}, strict)
	if ok {
		t.Error("单字 key 在 minMatchLen=2 时不应命中")
	}

	onlyContain := matchRules{minMatchLen: 2, maxEditDistance: 0}
	_, _, ok = matchBest("江苏", []string{"江苏省"}, onlyContain)
	if ok {
		t.Error("maxEditDistance=0 时不应通过编辑距离命中")
	}

	_, _, ok = matchBest("江苏", []string{"江苏省"}, matchRules{
		minMatchLen: 2, minOverlap: 2, maxEditDistance: 1,
	})
	if !ok {
		t.Error("minOverlap=2 时江苏/江苏省 应模糊命中")
	}

	_, _, ok = matchBest("哈哈", []string{"江苏省"}, matchRules{
		minMatchLen: 2, minOverlap: 2, maxEditDistance: -1,
	})
	if ok {
		t.Error("minOverlap=2 时哈哈/江苏省 无相同字不应命中")
	}

	_, _, ok = matchBest("abc", []string{"xyz"}, matchRules{
		minMatchLen: 1, minOverlap: 1, maxEditDistance: 1,
	})
	if ok {
		t.Error("minOverlap=1 时 abc/xyz 无相同字不应命中")
	}

	term, kind, ok := matchBest("深圳市", []string{"深圳"}, matchRules{
		minMatchLen: 2, minOverlap: 2, maxEditDistance: 0,
	})
	if !ok || term != "深圳" || kind != MatchContain {
		t.Errorf("子串包含+minOverlap=2: got %q %v %v", term, kind, ok)
	}
}

func TestCommonRuneOverlap(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"哈哈", "江苏省", 0},
		{"江苏", "江苏省", 2},
		{"哈哈", "哈哈啊", 2},
		{"aab", "aba", 3},
		{"abc", "xyz", 0},
	}
	for _, tt := range tests {
		got := commonRuneOverlap([]rune(tt.a), []rune(tt.b))
		if got != tt.want {
			t.Errorf("commonRuneOverlap(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
