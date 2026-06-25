package fuzzutil

import "strconv"

// VocabNode 表示词表中的一个节点，字段与 DB 表行一一对应：
// id → ID，parent_id → ParentID，name → Name。
// ParentID 指向上级节点 ID；词表中不存在该 ID 时，构建关系链在此停止。
// ID 使用 string，兼容数字主键（"1"）与未来 UUID。
type VocabNode struct {
	ID       string
	ParentID string
	Name     string
}

// NamePath 是一条关系链的 name 序列（链顶 → 终点），例如 {"广东省", "深圳"}。
type NamePath []string

// TreeNode 是嵌套树结构的节点，用于 JSON 配置等场景。
// ID 为空字符串时由 NewVocabularyFromTree 自动分配合成 ID；
// 自动分配的 ID 不会出现在 MatchResult.MatchedNodeID 中。
type TreeNode struct {
	ID       string
	Name     string
	Children []TreeNode
}

// MatchKind 的 String 便于测试与日志输出。
func (k MatchKind) String() string {
	switch k {
	case MatchContain:
		return "Contain"
	case MatchFuzzy:
		return "Fuzzy"
	default:
		return "None"
	}
}

// MatchResult 是 MatchFromText 的返回结果。
// Matched 为 false 时，MatchedNodeID、MatchKind、Path 均为零值。
type MatchResult struct {
	Matched bool
	// MatchedNodeID 为命中终点节点的 ID，是否非空取决于词表构建方式：
	// NewVocabulary(VocabNode) 时为业务主键；NewVocabularyFromPaths 时始终为空（请用 Path）；
	// NewVocabularyFromTree 时仅当终点 ID 由调用方显式提供时非空。
	MatchedNodeID string
	MatchKind     MatchKind
	Path          NamePath
}

// chain 是一条可匹配的关系链及其预计算权重（权重之和为 100）。
type chain struct {
	id      string
	path    []string
	weights []int
}

// Vocabulary 是预编译的关系链词表，初始化一次后可反复调用 MatchFromText。
type Vocabulary struct {
	chains       []chain
	exposableIDs map[string]bool // nil 表示全部可暴露；非 nil 时仅 map 中的 ID 可写入 MatchedNodeID
}

func (v *Vocabulary) matchedNodeID(chainID string) string {
	if v.exposableIDs != nil && !v.exposableIDs[chainID] {
		return ""
	}
	return chainID
}

// matchOpts 配置 MatchFromText 的匹配规则；通过 MatchOpts() 获取默认值后链式修改。
type matchOpts struct {
	minMatchLen        int
	minOverlap         int
	maxEditDistance    int
	nameAliases        map[string][]string
	stripAdminSuffixes bool
}

// aliasConfig 控制 scoreChain 是否对节点名启用别名匹配。
type aliasConfig struct {
	nameAliases        map[string][]string
	stripAdminSuffixes bool
}

func (a aliasConfig) enabled() bool {
	return a.nameAliases != nil || a.stripAdminSuffixes
}

// MatchOpts 返回默认匹配配置：MinMatchLen=2，MinOverlap=2，MaxEditDistance=1。
func MatchOpts() *matchOpts {
	return &matchOpts{
		minMatchLen:     2,
		minOverlap:      2,
		maxEditDistance: 1,
	}
}

// MinMatchLen 设置候选词最少 rune 数，低于此长度不算命中。默认 2。
func (o *matchOpts) MinMatchLen(n int) *matchOpts {
	o.minMatchLen = n
	return o
}

// MaxEditDistance 设置允许的最大编辑距离；0 表示仅子串包含。默认 1。
func (o *matchOpts) MaxEditDistance(n int) *matchOpts {
	o.maxEditDistance = n
	return o
}

// MinOverlap 设置 text 与候选至少相同的 rune 数（多重集，不要求连续）。默认 2。
func (o *matchOpts) MinOverlap(n int) *matchOpts {
	o.minOverlap = n
	return o
}

// WithDefaultRegionAliases 合并内建中国行政区划别名，并启用后缀剥离（省市区县自治区等）。
func (o *matchOpts) WithDefaultRegionAliases() *matchOpts {
	o.nameAliases = mergeAliasMap(o.nameAliases, DefaultRegionNameAliases())
	o.stripAdminSuffixes = true
	return o
}

// NameAliases 合并自定义节点名 → 额外匹配项映射；nil 时忽略。
// 与 WithDefaultRegionAliases 联用时，同 key 下自定义项追加在默认项之后并去重。
func (o *matchOpts) NameAliases(m map[string][]string) *matchOpts {
	if m != nil {
		o.nameAliases = mergeAliasMap(o.nameAliases, m)
	}
	return o
}

