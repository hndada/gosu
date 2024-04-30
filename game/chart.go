package game

import (
	"io/fs"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/util"
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
type Hash = string

func LoadChartFormat(fsys fs.FS, name string) (ChartFormat, Hash, error) {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		return nil, "", err
	}

	var format ChartFormat
	switch filepath.Ext(name) {
	case ".osu", ".OSU":
		format, err = osu.NewFormat(data)
		if err != nil {
			return nil, "", err
		}
	}
	return format, util.MD5(data), nil
}
