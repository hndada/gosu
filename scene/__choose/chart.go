package choose

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

type ChartSet struct {
	SetId            int
	ChildrenBeatmaps []*Chart
	RankedStatus     int
	ApprovedDate     string
	LastUpdate       string
	LastChecked      string
	Artist           string
	Title            string
	Creator          string
	Source           string
	Tags             string
	HasVideo         bool
	Genre            int
	Language         int
	Favourites       int
	Disabled         int

	Path string
}

type Chart struct {
	*ChartSet
	BeatmapId        int
	ParentSetId      int
	DiffName         string
	FileMD5          string
	Mode             int
	BPM              float64
	AR               float64
	OD               float64
	CS               float64
	HP               float64
	TotalLength      int
	HitLength        int
	Playcount        int
	Passcount        int
	MaxCombo         int
	DifficultyRating float64
	OsuFile          string
	DownloadPath     string
}

func LoadChartSets() (sets []ChartSet) {
	// read music dir
	const root = "./music"
	dirs, err := os.ReadDir(root)
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		set := ChartSet{
			ChildrenBeatmaps: make([]*Chart, 0, 4),
			Path:             dir.Name(),
			// Path:             filepath.Join("./", dir.Name()),
		}
		dpath := filepath.Join(root, dir.Name())
		osuFiles, err := os.ReadDir(dpath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, file := range osuFiles {
			if file.IsDir() {
				continue
			}
			if filepath.Ext(file.Name()) != ".osu" {
				continue
			}
			fpath := filepath.Join(root, dir.Name(), file.Name())
			b, err := os.ReadFile(fpath)
			if err != nil {
				fmt.Println(err)
				continue
			}
			f, err := osu.NewFormat(b)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if set.Title == "" {
				set.Title = f.Title
			}
			if set.Artist == "" {
				set.Artist = f.Artist
			}
			chart := Chart{
				ChartSet: &set,
				DiffName: f.Version,
				Mode:     f.Mode,
				CS:       f.CircleSize,
				OsuFile:  file.Name(),
				// OsuFile:  filepath.Join("./", dir.Name(), file.Name()),
				// DownloadPath: fpath,
			}
			set.ChildrenBeatmaps = append(set.ChildrenBeatmaps, &chart)
		}
		sets = append(sets, set)
	}
	return
}
