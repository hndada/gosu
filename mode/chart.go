package mode

import (
	"fmt"

	"github.com/hndada/gosu/format/osu"
)

type Chart interface {
	WindowTitle() string
}

// ChartHeader contains non-play information.
// Changing ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
type ChartHeader struct {
	SetID         int64 // Compatibility for osu.
	ID            int64
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
		if c.MusicFilename == "virtual" {
			c.MusicFilename = ""
		}
	}
	return c
}
