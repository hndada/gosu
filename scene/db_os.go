package scene

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/hndada/gosu/game"
)

func NewMusicDB(fsys fs.FS) ([]MusicRow, error) {
	dirs, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	var musicList []MusicRow
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		musicList = append(musicList, MusicRow{
			FileKey: FileKey{
				FS:   fsys,
				Name: dir.Name(),
			},
			// FolderName: dir.Name(),
		})
	}

	return musicList, nil
}
func NewChartDB(fsys fs.FS, dirName string) ([]ChartRow, error) {
	fs, err := fs.ReadDir(fsys, dirName)
	if err != nil {
		return nil, err
	}

	var chartList []ChartRow
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

		chartList = append(chartList, ChartRow{
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

	return chartList, nil
}
