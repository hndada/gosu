package choose

import (
	"archive/zip"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/mode"
)

// Chart contains information of a chart.
// Favorites and Played count needs to be checked frequently.
type Chart struct {
	mode.ChartHeader
	Duration int32
	MainBPM  float64
	MinBPM   float64
	MaxBPM   float64

	Dirname          string // for music ID
	AddAtTime        time.Time
	LastUpdateAtTime time.Time

	Level      float64
	NoteCounts []int

	// Attributes can be added by user, such as:
	// Genre, Language
	// Levels from game clients
	Attributes map[string]any
}

// This will be work as a key of music.
// Another possible way: MusicID = SetID + MusicFilename
func (c Chart) MusicPath() string { return filepath.Join(c.Dirname, c.MusicFilename) }

// newMusics reads only first depth of fsys.
func newMusics(root fs.FS) ([]Music, []error) {
	var musics []Music

	musicEntries, err := fs.ReadDir(root, ".")
	errs := make([]error, 0, 5)
	if err != nil {
		return nil, append(errs, err)
	}

	// Support format: directory/.osu, .osz
	for _, entry := range musicEntries {
		var musicFS fs.FS
		switch {
		case entry.IsDir():
			musicFS, err = fs.Sub(root, entry.Name())
		case ext(entry.Name()) == ".osz":
			musicFS, err = zipFS(entry.Name())
		}
		if err != nil {
			errs = append(errs, err)
			continue
		}

		cs, err := newCharts(musicFS)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		musics = append(musics, groupCharts(cs)...)
	}

	return musics, errs
}

// try parse. when no supporting mode,
// still enable to be shown, but not playable.
func newCharts(musicFS fs.FS) ([]Chart, error) {
	fs, err := fs.ReadDir(musicFS, ".")
	if err != nil {
		return nil, err
	}

	cs := make([]Chart, 0, 5)
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		// parse
	}
	return cs, nil
}

func groupCharts(cs []Chart) []Music {
	var ms []Music
	for _, c := range cs {
	}
	return ms
}

func ext(path string) string { return strings.ToLower(filepath.Ext(path)) }

func zipFS(path string) (fs.FS, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// Todo: re-wrap mode.Replay; include osr.Format's header part.
func newReplays(fsys fs.FS, charts map[[16]byte]*Chart) map[[16]byte]*osr.Format {
	m := make(map[[16]byte]*osr.Format)

	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || ext(path) != ".osr" {
			return nil
		}

		file, err := fsys.Open(path)
		if err != nil {
			return err
		}

		switch ext(path) {
		case ".osr":
			f, err := osr.NewFormat(file)
			if err != nil {
				return err
			}

			md5, err := f.MD5()
			if err != nil {
				return err
			}
			m[md5] = f
		}
		return nil
	})
	return m
}
