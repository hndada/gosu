package gosu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hndada/gosu/format/osu"
)

type Chart struct {
	ChartHeader
	TransPoints []*TransPoint
	Mode        int
	SubMode     int // e.g., KeyCount. // Todo: int -> float64; CircleSize may be float64
	Duration    int64
	NoteCounts  []int

	Notes []*Note
	Bars  []LaneSubject

	// SpeedScale float64
}

func NewChart(cpath string, mode, subMode int) (*Chart, error) {
	var c Chart
	dat, err := os.ReadFile(cpath)
	if err != nil {
		return nil, err
	}
	var f any
	switch strings.ToLower(filepath.Ext(cpath)) {
	case ".osu":
		f, err = osu.Parse(dat)
		if err != nil {
			return nil, err
		}
	}

	c.ChartHeader = NewChartHeader(f)
	c.TransPoints = NewTransPoints(f)
	c.Mode = mode
	c.SubMode = subMode
	c.Notes = NewNotes(f, c.TransPoints, mode, subMode)
	if len(c.Notes) > 0 {
		c.Duration = c.Notes[len(c.Notes)-1].Time
	}
	c.Bars = NewBars(c.TransPoints, c.Duration)
	c.NoteCounts = make([]int, 3) // Todo: general note counting
	for _, n := range c.Notes {
		switch n.Type {
		case Normal, Head:
			c.NoteCounts[n.Type]++
		}
	}
	// In Piano mode, Positions should be divided by main BPM for scaling.
	switch c.Mode {
	case ModePiano4, ModePiano7:
		mainBPM, _, _ := BPMs(c.TransPoints, c.Duration)
		for i := range c.TransPoints {
			c.TransPoints[i].Position /= mainBPM
		}
		for i := range c.Notes {
			c.Notes[i].Position /= mainBPM
		}
		for i := range c.Bars {
			c.Bars[i].Position /= mainBPM
		}
	}
	return &c, nil
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
	Path    string
	Mods    Mods
	Header  ChartHeader
	Mode    int
	SubMode int
	Level   float64

	Duration   int64
	NoteCounts []int
	MainBPM    float64
	MinBPM     float64
	MaxBPM     float64
	// Tags       []string // Auto-generated or User-defined
}

func NewChartInfo(c *Chart, cpath string, level float64) ChartInfo {
	mainBPM, minBPM, maxBPM := BPMs(c.TransPoints, c.Duration)
	cb := ChartInfo{
		Path:    cpath,
		Header:  c.ChartHeader,
		Mode:    c.Mode,
		SubMode: c.SubMode,
		Level:   level,

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
func (c *Chart) SetSpeedScale(speedScale float64) {
	for i := range c.Notes {
		c.Notes[i].Position /= c.SpeedScale
		c.Notes[i].Position *= speedScale
	}
	c.SpeedScale = speedScale

	// for i, bar := range c.Bars {
	// 	for tp.Next != nil && (tp.Time < bar.Time || tp.Time >= tp.Next.Time) {
	// 		tp = tp.Next
	// 	}
	// 	bpmRatio := tp.BPM / mainBPM
	// 	beatLength := bpmRatio * tp.BeatLengthScale
	// 	duration := float64(bar.Time - tp.Time)
	// 	position := tp.Position + duration*beatLength
	// 	c.Bars[i].Position = speedScale * position
	// }
	// }()

	// var distance float64 // Approaching notes have positive distance, vice versa.
	// tp := s.TransPoint
	// cursor := s.Time()
	// if time-s.Time() > 0 {
	// 	// When there are more than 2 TransPoint in bounded time.
	// 	for ; tp.Next != nil && tp.Next.Time < time; tp = tp.Next {
	// 		duration := tp.Next.Time - cursor
	// 		bpmRatio := tp.BPM / s.MainBPM
	// 		distance += s.SpeedScale * (bpmRatio * tp.BeatLengthScale) * float64(duration)
	// 		cursor += duration
	// 	}
	// } else {
	// 	for ; tp.Prev != nil && tp.Time > time; tp = tp.Prev {
	// 		duration := tp.Time - cursor // Negative value.
	// 		bpmRatio := tp.BPM / s.MainBPM
	// 		distance += s.SpeedScale * (bpmRatio * tp.BeatLengthScale) * float64(duration)
	// 		cursor += duration
	// 	}
	// }
	// bpmRatio := tp.BPM / s.MainBPM
	// // Calculate the remained (which is farthest from Hint within bound).
	// distance += s.SpeedScale * (bpmRatio * tp.BeatLengthScale) * float64(time-cursor)
	// return HitPosition - distance
}
