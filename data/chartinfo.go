package data

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
)

// ChartInfos are sorted with path, then mods.
// Todo: can the slice be sorted with Mode first, then MusicName?
var ChartInfos = make([]ChartInfo, 0)
var LastUpdateTime time.Time

// https://github.com/vmihailenco/msgpack
type ChartInfo struct {
	Path string
	// Mods mode.Mods
	Header  mode.ChartHeader
	Mode    int
	SubMode int
	Level   float64

	Duration   int64
	NoteCounts []int
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
	// Tags       []string // Auto-generated or User-defined
}

func (c ChartInfo) Text() string {
	return fmt.Sprintf("(%dK Lv %.1f) %s [%s]", c.SubMode, c.Level, c.Header.MusicName, c.Header.ChartName)
}

func NewChartInfo(c *mode.Chart, cpath string, level float64) ChartInfo {
	mainBPM, minBPM, maxBPM := mode.BPMs(c.TransPoints, c.Duration)
	cb := ChartInfo{
		Path:    cpath,
		Header:  c.ChartHeader,
		Mode:    c.Mode,
		SubMode: c.SubMode,
		Level:   level,

		Duration:   c.Duration,
		NoteCounts: c.NoteCounts,
		MainBPM:    mainBPM,
		MinBPM:     minBPM,
		MaxBPM:     maxBPM,
	}
	return cb
}

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
			cpath := filepath.Join(dpath, f.Name())
			switch mode.Mode(cpath) {
			case mode.ModePiano4, mode.ModePiano7:
				c, err := piano.NewChart(cpath)
				if err != nil {
					fmt.Printf("error at %s: %s\n", filepath.Base(cpath), err)
					continue
				}
				info := NewChartInfo(&c.Chart, cpath, mode.Level(c))
				ChartInfos = Put(ChartInfos, info)
			case mode.ModeDrum:
			}
		}
	}
	LastUpdateTime = time.Now()
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
	SaveData()
	return infos
}
