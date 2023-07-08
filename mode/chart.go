package mode

import (
	"fmt"
	"math"

	"github.com/hndada/gosu/format/osu"
)

// ChartHeader contains non-play information.
// Changing ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
type ChartHeader struct {
	SetID         int64 // Compatibility for osu.
	ID            int64
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

func (c ChartHeader) WindowTitle() string {
	return fmt.Sprintf("gosu | %s - %s [%s] (%s) ", c.Artist, c.MusicName, c.ChartName, c.Charter)
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
	}
	return c
}

// BPM with longest duration will be main BPM.
// When there are multiple BPMs with same duration, larger one will be main BPM.
func BPMs(transPoints []*TransPoint, duration int64) (main, min, max float64) {
	bpmDurations := make(map[float64]int64)
	for i, tp := range transPoints {
		if i == 0 {
			bpmDurations[tp.BPM] += tp.Time
		}
		if i < len(transPoints)-1 {
			bpmDurations[tp.BPM] += transPoints[i+1].Time - tp.Time
		} else {
			bpmDurations[tp.BPM] += duration - tp.Time // Bounds to final note time; confirmed with test.
		}
	}
	var maxDuration int64
	min = math.MaxFloat64
	for bpm, duration := range bpmDurations {
		if maxDuration < duration {
			maxDuration = duration
			main = bpm
		} else if maxDuration == duration && main < bpm {
			main = bpm
		}
		if min > bpm {
			min = bpm
		}
		if max < bpm {
			max = bpm
		}
	}
	return
}

type Chart struct{}

func (c Chart) FirstBeatNote() {

}
