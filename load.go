package gosu

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/mania"
)

var lastUpdateTime time.Time

// temp: currently only mania chart
// returns a slice of newly update charts
func updateCharts(cwd string) []*mania.Chart {
	root := filepath.Join(cwd, "music")
	dir, err := os.Stat(root)
	if err != nil {
		panic(err)
	}
	var cs []*mania.Chart
	if dir.ModTime().After(lastUpdateTime) {
		cs = loadCharts(cwd)
	}
	lastUpdateTime = time.Now()
	charts = append(charts, cs...)
	return cs
}

// todo:로드된 차트 데이터는 gob로 저장?
// 새로 만들어진 것만 추가로 로드
func loadCharts(cwd string) []*mania.Chart {
	cs := make([]*mania.Chart, 0, 40)
	root := filepath.Join(cwd, "music")
	dirs, err := ioutil.ReadDir(root) // music dirs
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		if dir.ModTime().Before(lastUpdateTime) {
			continue
		}
		dirPath := filepath.Join(root, dir.Name())
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			if f.ModTime().Before(lastUpdateTime) {
				continue
			}
			fpath := filepath.Join(dirPath, f.Name())
			switch strings.ToLower(filepath.Ext(fpath)) {
			case ".osu":
				switch common.OsuMode(fpath) {
				case common.ModeMania:
					c, err := mania.NewChart(fpath)
					if err != nil {
						panic(err)
					}
					cs = append(cs, c)
				}
			}
		}
	}
	sort.Slice(cs, func(i, j int) bool {
		if cs[i].KeyCount == cs[j].KeyCount {
			return cs[i].Level < cs[j].Level
		} else {
			return cs[i].KeyCount < cs[j].KeyCount
		}
	})
	return cs
}
