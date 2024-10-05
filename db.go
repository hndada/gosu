package gosu

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/scene"
)

func NewDatabases(fsys fs.FS) (*scene.Databases, error) {
	var dbs scene.Databases
	var err error

	dbs.Music, err = NewMusicDB(fsys)
	if err != nil {
		return nil, err
	}

	dbs.Chart = make(map[scene.FileKey][]scene.ChartRow)
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
func NewMusicDB(fsys fs.FS) ([]scene.MusicRow, error) {
	dirs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	var list []scene.MusicRow
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		list = append(list, scene.MusicRow{
			FileKey: scene.FileKey{
				FS:   fsys,
				Name: dir.Name(),
			},
		})
	}

	return list, nil
}

func NewChartDB(fsys fs.FS, dirName string) ([]scene.ChartRow, error) {
	fs, err := fs.ReadDir(fsys, dirName)
	if err != nil {
		return nil, err
	}

	var list []scene.ChartRow
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

		list = append(list, scene.ChartRow{
			FileKey: scene.FileKey{
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

func NewReplayDB(fsys fs.FS) ([]scene.ReplayRow, error) {
	fs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	var list []scene.ReplayRow
	for _, f := range fs {
		if f.IsDir() {
			continue
		}

		_, hash, err := game.NewReplay(fsys, f.Name(), 4)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		list = append(list, scene.ReplayRow{
			FileKey: scene.FileKey{
				FS:   fsys,
				Name: f.Name(),
			},
			ChartHash: hash,
		})
	}
	return list, nil
}
