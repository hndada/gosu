package mode

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/db"
	"github.com/hndada/gosu/format/osu"
)

type Header = db.Header

func NewHeader(f any) (c Header) {
	switch f := f.(type) {
	case *osu.Format:
		c = Header{
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
		// switch f.Mode {
		// case osu.ModeMania:
		// 	keyCount := int(f.CircleSize)
		// 	if keyCount <= 4 {
		// 		c.Mode = ModePiano4
		// 	} else {
		// 		c.Mode = ModePiano7
		// 	}
		// 	c.SubMode = keyCount
		// case osu.ModeTaiko:
		// 	c.Mode = ModeDrum
		// 	c.SubMode = 4
		// }
	}
	return c
}

func SetTitle(c Header) {
	title := fmt.Sprintf("gosu | %s - %s [%s] (%s) ", c.Artist, c.MusicName, c.ChartName, c.Charter)
	ebiten.SetWindowTitle(title)
}
