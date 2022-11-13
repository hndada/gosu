package chart

import (
	"github.com/hndada/gosu/game/format/osu"
)

// Header contains non-play information.
// Changing Header's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
type Header struct {
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

// func (c Header) MusicPath(cpath string) (string, bool) {
// 	if name := c.MusicFilename; name == "virtual" || name == "" {
// 		return "", false
// 	}
// 	return filepath.Join(filepath.Dir(cpath), c.MusicFilename), true
// }
// func (c Header) BackgroundPath(cpath string) string {
// 	return filepath.Join(filepath.Dir(cpath), c.ImageFilename)
// }
