package choose

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/scene"
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

	// Attributes can be added by user, such as:
	// Genre, Language
	// Levels from game clients
	Path             string
	Attributes       map[string]any
	LastUpdateAtTime time.Time

	Level      float64
	NoteCounts []int
}

// This will be work as a key of music.
// Another possible way: MusicID = SetID + MusicFilename
// func (c Chart) MusicPath() string { return filepath.Join(c.Dirname, c.MusicFilename) }

// 'name' is a officially used name as file path in io/fs.
// newMusics reads only first depth of root for music.
// Then it will read all charts in each music.
func newMusics(root fs.FS) ([]Chart, []error) {
	musicEntries, err := fs.ReadDir(root, ".")
	errs := make([]error, 0, 5)
	if err != nil {
		return nil, append(errs, err)
	}

	// Support format: directory/.osu, .osz
	for _, entry := range musicEntries {
		var musicFS fs.FS
		switch {
		case entry.IsDir():
			musicFS, err = fs.Sub(root, entry.Name())
		case ext(entry.Name()) == ".osz":
			musicFS, err = zipFS(entry.Name())
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}

		cs, err := newCharts(musicFS)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		musics = append(musics, groupCharts(cs)...)
	}

	return musics, errs
}

func ext(path string) string { return strings.ToLower(filepath.Ext(path)) }

func zipFS(path string) (fs.FS, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// try parse. when no supporting mode,
// still enable to be shown, but not playable.
func newCharts(musicFS fs.FS) ([]Chart, error) {
	fs, err := fs.ReadDir(musicFS, ".")
	if err != nil {
		return nil, err
	}

	cs := make([]Chart, 0, 5)
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		// parse
	}
	return cs, nil
}

type sortBy int

const (
	sortByMusicName sortBy = iota
	sortByLevel
	sortByTime
	sortByAddAtTime
)

// You can keep each order of slice when after copying slice
// even if the slice is a slice of pointers.
// https://go.dev/play/p/yhvMddwd2co
func sortCharts(src []Chart, sortBy sortBy) []Chart {
	cs := make([]Chart, len(src))
	copy(cs, src)

	var less func(i, j int) bool
	switch sortBy {
	case sortByMusicName:
		less = func(i, j int) bool {
			if cs[i].MusicName < cs[j].MusicName {
				return true
			} else if cs[i].MusicName > cs[j].MusicName {
				return false
			}
			if cs[i].Artist < cs[j].Artist {
				return true
			} else if cs[i].Artist > cs[j].Artist {
				return false
			}
			return cs[i].Level < cs[j].Level
		}
	case sortByLevel:
		// Currently all precision of level is used.
		// Usage of using a certain precision: int(cs[i].Level*10)
		less = func(i, j int) bool {
			return cs[i].Level < cs[j].Level
		}
	case sortByTime:
		less = func(i, j int) bool {
			if cs[i].Duration < cs[j].Duration {
				return true
			} else if cs[i].Duration > cs[j].Duration {
				return false
			}
			return cs[i].Level < cs[j].Level
		}
	case sortByAddAtTime:
		less = func(i, j int) bool {
			return cs[i].AddAtTime.Before(cs[j].AddAtTime)
		}
	}

	sort.Slice(cs, less)
	return cs
}

func newList(src []Chart, sortBy sortBy) *scene.List {
	root := &scene.List{Name: "root"}
	cs := sortCharts(src, sortBy)

	var isEqual func(c1, c2 Chart) bool
	switch sortBy {
	case sortByMusicName:
		// Music name itself may be duplicated.
		// Artist + Title (Music name) may be unique.
		isEqual = func(c1, c2 Chart) bool {
			return c1.MusicName == c2.MusicName && c1.Artist == c2.Artist
		}
	case sortByLevel:
		isEqual = func(c1, c2 Chart) bool {
			return int(c1.Level*10) == int(c2.Level*10)
		}
	case sortByTime:
		// Unit is 10 seconds.
		isEqual = func(c1, c2 Chart) bool {
			return c1.Duration/1e4 == c2.Duration/1e4
		}
		// case sortByAddAtTime:
	}

	list := &scene.List{}
	for _, c := range cs {
		s := fmt.Sprintf("%s", c.MusicName)
		if len(list.Children) == 0 || isEqual(c, list.Children[0]) {

		}
	}
	return root
}

func (c Chart) String() string {
	return fmt.Sprintf("[Lv. %.0f] %s", c.MusicName)
}
