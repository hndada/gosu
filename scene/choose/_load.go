package choose

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
)

// // ChartInfo is used at SceneSelect.
// type ChartInfo struct {
// 	Header
// 	// Tags       []string // Auto-generated or User-defined
// 	Path string
// 	// Mods    Mods

//		// Following fields are derived values.
//		Level      float64
//		NoteCounts []int
//		Duration   int64
//		MainBPM    float64
//		MinBPM     float64
//		MaxBPM     float64
//	}
var Charts = make([]Chart, 0, 50)

func loadChart(fs []fs.DirEntry) {
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
	}
}
func LoadCharts(fsys fs.FS, root string) (cs []Chart) {
	// defer sort
	musics, err := fs.ReadDir(fsys, root)
	if err != nil {
		return nil
	}
	for _, music := range musics {
		if music.IsDir() { // Directory
			fs, err := fs.ReadDir(fsys, path.Join(root, music.Name()))
			if err != nil {
				continue
			}
			loadChart(fs)
		} else { // Zip file
			info, err := music.Info()
			if err != nil {
				continue
			}
			switch ext := filepath.Ext(info.Name()); ext {
			case ".osz", ".OSZ":
				fs, err := fs.ReadDir()
				loadChart(music)
			}
		}
	}
	for _, dir := range dirs {
		for _, f := range fs {
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
	return
}
