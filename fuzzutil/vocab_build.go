package fuzzutil

type endpointFilterMode int

const (
	endpointFilterLeafOnly endpointFilterMode = iota + 1
	endpointFilterDepths
)

// endpointOpts 配置 Node/Tree 词表的链终点过滤；通过 EndpointOpts() 创建后链式设置。
// 深度 depth = len(path)，根为 1。
type endpointOpts struct {
	mode     endpointFilterMode
	depthSet map[int]struct{}
}

// EndpointOpts 返回默认 LeafOnly 的终点规则；需按层深过滤时链式调用 AtDepths。
func EndpointOpts() *endpointOpts {
	return &endpointOpts{mode: endpointFilterLeafOnly}
}

func resolveEndpointOpts(opts ...*endpointOpts) *endpointOpts {
	if len(opts) == 0 || opts[0] == nil {
		return EndpointOpts()
	}
	return opts[0]
}

// LeafOnly 仅叶子 node 为终点（词表中无其他 node 以该 ID 为 ParentID）。
// 无分支时等价于 NewVocabulary(NamePath{"a", "b", "c"}) 只传一条 path。
func (o *endpointOpts) LeafOnly() *endpointOpts {
	o.mode = endpointFilterLeafOnly
	return o
}

// AtDepths 仅 len(path) 在 depths 中的链为终点；depths 为空时词表 0 链。
func (o *endpointOpts) AtDepths(depths ...int) *endpointOpts {
	o.mode = endpointFilterDepths
	o.depthSet = nil
	if len(depths) == 0 {
		return o
	}
	o.depthSet = make(map[int]struct{}, len(depths))
	for _, d := range depths {
		o.depthSet[d] = struct{}{}
	}
	return o
}

func buildAllChains(nodes []VocabNode) []chain {
	if nodes == nil {
		nodes = []VocabNode{}
	}

	byID := make(map[string]VocabNode, len(nodes))
	for _, n := range nodes {
		byID[n.ID] = n
	}

	chains := make([]chain, 0, len(nodes))
	for _, n := range nodes {
		path, ok := buildPath(n.ID, byID)
		if !ok {
			continue
		}
		pathCopy := make([]string, len(path))
		copy(pathCopy, path)
		chains = append(chains, chain{
			id:      n.ID,
			path:    pathCopy,
			weights: chainWeights(len(path)),
		})
	}
	return chains
}

func leafNodeIDs(nodes []VocabNode) map[string]bool {
	hasChild := make(map[string]bool, len(nodes))
	for _, n := range nodes {
		if n.ParentID != "" {
			hasChild[n.ParentID] = true
		}
	}
	leaves := make(map[string]bool)
	for _, n := range nodes {
		if !hasChild[n.ID] {
			leaves[n.ID] = true
		}
	}
	return leaves
}

func filterChainsByEndpoint(chains []chain, nodes []VocabNode, opts *endpointOpts) []chain {
	if opts == nil {
		return nil
	}
	switch opts.mode {
	case endpointFilterLeafOnly:
		leaves := leafNodeIDs(nodes)
		filtered := make([]chain, 0, len(chains))
		for _, c := range chains {
			if leaves[c.id] {
				filtered = append(filtered, c)
			}
		}
		return filtered
	case endpointFilterDepths:
		if len(opts.depthSet) == 0 {
			return nil
		}
		filtered := make([]chain, 0, len(chains))
		for _, c := range chains {
			if _, ok := opts.depthSet[len(c.path)]; ok {
				filtered = append(filtered, c)
			}
		}
		return filtered
	default:
		return nil
	}
}

func filterChainsByEndpointIDs(chains []chain, endpointIDs map[string]bool) []chain {
	filtered := make([]chain, 0, len(chains))
	for _, c := range chains {
		if endpointIDs[c.id] {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// alignChains 将词表内各链右对齐到 maxDepth（前补 ""），并统一 weights。
func alignChains(chains []chain) []chain {
	if len(chains) == 0 {
		return chains
	}

	maxDepth := 0
	for _, c := range chains {
		if len(c.path) > maxDepth {
			maxDepth = len(c.path)
		}
	}

	weights := chainWeights(maxDepth)
	out := make([]chain, len(chains))
	for i, c := range chains {
		pad := maxDepth - len(c.path)
		aligned := make([]string, maxDepth)
		copy(aligned[pad:], c.path)
		out[i] = chain{
			id:      c.id,
			path:    c.path,
			aligned: aligned,
			weights: weights,
		}
	}
	return out
}
