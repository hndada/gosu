package mode

import (
	"crypto/md5"
	"io/fs"
	"math"
	"path/filepath"

	"github.com/hndada/gosu/format/osu"
)

const (
	ModePiano = iota
	ModeDrum
	ModeSing
)

type Chart interface {
	Duration() int32
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
func BPMs(dys []*Dynamic, duration int32) (main, min, max float64) {
	bpmDurations := make(map[float64]int32)
	for i, dy := range dys {
		if i == 0 {
			bpmDurations[dy.BPM] += dy.Time
		}
		if i < len(dys)-1 {
			bpmDurations[dy.BPM] += dys[i+1].Time - dy.Time
		} else {
			bpmDurations[dy.BPM] += duration - dy.Time // Bounds to final note time; confirmed with test.
		}
	}
	var maxDuration int32
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
