package listutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type treeNode struct {
	ID       int
	ParentID int
	Name     string
	Depth    int
	Children []*treeNode
}

func TestTreeBuilder_ToTree_Basic(t *testing.T) {
	as := assert.New(t)

	nodes := []*treeNode{
		{ID: 1, ParentID: 0, Name: "root"},
		{ID: 2, ParentID: 1, Name: "child-a"},
		{ID: 3, ParentID: 1, Name: "child-b"},
		{ID: 4, ParentID: 2, Name: "grandchild"},
	}

	roots := NewTreeBuilder[treeNode]().
		IsParentNode(func(a, b treeNode) bool {
			return a.ID == b.ParentID
		}).
		ToTree(nodes)

	require.Len(t, roots, 1)
	as.Equal(1, roots[0].ID)
	require.Len(t, roots[0].Children, 2)
	as.Equal(2, roots[0].Children[0].ID)
	as.Equal(3, roots[0].Children[1].ID)
	require.Len(t, roots[0].Children[0].Children, 1)
	as.Equal(4, roots[0].Children[0].Children[0].ID)
}

func TestTreeBuilder_ToTree_IsRootNode(t *testing.T) {
	as := assert.New(t)

	nodes := []*treeNode{
		{ID: 1, ParentID: 0, Name: "a"},
		{ID: 2, ParentID: 1, Name: "b"},
		{ID: 3, ParentID: 0, Name: "orphan-as-root"},
	}

	roots := NewTreeBuilder[treeNode]().
		IsParentNode(func(a, b treeNode) bool {
			return a.ID == b.ParentID
		}).
		IsRootNode(func(n treeNode) bool {
			return n.ParentID == 0
		}).
		ToTree(nodes)

	require.Len(t, roots, 2)
	as.Equal(1, roots[0].ID)
	as.Equal(3, roots[1].ID)
	require.Len(t, roots[0].Children, 1)
	as.Equal(2, roots[0].Children[0].ID)
	as.Empty(roots[1].Children)
}

func TestTreeBuilder_ToTree_SortAndTransforms(t *testing.T) {
	as := assert.New(t)

	nodes := []*treeNode{
		{ID: 1, ParentID: 0, Name: "root"},
		{ID: 3, ParentID: 1, Name: "c"},
		{ID: 2, ParentID: 1, Name: "b"},
	}

	roots := NewTreeBuilder[treeNode]().
		IsParentNode(func(a, b treeNode) bool {
			return a.ID == b.ParentID
		}).
		SortFun(func(a, b treeNode) int {
			return a.ID - b.ID
		}).
		TransformGlobal(func(n *treeNode) {
			n.Name = "g-" + n.Name
		}).
		TransformRoot(func(n *treeNode) {
			n.Name = "r-" + n.Name
		}).
		TransformNode(func(n *treeNode, parent treeNode) {
			n.Depth = parent.ID
		}).
		ToTree(nodes)

	require.Len(t, roots, 1)
	as.Equal("r-g-root", roots[0].Name)
	require.Len(t, roots[0].Children, 2)
	as.Equal(2, roots[0].Children[0].ID)
	as.Equal(3, roots[0].Children[1].ID)
	as.Equal("g-b", roots[0].Children[0].Name)
	as.Equal(1, roots[0].Children[0].Depth)
}

func TestTreeBuilder_ToTree_MaxDepth(t *testing.T) {
	as := assert.New(t)

	nodes := []*treeNode{
		{ID: 1, ParentID: 0, Name: "root"},
		{ID: 2, ParentID: 1, Name: "child"},
		{ID: 3, ParentID: 2, Name: "grandchild"},
	}

	roots := NewTreeBuilder[treeNode]().
		MaxDepth(1).
		IsParentNode(func(a, b treeNode) bool {
			return a.ID == b.ParentID
		}).
		ToTree(nodes)

	require.Len(t, roots, 1)
	require.Len(t, roots[0].Children, 1)
	as.Equal(2, roots[0].Children[0].ID)
	as.Empty(roots[0].Children[0].Children)
}

func TestTreeBuilder_ToTree_SetRootNode(t *testing.T) {
	as := assert.New(t)

	nodes := []*treeNode{
		{ID: 1, ParentID: 0, Name: "a"},
		{ID: 2, ParentID: 1, Name: "b"},
	}

	roots := NewTreeBuilder[treeNode]().
		IsParentNode(func(a, b treeNode) bool {
			return a.ID == b.ParentID
		}).
		SetRootNode(treeNode{ID: 0, Name: "virtual-root"}).
		ToTree(nodes)

	require.Len(t, roots, 1)
	as.Equal(0, roots[0].ID)
	as.Equal("virtual-root", roots[0].Name)
	require.Len(t, roots[0].Children, 1)
	as.Equal(1, roots[0].Children[0].ID)
}

func TestTreeBuilder_ToTree_PanicWithoutIsParentNode(t *testing.T) {
	as := assert.New(t)
	as.Panics(func() {
		NewTreeBuilder[treeNode]().ToTree([]*treeNode{{ID: 1}})
	})
}

func TestTreeBuilder_ChildrenField(t *testing.T) {
	as := assert.New(t)

	type altNode struct {
		ID   int
		Kids []*altNode
	}

	builder := NewTreeBuilder[altNode]().ChildrenField("Kids")
	as.Equal("Kids", builder.childrenField)

	as.Panics(func() {
		NewTreeBuilder[treeNode]().ChildrenField("Missing")
	})
}
