package fuzzutil

import (
	"reflect"
	"testing"

	"github.com/lontten/lutil/logutil"
)

func regionVocab() *Vocabulary {
	return NewVocabulary(
		NamePath{"广东省", "深圳"},
		NamePath{"上海", "浦东新区"},
		NamePath{"江苏省", "南京"},
		NamePath{"江苏省"},
		NamePath{"四川省", "成都"},
	)
}

func TestMatchResult_LastName_LastID(t *testing.T) {
	if (MatchResult{}).LastName() != "" || (MatchResult{}).LastID() != "" {
		t.Fatal("zero MatchResult")
	}
	got := regionVocab().Match("深圳市南山区科技园")
	if got.LastName() != "深圳" {
		t.Fatalf("LastName() = %q, want 深圳", got.LastName())
	}
	if got.LastID() != "" {
		t.Fatalf("LastID() = %q, want empty for FromPaths", got.LastID())
	}

	nodes := []VocabNode{
		{ID: "1", ParentID: "", Name: "广东省"},
		{ID: "2", ParentID: "1", Name: "深圳"},
	}
	got2 := NewVocabularyFromNodes(nodes, EndpointOpts().AtDepths(2)).Match("深圳市南山区")
	if got2.LastName() != "深圳" || got2.LastID() != "2" {
		t.Fatalf("FromNodes: LastName=%q LastID=%q", got2.LastName(), got2.LastID())
	}
}

func TestNewVocabularyFromNodes(t *testing.T) {
	nodes := []VocabNode{
		{ID: "1", ParentID: "", Name: "广东省"},
		{ID: "2", ParentID: "1", Name: "深圳"},
		{ID: "3", ParentID: "", Name: "江苏省"},
	}
	vocab := NewVocabularyFromNodes(nodes, EndpointOpts().AtDepths(2))
	got := vocab.Match("深圳市南山区")
	wantPath := NamePath{"广东省", "深圳"}
	if !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("Path = %v, want %v", got.Path, wantPath)
	}
	wantPathIDs := IDPath{"1", "2"}
	if !reflect.DeepEqual(got.PathIDs, wantPathIDs) {
		t.Fatalf("PathIDs = %v, want %v", got.PathIDs, wantPathIDs)
	}
	if got.LastID() != "2" {
		t.Fatalf("LastID() = %q, want %q", got.LastID(), "2")
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
	got := vocab.Match("广州市天河区")
	if !got.Matched || got.LastName() != "广州" {
		t.Fatalf("unexpected result: %+v", got)
	}
	if got.LastID() != "" {
		t.Fatalf("LastID() = %q, want empty for auto-assigned IDs", got.LastID())
	}
}

func TestVocabulary_Match_FromPaths_NoNodeID(t *testing.T) {
	got := regionVocab().Match("深圳市南山区科技园")
	if !got.Matched {
		t.Fatal("expected match")
	}
	if got.LastID() != "" {
		t.Fatalf("LastID() = %q, want empty for FromPaths", got.LastID())
	}
	if len(got.Path) == 0 {
		t.Fatal("Path should be non-empty")
	}
	wantPathIDs := IDPath{"", ""}
	if !reflect.DeepEqual(got.PathIDs, wantPathIDs) {
		t.Fatalf("PathIDs = %v, want %v", got.PathIDs, wantPathIDs)
	}
}

func TestVocabulary_Match_FromTree_WithUserID(t *testing.T) {
	vocab := NewVocabularyFromTree(TreeNode{
		ID:   "province-1",
		Name: "广东省",
		Children: []TreeNode{
			{Name: "深圳"},
		},
	})
	got := vocab.Match("深圳市")
	if got.LastID() != "" {
		t.Fatalf("synthetic child ID: LastID() = %q, want empty", got.LastID())
	}
	wantPathIDs := IDPath{"province-1", ""}
	if !reflect.DeepEqual(got.PathIDs, wantPathIDs) {
		t.Fatalf("PathIDs = %v, want %v", got.PathIDs, wantPathIDs)
	}

	vocab2 := NewVocabularyFromTree(TreeNode{
		Name: "广东省",
		Children: []TreeNode{
			{ID: "city-sz", Name: "深圳"},
		},
	})
	got2 := vocab2.Match("深圳市")
	if got2.LastID() != "city-sz" {
		t.Fatalf("user-provided ID: LastID() = %q, want city-sz", got2.LastID())
	}
	wantPathIDs2 := IDPath{"", "city-sz"}
	if !reflect.DeepEqual(got2.PathIDs, wantPathIDs2) {
		t.Fatalf("PathIDs = %v, want %v", got2.PathIDs, wantPathIDs2)
	}
}

