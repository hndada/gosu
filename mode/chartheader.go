package mode

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/hndada/gosu/format/osu"
)

// ChartHeader contains non-play information.
// Changing ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
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

	// Hash works as id in database.
	// Hash is not exported to file.
	ChartHash [16]byte // MD5

	// MusicHash is used to check for music updates.
	// A player may replace the music file with another,
	// such as a higher-quality version.
	MusicHash [16]byte // MD5
}

func NewChartHeader(f any) (c ChartHeader) {
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

func Hash(r io.Reader) ([16]byte, error) {
	dat, err := io.ReadAll(r)
	if err != nil {
		return [16]byte{}, err
	}
	return md5.Sum(dat), nil
}
