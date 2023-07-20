package choose

type NodeType int

const (
	RootNode NodeType = iota
	FolderNode
	ChartNote // Chart
	PathNode  // leaf; Chart.Path
)

// Inspired by html.Node.
// But html.Node is too complicated for constructing chart tree.
type Node struct {
	Parent                   *Node
	PrevSibling, NextSibling *Node
	FirstChild, LastChild    *Node

	Type NodeType
	Data string
}

// GPT-4
func (n *Node) AppendChild(child *Node) {
	child.Parent = n
	if n.FirstChild == nil {
		// If this node has no children, just set the new child as the first and last child
		n.FirstChild = child
		n.LastChild = child
	} else {
		// If this node does have children, append the new child to the end of the list
		n.LastChild.NextSibling = child
		child.PrevSibling = n.LastChild
		n.LastChild = child
	}
}

func (n *Node) Prev() *Node {
	if n.PrevSibling != nil {
		return n.PrevSibling
	}
	return n.Parent
}

func (n *Node) Next() *Node {
	if n.NextSibling != nil {
		return n.NextSibling
	}
	return n.Parent.NextSibling
}
