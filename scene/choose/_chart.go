package choose

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
)

type ChartKey struct {
	MD5  [16]byte
	Mods interface{} // Todo: use mods code for mode-specific mods?
}
type Chart struct {
	ChartKey

	mode.ChartHeader
	// Following fields are derived values.
	Level      float64
	NoteCounts []int
	Duration   int64
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64

	// Todo: should be separated as different struct?
	LastUpdateTime time.Time
	AddedTime      time.Time

	Genre    int //string
	Language int //string
	NSFW     bool
	// Tags can be added by user.
	Tags []string
	// Dropped Favorites and Played count since it
	// needs to be checked frequently.

	Pitch bool
}

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

// download.go
func (r ChimuResult) Filename() string {
	name := fmt.Sprintf("%d %s - %s.osz", r.SetId, r.Artist, r.Title)
	for _, letter := range []string{"<", ">", ":", "\"", "/", "\\", "|", "?", "*"} {
		name = strings.ReplaceAll(name, letter, "-")
	}
	return name
}

func ChartSetList(root string) map[int]bool {
	l := make(map[int]bool)
	dirs, err := os.ReadDir(root) // music dirs
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if dir.IsDir() || filepath.Ext(dir.Name()) == ".osz" {
			s := strings.Split(dir.Name(), " ")
			setId, err := strconv.Atoi(s[0])
			if err != nil {
				// fmt.Printf("%s: %s\n", err, s)
				continue
			}
			l[setId] = true
		}
	}
	return l
}

func BanList(path string) map[int]bool {
	ban := make(map[int]bool)
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()
		if s == "" || s[0] == '#' {
			continue
		}
		id, err := strconv.Atoi(s)
		if err != nil {
			continue
		}
		ban[id] = true
	}
	return ban
}
