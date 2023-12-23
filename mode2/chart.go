package mode

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

const (
	GameModePiano = iota
	GameModeDrum
	GameModeSing
	GameModeAll = -1
)

type Chart interface {
	Duration() int32
	Difficulties() []float64
}

func LoadChartFile(fsys fs.FS, name string) (format any, hash [16]byte, err error) {
	f, err := fsys.Open(name)
	if err != nil {
		return
	}

	switch filepath.Ext(name) {
	case ".osu", ".OSU":
		format, err = osu.NewFormat(f)
		if err != nil {
			return
		}
	}
	hash, _ = Hash(f)
	return
}

// ChartHeader contains non-play information.
// Changing ChartHeader's data will not affect integrity of the chart.
// GameMode-specific fields are located to each Chart struct.
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

	PreviewTime        int32
	MusicFilename      string // Filename is fine to use. (cf. FileName; Filepath)
	BackgroundFilename string
	VideoFilename      string
	VideoTimeOffset    int32

	GameMode    int
	SubGameMode int

	// Hash works as id in database.
	// Hash is not exported to file.
	ChartHash [16]byte // MD5

	// MusicHash is used to check for music updates.
	// A player may replace the music file with another,
	// such as a higher-quality version.
	MusicHash [16]byte // MD5
}

func NewChartHeader(format any, hash [16]byte) (c ChartHeader) {
	switch format := format.(type) {
	case *osu.Format:
		c = newChartHeaderFromOsu(format)
	}
	c.ChartHash = hash
	return
}
func newChartHeaderFromOsu(format *osu.Format) (c ChartHeader) {
	const unknownID = -1
	c = ChartHeader{
		SetID: int32(format.BeatmapSetID),
		ID:    int32(format.BeatmapID),

		MusicName:     format.Title,
		MusicUnicode:  format.TitleUnicode,
		Artist:        format.Artist,
		ArtistUnicode: format.ArtistUnicode,
		MusicSource:   format.Source,
		ChartName:     format.Version,
		Charter:       format.Creator,
		CharterID:     unknownID,
		HolderID:      unknownID,
		Tags:          format.Tags,

		PreviewTime:   int32(format.PreviewTime),
		MusicFilename: format.AudioFilename,
	}

	var e osu.Event
	e, _ = format.Background()
	c.BackgroundFilename = e.Filename
	e, _ = format.Video()
	c.VideoFilename = e.Filename
	c.VideoTimeOffset = int32(e.StartTime)
	if c.MusicFilename == "virtual" {
		c.MusicFilename = ""
	}

	c.GameMode = -1
	switch format.Mode {
	case osu.ModeStandard:
	case osu.ModeTaiko:
		c.GameMode = GameModeDrum
	case osu.ModeCatch:
	case osu.ModeMania:
		c.GameMode = GameModePiano
		c.SubGameMode = int(format.CircleSize)
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
