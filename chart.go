package gosu

import (
	"path/filepath"
	"sort"

	"github.com/hndada/gosu/parse/osu"
)

const (
	ModeStandard = iota
	ModeTaiko
	ModeCatch
	ModeMania
)

const DefaultMode = ModeStandard

// ChartHeader contains non-play information.
// Chaning ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
type ChartHeader struct {
	ChartID       int64 // 6byte: setID, 2byte: subID
	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string // diff name
	Producer      string // Name of field may change
	HolderID      int64  // 0: gosu Chart Management

	AudioFilename   string
	PreviewTime     int64
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64
}

// Chart should avoid redundant data as much as possible
type Chart struct {
	ChartHeader
	TransPoints []*TransPoint
	KeyCount    int
	Notes       []Note
	// Level    float64
	// ScratchMode int
}

func NewChartFromOsu(o *osu.Format) (*Chart, error) {
	c := &Chart{
		ChartHeader: NewChartHeaderFromOsu(o),
		TransPoints: NewTransPointsFromOsu(o),
		KeyCount:    int(o.CircleSize),
		Notes:       make([]Note, 0, len(o.HitObjects)*2),
	}
	for _, ho := range o.HitObjects {
		c.Notes = append(c.Notes, NewNoteFromOsu(ho, c.KeyCount)...)
	}
	sort.Slice(c.Notes, func(i, j int) bool {
		if c.Notes[i].Time == c.Notes[j].Time {
			return c.Notes[i].Key < c.Notes[j].Key
		}
		return c.Notes[i].Time < c.Notes[j].Time
	})
	return c, nil
}
func NewChartHeaderFromOsu(o *osu.Format) ChartHeader {
	c := ChartHeader{
		MusicName:     o.Title,
		MusicUnicode:  o.TitleUnicode,
		Artist:        o.Artist,
		ArtistUnicode: o.ArtistUnicode,
		MusicSource:   o.Source,
		ChartName:     o.Version,
		Producer:      o.Creator,

		AudioFilename: o.AudioFilename,
		PreviewTime:   int64(o.PreviewTime),
	}
	var e osu.Event
	e, _ = o.Background()
	c.ImageFilename = e.Filename
	e, _ = o.Video()
	c.VideoFilename, c.VideoTimeOffset = e.Filename, int64(e.StartTime)
	return c
}

func (c ChartHeader) MusicPath(cpath string) string {
	return filepath.Join(filepath.Dir(cpath), c.AudioFilename)
}
func (c ChartHeader) BackgroundPath(cpath string) string {
	return filepath.Join(filepath.Dir(cpath), c.ImageFilename)
}
