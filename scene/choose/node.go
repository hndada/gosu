package choose

import (
	"fmt"
	"io/fs"
	"path"
	"sort"

	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/play"
)

type NodeType int

const (
	RootNode NodeType = iota
	FolderNode
	ChartNote // Chart
	PathNode  // leaf; Chart.Path
)

type Node struct {
	Parent *Node
	// Children []*Node
	PrevSibling *Node
	NextSibling *Node
	FirstChild  *Node
	LastChild   *Node

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

// type Item struct {
// 	Text     string
// 	Parent   *Item
// 	Children []*Item
// 	IsLeaf   bool
// }

// func newTree(src []*Chart) *Node {
// 	cs := make([]*Chart, len(src))
// 	copy(cs, src)
// 	sort.Slice(cs, func(i, j int) bool {
// 		return cs[i].Level < cs[j].Level
// 	})

// 	folders := make(map[string][]*Item)
// 	for _, c := range cs {
// 		// Sort first since items in a folder has few information about charts.
// 		leaf := &Item{Text: c.Path, IsLeaf: true}
// 		item := &Item{
// 			Text:     fmt.Sprintf("[Lv. %.0f] %s [%s]", c.Level, c.MusicName, c.ChartName),
// 			Children: []*Item{leaf},
// 		}
// 		leaf.Parent = item

// 		gname := fmt.Sprintf("%s - %s", c.MusicName, c.Artist) // folder name
// 		folders[gname] = append(folders[gname], item)
// 	}

// 	root := &Item{Children: make([]*Item, 0, len(folders))}
// 	for name, folder := range folders {
// 		g := &Item{Text: name, Parent: root, Children: folder}
// 		root.Children = append(root.Children, g)
// 	}
// 	return root
// }

func newTree(src []*Chart) *Node {
	cs := make([]*Chart, len(src))
	copy(cs, src)

	folders := make(map[string][]*Chart)
	for _, c := range cs {
		fdname := c.FolderName()
		folders[fdname] = append(folders[fdname], c)
	}

	keys := make([]string, len(folders))
	for k, cs := range folders {
		sort.Slice(cs, func(i, j int) bool {
			return cs[i].Level < cs[j].Level
		})
		folders[k] = cs
		keys = append(keys, k)
	}
	sort.Strings(keys)

	root := &Node{Type: RootNode}
	for _, name := range keys {
		folder := &Node{Type: FolderNode, Data: name, Parent: root}
		for _, c := range folders[name] {
			chart := &Node{Type: ChartNote, Data: c.NodeName(), Parent: folder}
			path := &Node{Type: PathNode, Data: c.Path, Parent: chart}
			chart.AppendChild(path)
			folder.AppendChild(chart)
		}
		root.AppendChild(folder)
	}
	return root
}

// Todo: pass folder mode?
func (c Chart) FolderName() string {
	return fmt.Sprintf("%s - %s", c.MusicName, c.Artist)
}

// Todo: NodeName vs String vs Name vs ChartName?
// Cosider: Currently, there is a field named "ChartName" in Chart.
func (c Chart) NodeName() string {
	return fmt.Sprintf("[Lv. %.0f] %s [%s]", c.Level, c.MusicName, c.ChartName)
}

func (s *Scene) setMusicPlayer(fsys fs.FS, name string) {

}
func (s *Scene) setBackground(fsys fs.FS, name string) {
	scene.NewBackgroundDrawer()
}
func (s Scene) Open(fsys fs.FS, name string) fs.File {
	//
}

type ChartFileInfo struct {
	Root string
	Dir  string
	Name string
}

// music root, music folder, chart
// music folder can be two types: .osz, directory

type Chart2 struct {
	// Base string // Root + Dir
	Path     string
	IsZipped bool
}

// s.currentNode deserves to be stored, during the game. -> map[string]scene
// func (s *Scene) PlayChart(path string) (fsys fs.FS, name string) {}

func (g *game) Update() any {
	args := g.scene.Update()
	switch g.scene {
	case *choose.Scene:
		switch args.(type) {
		case string:
			path := args.(string)
		}
	case *play.Scene:
		switch args.(type) {
		case Exit:
		case piano.Scorer: // , drum.Scorer
			// Result scene
		}
	}
}

func (c Chart2) FSPath(root fs.FS) (fsys fs.FS, name string, err error) {
	dir, name := path.Split(c.Path)
	if c.IsZipped {
		fsys, err = zipFS(dir)
	} else {
		fsys, err = fs.Sub(root, dir)
	}
	// try .osz
	if err != nil {

	}

	// handle .osz or dir

}

// let's drop .osz feature for now.
// zipfs: experimental
// fs.FS: Open(name string) (File, error)
// zip.OpenReader returns ReadSeeker
// ReadSeeker implements Read
// Read(name string) (fs.File, error)
