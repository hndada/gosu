package choose

import (
	"time"

	"github.com/hndada/gosu/mode"
)

//	func (c ChartInfo) Text() string {
//		switch c.Mode {
//		case ModePiano4, ModePiano7:
//			return fmt.Sprintf("(%dK Level %3.1f) %s [%s]", c.SubMode, c.Level, c.MusicName, c.ChartName)
//		case ModeDrum:
//			return fmt.Sprintf("(Level %3.1f) %s [%s]", c.Level, c.MusicName, c.ChartName)
//		}
//		return ""
//	}
//
//	func (c ChartInfo) TimeString() string {
//		c.Duration /= 1000
//		return fmt.Sprintf("%02d:%02d", c.Duration/60, c.Duration%60)
//	}
//
// ChartInfo is used at SceneSelect.
type ChartInfo struct {
	mode.ChartHeader
	// Tags       []string // Auto-generated or User-defined
	Path string
	// Mods    Mods

	// Following fields are derived values.
	Level      float64
	NoteCounts []int
	Duration   int64
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
}
type ChartKey struct {
	MD5  [16]byte
	Mods interface{} // Todo: use mods code for mode-specific mods?
}
type Chart struct {
	ChartKey

	mode.ChartHeader
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

	Genre    int //string
	Language int //string
	NSFW     bool
	// Tags can be added by user.
	Tags []string
	// Dropped Favorites and Played count since it
	// needs to be checked frequently.

	Pitch bool
}
