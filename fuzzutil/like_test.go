package fuzzutil

import (
	"testing"
)

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
}
