package db

import (
	"fmt"
	"time"
)

// ChartHeader contains non-play information.
// Changing ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
type ChartHeader struct {
	ChartSetID    int64 // Compatibility for osu.
	ChartID       int64 // Todo: ChartID -> ID
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string
	Charter       string
	HolderID      int64

	PreviewTime     int64
	MusicFilename   string // Filename is fine to use (cf. Filepath)
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64

	Mode    int
	SubMode int
}

func (c ChartHeader) WindowTitle() string {
	return fmt.Sprintf("gosu | %s - %s [%s] (%s) ",
		c.Artist, c.MusicName, c.ChartName, c.Charter)
}

// Todo: should MD5 be substituded with SHA256?
type ChartKey struct {
	MD5  [16]byte
	Mods interface{} // Todo: use mods code for mode-specific mods?
}
type Chart struct {
	ChartKey

	ChartHeader
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

// func (c ChartHeader) MusicPath(cpath string) (string, bool) {
// 	if name := c.MusicFilename; name == "virtual" || name == "" {
// 		return "", false
// 	}
// 	return filepath.Join(filepath.Dir(cpath), c.MusicFilename), true
// }
// func (c ChartHeader) BackgroundPath(cpath string) string {
// 	return filepath.Join(filepath.Dir(cpath), c.ImageFilename)
// }

// https://osu.ppy.sh/docs/index.html#beatmapsetcompact-covers
// cover, card, list, slimcover

// rnkaed status
// graveyard, wip, pending
// ranked, approved, qualified, loved
