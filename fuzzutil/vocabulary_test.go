package fuzzutil

import (
	"reflect"
	"testing"

	"github.com/lontten/lutil/logutil"
)

func regionVocab() *Vocabulary {
	return NewVocabularyFromPaths(
		NamePath{"广东省", "深圳"},
		NamePath{"上海", "浦东新区"},
		NamePath{"江苏省", "南京"},
		NamePath{"江苏省"},
		NamePath{"四川省", "成都"},
	)
}

func TestNewVocabularyFromNodes(t *testing.T) {
	nodes := []VocabNode{
		{ID: "1", ParentID: "", Name: "广东省"},
		{ID: "2", ParentID: "1", Name: "深圳"},
		{ID: "3", ParentID: "", Name: "江苏省"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.MatchFromText("深圳市南山区")
	wantPath := NamePath{"广东省", "深圳"}
	if !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("Path = %v, want %v", got.Path, wantPath)
	}
	if got.MatchedNodeID != "2" {
		t.Fatalf("MatchedNodeID = %q, want %q", got.MatchedNodeID, "2")
	}
}

func TestNewVocabularyFromTree(t *testing.T) {
	vocab := NewVocabularyFromTree(TreeNode{
		Name: "广东省",
		Children: []TreeNode{
			{Name: "深圳"},
			{Name: "广州"},
		},
	})
	got := vocab.MatchFromText("广州市天河区")
	if !got.Matched || got.Path[len(got.Path)-1] != "广州" {
		t.Fatalf("unexpected result: %+v", got)
	}
	if got.MatchedNodeID != "" {
		t.Fatalf("MatchedNodeID = %q, want empty for auto-assigned IDs", got.MatchedNodeID)
	}
}

func TestMatchFromText_FromPaths_NoNodeID(t *testing.T) {
	got := regionVocab().MatchFromText("深圳市南山区科技园")
	if !got.Matched {
		t.Fatal("expected match")
	}
	if got.MatchedNodeID != "" {
		t.Fatalf("MatchedNodeID = %q, want empty for FromPaths", got.MatchedNodeID)
	}
	if len(got.Path) == 0 {
		t.Fatal("Path should be non-empty")
	}
}

func TestMatchFromText_FromTree_WithUserID(t *testing.T) {
	vocab := NewVocabularyFromTree(TreeNode{
		ID:   "province-1",
		Name: "广东省",
		Children: []TreeNode{
			{Name: "深圳"},
		},
	})
	got := vocab.MatchFromText("深圳市")
	if got.MatchedNodeID != "" {
		t.Fatalf("synthetic child ID: MatchedNodeID = %q, want empty", got.MatchedNodeID)
	}

	vocab2 := NewVocabularyFromTree(TreeNode{
		Name: "广东省",
		Children: []TreeNode{
			{ID: "city-sz", Name: "深圳"},
		},
	})
	got2 := vocab2.MatchFromText("深圳市")
	if got2.MatchedNodeID != "city-sz" {
		t.Fatalf("user-provided ID: MatchedNodeID = %q, want city-sz", got2.MatchedNodeID)
	}
}

func TestMatchFromText(t *testing.T) {
	vocab := regionVocab()

	tests := []struct {
		name     string
		text     string
		matched  bool
		wantPath NamePath
		wantKind MatchKind
		opts     *matchOpts
	}{
		{
			name:     "省市-市命中",
			text:     "深圳市南山区科技园",
			matched:  true,
			wantPath: NamePath{"广东省", "深圳"},
			wantKind: MatchContain,
		},
		{
			name:     "省市-区命中",
			text:     "沪上海市浦东新区",
			matched:  true,
			wantPath: NamePath{"上海", "浦东新区"},
			wantKind: MatchContain,
		},
		{
			name:     "省市-仅省-模糊",
			text:     "江苏",
			matched:  true,
			wantPath: NamePath{"江苏省"},
			wantKind: MatchFuzzy,
		},
		{
			name:    "无匹配",
			text:    "abcxyz",
			matched: false,
		},
		{
			name:    "MinMatchLen 单字不命中",
			text:    "京",
			matched: false,
			opts:    MatchOpts().MinMatchLen(2),
		},
		{
			name:    "MaxEditDistance=0 无包含则不命中",
			text:    "江苏",
			matched: false,
			opts:    MatchOpts().MaxEditDistance(0),
		},
		{
			name:     "成都",
			text:     "成都市金牛区金府路111号",
			matched:  true,
			wantPath: NamePath{"四川省", "成都"},
			wantKind: MatchContain,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vocab.MatchFromText(tt.text, tt.opts)
			if got.Matched != tt.matched {
				t.Fatalf("Matched = %v, want %v (result=%+v)", got.Matched, tt.matched, got)
			}
			if !tt.matched {
				return
			}
			if !reflect.DeepEqual(got.Path, tt.wantPath) {
				t.Errorf("Path = %v, want %v", got.Path, tt.wantPath)
			}
			if got.MatchKind != tt.wantKind {
				t.Errorf("MatchKind = %v, want %v", got.MatchKind, tt.wantKind)
			}
		})
	}
}

