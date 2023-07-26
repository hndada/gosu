package choose

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
)

// Chart contains information of a chart.
// Favorites and Played count needs to be checked frequently.
type Chart struct {
	MusicFS  fs.FS
	Base     string // For comparing fsys.
	Filename string
	Hash     string // md5 with 16 bytes

	mode.ChartHeader
	Duration         int32
	MainBPM          float64
	MinBPM           float64
	MaxBPM           float64
	AddAtTime        time.Time
	LastUpdateAtTime time.Time
	Level            float64
	NoteCounts       []int
	// Attributes can be added by user, such as
	// Genre, Language, Levels from game clients
	Attributes map[string]any
}

// Music name itself may be duplicated.
// Artist + Title (Music name) may be unique.
func (c Chart) FolderNodeName() string {
	return fmt.Sprintf("%s - %s", c.MusicName, c.Artist)
}

func (c Chart) NodeName() string {
	return fmt.Sprintf("[Lv. %.0f] %s [%s]", c.Level, c.MusicName, c.ChartName) // [Lv. %4.2f]
}

// newCharts reads only first depth of root for directory.
// Then it will read all charts in each directory.
// Memo: 'name' is a officially used name as file path in io/fs.
func newCharts(root fs.FS) (map[string]*Chart, []error) {
	cs := make(map[string]*Chart)
	musicEntries, err := fs.ReadDir(root, ".")
	errs := make([]error, 0, 5)
	if err != nil {
		return nil, append(errs, err)
	}

	// Support format: directory/.osu
	for _, me := range musicEntries {
		// Todo: support .osz as music folder
		if !me.IsDir() {
			continue
		}

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
				defer file.Close()

				dat, err := fs.ReadFile(musicFS, ce.Name())
				if err != nil {
					errs = append(errs, err)
					continue
				}
				hash := md5.Sum(dat)

				f, err := osu.NewFormat(file)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				c := &Chart{
					MusicFS:     musicFS,
					Base:        me.Name(),
					Filename:    ce.Name(),
					Hash:        string(hash[:]),
					ChartHeader: mode.NewChartHeader(f),
				}

				switch f.Mode {
				case osu.ModeMania:
					cp, err := piano.NewChart(musicFS, ce.Name())
					if err != nil {
						errs = append(errs, err)
						continue
					}
					c.Duration = cp.Duration()
					c.Level = cp.Level
				}

				cs[c.Hash] = c
			}
		}
	}
	return cs, errs
}

func ext(path string) string {
	return strings.ToLower(filepath.Ext(path))
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

// func (c Chart) FSPath(root fs.FS) (fs.FS, string) {
// 	dir, name := path.Split(c.Path)
// 	fsys, _ := fs.Sub(root, dir)
// 	return fsys, name
// }