func TestVocabulary_Match(t *testing.T) {
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
			got := vocab.Match(tt.text, tt.opts)
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

func TestVocabulary_Match_CategoryThreeLevel(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"服饰", "运动鞋", "跑步鞋"},
		NamePath{"服饰", "运动鞋", "篮球鞋"},
	)
	got := vocab.Match("男士运动跑步鞋")
	want := NamePath{"服饰", "运动鞋", "跑步鞋"}
	if !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("Path = %v, want %v", got.Path, want)
	}
	if got.MatchKind != MatchContain {
		t.Fatalf("MatchKind = %v, want Contain", got.MatchKind)
	}
}

func TestVocabulary_Match_LongestWins(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"江苏省", "南京"},
		NamePath{"江苏省", "南京市"},
	)
	got := vocab.Match("南京市鼓楼区")
	if got.Path[len(got.Path)-1] != "南京市" {
		t.Fatalf("want 南京市, got Path=%v", got.Path)
	}
}

func TestVocabulary_Match_EmptyVocabulary(t *testing.T) {
	vocab := NewVocabulary()
	got := vocab.Match("深圳市")
	if got.Matched {
		t.Fatal("empty vocabulary should not match")
	}
}

func TestVocabulary_Match_SingleRoot(t *testing.T) {
	vocab := NewVocabularyFromNodes(
		[]VocabNode{{ID: "1", ParentID: "", Name: "江苏省"}},
		EndpointOpts().AtDepths(1),
	)
	got := vocab.Match("江苏")
	if !got.Matched || !reflect.DeepEqual(got.Path, NamePath{"江苏省"}) {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestNewVocabulary_ParentNotInVocab(t *testing.T) {
	nodes := []VocabNode{
		{ID: "1", ParentID: "99", Name: "孤儿"},
		{ID: "2", ParentID: "", Name: "江苏省"},
	}
	vocab := NewVocabularyFromNodes(nodes, EndpointOpts().AtDepths(1))
	got := vocab.Match("江苏")
	if !got.Matched || got.LastID() != "2" {
		t.Fatalf("江苏: unexpected %+v", got)
	}
	got = vocab.Match("孤儿")
	if !got.Matched || got.LastID() != "1" || !reflect.DeepEqual(got.Path, NamePath{"孤儿"}) {
		t.Fatalf("孤儿: unexpected %+v", got)
	}
	if !reflect.DeepEqual(got.PathIDs, IDPath{"1"}) {
		t.Fatalf("PathIDs = %v, want %v", got.PathIDs, IDPath{"1"})
	}
}

func TestNewVocabulary_ParentIDZero(t *testing.T) {
	nodes := []VocabNode{
		{ID: "1", ParentID: "0", Name: "广东省"},
		{ID: "2", ParentID: "1", Name: "深圳"},
	}
	vocab := NewVocabularyFromNodes(nodes, EndpointOpts().AtDepths(2))
	got := vocab.Match("深圳市南山区")
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
	vocab := NewVocabularyFromNodes(nodes, EndpointOpts().AtDepths(1))
	got := vocab.Match("江苏")
	if !got.Matched || got.LastID() != "3" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestNewVocabulary_SharedPrefix(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"江苏省"},
		NamePath{"江苏省", "南京"},
	)
	got := vocab.Match("南京市")
	want := NamePath{"江苏省", "南京"}
	if !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("Path = %v, want %v", got.Path, want)
	}
}

func linearABCVocabNodes() []VocabNode {
	return []VocabNode{
		{ID: "1", ParentID: "", Name: "a"},
		{ID: "2", ParentID: "1", Name: "b"},
		{ID: "3", ParentID: "2", Name: "c"},
	}
}

func branchABCvocabNodes() []VocabNode {
	return []VocabNode{
		{ID: "1", ParentID: "", Name: "a"},
		{ID: "2", ParentID: "1", Name: "b"},
		{ID: "3", ParentID: "2", Name: "c"},
		{ID: "4", ParentID: "1", Name: "d"},
	}
}

