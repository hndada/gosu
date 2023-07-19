package scene

// Tree structure.
type List struct {
	Value ListValue
	// Title is derived from Value.String().
	// It is stored for drawing efficiently.
	Title string

	Parent   *List
	Children []*List
}

type ListValue interface{ String() string }

func NewList(v ListValue, parent *List, children []*List) *List {
	list := &List{
		Value:    v,
		Title:    v.String(),
		Parent:   parent,
		Children: children,
	}
	return list
}

// It is possible that non-leaf node has no children.
// To make sure that a node is a leaf, check its Children field has initialized.
// By the way, It is safe to call len() at nil slice.
// Yet, I explicitly check whether list is leaf or not for readability.
// https://go.dev/play/p/-1VWc9iDgMl
func (l List) IsLeaf() bool  { return l.Children == nil }
func (l List) IsEmpty() bool { return l.IsLeaf() || len(l.Children) == 0 }
