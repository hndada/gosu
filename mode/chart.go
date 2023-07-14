package mode

import "math"

const (
	ModePiano = iota
	ModeDrum
	ModeSing
)

type Chart interface {
	Duration() int32
	Difficulties() []float64
}

// BPM with longest duration will be main BPM.
// When there are multiple BPMs with same duration, larger one will be main BPM.
func BPMs(ds []*Dynamic, duration int32) (main, min, max float64) {
	bpmDurations := make(map[float64]int32)
	for i, d := range ds {
		if i == 0 {
			bpmDurations[d.BPM] += d.Time
		}
		if i < len(ds)-1 {
			bpmDurations[d.BPM] += ds[i+1].Time - d.Time
		} else {
			bpmDurations[d.BPM] += duration - d.Time // Bounds to final note time; confirmed with test.
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
