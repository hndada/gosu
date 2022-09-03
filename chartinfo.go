package gosu

import (
	"fmt"
)

// ChartInfo is used at SceneSelect.
type ChartInfo struct {
	Path string
	// Mods    Mods
	Header  ChartHeader
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
	return fmt.Sprintf("(%dK Lv %.1f) %s [%s]", c.SubMode, c.Level, c.Header.MusicName, c.Header.ChartName)
}
func (c ChartInfo) BackgroundPath() string {
	return c.Header.BackgroundPath(c.Path)
}
