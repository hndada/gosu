package mode

import (
	"github.com/hndada/gosu/db"
	"github.com/hndada/gosu/format/osu"
)

type Chart interface {
	WindowTitle() string
}
type ChartHeader = db.ChartHeader

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
		switch f.Mode {
		case osu.ModeMania:
			c.Mode = ModePiano
			c.SubMode = int(f.CircleSize)
		case osu.ModeTaiko:
			c.Mode = ModeDrum
			c.SubMode = 4
		}
	}
	return c
}