func (o *matchOpts) aliasConfig() aliasConfig {
	return aliasConfig{
		nameAliases:        o.nameAliases,
		stripAdminSuffixes: o.stripAdminSuffixes,
	}
}

func matchCandidatesForName(name string, aliases aliasConfig) []string {
	if !aliases.enabled() {
		return []string{name}
	}

	candidates := []string{name}
	if aliases.nameAliases != nil {
		if extra, ok := aliases.nameAliases[name]; ok {
			candidates = append(candidates, extra...)
		}
	}
	if aliases.stripAdminSuffixes {
		candidates = append(candidates, adminSuffixAliases(name)...)
	}
	candidates = dedupeStrings(candidates)
	sortNamesByRuneDesc(candidates)
	return candidates
}

// NewVocabulary 从 DB 扁平节点列表构建词表。
// 默认每个 node 都是一条关系链的终点；成环的 node 会被跳过。
// 可选 WithLeafEndpointsOnly / WithEndpointDepths 限制哪些 node 成为链终点，
// 以便与 NewVocabularyFromPaths 的显式 path 语义对齐。
func NewVocabulary(nodes []VocabNode, opts ...VocabBuildOption) *Vocabulary {
	buildOpts := applyVocabBuildOpts(opts)
	if nodes == nil {
		nodes = []VocabNode{}
	}
	all := buildAllChains(nodes)
	return &Vocabulary{chains: filterChainsByEndpoint(all, nodes, buildOpts)}
}

// buildPath 沿 ParentID 向上追溯，返回链顶→终点的 name 链；仅成环返回 false。
func buildPath(id string, byID map[string]VocabNode) ([]string, bool) {
	var names []string
	visited := make(map[string]bool)

	for {
		if visited[id] {
			return nil, false
		}
		visited[id] = true

		n, ok := byID[id]
		if !ok {
			return nil, false
		}
		names = append(names, n.Name)

		parentID := n.ParentID
		if _, exists := byID[parentID]; !exists {
			break
		}
		id = parentID
	}

	for i, j := 0, len(names)-1; i < j; i, j = i+1, j-1 {
		names[i], names[j] = names[j], names[i]
	}
	return names, true
}

// chainWeights 返回长度为 n 的权重切片，链内权重之和为 100，链尾最重。
func chainWeights(n int) []int {
	if n <= 0 {
		return nil
	}
	if n == 1 {
		return []int{100}
	}
	sumRaw := n * (n + 1) / 2
	weights := make([]int, n)
	allocated := 0
	for i := 0; i < n; i++ {
		weights[i] = (i + 1) * 100 / sumRaw
		allocated += weights[i]
	}
	weights[n-1] += 100 - allocated
	return weights
}

// chainScoreResult 是一条链的计分结果。
type chainScoreResult struct {
	total       int
	tailKind    MatchKind
	tailMatched bool
	kind        MatchKind // 返回用：链尾命中方式，否则链内最高分节点的方式
}

// scoreChain 对一条关系链逐节点计分。
func scoreChain(text string, path []string, weights []int, rules matchRules, aliases aliasConfig) chainScoreResult {
	var res chainScoreResult
	bestNodeScore := 0

	for i, name := range path {
		candidates := matchCandidatesForName(name, aliases)
		_, kind, ok := matchBest(text, candidates, rules)
		if !ok {
			continue
		}
		pts := weights[i]
		if kind == MatchFuzzy {
			pts = weights[i] / 2
		}
		res.total += pts

		if pts > bestNodeScore {
			bestNodeScore = pts
			res.kind = kind
		}

		if i == len(path)-1 {
			res.tailMatched = true
			res.tailKind = kind
		}
	}

	if res.tailMatched {
		res.kind = res.tailKind
	}

	return res
}

func kindRank(k MatchKind) int {
	switch k {
	case MatchContain:
		return 2
	case MatchFuzzy:
		return 1
	default:
		return 0
	}
}

// betterChain 判断候选链是否优于当前最佳链（用于同分决胜）。
func betterChain(total int, tailKind MatchKind, pathLen, tailRuneLen int, bestTotal int, bestTailKind MatchKind, bestPathLen, bestTailRuneLen int) bool {
	if total != bestTotal {
		return total > bestTotal
	}
	if kindRank(tailKind) != kindRank(bestTailKind) {
		return kindRank(tailKind) > kindRank(bestTailKind)
	}
	if pathLen != bestPathLen {
		return pathLen > bestPathLen
	}
	return tailRuneLen > bestTailRuneLen
}

