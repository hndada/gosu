package choose

import (
	"sort"
)

// type LevelFolder struct {
// 	Level10 int // scale is multiplied by 10.
// }
// type TimeFolder struct {
// 	Duration int // unit is 10 seconds.
// }

func a() {

	// type MusicNameFolder struct {
	// 	MusicName string
	// 	Artist    string
	// }

	// var charts map[string]Chart // name
	// var folders map[MusicNameFolder][]Chart
	// for _, c := range charts {
	// 	f := MusicNameFolder{MusicName: c.MusicName, Artist: c.Artist}
	// 	folders[f] = append(folders[f], c)
	// }
	// for name, cs := range folders {
	// 	n := &Node{Name: name.MusicName + name.Artist}
	// 	for _, c := range cs {
	// 		n.Children = append(n.Children, &Node{Name: c.String()})
	// 	}
	// }

	// It is fine to append first element at map without make().
	// https://go.dev/play/p/nXBtGxBIh1p
	var m = make(map[string][]*Node)
	for _, c := range []Chart{} {
		var path string
		name := &Node{Name: path}
		n := &Node{Name: c.LeafName(), Children: []*Node{name}}
		m[c.MusicArtistName()] = append(m[c.MusicArtistName()], n)
	}
	root := &Node{}
	for name, ns := range m {
		sort.Slice(ns, func(i, j int) bool {
			return ns[i].Name < ns[j].Name
		})
		n := &Node{Name: name, Children: ns}
		root.Children = append(root.Children, n)
	}
}
