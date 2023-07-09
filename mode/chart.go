package mode

import (
	"crypto/md5"
	"io/fs"
	"math"
	"path/filepath"
	"sort"

	"github.com/hndada/gosu/format/osu"
)

const (
	ModePiano = iota
	ModeDrum
	ModeSing
)

type Chart interface {
	Duration() int64
	Difficulties() []float64
}

func ParseChartFile(fsys fs.FS, name string) (format any, hash [16]byte, err error) {
	var dat []byte
	dat, err = fs.ReadFile(fsys, name)
	if err != nil {
		return
	}
	hash = md5.Sum(dat)

	switch filepath.Ext(name) {
	case ".osu", ".OSU":
		format, err = osu.Parse(dat)
		if err != nil {
			return
		}
	}
	return
}

// BPM with longest duration will be main BPM.
// When there are multiple BPMs with same duration, larger one will be main BPM.
func BPMs(transPoints []*Dynamic, duration int64) (main, min, max float64) {
	bpmDurations := make(map[float64]int64)
	for i, dy := range transPoints {
		if i == 0 {
			bpmDurations[dy.BPM] += dy.Time
		}
		if i < len(transPoints)-1 {
			bpmDurations[dy.BPM] += transPoints[i+1].Time - dy.Time
		} else {
			bpmDurations[dy.BPM] += duration - dy.Time // Bounds to final note time; confirmed with test.
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
