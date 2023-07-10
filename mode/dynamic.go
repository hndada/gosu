package mode

import (
	"sort"

	"github.com/hndada/gosu/format/osu"
)

// First BPM is used as temporary main BPM.
// No two Dynamics have same Time.
type Dynamic struct {
	Time      int32
	BPM       float64
	Speed     float64
	Meter     int
	NewBeat   bool    // NewBeat draws a bar.
	Volume    float64 // Range is [0, 1].
	Highlight bool

	Position float64
	Next     *Dynamic
	Prev     *Dynamic
}

func NewDynamics(f any) []*Dynamic {
	var transPoints []*Dynamic
	switch f := f.(type) {
	case *osu.Format:
		transPoints = newDynamicsFromOsu(f)
	}

	var prev *Dynamic
	for _, dy := range transPoints {
		dy.Prev = prev
		if prev != nil {
			prev.Next = dy
		}
		prev = dy
	}
	return nil
}

// When gathering Dynamics from osu.Format, it should input the whole slice.
// It is because osu.Format.TimingPoints brings some value from previous TimingPoint.
func newDynamicsFromOsu(f *osu.Format) []*Dynamic {
	var transPoints []*Dynamic
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
		return transPoints
	}
	tempMainBPM := f.TimingPoints[0].BPM()
	transPoints = make([]*Dynamic, 0, len(f.TimingPoints))
	prevBPM := tempMainBPM
	for _, timingPoint := range f.TimingPoints {
		dy := &Dynamic{
			Time:      int32(timingPoint.Time),
			BPM:       prevBPM,
			Speed:     prevBPM / tempMainBPM,
			Meter:     timingPoint.Meter,
			NewBeat:   timingPoint.Uninherited,
			Volume:    float64(timingPoint.Volume) / 100,
			Highlight: timingPoint.IsKiai(),
		}
		if timingPoint.Uninherited {
			dy.BPM = timingPoint.BPM()
			dy.Speed = dy.BPM / tempMainBPM
		} else {
			dy.Speed *= timingPoint.BeatLengthScale()
		}
		// Drop a Dynamic with a same time
		if len(transPoints) > 0 && transPoints[len(transPoints)-1].Time == dy.Time {
			// Either one makes Dynamic a NewBeat
			dy.NewBeat = transPoints[len(transPoints)-1].NewBeat || dy.NewBeat
			transPoints = transPoints[:len(transPoints)-1]
		}
		transPoints = append(transPoints, dy)
		prevBPM = dy.BPM
	}
	return transPoints
}

// 0: Use default meter.
func (dy Dynamic) BeatDuration(meter int) float64 {
	m := float64(dy.Meter)
	if meter > 0 {
		m = float64(meter)
	}
	return m * (60000 / dy.BPM)
}

func BeatTimes(dys []*Dynamic, duration int32, meter int) (times []int32) {
	// These variables are for iterating over the Time.
	var start, end, step float64
	const bufferTime = 5000

	// times before first Dynamic
	start = float64(dys[0].Time)
	end = start
	if end > -bufferTime {
		end = -bufferTime
	}
	step = dys[0].BeatDuration(meter)
	for t := start; t >= end; t -= step {
		times = append([]int32{int32(t)}, times...)
	}
	// Need to drop a last element because it will be duplicated.
	times = times[:len(times)-1]

	// times after first Dynamic
	var newDys []*Dynamic
	for _, dy := range dys {
		if dy.NewBeat {
			newDys = append(newDys, dy)
		}
	}

	for i, ndy := range newDys {
		start = float64(ndy.Time)
		if i == len(newDys)-1 {
			end = float64(duration + bufferTime)
		} else {
			end = float64(newDys[i+1].Time)
		}
		step = ndy.BeatDuration(meter)
		for t := start; t < end; t += step {
			times = append(times, int32(t))
		}
	}
	return
}