// NewVocabularyFromPaths 从路径列表构建词表（测试与手工配置用）。
// 每条 NamePath 对应一条关系链（终点为该路径最后一段）；共享前缀自动合并节点。
func NewVocabularyFromPaths(paths ...NamePath) *Vocabulary {
	var nodes []VocabNode
	nextID := int64(1)
	seen := make(map[string]string)
	endpointIDs := make([]string, 0, len(paths))

	for _, p := range paths {
		parentID := ""
		var lastID string
		for _, seg := range p {
			key := nodeKey(parentID, seg)
			id, exists := seen[key]
			if !exists {
				id = strconv.FormatInt(nextID, 10)
				nextID++
				seen[key] = id
				nodes = append(nodes, VocabNode{
					ID:       id,
					ParentID: parentID,
					Name:     seg,
				})
			}
			lastID = id
			parentID = id
		}
		if lastID != "" {
			endpointIDs = append(endpointIDs, lastID)
		}
	}

	full := NewVocabulary(nodes)
	endpoints := make(map[string]bool, len(endpointIDs))
	for _, id := range endpointIDs {
		endpoints[id] = true
	}
	chains := make([]chain, 0, len(endpointIDs))
	for _, c := range full.chains {
		if endpoints[c.id] {
			chains = append(chains, c)
		}
	}
	return &Vocabulary{chains: chains, exposableIDs: map[string]bool{}}
}

func nodeKey(parentID, name string) string {
	return parentID + ":" + name
}

// NewVocabularyFromTree 从嵌套树构建词表。
// 默认每个 node 都是链终点；可选 WithLeafEndpointsOnly / WithEndpointDepths 过滤终点，
// 使行为更接近 NewVocabularyFromPaths 只传叶子 path 的场景。
func NewVocabularyFromTree(root TreeNode, opts ...VocabBuildOption) *Vocabulary {
	return newVocabularyFromTree(opts, root)
}

func newVocabularyFromTree(opts []VocabBuildOption, roots ...TreeNode) *Vocabulary {
	var nodes []VocabNode
	exposable := make(map[string]bool)
	nextID := int64(1)
	var walk func(children []TreeNode, parentID string)
	walk = func(children []TreeNode, parentID string) {
		for _, child := range children {
			id := child.ID
			if id == "" {
				id = strconv.FormatInt(nextID, 10)
				nextID++
			} else {
				bumpNextID(&nextID, id)
				exposable[id] = true
			}
			nodes = append(nodes, VocabNode{
				ID:       id,
				ParentID: parentID,
				Name:     child.Name,
			})
			if len(child.Children) > 0 {
				walk(child.Children, id)
			}
		}
	}
	walk(roots, "")
	buildOpts := applyVocabBuildOpts(opts)
	all := buildAllChains(nodes)
	return &Vocabulary{
		chains:       filterChainsByEndpoint(all, nodes, buildOpts),
		exposableIDs: exposable,
	}
}

// bumpNextID 在用户提供数字字符串 ID 时，推进自增计数器避免冲突。
func bumpNextID(nextID *int64, id string) {
	if n, err := strconv.ParseInt(id, 10, 64); err == nil && n >= *nextID {
		*nextID = n + 1
	}
}

// MatchFromText 从 text 中匹配得分最高的关系链终点节点。
// 对每条链的每个节点独立匹配并加权求和（链内权重之和为 100，链尾最重）。
// 默认 MinMatchLen=2，MinOverlap=2，MaxEditDistance=1。
func (v *Vocabulary) MatchFromText(text string, opts ...*matchOpts) MatchResult {
	var o *matchOpts
	if len(opts) > 0 {
		o = opts[0]
	}
	if o == nil {
		o = MatchOpts()
	}

	rules := matchRules{
		minMatchLen:     o.minMatchLen,
		minOverlap:      o.minOverlap,
		maxEditDistance: o.maxEditDistance,
	}
	aliases := o.aliasConfig()

	var (
		found         bool
		bestTotal     int
		bestTailKind  MatchKind
		bestPathLen   int
		bestTailRunes int
		bestChain     chain
		bestKind      MatchKind
	)

	for _, c := range v.chains {
		scored := scoreChain(text, c.path, c.weights, rules, aliases)
		if scored.total == 0 {
			continue
		}

		tailRunes := 0
		if len(c.path) > 0 {
			tailRunes = len([]rune(c.path[len(c.path)-1]))
		}

		if !found || betterChain(scored.total, scored.tailKind, len(c.path), tailRunes, bestTotal, bestTailKind, bestPathLen, bestTailRunes) {
			found = true
			bestTotal = scored.total
			bestTailKind = scored.tailKind
			bestPathLen = len(c.path)
			bestTailRunes = tailRunes
			bestChain = c
			bestKind = scored.kind
		}
	}

	if !found {
		return MatchResult{MatchKind: MatchNone}
	}

	path := make(NamePath, len(bestChain.path))
	copy(path, bestChain.path)
	return MatchResult{
		Matched:       true,
		MatchedNodeID: v.matchedNodeID(bestChain.id),
		MatchKind:     bestKind,
		Path:          path,
	}
}
