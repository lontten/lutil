package fuzzutil

type endpointFilterMode int

const (
	endpointFilterAll endpointFilterMode = iota
	endpointFilterLeafOnly
	endpointFilterDepths
)

type vocabBuildOpts struct {
	mode     endpointFilterMode
	depthSet map[int]struct{}
}

// VocabBuildOption 配置 NewVocabulary / NewVocabularyFromTree 的链终点过滤规则。
type VocabBuildOption func(*vocabBuildOpts)

// WithLeafEndpointsOnly 仅保留叶子 node 对应的链作为匹配终点。
// 与 NewVocabularyFromPaths 只传 {a,b,c} 一条 path 时的链集合类似（无分支时等价）。
func WithLeafEndpointsOnly() VocabBuildOption {
	return func(o *vocabBuildOpts) {
		o.mode = endpointFilterLeafOnly
	}
}

// WithEndpointDepths 仅保留 path 深度（len(path)，根为 1）在集合内的链。
// depths 为空时忽略，等同默认（全部 node 为终点）。
func WithEndpointDepths(depths ...int) VocabBuildOption {
	return func(o *vocabBuildOpts) {
		if len(depths) == 0 {
			return
		}
		o.mode = endpointFilterDepths
		o.depthSet = make(map[int]struct{}, len(depths))
		for _, d := range depths {
			o.depthSet[d] = struct{}{}
		}
	}
}

func applyVocabBuildOpts(opts []VocabBuildOption) vocabBuildOpts {
	merged := vocabBuildOpts{mode: endpointFilterAll}
	for _, opt := range opts {
		if opt != nil {
			opt(&merged)
		}
	}
	return merged
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

func filterChainsByEndpoint(chains []chain, nodes []VocabNode, opts vocabBuildOpts) []chain {
	switch opts.mode {
	case endpointFilterAll:
		return chains
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
		filtered := make([]chain, 0, len(chains))
		for _, c := range chains {
			if _, ok := opts.depthSet[len(c.path)]; ok {
				filtered = append(filtered, c)
			}
		}
		return filtered
	default:
		return chains
	}
}