func vocabChainPaths(v *Vocabulary) []NamePath {
	paths := make([]NamePath, len(v.chains))
	for i, c := range v.chains {
		paths[i] = append(NamePath(nil), c.path...)
	}
	return paths
}

func chainPathsEqual(got, want []NamePath) bool {
	if len(got) != len(want) {
		return false
	}
	seen := make(map[string]int, len(want))
	for _, p := range want {
		seen[pathKey(p)]++
	}
	for _, p := range got {
		k := pathKey(p)
		if seen[k] == 0 {
			return false
		}
		seen[k]--
	}
	return true
}

func pathKey(p NamePath) string {
	key := ""
	for i, s := range p {
		if i > 0 {
			key += "|"
		}
		key += s
	}
	return key
}

func TestEndpointOpts_EmptyAtDepths(t *testing.T) {
	vocab := NewVocabularyFromNodes(linearABCVocabNodes(), EndpointOpts().AtDepths())
	if len(vocab.chains) != 0 {
		t.Fatalf("empty depths: want 0 chains, got %d", len(vocab.chains))
	}
}

func TestEndpointOpts_DefaultLeafOnly(t *testing.T) {
	omitted := NewVocabularyFromNodes(linearABCVocabNodes())
	defaultOpts := NewVocabularyFromNodes(linearABCVocabNodes(), EndpointOpts())
	explicitLeaf := NewVocabularyFromNodes(linearABCVocabNodes(), EndpointOpts().LeafOnly())
	if !chainPathsEqual(vocabChainPaths(omitted), vocabChainPaths(defaultOpts)) {
		t.Fatalf("omitted = %v, default = %v", vocabChainPaths(omitted), vocabChainPaths(defaultOpts))
	}
	if !chainPathsEqual(vocabChainPaths(defaultOpts), vocabChainPaths(explicitLeaf)) {
		t.Fatalf("default = %v, LeafOnly = %v", vocabChainPaths(defaultOpts), vocabChainPaths(explicitLeaf))
	}
	want := []NamePath{{"a", "b", "c"}}
	if !chainPathsEqual(vocabChainPaths(omitted), want) {
		t.Fatalf("chains = %v, want %v", vocabChainPaths(omitted), want)
	}
}

func TestEndpointOpts_NilIsLeafOnly(t *testing.T) {
	nilOpts := NewVocabularyFromNodes(linearABCVocabNodes(), nil)
	omitted := NewVocabularyFromNodes(linearABCVocabNodes())
	if !chainPathsEqual(vocabChainPaths(nilOpts), vocabChainPaths(omitted)) {
		t.Fatalf("nil = %v, omitted = %v", vocabChainPaths(nilOpts), vocabChainPaths(omitted))
	}
}

func TestNewVocabularyFromNodes_LeafEndpoints(t *testing.T) {
	vocab := NewVocabularyFromNodes(linearABCVocabNodes())
	want := []NamePath{{"a", "b", "c"}}
	if !chainPathsEqual(vocabChainPaths(vocab), want) {
		t.Fatalf("chains = %v, want %v", vocabChainPaths(vocab), want)
	}
}

func TestNewVocabularyFromNodes_LeafEndpoints_Branch(t *testing.T) {
	vocab := NewVocabularyFromNodes(branchABCvocabNodes())
	want := []NamePath{
		{"a", "b", "c"},
		{"a", "d"},
	}
	if !chainPathsEqual(vocabChainPaths(vocab), want) {
		t.Fatalf("chains = %v, want %v", vocabChainPaths(vocab), want)
	}
}

func TestNewVocabularyFromNodes_EndpointDepths(t *testing.T) {
	nodes := linearABCVocabNodes()

	vocabDepth2 := NewVocabularyFromNodes(nodes, EndpointOpts().AtDepths(2))
	wantDepth2 := []NamePath{{"a", "b"}}
	if !chainPathsEqual(vocabChainPaths(vocabDepth2), wantDepth2) {
		t.Fatalf("depth 2: chains = %v, want %v", vocabChainPaths(vocabDepth2), wantDepth2)
	}

	vocabDepth13 := NewVocabularyFromNodes(nodes, EndpointOpts().AtDepths(1, 3))
	wantDepth13 := []NamePath{{"a"}, {"a", "b", "c"}}
	if !chainPathsEqual(vocabChainPaths(vocabDepth13), wantDepth13) {
		t.Fatalf("depth 1,3: chains = %v, want %v", vocabChainPaths(vocabDepth13), wantDepth13)
	}
}

