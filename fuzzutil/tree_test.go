package fuzzutil

import (
	"reflect"
	"testing"
)

func regionVocab() *Vocabulary {
	return NewVocabularyFromPaths(
		Path{"广东省", "深圳"},
		Path{"上海", "浦东新区"},
		Path{"江苏省", "南京"},
		Path{"江苏省"},
		Path{"四川省", "成都"},
	)
}

func TestNewVocabularyFromNodes(t *testing.T) {
	nodes := []Node{
		{ID: "1", ParentID: "", Name: "广东省"},
		{ID: "2", ParentID: "1", Name: "深圳"},
		{ID: "3", ParentID: "", Name: "江苏省"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.ExtractFromText("深圳市南山区")
	wantPath := []string{"广东省", "深圳"}
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
	got := vocab.ExtractFromText("广州市天河区")
	if !got.Matched || got.Path[len(got.Path)-1] != "广州" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestExtractFromText(t *testing.T) {
	vocab := regionVocab()

	tests := []struct {
		name     string
		text     string
		matched  bool
		wantPath []string
		wantKind MatchKind
		opts     []ExtractOption
	}{
		{
			name:     "省市-市命中",
			text:     "深圳市南山区科技园",
			matched:  true,
			wantPath: []string{"广东省", "深圳"},
			wantKind: MatchContain,
		},
		{
			name:     "省市-区命中",
			text:     "沪上海市浦东新区",
			matched:  true,
			wantPath: []string{"上海", "浦东新区"},
			wantKind: MatchContain,
		},
		{
			name:     "省市-仅省-模糊",
			text:     "江苏",
			matched:  true,
			wantPath: []string{"江苏省"},
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
			opts:    []ExtractOption{WithMinMatchLen(2)},
		},
		{
			name:    "MaxEditDistance=0 无包含则不命中",
			text:    "江苏",
			matched: false,
			opts:    []ExtractOption{WithMaxEditDistance(0)},
		},
		{
			name:     "成都",
			text:     "成都市金牛区金府路111号",
			matched:  true,
			wantPath: []string{"四川省", "成都"},
			wantKind: MatchContain,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vocab.ExtractFromText(tt.text, tt.opts...)
			if got.Matched != tt.matched {
				t.Fatalf("Matched = %v, want %v (result=%+v)", got.Matched, tt.matched, got)
			}
			if !tt.matched {
				return
			}
			if !reflect.DeepEqual(got.Path, tt.wantPath) {
				t.Errorf("Path = %v, want %v", got.Path, tt.wantPath)
			}
			if got.Kind != tt.wantKind {
				t.Errorf("Kind = %v, want %v", got.Kind, tt.wantKind)
			}
		})
	}
}

func TestExtractFromText_CategoryThreeLevel(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		Path{"服饰", "运动鞋", "跑步鞋"},
		Path{"服饰", "运动鞋", "篮球鞋"},
	)
	got := vocab.ExtractFromText("男士运动跑步鞋")
	want := []string{"服饰", "运动鞋", "跑步鞋"}
	if !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("Path = %v, want %v", got.Path, want)
	}
	if got.Kind != MatchContain {
		t.Fatalf("Kind = %v, want Contain", got.Kind)
	}
}

func TestExtractFromText_LongestWins(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		Path{"江苏省", "南京"},
		Path{"江苏省", "南京市"},
	)
	got := vocab.ExtractFromText("南京市鼓楼区")
	if got.Path[len(got.Path)-1] != "南京市" {
		t.Fatalf("want 南京市, got Path=%v", got.Path)
	}
}

func TestExtractFromText_EmptyVocabulary(t *testing.T) {
	vocab := NewVocabulary(nil)
	got := vocab.ExtractFromText("深圳市")
	if got.Matched {
		t.Fatal("empty vocabulary should not match")
	}
}

