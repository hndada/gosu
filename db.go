package gosu

import (
	"github.com/hndada/gosu/mode"
)

// Its structure resembles mode.Chart.
// Todo: Relational Database
// https://go.dev/doc/database/querying
// https://github.com/osuripple/cheesegull
// https://www.cockroachlabs.com/docs/stable/start-a-local-cluster.html
type ChartInfo struct {
	Header mode.ChartHeader
	Mode   int
	Level  float64
	Path   string

	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
	Duration   int64
	NoteCounts [3]int
	// Tags    []string // Auto-generated or User-defined
}
