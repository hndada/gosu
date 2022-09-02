package gosu

import (
	"math"
	"sort"

	"github.com/hndada/gosu/format/osu"
)

// BPM means Beats Per Minute. Higher BPM means more beats in unit time.
// Scroll goes faster proportional to BPM (or BPM ratio), since
// length of beat is fixed by default, which can be scaled by BeatLengthScale.
type TransPoint struct {
	Time            int64
	BPM             float64
	BeatLengthScale float64
	Meter           uint8
	Volume          float64 // Range is 0 to 1.
	Highlight       bool
	NewBPM          bool
	Position        float64
	Prev            *TransPoint
	Next            *TransPoint
	NextBPMPoint    *TransPoint // For performance
}

// In osu!, Uninherited point is base point. While Inherited point inherits
// starting from previous Uninherited point, then possible following Inherited points.
// BPM and BeatLengthScale are derived only from Uninherited and Inherited each.

// All .osu chart have at least one Uninherited point.
// Initial BPM is derived from the first Uninherited point.
// When there are multiple timing points with same time, the last one will overwrites all precedings.
func NewTransPoints(f any, fixed bool) []*TransPoint {
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

		transPoints = make([]*TransPoint, 0, len(f.TimingPoints))
		var prev = &TransPoint{BPM: f.TimingPoints[0].BPM()}
		var prevBPMPoint = prev
		for _, timingPoint := range f.TimingPoints {
			if timingPoint.Uninherited {
				tp := &TransPoint{
					Time:            int64(timingPoint.Time),
					BPM:             timingPoint.BPM(),
					BeatLengthScale: 1,
					Meter:           uint8(timingPoint.Meter),
					Volume:          float64(timingPoint.Volume) / 100,
					Highlight:       timingPoint.IsKiai(),
					NewBPM:          true,
					// Position:        prev.Position,
					Prev: prev,
					// Next:            nil,
					// NextBPMPoint:    nil,
				}
				if fixed {
					tp.Position = prev.Position + prev.Speed()*float64(tp.Time-prev.Time)
				} else {
					tp.Position = tp.Speed() * float64(tp.Time)
				}
				prev.Next = tp
				prev = tp
				prevBPMPoint.NextBPMPoint = tp // This was hard to find the bug to me.
				prevBPMPoint = tp
				transPoints = append(transPoints, tp)
			} else {
				tp := &TransPoint{
					Time:            int64(timingPoint.Time),
					BPM:             prevBPMPoint.BPM,
					BeatLengthScale: timingPoint.BeatLengthScale(),
					Meter:           uint8(timingPoint.Meter),
					Volume:          float64(timingPoint.Volume) / 100,
					Highlight:       timingPoint.IsKiai(),
					NewBPM:          false,
					// Position:        prev.Position,
					Prev: prev,
					// Next:            nil,
					// NextBPMPoint:    nil,
				}
				if fixed {
					tp.Position = prev.Position + prev.Speed()*float64(tp.Time-prev.Time)
				} else {
					tp.Position = tp.Speed() * float64(tp.Time)
				}
				prev.Next = tp // Inherited point is never the first.
				prev = tp
				transPoints = append(transPoints, tp)
			}
		}
	}
	transPoints[0].Prev = nil // First TransPoint's Prev is just a dummy.
	return transPoints
}

func (tp *TransPoint) FetchByTime(time int64) *TransPoint {
	for tp.Next != nil && time >= tp.Next.Time {
		tp = tp.Next
	}
	return tp
}

// FetchPresent is useful for NewBPM TransPoint fetching latest TransPoint.
func (tp *TransPoint) FetchPresent() *TransPoint { return tp.FetchByTime(tp.Time) }

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

func (tp TransPoint) Speed() float64 {
	return tp.BPM * tp.BeatLengthScale
}
func (tp TransPoint) BeatDuration() float64 {
	return float64(tp.Meter) * (60000 / tp.BPM)
}
