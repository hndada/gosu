package plays

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/util"
)

const (
	ModePiano = iota
	ModeDrum
	ModeSing
	ModeAll = -1
)

type ChartFormat any

// *osu.Format
type Chart interface {
	// chart header
	WindowTitle() string

	// dynamics
	Dynamics() []Dynamic
	BPMs() (main, min, max float64)

	// notes
	NoteCounts() []int
	TotalDuration() int32 // Span()
}

// ChartHeader contains non-play information.
// Changing ChartHeader's data will not affect integrity of the chart.
// Play mode-specific fields are located to each Chart struct.
type ChartHeader struct {
	SetID int32 // Compatibility with osu.
	ID    int32 // Compatibility with osu.

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

	Mode    int
	SubMode int

	// Hash works as id in database.
	// Hash is not exported to file.
	// ChartHash [16]byte // MD5
	ChartHash string

	// MusicHash is used to check for music updates.
	// A player may replace the music file with another,
	// such as a higher-quality version.
	MusicHash string
}

// internal game packages use chart format.
func LoadChartFormat(fsys fs.FS, name string) (any, string, error) {
	data, err := util.ReadFile(fsys, name)
	if err != nil {
		return nil, "", err
	}
	hash := util.MD5(data)

	switch filepath.Ext(name) {
	case ".osu", ".OSU":
		format, err := osu.NewFormat(data)
		if err != nil {
			return nil, "", err
		}
		return format, hash, nil
	}
	return nil, "", fmt.Errorf("unsupported file format")
}

// scene select use NewChartHeaderFromFile.
func NewChartHeaderFromFile(fsys fs.FS, name string) (*ChartHeader, error) {
	format, hash, err := LoadChartFormat(fsys, name)
	if err != nil {
		return nil, err
	}

	switch format := format.(type) {
	case *osu.Format:
		c := newChartHeaderFromOsu(format)
		c.ChartHash = hash
		return c, nil
	}

	return nil, fmt.Errorf("unsupported file format")
}

// game piano use NewChartHeaderFromFormat.
func NewChartHeaderFromFormat(format any, hash string) *ChartHeader {
	switch format := format.(type) {
	case *osu.Format:
		c := newChartHeaderFromOsu(format)
		c.ChartHash = hash
		return c
	}
	return nil
}

func newChartHeaderFromOsu(format *osu.Format) (c *ChartHeader) {
	const unknownID = -1
	c = &ChartHeader{
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

	c.Mode = -1
	switch format.Mode {
	case osu.ModeStandard:
	case osu.ModeTaiko:
		c.Mode = ModeDrum
	case osu.ModeCatch:
	case osu.ModeMania:
		c.Mode = ModePiano
		c.SubMode = int(format.CircleSize)
	}
	return
}

func (c ChartHeader) WindowTitle() string {
	return fmt.Sprintf("gosu | %s - %s [%s] (%s) ", c.Artist, c.MusicName, c.ChartName, c.Charter)
}
