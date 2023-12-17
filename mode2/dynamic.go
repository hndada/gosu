package mode

import (
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

	Position float64
	// Next and Prev are used in each mode too, hence exported.
	Next *Dynamic
	Prev *Dynamic
}

func NewDynamics(f any) []*Dynamic {
	var ds []*Dynamic
	switch f := f.(type) {
	case *osu.Format:
		ds = newDynamicsFromOsu(f)
	}

	// linking
	var prev *Dynamic
	for _, d := range ds {
		d.Prev = prev
		if prev != nil {
			prev.Next = d
		}
		prev = d
	}
	return ds
}

// When gathering Dynamics from osu.Format, it should input the whole slice.
// It is because osu.Format.TimingPoints brings some value from previous TimingPoint.

// First BPM is used as temporary main BPM.
// No two Dynamics have same Time.
func newDynamicsFromOsu(f *osu.Format) []*Dynamic {
	var ds []*Dynamic
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
	ds = make([]*Dynamic, 0, len(f.TimingPoints))
	prevBPM := tempMainBPM
	for _, timingPoint := range f.TimingPoints {
		d := &Dynamic{
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
		// Drop a Dynamic with a same time
		if len(ds) > 0 && ds[len(ds)-1].Time == d.Time {
			// Either one makes Dynamic a NewBeat
			d.NewBeat = ds[len(ds)-1].NewBeat || d.NewBeat
			ds = ds[:len(ds)-1]
		}
		ds = append(ds, d)
		prevBPM = d.BPM
	}
	return ds
}

// 0: Use default meter.
func (d Dynamic) BeatDuration(meter int) float64 {
	m := float64(d.Meter)
	if meter > 0 {
		m = float64(meter)
	}
	return m * (60000 / d.BPM)
}

func BeatTimes(ds []*Dynamic, duration int32, meter int) (times []int32) {
	// These variables are for iterating over the Time.
	var start, end, step float64
	const bufferTime = 5000

	// times before first Dynamic
	start = float64(ds[0].Time)
	end = start
	if end > -bufferTime {
		end = -bufferTime
	}
	step = ds[0].BeatDuration(meter)
	for t := start; t >= end; t -= step {
		times = append([]int32{int32(t)}, times...)
	}
	// Need to drop a last element because it will be duplicated.
	times = times[:len(times)-1]

	// times after first Dynamic
	var newDys []*Dynamic
	for _, d := range ds {
		if d.NewBeat {
			newDys = append(newDys, d)
		}
	}

	for i, nd := range newDys {
		start = float64(nd.Time)
		if i == len(newDys)-1 {
			end = float64(duration + bufferTime)
		} else {
			end = float64(newDys[i+1].Time)
		}
		step = nd.BeatDuration(meter)
		for t := start; t < end; t += step {
			times = append(times, int32(t))
		}
	}
	return
}

// Used in ScenePlay for fetching next Dynamic.
func NextDynamics(d *Dynamic, now int32) *Dynamic {
	for d.Next != nil && now >= d.Next.Time {
		d = d.Next
	}
	return d
}