func TestNewVocabularyFromTree_LeafEndpoints(t *testing.T) {
	root := TreeNode{
		Name: "a",
		Children: []TreeNode{{
			Name: "b",
			Children: []TreeNode{{
				Name: "c",
			}},
		}},
	}
	fromTree := NewVocabularyFromTree(root)
	fromPaths := NewVocabulary(NamePath{"a", "b", "c"})
	if !chainPathsEqual(vocabChainPaths(fromTree), vocabChainPaths(fromPaths)) {
		t.Fatalf("tree = %v, paths = %v", vocabChainPaths(fromTree), vocabChainPaths(fromPaths))
	}
}

func TestNewVocabularyFromNodes_EndpointDepths_MatchBehavior(t *testing.T) {
	nodes := []VocabNode{
		{ID: "1", ParentID: "", Name: "安徽省"},
		{ID: "2", ParentID: "1", Name: "滁州市"},
		{ID: "3", ParentID: "2", Name: "天长市"},
	}
	vocab := NewVocabularyFromNodes(nodes, EndpointOpts().AtDepths(3))
	got := vocab.Match(
		"安徽省天长市铜城工业园区",
		MatchOpts().WithDefaultRegionAliases(),
	)
	want := NamePath{"安徽省", "滁州市", "天长市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestNewVocabulary_UUIDIDs(t *testing.T) {
	nodes := []VocabNode{
		{ID: "550e8400-e29b-41d4-a716-446655440000", ParentID: "", Name: "广东省"},
		{ID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", ParentID: "550e8400-e29b-41d4-a716-446655440000", Name: "深圳"},
	}
	vocab := NewVocabularyFromNodes(nodes, EndpointOpts().AtDepths(2))
	got := vocab.Match("深圳市南山区")
	if got.LastID() != "6ba7b810-9dad-11d1-80b4-00c04fd430c8" {
		t.Fatalf("LastID() = %q, want UUID", got.LastID())
	}
	wantPathIDs := IDPath{
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
	}
	if !reflect.DeepEqual(got.PathIDs, wantPathIDs) {
		t.Fatalf("PathIDs = %v, want %v", got.PathIDs, wantPathIDs)
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
	vocab := NewVocabulary(
		NamePath{"四川省", "成都市", "武侯区"},
		NamePath{"四川省", "绵阳市", "德阳市", "乐山市", "宜宾市", "武侯区"},
	)
	got := vocab.Match("四川省成都市武侯区")
	wantPath := NamePath{"四川省", "成都市", "武侯区"}
	if !got.Matched || !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("got %+v, want Path %v", got, wantPath)
	}
}

func TestMatch_ChainScoringTieBreak(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"中国", "北京市"},
		NamePath{"美国", "旧金山"},
		NamePath{"旧金山", "德克萨斯州", "中国"},
	)
	got := vocab.Match("中国旧金山")
	// 链 2/3 同分 67，链 3 更长 → 胜出
	wantPath := NamePath{"旧金山", "德克萨斯州", "中国"}
	if !got.Matched || !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("got %+v, want Path %v", got, wantPath)
	}
}

func TestMatch_AddressPartialMiddle(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"亚洲", "云南省", "昆明市"},
	)
	got := vocab.Match("云南省昆明市人民路")
	wantPath := NamePath{"亚洲", "云南省", "昆明市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, wantPath) {
		t.Fatalf("got %+v, want Path %v", got, wantPath)
	}
}

func TestMatch_1(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"四川省", "成都市", "武侯区"},
		NamePath{"四川省", "绵阳市", "德阳市", "乐山市", "宜宾市", "武侯区"},
	)
	got := vocab.Match("四川省成都市武侯区")
	logutil.Log(got)
}

func xinjiangVocab() *Vocabulary {
	return NewVocabulary(
		NamePath{"新疆维吾尔自治区", "乌鲁木齐市", "天山区"},
	)
}

func tibetVocab() *Vocabulary {
	return NewVocabulary(
		NamePath{"西藏自治区", "拉萨市"},
	)
}

func anhuiMultiDepthVocab() *Vocabulary {
	return NewVocabulary(
		NamePath{"安徽省"},
		NamePath{"安徽省", "滁州市"},
		NamePath{"安徽省", "滁州市", "天长市"},
	)
}

