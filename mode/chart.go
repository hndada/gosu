package mode

import (
	"fmt"
	"math"
	"sort"

	"github.com/hndada/gosu/format/osu"
)

type Chart interface {
	Difficulties() []float64
}

// ChartHeader contains non-play information.
// Changing ChartHeader's data will not affect integrity of the chart.
// Mode-specific fields are located to each Chart struct.
// Todo: add AudioHash
type ChartHeader struct {
	SetID int64 // Compatibility for osu.
	ID    int64 // Compatibility for osu.

	MusicName     string
	MusicUnicode  string
	Artist        string
	ArtistUnicode string
	MusicSource   string
	ChartName     string
	Charter       string
	CharterID     int64
	HolderID      int64 // When the chart is uploaded by non-charter.
	Tags          []string

	PreviewTime     int64
	MusicFilename   string // Filename is fine to use. (cf. FileName; Filepath)
	ImageFilename   string
	VideoFilename   string
	VideoTimeOffset int64

	Mode    int
	SubMode int
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
		SetID: int64(f.BeatmapSetID),
		ID:    int64(f.BeatmapID),

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

		PreviewTime:   int64(f.PreviewTime),
		MusicFilename: f.AudioFilename,
	}

	var e osu.Event
	e, _ = f.Background()
	c.ImageFilename = e.Filename
	e, _ = f.Video()
	c.VideoFilename = e.Filename
	c.VideoTimeOffset = int64(e.StartTime)
	if c.MusicFilename == "virtual" {
		c.MusicFilename = ""
	}

	c.Mode = -1
	switch f.Mode {
	case osu.ModeStandard:
	case osu.ModeTaiko:
		c.Mode = 1
	case osu.ModeCatch:
	case osu.ModeMania:
		c.Mode = 0
		c.SubMode = int(f.CircleSize)
	}
	return
}

func (c ChartHeader) WindowTitle() string {
	return fmt.Sprintf("gosu | %s - %s [%s] (%s) ", c.Artist, c.MusicName, c.ChartName, c.Charter)
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

func Level(c Chart) float64 {
	const decayFactor = 0.95

	ds := c.Difficulties()
	sort.Slice(ds, func(i, j int) bool { return ds[i] > ds[j] })

	sum, weight := 0.0, 1.0
	for _, term := range ds {
		sum += weight * term
		weight *= decayFactor
	}

	// No additional Math.Pow; it would make a little change.
	return sum
}

// 오프셋 부터
// Bar 그리는 것을 체크
func BeatDurations(tps []*TransPoint) []float64 {
	durations := make([]float64, 0, 10)
	for i, tp := range tps {
		if i < len(tps)-1 {
			durations[i] = tps[i+1].Time - tp.Time
		} else {
			durations[i] = 0
		}
	}
	return durations
}
