package mode

import (
	"fmt"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

type BaseChart struct {
	ChartHeader
	TransPoints []*TransPoint
	ModeType
	SubMode    int // e.g., KeyCount. // Todo: int -> float64; CircleSize may be float64
	Duration   int64
	NoteCounts []int
}

// ChartHeader contains non-play information.
// Chaning ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
// Music: when treating it as media
// Audio: when considering as programming aspect
type ChartHeader struct {
	ChartID       int64
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string
	Producer      string // Name of field may change.
	HolderID      int64

	AudioFilename   string
	PreviewTime     int64
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64
}

func NewChartHeader(f any) ChartHeader {
	var c ChartHeader
	switch f := f.(type) {
	case *osu.Format:
		c = ChartHeader{
			MusicName:     f.Title,
			MusicUnicode:  f.TitleUnicode,
			Artist:        f.Artist,
			ArtistUnicode: f.ArtistUnicode,
			MusicSource:   f.Source,
			ChartName:     f.Version,
			Producer:      f.Creator,

			AudioFilename: f.AudioFilename,
			PreviewTime:   int64(f.PreviewTime),
		}
		var e osu.Event
		e, _ = f.Background()
		c.ImageFilename = e.Filename
		e, _ = f.Video()
		c.VideoFilename, c.VideoTimeOffset = e.Filename, int64(e.StartTime)
	}
	return c
}

func (c ChartHeader) MusicPath(cpath string) string {
	return filepath.Join(filepath.Dir(cpath), c.AudioFilename)
}
func (c ChartHeader) BackgroundPath(cpath string) string {
	return filepath.Join(filepath.Dir(cpath), c.ImageFilename)
}

// ChartInfo is used at SceneSelect.
type ChartInfo struct {
	Path   string
	Mods   Mods
	Header ChartHeader
	ModeType
	SubMode int
	Level   float64

	Duration   int64
	NoteCounts []int
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
	// Tags       []string // Auto-generated or User-defined
}

func NewChartInfo(c *BaseChart, cpath string, level float64) ChartInfo {
	mainBPM, minBPM, maxBPM := BPMs(c.TransPoints, c.Duration)
	cb := ChartInfo{
		Path:     cpath,
		Header:   c.ChartHeader,
		ModeType: c.ModeType,
		SubMode:  c.SubMode,
		Level:    level,

		Duration:   c.Duration,
		NoteCounts: c.NoteCounts,
		MainBPM:    mainBPM,
		MinBPM:     minBPM,
		MaxBPM:     maxBPM,
	}
	return cb
}

func (c ChartInfo) Text() string {
	return fmt.Sprintf("(%dK Lv %.1f) %s [%s]", c.SubMode, c.Level, c.Header.MusicName, c.Header.ChartName)
}