func TestMatchFromText_CategoryThreeLevel(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"服饰", "运动鞋", "跑步鞋"},
		NamePath{"服饰", "运动鞋", "篮球鞋"},
	)
	got := vocab.MatchFromText("男士运动跑步鞋")
	want := NamePath{"服饰", "运动鞋", "跑步鞋"}
	if !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("Path = %v, want %v", got.Path, want)
	}
	if got.MatchKind != MatchContain {
		t.Fatalf("MatchKind = %v, want Contain", got.MatchKind)
	}
}

func TestMatchFromText_LongestWins(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"江苏省", "南京"},
		NamePath{"江苏省", "南京市"},
	)
	got := vocab.MatchFromText("南京市鼓楼区")
	if got.Path[len(got.Path)-1] != "南京市" {
		t.Fatalf("want 南京市, got Path=%v", got.Path)
	}
}

func TestMatchFromText_EmptyVocabulary(t *testing.T) {
	vocab := NewVocabulary(nil)
	got := vocab.MatchFromText("深圳市")
	if got.Matched {
		t.Fatal("empty vocabulary should not match")
	}
}

func TestMatchFromText_SingleRoot(t *testing.T) {
	vocab := NewVocabulary([]VocabNode{{ID: "1", ParentID: "", Name: "江苏省"}})
	got := vocab.MatchFromText("江苏")
	if !got.Matched || !reflect.DeepEqual(got.Path, NamePath{"江苏省"}) {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestNewVocabulary_ParentNotInVocab(t *testing.T) {
	nodes := []VocabNode{
		{ID: "1", ParentID: "99", Name: "孤儿"},
		{ID: "2", ParentID: "", Name: "江苏省"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.MatchFromText("江苏")
	if !got.Matched || got.MatchedNodeID != "2" {
		t.Fatalf("江苏: unexpected %+v", got)
	}
	got = vocab.MatchFromText("孤儿")
	if !got.Matched || got.MatchedNodeID != "1" || !reflect.DeepEqual(got.Path, NamePath{"孤儿"}) {
		t.Fatalf("孤儿: unexpected %+v", got)
	}
}

func TestNewVocabulary_ParentIDZero(t *testing.T) {
	nodes := []VocabNode{
		{ID: "1", ParentID: "0", Name: "广东省"},
		{ID: "2", ParentID: "1", Name: "深圳"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.MatchFromText("深圳市南山区")
	wantPath := NamePath{"广东省", "深圳"}
	if !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("Path = %v, want %v", got.Path, wantPath)
	}
}

func TestNewVocabulary_CycleSkipped(t *testing.T) {
	nodes := []VocabNode{
		{ID: "1", ParentID: "2", Name: "A"},
		{ID: "2", ParentID: "1", Name: "B"},
		{ID: "3", ParentID: "", Name: "江苏省"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.MatchFromText("江苏")
	if !got.Matched || got.MatchedNodeID != "3" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestNewVocabularyFromPaths_SharedPrefix(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"江苏省"},
		NamePath{"江苏省", "南京"},
	)
	got := vocab.MatchFromText("南京市")
	want := NamePath{"江苏省", "南京"}
	if !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("Path = %v, want %v", got.Path, want)
	}
}

func TestNewVocabulary_UUIDIDs(t *testing.T) {
	nodes := []VocabNode{
		{ID: "550e8400-e29b-41d4-a716-446655440000", ParentID: "", Name: "广东省"},
		{ID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", ParentID: "550e8400-e29b-41d4-a716-446655440000", Name: "深圳"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.MatchFromText("深圳市南山区")
	if got.MatchedNodeID != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Fatalf("MatchedNodeID = %q, want UUID", got.MatchedNodeID)
	}
	if !reflect.DeepEqual(got.Path, NamePath{"广东省", "深圳"}) {
		t.Fatalf("Path = %v", got.Path)
	}
}

func TestChainWeights_Sum100(t *testing.T) {
	for _, n := range []int{1, 2, 3, 6, 10} {
		w := chainWeights(n)
		if len(w) != n {
			t.Fatalf("n=%d: len=%d", n, len(w))
		}
		sum := 0
		for i, v := range w {
			sum += v
			if i > 0 && w[i] <= w[i-1] {
				t.Fatalf("n=%d: weights not increasing: %v", n, w)
			}
		}
		if sum != 100 {
			t.Fatalf("n=%d: sum=%d, want 100, weights=%v", n, sum, w)
		}
	}
}

func TestMatch_ChainScoringFullMatch(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"四川省", "成都市", "武侯区"},
		NamePath{"四川省", "绵阳市", "德阳市", "乐山市", "宜宾市", "武侯区"},
	)
	got := vocab.MatchFromText("四川省成都市武侯区")
	wantPath := NamePath{"四川省", "成都市", "武侯区"}
	if !got.Matched || !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("got %+v, want Path %v", got, wantPath)
	}
}

func TestMatch_ChainScoringTieBreak(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"中国", "北京市"},
		NamePath{"美国", "旧金山"},
		NamePath{"旧金山", "德克萨斯州", "中国"},
	)
	got := vocab.MatchFromText("中国旧金山")
	// 链 2/3 同分 67，链 3 更长 → 胜出
	wantPath := NamePath{"旧金山", "德克萨斯州", "中国"}
	if !got.Matched || !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("got %+v, want Path %v", got, wantPath)
	}
}