func anhuiMultiDepthNodes() []VocabNode {
	return []VocabNode{
		{ID: "1", ParentID: "", Name: "安徽省"},
		{ID: "2", ParentID: "1", Name: "滁州市"},
		{ID: "3", ParentID: "2", Name: "天长市"},
	}
}

func anhuiMultiDepthAlignCases() []struct {
	name     string
	text     string
	matched  bool
	wantPath NamePath
	wantKind MatchKind
} {
	return []struct {
		name     string
		text     string
		matched  bool
		wantPath NamePath
		wantKind MatchKind
	}{
		{
			name:     "skip_middle_level",
			text:     "安徽省天长市铜城工业园区纬三大道一号",
			matched:  true,
			wantPath: NamePath{"安徽省", "滁州市", "天长市"},
			wantKind: MatchContain,
		},
		{
			name:     "province_only",
			text:     "安徽省",
			matched:  true,
			wantPath: NamePath{"安徽省"},
			wantKind: MatchContain,
		},
		{
			name:     "full_hierarchy",
			text:     "安徽省滁州市天长市",
			matched:  true,
			wantPath: NamePath{"安徽省", "滁州市", "天长市"},
			wantKind: MatchContain,
		},
		{
			name:     "county_only",
			text:     "天长市",
			matched:  true,
			wantPath: NamePath{"安徽省", "滁州市", "天长市"},
			wantKind: MatchContain,
		},
		{
			name:     "city_and_county_no_province",
			text:     "滁州市天长市",
			matched:  true,
			wantPath: NamePath{"安徽省", "滁州市", "天长市"},
			wantKind: MatchContain,
		},
		{
			name:     "city_only",
			text:     "滁州市",
			matched:  true,
			wantPath: NamePath{"安徽省", "滁州市"},
			wantKind: MatchContain,
		},
		{
			name:    "no_match",
			text:    "abcxyz",
			matched: false,
		},
	}
}

func runMultiDepthAlignCases(t *testing.T, vocab *Vocabulary, cases []struct {
	name     string
	text     string
	matched  bool
	wantPath NamePath
	wantKind MatchKind
}) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := vocab.Match(tt.text)
			if got.Matched != tt.matched {
				t.Fatalf("Matched = %v, want %v (result=%+v)", got.Matched, tt.matched, got)
			}
			if !tt.matched {
				return
			}
			if !reflect.DeepEqual(got.Path, tt.wantPath) {
				t.Fatalf("Path = %v, want %v", got.Path, tt.wantPath)
			}
			if got.MatchKind != tt.wantKind {
				t.Fatalf("MatchKind = %v, want %v", got.MatchKind, tt.wantKind)
			}
			if tt.name == "skip_middle_level" {
				for _, seg := range got.Path {
					if seg == "" {
						t.Fatal("Path must not contain padded empty segments")
					}
				}
			}
		})
	}
}

func TestAlignChains(t *testing.T) {
	vocab := anhuiMultiDepthVocab()
	if len(vocab.chains) != 3 {
		t.Fatalf("want 3 chains, got %d", len(vocab.chains))
	}

	byDepth := make(map[int]chain, len(vocab.chains))
	for _, c := range vocab.chains {
		byDepth[len(c.path)] = c
	}

	shallow, ok := byDepth[1]
	if !ok {
		t.Fatal("missing depth-1 chain")
	}
	wantShallowAligned := []string{"", "", "安徽省"}
	if !reflect.DeepEqual(shallow.aligned, wantShallowAligned) {
		t.Fatalf("depth 1 aligned = %v, want %v", shallow.aligned, wantShallowAligned)
	}

	mid, ok := byDepth[2]
	if !ok {
		t.Fatal("missing depth-2 chain")
	}
	wantMidAligned := []string{"", "安徽省", "滁州市"}
	if !reflect.DeepEqual(mid.aligned, wantMidAligned) {
		t.Fatalf("depth 2 aligned = %v, want %v", mid.aligned, wantMidAligned)
	}

	deep, ok := byDepth[3]
	if !ok {
		t.Fatal("missing depth-3 chain")
	}
	if !reflect.DeepEqual(deep.aligned, deep.path) {
		t.Fatalf("depth 3 aligned = %v, want path %v", deep.aligned, deep.path)
	}

	baseWeights := vocab.chains[0].weights
	if len(baseWeights) != 3 {
		t.Fatalf("weights len = %d, want 3", len(baseWeights))
	}
	sum := 0
	for _, w := range baseWeights {
		sum += w
	}
	if sum != 100 {
		t.Fatalf("weights sum = %d, want 100", sum)
	}
	for i, c := range vocab.chains[1:] {
		if len(c.weights) != len(baseWeights) {
			t.Fatalf("chain %d weights len = %d, want %d", i+1, len(c.weights), len(baseWeights))
		}
		if !reflect.DeepEqual(c.weights, baseWeights) {
			t.Fatalf("chain %d weights = %v, want shared %v", i+1, c.weights, baseWeights)
		}
	}

	if alignChains(nil) != nil {
		t.Fatal("alignChains(nil) should return nil")
	}
	if got := alignChains([]chain{}); len(got) != 0 {
		t.Fatalf("alignChains(empty) len = %d, want 0", len(got))
	}

	singleDepth := NewVocabulary(NamePath{"江苏省", "南京"})
	if len(singleDepth.chains) != 1 {
		t.Fatalf("single depth vocab: want 1 chain, got %d", len(singleDepth.chains))
	}
	c := singleDepth.chains[0]
	if !reflect.DeepEqual(c.aligned, c.path) {
		t.Fatalf("same-depth aligned = %v, want path %v", c.aligned, c.path)
	}
}

