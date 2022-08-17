package db

import (
	"github.com/hndada/gosu/mode"
)

// https://github.com/vmihailenco/msgpack
// https://github.com/osuripple/cheesegull
type ChartInfo struct {
	Path   string
	Header mode.ChartHeader
	Mode   int
	Level  float64

	Duration   int64
	NoteCounts []int
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
	// Tags       []string // Auto-generated or User-defined
}

func NewChartInfo(c mode.Chart, fpath string, level float64) ChartInfo {
	mainBPM, minBPM, maxBPM := mode.BPMs(c.TransPoints, c.Duration)
	return ChartInfo{
		Path:   fpath,
		Header: c.ChartHeader,
		Mode:   c.Mode,
		Level:  level,

		Duration:   c.Duration,
		NoteCounts: c.NoteCounts,
		MainBPM:    mainBPM,
		MinBPM:     minBPM,
		MaxBPM:     maxBPM,
	}
}
