package choose

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
)

// Chart contains information of a chart.
// Favorites and Played count needs to be checked frequently.
type Chart struct {
	mode.ChartHeader
	Duration  int32
	MainBPM   float64
	MinBPM    float64
	MaxBPM    float64
	AddAtTime time.Time

	Path string
	// Attributes can be added by user, such as
	// Genre, Language, Levels from game clients
	Attributes       map[string]any
	LastUpdateAtTime time.Time

	Level      float64
	NoteCounts []int
}

// Music name itself may be duplicated.
// Artist + Title (Music name) may be unique.
func (c Chart) FolderNodeName() string {
	return fmt.Sprintf("%s - %s", c.MusicName, c.Artist)
}

func (c Chart) NodeName() string {
	return fmt.Sprintf("[Lv. %.0f] %s [%s]", c.Level, c.MusicName, c.ChartName)
}

// Todo: support .osz as music folder
// Memo: archive/zip.OpenReader returns ReadSeeker, which implements Read.
// Both Read and fs.Open are same in type: (name string) (fs.File, error)
func (c Chart) FSPath(root fs.FS) (fs.FS, string) {
	dir, name := path.Split(c.Path)
	fsys, _ := fs.Sub(root, dir)
	return fsys, name
}

// newCharts reads only first depth of root for directory.
// Then it will read all charts in each directory.
// Memo: 'name' is a officially used name as file path in io/fs.
func newCharts(root fs.FS) ([]*Chart, []error) {
	var cs []*Chart
	musicEntries, err := fs.ReadDir(root, ".")
	errs := make([]error, 0, 5)
	if err != nil {
		return nil, append(errs, err)
	}

	// Support format: directory/.osu, .osz
	for _, me := range musicEntries {
		var musicFS fs.FS
		switch {
		case me.IsDir():
			musicFS, err = fs.Sub(root, me.Name())
			// case ext(entry.Name()) == ".osz":
			// 	musicFS, err = zipFS(entry.Name())
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}

		chartEntries, err := fs.ReadDir(musicFS, ".")
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for _, ce := range chartEntries {
			if ce.IsDir() {
				continue
			}
			// Charts of unsupporting mode are still shown.
			switch ext(ce.Name()) {
			case ".osu":
				// read chart
				file, err := musicFS.Open(ce.Name())
				if err != nil {
					errs = append(errs, err)
					continue
				}

				f, err := osu.NewFormat(file)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				c := &Chart{
					ChartHeader: mode.NewChartHeader(f),
					Path:        path.Join(me.Name(), ce.Name()),
				}
				cs = append(cs, c)
			}
		}
	}
	return cs, errs
}

func ext(path string) string {
	return strings.ToLower(filepath.Ext(path))
}

func zipFS(path string) (fs.FS, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// You can keep each order of slice when after copying slice
// even if the slice is a slice of pointers.
// https://go.dev/play/p/yhvMddwd2co
func newChartTree(src []*Chart) *Node {
	cs := make([]*Chart, len(src))
	copy(cs, src)

	folders := make(map[string][]*Chart)
	for _, c := range cs {
		fdname := c.FolderNodeName()
		folders[fdname] = append(folders[fdname], c)
	}

	// Sort folders by name, sort charts by level.
	keys := make([]string, len(folders))
	for k, cs := range folders {
		// Currently all precision of level is used.
		// Usage of using a certain precision: int(cs[i].Level*10)
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