func TestExtractFromText_SingleRoot(t *testing.T) {
	vocab := NewVocabulary([]Node{{ID: "1", ParentID: "", Name: "江苏省"}})
	got := vocab.ExtractFromText("江苏")
	if !got.Matched || !reflect.DeepEqual(got.Path, []string{"江苏省"}) {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestExtractResult_Ancestors(t *testing.T) {
	r := ExtractResult{Path: []string{"广东省", "深圳"}}
	if !reflect.DeepEqual(r.Ancestors(), r.Path) {
		t.Fatal("Ancestors should equal Path")
	}
}

func TestNewVocabulary_ParentNotInVocab(t *testing.T) {
	nodes := []Node{
		{ID: "1", ParentID: "99", Name: "孤儿"},
		{ID: "2", ParentID: "", Name: "江苏省"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.ExtractFromText("江苏")
	if !got.Matched || got.MatchedNodeID != "2" {
		t.Fatalf("江苏: unexpected %+v", got)
	}
	got = vocab.ExtractFromText("孤儿")
	if !got.Matched || got.MatchedNodeID != "1" || !reflect.DeepEqual(got.Path, []string{"孤儿"}) {
		t.Fatalf("孤儿: unexpected %+v", got)
	}
}

func TestNewVocabulary_ParentIDZero(t *testing.T) {
	nodes := []Node{
		{ID: "1", ParentID: "0", Name: "广东省"},
		{ID: "2", ParentID: "1", Name: "深圳"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.ExtractFromText("深圳市南山区")
	wantPath := []string{"广东省", "深圳"}
	if !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("Path = %v, want %v", got.Path, wantPath)
	}
}

func TestNewVocabulary_CycleSkipped(t *testing.T) {
	nodes := []Node{
		{ID: "1", ParentID: "2", Name: "A"},
		{ID: "2", ParentID: "1", Name: "B"},
		{ID: "3", ParentID: "", Name: "江苏省"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.ExtractFromText("江苏")
	if !got.Matched || got.MatchedNodeID != "3" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestNewVocabularyFromPaths_SharedPrefix(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		Path{"江苏省"},
		Path{"江苏省", "南京"},
	)
	got := vocab.ExtractFromText("南京市")
	want := []string{"江苏省", "南京"}
	if !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("Path = %v, want %v", got.Path, want)
	}
}

func TestNewVocabulary_UUIDIDs(t *testing.T) {
	nodes := []Node{
		{ID: "550e8400-e29b-41d4-a716-446655440000", ParentID: "", Name: "广东省"},
		{ID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", ParentID: "550e8400-e29b-41d4-a716-446655440000", Name: "深圳"},
	}
	vocab := NewVocabulary(nodes)
	got := vocab.ExtractFromText("深圳市南山区")
	if got.MatchedNodeID != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Fatalf("MatchedNodeID = %q, want UUID", got.MatchedNodeID)
	}
	if !reflect.DeepEqual(got.Path, []string{"广东省", "深圳"}) {
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

func TestExtract_ChainScoringFullMatch(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		Path{"四川省", "成都市", "武侯区"},
		Path{"四川省", "绵阳市", "德阳市", "乐山市", "宜宾市", "武侯区"},
	)
	got := vocab.ExtractFromText("四川省成都市武侯区")
	wantPath := []string{"四川省", "成都市", "武侯区"}
	if !got.Matched || !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("got %+v, want Path %v", got, wantPath)
	}
}

func TestExtract_ChainScoringTieBreak(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		Path{"中国", "北京市"},
		Path{"美国", "旧金山"},
		Path{"旧金山", "德克萨斯州", "中国"},
	)
	got := vocab.ExtractFromText("中国旧金山")
	// 链 2/3 同分 67，链 3 更长 → 胜出
	wantPath := []string{"旧金山", "德克萨斯州", "中国"}
	if !got.Matched || !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("got %+v, want Path %v", got, wantPath)
	}
}

func TestExtract_AddressPartialMiddle(t *testing.T) {
	vocab := NewVocabularyFromPaths(
		Path{"亚洲", "云南省", "昆明市"},
	)
	got := vocab.ExtractFromText("云南省昆明市人民路")
	wantPath := []string{"亚洲", "云南省", "昆明市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("got %+v, want Path %v", got, wantPath)
	}
}
