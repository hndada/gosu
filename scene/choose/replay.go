package choose

import (
	"io/fs"

	"github.com/hndada/gosu/format/osr"
)

// Todo: implement non-playing score simulator
// Todo: re-wrap mode.Replay; include osr.Format's header part.
func newReplays(fsys fs.FS, charts map[string]*Chart) map[string]*osr.Format {
	m := make(map[string]*osr.Format)

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
			m[f.BeatmapMD5] = f
		}
		return nil
	})

	return m
}