func TestMatch_AddressPartialMiddle(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"亚洲", "云南省", "昆明市"},
	)
	got := vocab.MatchFromText("云南省昆明市人民路")
	wantPath := NamePath{"亚洲", "云南省", "昆明市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("got %+v, want Path %v", got, wantPath)
	}
}

func TestMatch_1(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"四川省", "成都市", "武侯区"},
		NamePath{"四川省", "绵阳市", "德阳市", "乐山市", "宜宾市", "武侯区"},
	)
	got := vocab.MatchFromText("四川省成都市武侯区")
	logutil.Log(got)
}

func xinjiangVocab() *Vocabulary {
	return NewVocabularyFromPaths(
		NamePath{"新疆维吾尔自治区", "乌鲁木齐市", "天山区"},
	)
}

func tibetVocab() *Vocabulary {
	return NewVocabularyFromPaths(
		NamePath{"西藏自治区", "拉萨市"},
	)
}

func TestMatchFromText_RegionAliases_NoOptsNoMatch(t *testing.T) {
	vocab := NewVocabularyFromPaths(NamePath{"新疆维吾尔自治区"})
	got := vocab.MatchFromText("新疆")
	if got.Matched {
		t.Fatalf("without region aliases 新疆 should not match 新疆维吾尔自治区, got %+v", got)
	}
}

func TestMatchFromText_WithDefaultRegionAliases_Xinjiang(t *testing.T) {
	got := xinjiangVocab().MatchFromText(
		"新疆乌鲁木齐市天山区",
		MatchOpts().WithDefaultRegionAliases(),
	)
	want := NamePath{"新疆维吾尔自治区", "乌鲁木齐市", "天山区"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestMatchFromText_WithDefaultRegionAliases_Tibet(t *testing.T) {
	got := tibetVocab().MatchFromText(
		"西藏拉萨市",
		MatchOpts().WithDefaultRegionAliases(),
	)
	want := NamePath{"西藏自治区", "拉萨市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestMatchFromText_WithDefaultRegionAliases_ExistingCityMatch(t *testing.T) {
	got := regionVocab().MatchFromText(
		"深圳市南山区科技园",
		MatchOpts().WithDefaultRegionAliases(),
	)
	want := NamePath{"广东省", "深圳"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestMatchFromText_NameAliases_CustomOnly(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"新疆维吾尔自治区", "乌鲁木齐市"},
	)
	got := vocab.MatchFromText(
		"新疆乌鲁木齐市",
		MatchOpts().NameAliases(map[string][]string{
			"新疆维吾尔自治区": {"新疆"},
		}),
	)
	want := NamePath{"新疆维吾尔自治区", "乌鲁木齐市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestMatchFromText_NameAliases_CustomOnly2(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"新疆", "乌鲁木齐市"},
	)
	got := vocab.MatchFromText(
		"新疆维吾尔自治区乌鲁木齐市",
		MatchOpts().NameAliases(map[string][]string{
			"新疆维吾尔自治区": {"新疆"},
		}),
	)
	want := NamePath{"新疆", "乌鲁木齐市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestMatchFromText_WithDefaultRegionAliases_CategoryUnaffected(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"服饰", "运动鞋", "跑步鞋"},
		NamePath{"服饰", "运动鞋", "篮球鞋"},
	)
	got := vocab.MatchFromText(
		"男士运动跑步鞋",
		MatchOpts().WithDefaultRegionAliases(),
	)
	want := NamePath{"服饰", "运动鞋", "跑步鞋"}
	if !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("Path = %v, want %v", got.Path, want)
	}
}

func TestMatchFromText_DefaultPlusCustomNameAliases(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		NamePath{"广东省", "深圳市"},
	)
	got := vocab.MatchFromText(
		"鹏城南山区",
		MatchOpts().WithDefaultRegionAliases().NameAliases(map[string][]string{
			"深圳市": {"鹏城"},
		}),
	)
	want := NamePath{"广东省", "深圳市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}
