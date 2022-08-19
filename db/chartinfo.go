package db

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
	"github.com/vmihailenco/msgpack/v5"
)

// var ChartBoxs = make(map[string]ChartBox)
// ChartInfos are sorted with path, then mods.
// Todo: can the slice be sorted with Mode first, then MusicName?
var ChartInfos = make([]ChartInfo, 0)

// https://github.com/vmihailenco/msgpack
// https://github.com/osuripple/cheesegull
type ChartInfo struct {
	Path string
	// Mods mode.Mods
	Header mode.ChartHeader
	Mode   int
	Mode2  int
	Level  float64

	Duration   int64
	NoteCounts []int
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
	// Tags       []string // Auto-generated or User-defined
	// Box render.Sprite
}

func NewChartInfo(c *mode.Chart, fpath string, level float64) ChartInfo {
	mainBPM, minBPM, maxBPM := mode.BPMs(c.TransPoints, c.Duration)
	cb := ChartInfo{
		Path:   fpath,
		Header: c.ChartHeader,
		Mode:   c.Mode,
		Mode2:  c.Mode2,
		Level:  level,

		Duration:   c.Duration,
		NoteCounts: c.NoteCounts,
		MainBPM:    mainBPM,
		MinBPM:     minBPM,
		MaxBPM:     maxBPM,
	}
	// cb.Box = NewBoxSprite(c, level)
	return cb
}

// Todo: should deleted charts be checked after Unmarshal ChartInfos?
// Todo: MessagePack when tags=release, JSON when tags=debug
func LoadCharts(musicPath string) {
	const fname = "chart.db"
	b, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
		os.Rename(fname, fname+".crashed") // Fail if not existed.
	}
	msgpack.Unmarshal(b, &ChartInfos)
	LoadNewCharts(musicPath)
}

var LastUpdateTime time.Time

// Todo: ChartBoxDB, ScoreDB
func LoadNewCharts(musicPath string) {
	isNew := func(e fs.DirEntry) bool {
		info, err := e.Info()
		if err != nil {
			fmt.Println(err)
			return false
		}
		// Skip when modified time is equal or former than last update time.
		if !info.ModTime().After(LastUpdateTime) {
			return false
		}
		return true
	}

	dirs, err := os.ReadDir(musicPath)
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() || !isNew(dir) {
			continue
		}
		dpath := filepath.Join(musicPath, dir.Name())
		fs, err := os.ReadDir(dpath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, f := range fs {
			if f.IsDir() || !isNew(f) { // There may be directory e.g., SB
				continue
			}
			fpath := filepath.Join(dpath, f.Name())
			// dat, err := os.ReadFile(fpath)
			// if err != nil {
			// 	fmt.Println(err)
			// 	continue
			// }
			switch mode.Mode(fpath) {
			case mode.ModePiano4, mode.ModePiano7:
				// o, err := osu.Parse(dat)
				// if err != nil {
				// 	fmt.Println(err)
				// 	continue
				// }
				c, err := piano.NewChart(fpath)
				if err != nil {
					fmt.Printf("error at %s: %s\n", filepath.Base(fpath), err)
					continue
				}
				info := NewChartInfo(&c.Chart, fpath, mode.Level(c))
				ChartInfos = Put(ChartInfos, info)
				// ChartBoxs[fpath] = info
				// for _, sort := range []int{SortByName, SortByLevel} {
				// 	Insert(ChartViews[ViewMode(c.Mode, sort)], info, sort)
				// }
			case mode.ModeDrum:
			}
		}
	}
	LastUpdateTime = time.Now()
	fmt.Println(len(ChartInfos))
	SaveChartInfos()
	// sort.Slice(s.ChartBoxs, func(i, j int) bool {
	// 	if s.ChartBoxs[i].Chart.MusicName == s.ChartBoxs[j].Chart.MusicName {
	// 		return s.ChartBoxs[i].Level < s.ChartBoxs[j].Level
	// 	}
	// 	return s.ChartBoxs[i].Chart.MusicName < s.ChartBoxs[j].Chart.MusicName
	// })
}

// Todo: mode of ChartSet as a move unit
func Put(infos []ChartInfo, info ChartInfo) []ChartInfo {
	// i := sort.Search(len(view), func(i int) bool {
	// 	switch sortBy {
	// 	case SortByLevel:
	// 		return view[i].Level >= info.Level
	// 	case SortByName:
	// 		if view[i].Header.MusicName == info.Header.MusicName {
	// 			return view[i].Level >= info.Level
	// 		}
	// 		return view[i].Header.MusicName >= info.Header.MusicName
	// 	default:
	// 		return view[i].Level >= info.Level
	// 	}
	// })
	i := sort.Search(len(infos), func(i int) bool {
		return infos[i].Path >= info.Path
	})
	// Append if not existed, update otherwise.
	if i == len(ChartInfos) || ChartInfos[i].Path != info.Path {
		infos = append(infos, ChartInfo{})
		copy(infos[i+1:], infos[i:])
		infos[i] = info
	} else {
		infos[i] = info
	}
	SaveChartInfos()
	return infos
}
func SaveChartInfos() {
	b, err := msgpack.Marshal(&ChartInfos)
	if err != nil {
		fmt.Printf("Failed to save by an error: %s", err)
		return
	}
	err = os.WriteFile("chart.db", b, 0644)
	if err != nil {
		fmt.Printf("Failed to save by an error: %s", err)
		return
	}
}
func (c ChartInfo) Text() string {
	return fmt.Sprintf("(%dK Lv %.1f) %s [%s]", c.Mode2, c.Level, c.Header.MusicName, c.Header.ChartName)
}