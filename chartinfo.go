package gosu

import (
	"fmt"
)

// ChartInfo is used at SceneSelect.
type ChartInfo struct {
	Path string
	// Mods    Mods
	// Header  ChartHeader
	ChartHeader
	Mode    int
	SubMode int
	Level   float64

	Duration   int64
	NoteCounts []int
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
	// Tags       []string // Auto-generated or User-defined
}

func (c ChartInfo) Text() string {
	return fmt.Sprintf("(%dK Lv %.1f) %s [%s]", c.SubMode, c.Level, c.MusicName, c.ChartName)
}
func (c ChartInfo) BackgroundPath() string {
	return c.ChartHeader.BackgroundPath(c.Path)
}
func (c ChartInfo) TimeString() string {
	c.Duration /= 1000
	return fmt.Sprintf("%02d:%02d", c.Duration/60, c.Duration%60)
}
func (c ChartInfo) BPMString() string {
	return fmt.Sprintf("%.0f BPM (%.0f ~ %.0f)", c.MainBPM, c.MinBPM, c.MaxBPM)
}
