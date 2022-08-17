package gosu

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/render"
)

func LoadDB() {}
func LoadNewMusic() {
	dirs, err := os.ReadDir(MusicPath)
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dpath := filepath.Join(MusicPath, dir.Name())
		fs, err := os.ReadDir(dpath)
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			if f.IsDir() { // There may be directory e.g., SB
				continue
			}
			fpath := filepath.Join(dpath, f.Name())
			b, err := os.ReadFile(fpath)
			if err != nil {
				panic(err)
			}
			switch strings.ToLower(filepath.Ext(fpath)) {
			case ".osu":
				o, err := osu.Parse(b)
				if err != nil {
					panic(err)
				}
				switch o.Mode {
				case mode.ModeMania:
					c, err := piano.NewChartFromOsu(o)
					if err != nil {
						panic(err)
					}
					info := ChartInfo{
						Chart: c,
						Path:  fpath,
						Level: mode.Level(c.Difficulties()),
						Box: render.Sprite{
							I: NewBox(c, mode.Level(c.Difficulties())),
							W: bw,
							H: bh,
						},
					}
					s.ChartInfos = append(s.ChartInfos, info)
					// box's x value is not fixed.
					// box's y value is not fixed.
				case mode.ModeTaiko:
				}
			}
		}
	}
	sort.Slice(s.ChartInfos, func(i, j int) bool {
		if s.ChartInfos[i].Chart.MusicName == s.ChartInfos[j].Chart.MusicName {
			return s.ChartInfos[i].Level < s.ChartInfos[j].Level
		}
		return s.ChartInfos[i].Chart.MusicName < s.ChartInfos[j].Chart.MusicName
	})
}
