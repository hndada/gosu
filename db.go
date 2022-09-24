package gosu

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hndada/gosu/db"
)

// LoadChartInfos supposes Game's Modes has already set.
// ChartInfos are sorted with path, then mods.
// Todo: can the slice be sorted with Mode first, then MusicName?
func LoadChartInfosSet(modeProps []ModeProp) error {
	var err error
	set := make([][]ChartInfo, 0)
	err = db.LoadData("chart.db", &set)
	if err != nil {
		return err
	}
	if len(modeProps) != len(set) {
		return fmt.Errorf("mismatch game's modes length and db's modes length")
	}
	for mode, infos := range set {
		modeProps[mode].ChartInfos = infos
	}
	return nil
}

// TidyChartInfosSet drops unavailable chart infos from games.
func TidyChartInfosSet(modeProps []ModeProp) {
	for i, prop := range modeProps {
		for j, info := range prop.ChartInfos {
			if _, err := os.Stat(info.Path); err != nil {
				info1 := prop.ChartInfos[:j]
				info2 := prop.ChartInfos[j+1:]
				modeProps[i].ChartInfos = append(info1, info2...)
			}
		}
	}
}

// Todo: multiple music root. Would be not that hard.
// func LoadNewChartInfos(musicRoot string, prop *ModeProp) []ChartInfo {
func (prop ModeProp) LoadNewChartInfos(musicRoot string) []ChartInfo {
	chartInfos := make([]ChartInfo, 0)
	isNew := func(e fs.DirEntry) bool {
		info, err := e.Info()
		if err != nil {
			fmt.Println(err)
			return false
		}
		// Skip when modified time is equal or former than last update time.
		if !info.ModTime().After(prop.LastUpdateTime) {
			return false
		}
		return true
	}

	dirs, err := os.ReadDir(musicRoot)
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() || !isNew(dir) {
			continue
		}
		dpath := filepath.Join(musicRoot, dir.Name())
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
			if ChartFileMode(cpath) != prop.Mode {
				continue
			}
			info, err := prop.NewChartInfo(cpath) // First load should be done with no mods
			if err != nil {
				fmt.Printf("error at %s: %s\n", filepath.Base(cpath), err)
				continue
			}
			chartInfos = PutChartInfo(chartInfos, info)
		}
	}
	prop.LastUpdateTime = time.Now()
	return chartInfos
}

func SaveChartInfosSet(modeProps []ModeProp) {
	set := make([][]ChartInfo, len(modeProps))
	for i, prop := range modeProps {
		set[i] = prop.ChartInfos
	}
	var fname string
	switch db.MarshalType {
	case "json":
		fname = "chart.json"
	case "msgpack":
		fname = "chart.db"
	}
	db.SaveData(fname, &set)
}

// Todo: mode of ChartSet as a move unit
func PutChartInfo(infos []ChartInfo, info ChartInfo) []ChartInfo {
	i := sort.Search(len(infos), func(i int) bool {
		return infos[i].Path >= info.Path
	})
	// Append if not existed, update otherwise.
	if i == len(infos) || infos[i].Path != info.Path {
		infos = append(infos, ChartInfo{})
		copy(infos[i+1:], infos[i:])
		infos[i] = info
	} else {
		infos[i] = info
	}
	return infos
}
func Sort(sortBy int) {
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

	// sort.Slice(s.ChartBoxs, func(i, j int) bool {
	// 	if s.ChartBoxs[i].Chart.MusicName == s.ChartBoxs[j].Chart.MusicName {
	// 		return s.ChartBoxs[i].Level < s.ChartBoxs[j].Level
	// 	}
	// 	return s.ChartBoxs[i].Chart.MusicName < s.ChartBoxs[j].Chart.MusicName
	// })
}
