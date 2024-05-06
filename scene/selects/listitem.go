package selects

import (
	"fmt"
	"sort"

	"github.com/hndada/gosu/game"
)

// Music name itself may be duplicated.
// Artist + Title (Music name) may be unique.
func FolderText(c *game.ChartHeader) string {
	return fmt.Sprintf("%s - %s", c.MusicName, c.Artist)
}

// Todo: add level database. attach level info to the text.
func ItemText(c *game.ChartHeader) string {
	// return fmt.Sprintf("[Lv. %.0f] %s [%s]", c.Level, c.MusicName, c.ChartName) // [Lv. %4.2f]
	return fmt.Sprintf("%s [%s]", c.MusicName, c.ChartName)
}

// You can keep each order of slice when after copying slice
// even if the slice is a slice of pointers.
// https://go.dev/play/p/yhvMddwd2co
func newChartTree(src map[string]*Chart) *Node { // key: c.Hash
	cs := make([]*Chart, len(src))
	i := 0
	for _, c := range src {
		cs[i] = c
		i++
	}

	folders := make(map[string][]*Chart) // key: c.FolderNodeName()
	for _, c := range cs {
		fdname := c.FolderNodeName()
		folders[fdname] = append(folders[fdname], c)
	}

	// Sort folders by name, sort charts by level.
	// Memo: make([]T, len) and make([]T, 0, len) is prone to be erroneous.
	keys := make([]string, 0, len(folders))
	for k, cs := range folders {
		// Currently all precision of level is used.
		// Usage of using a certain precision: int(cs[i].Level*10)
		sort.Slice(cs, func(i, j int) bool {
			return cs[i].Level < cs[j].Level
		})
		folders[k] = cs
		keys = append(keys, k)
	}
	// Todo: add sort criteria to config.
	// Group1, Group2, Sort, Filter int
	// sortByMusicName, Level, Time, AddAtTime
	// func(i, j int) bool { return cs[i].AddAtTime.Before(cs[j].AddAtTime) }
	sort.Strings(keys)

	root := &Node{Type: RootNode}
	for _, name := range keys {
		folder := &Node{Type: FolderNode, Data: name}
		for _, c := range folders[name] {
			chart := &Node{Type: ChartNode, Data: c.NodeName()}
			path := &Node{Type: LeafNode, Data: c.Hash}
			chart.AppendChild(path)
			folder.AppendChild(chart)
		}
		root.AppendChild(folder)
	}
	return root
}

// Memo: archive/zip.OpenReader returns ReadSeeker, which implements Read.
// Both Read and fs.Open are same in type: (name string) (fs.File, error)
// func zipFS(path string) (fs.FS, error) {
// 	r, err := zip.OpenReader(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return r, nil
// }
