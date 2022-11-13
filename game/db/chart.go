package db

import (
	"time"

	"github.com/hndada/gosu/game/chart"
)

// Todo: should MD5 be substituded with SHA256?
type ChartKey struct {
	MD5  [16]byte
	Mods interface{} // Todo: use mods code for mode-specific mods?
}
type Chart struct {
	ChartKey

	chart.Header
	// Following fields are derived values.
	Level      float64
	NoteCounts []int
	Duration   int64
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64

	// Todo: should be separated as different struct?
	LastUpdateTime time.Time
	AddedTime      time.Time

	// Tags can be added by user.
	Tags []string
}
