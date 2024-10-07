package scene

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/hndada/gosu/game"
)

// A function which load database should not load the entire file system into memory.
// Instead, they should load the minimum information only, such as
// path to the file, and read the file when needed.

type Databases struct {
	Music  []MusicRow
	Chart  map[FileKey][]ChartRow
	Replay []ReplayRow
}

type FileKey struct {
	FS   fs.FS
	Name string
}

type MusicRow struct {
	FileKey
	// FolderName string
}

func (r MusicRow) String() string {
	return r.Name
}

// Pseudo database
type ChartRow struct {
	FileKey
	MusicName          string
	Artist             string
	ChartName          string
	BackgroundFilename string
	Mode               int
	SubMode            int
	ChartHash          string
	Level              float64
}

func (r ChartRow) MusicString() string {
	return fmt.Sprintf("%s - %s", r.Artist, r.MusicName)
}

func (r ChartRow) ChartString() string {
	// return fmt.Sprintf("[Lv. %.0f] %s [%s]", c.Level, c.MusicName, c.ChartName) // [Lv. %4.2f]
	return fmt.Sprintf("%s [%s]", r.MusicName, r.ChartName)
}

type ReplayRow struct {
	FileKey
	ChartHash string
}

// Memo: archive/zip.OpenReader returns ReadSeeker, which implements Read.
// Both Read and fs.Open are same in type: (name string) (fs.File, error)
// func zipFS(path string) (fs.FS, error) {
// 	r, err := zip.OpenReader(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return r, nil
// }

func NewDatabases(fsys fs.FS) (*Databases, error) {
	var dbs Databases
	var err error

	dbs.Music, err = NewMusicDB(fsys)
	if err != nil {
		return nil, err
	}

	dbs.Chart = make(map[FileKey][]ChartRow)
	for _, m := range dbs.Music {
		cdb, err := NewChartDB(m.FS, m.Name)
		if err != nil {
			return nil, err
		}

		dbs.Chart[m.FileKey] = cdb
	}

	dbs.Replay, err = NewReplayDB(fsys)
	if err != nil {
		return nil, err
	}

	return &dbs, nil
}

// NewMusicDB reads only first depth of root for directory.
// Then it will read all charts in each directory.
func NewMusicDB(fsys fs.FS) ([]MusicRow, error) {
	dirs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	var list []MusicRow
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		list = append(list, MusicRow{
			FileKey: FileKey{
				FS:   fsys,
				Name: dir.Name(),
			},
		})
	}

	return list, nil
}

func NewChartDB(fsys fs.FS, dirName string) ([]ChartRow, error) {
	fs, err := fs.ReadDir(fsys, dirName)
	if err != nil {
		return nil, err
	}

	var list []ChartRow
	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		ext := filepath.Ext(f.Name())
		if ext != ".osu" {
			continue
		}

		c, err := game.NewChartHeaderFromFile(fsys, filepath.Join(dirName, f.Name()))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		list = append(list, ChartRow{
			FileKey: FileKey{
				FS:   fsys,
				Name: filepath.Join(dirName, f.Name()),
			},
			MusicName: c.MusicName,
			Artist:    c.Artist,
			ChartName: c.ChartName,
			Mode:      c.Mode,
			SubMode:   c.SubMode,
			ChartHash: c.ChartHash,
			// Level:     c.Level,
		})
	}

	return list, nil
}

func NewReplayDB(fsys fs.FS) ([]ReplayRow, error) {
	fs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	var list []ReplayRow
	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		_, hash, err := game.NewReplay(fsys, f.Name(), 4)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		list = append(list, ReplayRow{
			FileKey: FileKey{
				FS:   fsys,
				Name: f.Name(),
			},
			ChartHash: hash,
		})
	}
	return list, nil
}
