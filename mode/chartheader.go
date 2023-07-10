package mode

import (
	"fmt"

	"github.com/hndada/gosu/format/osu"
)

// ChartHeader contains non-play information.
// Changing ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.

// If AudioHash was a thing, a player must have to play a chart
// with the certain audio file, which is too strict.
// Hence, AudioHash should be handled via outer way.
type ChartHeader struct {
	SetID int32 // Compatibility for osu.
	ID    int32 // Compatibility for osu.

	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string
	Charter       string
	CharterID     int32
	HolderID      int32 // When the chart is uploaded by non-charter.
	Tags          []string

	PreviewTime     int32
	MusicFilename   string // Filename is fine to use. (cf. FileName; Filepath)
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int32

	Mode    int
	SubMode int
}

func LoadChartHeader(f any) (c ChartHeader) {
	switch f := f.(type) {
	case *osu.Format:
		return newChartHeaderFromOsu(f)
	}
	return
}
func newChartHeaderFromOsu(f *osu.Format) (c ChartHeader) {
	const unknownID = -1
	c = ChartHeader{
		SetID: int32(f.BeatmapSetID),
		ID:    int32(f.BeatmapID),

		MusicName:     f.Title,
		MusicUnicode:  f.TitleUnicode,
		Artist:        f.Artist,
		ArtistUnicode: f.ArtistUnicode,
		MusicSource:   f.Source,
		ChartName:     f.Version,
		Charter:       f.Creator,
		CharterID:     unknownID,
		HolderID:      unknownID,
		Tags:          f.Tags,

		PreviewTime:   int32(f.PreviewTime),
		MusicFilename: f.AudioFilename,
	}

	var e osu.Event
	e, _ = f.Background()
	c.ImageFilename = e.Filename
	e, _ = f.Video()
	c.VideoFilename = e.Filename
	c.VideoTimeOffset = int32(e.StartTime)
	if c.MusicFilename == "virtual" {
		c.MusicFilename = ""
	}

	c.Mode = -1
	switch f.Mode {
	case osu.ModeStandard:
	case osu.ModeTaiko:
		c.Mode = ModeDrum
	case osu.ModeCatch:
	case osu.ModeMania:
		c.Mode = ModePiano
		c.SubMode = int(f.CircleSize)
	}
	return
}

func (c ChartHeader) WindowTitle() string {
	return fmt.Sprintf("gosu | %s - %s [%s] (%s) ", c.Artist, c.MusicName, c.ChartName, c.Charter)
}