func TestVocabulary_Match_MultiDepthChainAlignment(t *testing.T) {
	runMultiDepthAlignCases(t, anhuiMultiDepthVocab(), anhuiMultiDepthAlignCases())
}

func TestVocabulary_Match_MultiDepthChainAlignment_FromNodes(t *testing.T) {
	vocab := NewVocabularyFromNodes(
		anhuiMultiDepthNodes(),
		EndpointOpts().AtDepths(1, 2, 3),
	)
	runMultiDepthAlignCases(t, vocab, anhuiMultiDepthAlignCases())
}

func TestVocabulary_Match_RegionAliases_NoOptsNoMatch(t *testing.T) {
	vocab := NewVocabulary(NamePath{"新疆维吾尔自治区"})
	got := vocab.Match("新疆")
	if got.Matched {
		t.Fatalf("without region aliases 新疆 should not match 新疆维吾尔自治区, got %+v", got)
	}
}

func TestVocabulary_Match_WithDefaultRegionAliases_Xinjiang(t *testing.T) {
	got := xinjiangVocab().Match(
		"新疆乌鲁木齐市天山区",
		MatchOpts().WithDefaultRegionAliases(),
	)
	want := NamePath{"新疆维吾尔自治区", "乌鲁木齐市", "天山区"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestVocabulary_Match_WithDefaultRegionAliases_Tibet(t *testing.T) {
	got := tibetVocab().Match(
		"西藏拉萨市",
		MatchOpts().WithDefaultRegionAliases(),
	)
	want := NamePath{"西藏自治区", "拉萨市"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestVocabulary_Match_WithDefaultRegionAliases_ExistingCityMatch(t *testing.T) {
	got := regionVocab().Match(
		"深圳市南山区科技园",
		MatchOpts().WithDefaultRegionAliases(),
	)
	want := NamePath{"广东省", "深圳"}
	if !got.Matched || !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("got %+v, want Path %v", got, want)
	}
}

func TestVocabulary_Match_NameAliases_CustomOnly(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"新疆维吾尔自治区", "乌鲁木齐市"},
	)
	got := vocab.Match(
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

func TestVocabulary_Match_NameAliases_CustomOnly2(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"新疆", "乌鲁木齐市"},
	)
	got := vocab.Match(
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

func TestVocabulary_Match_WithDefaultRegionAliases_CategoryUnaffected(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"服饰", "运动鞋", "跑步鞋"},
		NamePath{"服饰", "运动鞋", "篮球鞋"},
	)
	got := vocab.Match(
		"男士运动跑步鞋",
		MatchOpts().WithDefaultRegionAliases(),
	)
	want := NamePath{"服饰", "运动鞋", "跑步鞋"}
	if !reflect.DeepEqual(got.Path, want) {
		t.Fatalf("Path = %v, want %v", got.Path, want)
	}
}

func TestVocabulary_Match_DefaultPlusCustomNameAliases(t *testing.T) {
	vocab := NewVocabulary(
		NamePath{"广东省", "深圳市"},
	)
	got := vocab.Match(
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
