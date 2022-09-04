package gosu

import (
	"math"
	"sort"

	"github.com/hndada/gosu/format/osu"
)

type TransPoint struct {
	Time      int64
	BPM       float64
	Speed     float64
	Meter     int
	NewBeat   bool    // NewBeat draws a bar.
	Volume    float64 // Range is [0, 1].
	Highlight bool

	Position float64
	Next     *TransPoint
	Prev     *TransPoint
}

// First BPM is used as temporary main BPM.
// No two TransPoints have same Time.
func NewTransPoints(f any) []*TransPoint {
	var transPoints []*TransPoint
	switch f := f.(type) {
	case *osu.Format:
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
		transPoints = make([]*TransPoint, 0, len(f.TimingPoints))
		prevBPM := tempMainBPM
		for _, timingPoint := range f.TimingPoints {
			tp := &TransPoint{
				Time:      int64(timingPoint.Time),
				BPM:       prevBPM,
				Speed:     prevBPM / tempMainBPM,
				Meter:     timingPoint.Meter,
				NewBeat:   timingPoint.Uninherited,
				Volume:    float64(timingPoint.Volume) / 100,
				Highlight: timingPoint.IsKiai(),
			}
			if timingPoint.Uninherited {
				tp.BPM = timingPoint.BPM()
				tp.Speed = tp.BPM / tempMainBPM
			} else {
				tp.Speed *= timingPoint.BeatLengthScale()
			}
			if len(transPoints) > 0 && transPoints[len(transPoints)-1].Time == tp.Time { // Drop a TransPoint with a same time
				tp.NewBeat = transPoints[len(transPoints)-1].NewBeat || tp.NewBeat
				transPoints = transPoints[:len(transPoints)-1]
			}
			transPoints = append(transPoints, tp)
			prevBPM = tp.BPM
		}
	}
	var prev *TransPoint
	for _, tp := range transPoints {
		tp.Prev = prev
		if prev != nil {
			prev.Next = tp
		}
		prev = tp
	}
	return transPoints
}

func (tp TransPoint) BeatDuration() float64 {
	return float64(tp.Meter) * (60000 / tp.BPM)
}
func (tp *TransPoint) FetchByTime(time int64) *TransPoint {
	for tp.Next != nil && time >= tp.Next.Time {
		tp = tp.Next
	}
	return tp
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
