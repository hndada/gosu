package game

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/plays"
)

// A function which load database should not load the entire file system into memory.
// Instead, they should load the minimum information only, such as
// path to the file, and read the file when needed.

// Database can cover multiple file systems.
type Database struct {
	Chart  []ChartRow
	Replay []ReplayRow
}

type FSFile struct {
	FS   fs.FS
	Name string // Includes path
}

type ChartRow struct {
	FSFile
	MusicName          string
	Artist             string
	ChartName          string
	BackgroundFilename string
	Mode               int
	SubMode            int
	ChartHash          string
	Level              float64
}

func (c ChartRow) MusicString() string {
	return fmt.Sprintf("%s - %s", c.MusicName, c.Artist)
}

func (c ChartRow) LevelString() string {
	return fmt.Sprintf("Level %.0f", c.Level)
}

func (c ChartRow) ChartString() string {
	// return fmt.Sprintf("[Lv. %.0f] %s [%s]", c.Level, c.MusicName, c.ChartName) // [Lv. %4.2f]
	return fmt.Sprintf("%s [%s]", c.MusicName, c.ChartName)
}

func (c ChartRow) IsMatch(query string) bool {
	for _, s := range []string{c.MusicName, c.Artist, c.ChartName} {
		if strings.Contains(s, query) {
			return true
		}
	}
	return false
}

type ReplayRow struct {
	FSFile
	ChartHash string
}

// TODO: handle other music and replay directories
func NewDatabase(root fs.FS) (*Database, error) {
	var dbs Database

	// fs.Sub returns nil even if the directory does not exist.
	if _, err := fs.Stat(root, "music"); err == nil {
		fsys, err := fs.Sub(root, "music")
		if err != nil {
			return nil, fmt.Errorf("NewDatabase music: %w", err)
		}
		db, err := newChartDB(fsys)
		if err != nil {
			return nil, err
		}
		dbs.Chart = db
	}

	if _, err := fs.Stat(root, "replays"); err == nil {
		fsys, err := fs.Sub(root, "replays")
		if err != nil {
			return nil, fmt.Errorf("NewDatabase replays: %w", err)
		}
		db, err := newReplayDB(fsys)
		if err != nil {
			return nil, fmt.Errorf("NewDatabase replays: %w", err)
		}
		dbs.Replay = db
	}
	return &dbs, nil
}

// NewMusicDB reads only first depth of root for directory.
// Then it will read all charts in each directory.
func newChartDB(fsys fs.FS) ([]ChartRow, error) {
	dirs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("newChartDB dirs: %w", err)
	}

	var db []ChartRow
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dname := dir.Name()
		fs, err := fs.ReadDir(fsys, dname)
		if err != nil {
			return nil, fmt.Errorf("newChartDB dir: %w", err)
		}

		for _, f := range fs {
			if f.IsDir() {
				continue
			}

			ext := filepath.Ext(f.Name())
			if ext != ".osu" {
				continue
			}

			fname := path.Join(dname, f.Name())
			c, err := plays.NewChartHeaderFromFile(fsys, fname)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			db = append(db, ChartRow{
				FSFile: FSFile{
					FS:   fsys,
					Name: fname,
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
	}
	return db, nil
}

func newReplayDB(fsys fs.FS) ([]ReplayRow, error) {
	const maxKeyCount = 10

	fs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	var db []ReplayRow
	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		_, hash, err := plays.NewReplay(fsys, f.Name(), maxKeyCount)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		db = append(db, ReplayRow{
			FSFile: FSFile{
				FS:   fsys,
				Name: f.Name(),
			},
			ChartHash: hash,
		})
	}
	return db, nil
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
