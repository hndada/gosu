package gosu

import (
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

// ChartHeader contains non-play information.
// Chaning ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
type ChartHeader struct {
	ChartID       int64
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string
	Charter       string
	HolderID      int64

	PreviewTime     int64
	MusicFilename   string
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64
}

func NewChartHeader(f any) (c ChartHeader) {
	switch f := f.(type) {
	case *osu.Format:
		c = ChartHeader{
			MusicName:     f.Title,
			MusicUnicode:  f.TitleUnicode,
			Artist:        f.Artist,
			ArtistUnicode: f.ArtistUnicode,
			MusicSource:   f.Source,
			ChartName:     f.Version,
			Charter:       f.Creator,

			PreviewTime:   int64(f.PreviewTime),
			MusicFilename: f.AudioFilename,
		}
		var e osu.Event
		e, _ = f.Background()
		c.ImageFilename = e.Filename
		e, _ = f.Video()
		c.VideoFilename, c.VideoTimeOffset = e.Filename, int64(e.StartTime)
	}
	return c
}

func (c ChartHeader) MusicPath(cpath string) (string, bool) {
	if name := c.MusicFilename; name == "virtual" || name == "" {
		return "", false
	}
	return filepath.Join(filepath.Dir(cpath), c.MusicFilename), true
}
func (c ChartHeader) BackgroundPath(cpath string) string {
	return filepath.Join(filepath.Dir(cpath), c.ImageFilename)
}
