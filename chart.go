package gosu

import (
	"fmt"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

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

func (c ChartHeader) MusicPath(cpath string) (string, bool) {
	if name := c.AudioFilename; name != "virtual" && name != "" {
		return "", false
	}
	return filepath.Join(filepath.Dir(cpath), c.AudioFilename), true
}
func (c ChartHeader) BackgroundPath(cpath string) string {
	return filepath.Join(filepath.Dir(cpath), c.ImageFilename)
}

// ChartInfo is used at SceneSelect.
type ChartInfo struct {
	Path string
	// Mods    Mods
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

func (c ChartInfo) Text() string {
	return fmt.Sprintf("(%dK Lv %.1f) %s [%s]", c.SubMode, c.Level, c.Header.MusicName, c.Header.ChartName)
}
func (c ChartInfo) BackgroundPath() string {
	return c.Header.BackgroundPath(c.Path)
}

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
