package mode

import (
	"math"
	"sort"

	"github.com/hndada/gosu/format/osu"
)

// Except Volume, all fields in Dynamic are related in beat, or 'Pace'.
// Tempo: Allergo, Adagio
// Rhythm: confusing with pattern
// Measure: aka BPM
// Meter: confusing with field 'Meter'
type Dynamic struct {
	Time    int32
	BPM     float64
	Speed   float64
	Meter   int
	NewBeat bool // NewBeat draws a bar.

	Volume    float64 // Used when sample volume is 0.
	Highlight bool
	Position  float64
}

func (d Dynamic) BeatDuration() float64 {
	return float64(d.Meter) * (60000 / d.BPM)
}

// No two Dynamics have same Time.
type Dynamics struct {
	Dynamics []Dynamic
	idx      int
}

func NewDynamics(format any) (Dynamics, error) {
	var ds []Dynamic
	switch format := format.(type) {
	case *osu.Format:
		ds = newDynamicListFromOsu(format)
	}
	return Dynamics{ds, 0}, nil
}

// When gathering Dynamics from osu.Format, it should input the whole slice.
// It is because osu.Format.TimingPoints brings some value from previous TimingPoint.
// First BPM is used as temporary main BPM.

func newDynamicListFromOsu(f *osu.Format) (ds []Dynamic) {
	sort.SliceStable(f.TimingPoints, func(i int, j int) bool {
		if f.TimingPoints[i].Time == f.TimingPoints[j].Time {
			return f.TimingPoints[i].Uninherited
		}
		return f.TimingPoints[i].Time < f.TimingPoints[j].Time
	})
	// Inherited points without Uninherited points will go dropped.
	for len(f.TimingPoints) > 0 && !f.TimingPoints[0].Uninherited {
		f.TimingPoints = f.TimingPoints[1:]
	}
	if len(f.TimingPoints) == 0 {
		return ds
	}

	tempMainBPM := f.TimingPoints[0].BPM()
	ds = make([]Dynamic, 0, len(f.TimingPoints))
	prevBPM := tempMainBPM
	for _, timingPoint := range f.TimingPoints {
		d := Dynamic{
			Time:      int32(timingPoint.Time),
			BPM:       prevBPM,
			Speed:     prevBPM / tempMainBPM,
			Meter:     timingPoint.Meter,
			NewBeat:   timingPoint.Uninherited,
			Volume:    float64(timingPoint.Volume) / 100,
			Highlight: timingPoint.IsKiai(),
		}
		if timingPoint.Uninherited {
			d.BPM = timingPoint.BPM()
			d.Speed = d.BPM / tempMainBPM
		} else {
			d.Speed *= timingPoint.BeatLengthScale()
		}

		// Drop a Dynamic with a same time.
		if len(ds) > 0 && ds[len(ds)-1].Time == d.Time {
			// Either one makes Dynamic a NewBeat.
			d.NewBeat = ds[len(ds)-1].NewBeat || d.NewBeat
			ds = ds[:len(ds)-1]
		}
		ds = append(ds, d)
		prevBPM = d.BPM
	}
	return ds
}

func (ds Dynamics) Speed() float64  { return ds.Dynamics[ds.idx].Speed }
func (ds Dynamics) Volume() float64 { return ds.Dynamics[ds.idx].Volume }
func (ds Dynamics) Highlight() bool { return ds.Dynamics[ds.idx].Highlight }

// Used in ScenePlay for fetching next Dynamic.
func (ds *Dynamics) Update(now int32) {
	for i := ds.idx; i < len(ds.Dynamics); i++ {
		if ds.Dynamics[i].Time > now {
			break
		}
		// index should be updated outside of if block.
		ds.idx = i
	}
}

func (ds Dynamics) BeatTimes(chartDuration int32) (times []int32) {
	// These variables are for iterating over the Time.
	var start, end, step float64
	const bufferTime = 5000

	// times before first Dynamic
	start = float64(ds.Dynamics[0].Time)
	end = start
	if end > -bufferTime {
		end = -bufferTime
	}
	step = ds.Dynamics[0].BeatDuration()
	for t := start; t >= end; t -= step {
		times = append([]int32{int32(t)}, times...)
	}
	// Need to drop a last element because it will be duplicated.
	times = times[:len(times)-1]

	// times after first Dynamic
	var newDys []Dynamic
	for _, d := range ds.Dynamics {
		if d.NewBeat {
			newDys = append(newDys, d)
		}
	}

	for i, nd := range newDys {
		start = float64(nd.Time)
		if i == len(newDys)-1 {
			end = float64(chartDuration + bufferTime)
		} else {
			end = float64(newDys[i+1].Time)
		}
		step = nd.BeatDuration()
		for t := start; t < end; t += step {
			times = append(times, int32(t))
		}
	}
	return
}

// BPM with longest duration will be main BPM.
// When there are multiple BPMs with same duration, larger one will be main BPM.
func (ds Dynamics) BPMs(duration int32) (main, min, max float64) {
	bpmDurations := make(map[float64]int32)
	for i, d := range ds.Dynamics {
		if i == 0 {
			bpmDurations[d.BPM] += d.Time
		}
		if i < len(ds.Dynamics)-1 {
			bpmDurations[d.BPM] += ds.Dynamics[i+1].Time - d.Time
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
