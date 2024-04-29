package game

import (
	"crypto/md5"
	"io/fs"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

const (
	ModePiano = iota
	ModeDrum
	ModeSing
	ModeAll = -1
)

// *osu.Format
type Chart interface {
	// chart header
	WindowTitle() string

	// dynamics
	Dynamics() []Dynamic
	BPMs() (main, min, max float64)

	// notes
	NoteCounts() []int
	TotalDuration() int32 // Span()
}

type ChartFormat any

func LoadChartFormat(fsys fs.FS, name string) (ChartFormat, [16]byte, error) {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, [16]byte{}, err
	}

	var format ChartFormat
	switch filepath.Ext(name) {
	case ".osu", ".OSU":
		format, err = osu.NewFormat(data)
		if err != nil {
			return nil, [16]byte{}, err
		}
	}

	hash := md5.Sum(data)
	return format, hash, nil
}
