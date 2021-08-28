package gosu

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"
)

// todo: 새 것만 로드하기
func updateCharts(cwd string) {
	// loadCharts(cwd)
}

// 로드된 차트 데이터는 gob로 저장
func loadCharts(cwd string) []*mania.Chart { // temp: currently only mania chart
	charts := make([]*mania.Chart, 0, 100)
	root := filepath.Join(cwd, "music")
	dirs, err := ioutil.ReadDir(root) // music dirs
	if err != nil {
		log.Fatal(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dirPath := filepath.Join(root, dir.Name())
		files, err := ioutil.ReadDir(dirPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			fpath := filepath.Join(dirPath, f.Name())
			switch strings.ToLower(filepath.Ext(fpath)) {
			case ".osu":
				switch game.OsuMode(fpath) {
				case game.ModeMania:
					c, err := mania.NewChart(fpath)
					if err != nil {
						log.Fatal(err)
					}
					charts = append(charts, c)
				}
			}
		}
	}
	sort.Slice(charts, func(i, j int) bool {
		if charts[i].KeyCount == charts[j].KeyCount {
			return charts[i].Level < charts[j].Level
		} else {
			return charts[i].KeyCount < charts[j].KeyCount
		}
	})
	return charts
}
