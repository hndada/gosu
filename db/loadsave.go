package db

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
)

var LastUpdateTime time.Time

// Todo: ChartInfoDB, ScoreDB
func LoadNewMusic(musicPath string) {
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
			dat, err := os.ReadFile(fpath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			switch mode.Mode(fpath) {
			case mode.ModePiano:
				o, err := osu.Parse(dat)
				if err != nil {
					fmt.Println(err)
					continue
				}
				c := piano.NewChart(o)
				info := NewChartInfo(c.Chart, fpath, mode.Level(c))
				ChartInfos[fpath] = info // Append if not existed, update otherwise.
				for _, sort := range []int{SortByName, SortByLevel} {
					Insert(ChartViews[ViewMode(c.Mode, sort)], info, sort)
				}
			case mode.ModeDrum:
			}
		}
	}
	LastUpdateTime = time.Now()
	// sort.Slice(s.ChartInfos, func(i, j int) bool {
	// 	if s.ChartInfos[i].Chart.MusicName == s.ChartInfos[j].Chart.MusicName {
	// 		return s.ChartInfos[i].Level < s.ChartInfos[j].Level
	// 	}
	// 	return s.ChartInfos[i].Chart.MusicName < s.ChartInfos[j].Chart.MusicName
	// })
}

// Todo: MessagePack when tags=release, JSON when tags=debug
func LoadDB() {

}

// Todo: mode of ChartSet as a move unit
func Insert(view []ChartInfo, info ChartInfo, sortBy int) []ChartInfo {
	i := sort.Search(len(view), func(i int) bool {
		switch sortBy {
		case SortByLevel:
			return view[i].Level >= info.Level
		case SortByName:
			if view[i].Header.MusicName == info.Header.MusicName {
				return view[i].Level >= info.Level
			}
			return view[i].Header.MusicName >= info.Header.MusicName
		default:
			return view[i].Level >= info.Level
		}
	})
	view = append(view, ChartInfo{})
	copy(view[i+1:], view[i:])
	view[i] = info
	return view
}
