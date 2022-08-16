package mode

import (
	"math"
	"sort"

	"github.com/hndada/gosu/format/osu"
)

type TransPoint struct {
	Time         int64
	BPM          float64
	BeatScale    float64
	Meter        uint8
	Volume       uint8
	Highlight    bool
	NewBPM       bool
	Prev         *TransPoint
	Next         *TransPoint
	NextBPMPoint *TransPoint // For performance
}

// Uninherited point is base point. whereas, Inherited point 'inherits'
// values from previous Uninherited point first, then after.
// All osu format have at least one Uninherited timing point.
// Uninherited: BPM, Inherited: beat scale
// Initial BPM is derived from the first timing point's BPM.
// When there are multiple timing points with same time, the last one will overwrites all precedings.
func NewTransPointsFromOsu(o *osu.Format) []*TransPoint {
	sort.SliceStable(o.TimingPoints, func(i int, j int) bool {
		if o.TimingPoints[i].Time == o.TimingPoints[j].Time {
			return o.TimingPoints[i].Uninherited
		}
		return o.TimingPoints[i].Time < o.TimingPoints[j].Time
	})
	tps := make([]*TransPoint, 0, len(o.TimingPoints))
	lastBPM, _ := o.TimingPoints[0].BPM()
	var prev *TransPoint
	var prevBPMPoint *TransPoint
	for _, timingPoint := range o.TimingPoints {
		if timingPoint.Uninherited {
			tp := &TransPoint{
				Time: int64(timingPoint.Time),
				// BPM:,
				BeatScale: 1,
				Meter:     uint8(timingPoint.Meter),
				Volume:    uint8(timingPoint.Volume),
				Highlight: timingPoint.IsKiai(),
				NewBPM:    true,
				Prev:      prev,
			}
			tp.BPM, _ = timingPoint.BPM()

			if prev != nil {
				prev.Next = tp
			}
			if prevBPMPoint != nil {
				prevBPMPoint.NextBPMPoint = tp // This was hard to find the bug to me.
			}
			prev = tp
			prevBPMPoint = tp
			lastBPM = tp.BPM
			tps = append(tps, tp)
		} else {
			tp := &TransPoint{
				Time: int64(timingPoint.Time),
				BPM:  lastBPM,
				// BeatScale: ,
				Meter:     uint8(timingPoint.Meter),
				Volume:    uint8(timingPoint.Volume),
				Highlight: timingPoint.IsKiai(),
				NewBPM:    false,
				Prev:      prev,
			}
			tp.BeatScale, _ = timingPoint.BeatScale()

			prev.Next = tp // Inherited point is never the first.
			prev = tp
			tps = append(tps, tp)
		}
	}
	return tps
}

// BPM with longest duration will be main BPM.
// Suppose when there are multiple BPMs with same duration, larger one will be main.
func BPMs(tps []*TransPoint, endTime int64) (main, min, max float64) {
	bpmDurations := make(map[float64]int64)
	for i, tp := range tps {
		if i == 0 {
			bpmDurations[tp.BPM] += tp.Time
		}
		if i < len(tps)-1 {
			bpmDurations[tp.BPM] += tps[i+1].Time - tp.Time
		} else {
			bpmDurations[tp.BPM] += endTime - tp.Time // Bounds to final note time; confirmed with test.
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

// wb, wa stands for buffer times: wait before, wait after.
// Multiply wa with 2 for preventing indexing a time slice over length.
func BarLineTimes(tps []*TransPoint, endTime int64, wb, wa int64) []int64 {
	ts := make([]int64, 0)
	tp0 := tps[0]
	for t := float64(tp0.Time); t >= float64(wb); t -= float64(tp0.Meter) * 60000 / tp0.BPM {
		ts = append([]int64{int64(t)}, ts...)
	}
	ts = ts[:len(ts)-1] // Drop bar line for tp0 for avoiding duplicated
	for tp := tps[0]; tp != nil; tp = tp.NextBPMPoint {
		next := endTime + 2*wa
		if tp.NextBPMPoint != nil {
			next = tp.NextBPMPoint.Time
		}
		unit := float64(tp.Meter) * 60000 / tp.BPM
		for t := float64(tp.Time); t < float64(next); t += unit {
			ts = append(ts, int64(t))
		}
	}
	// for i, tp := range c.TransPoints {
	// 	if !tp.NewBPM {
	// 		continue
	// 	}
	// 	next := c.EndTime() + 2*wa
	// 	if i < len(c.TransPoints)-1 {
	// 		for _, tp2 := range c.TransPoints[i:] {
	// 			if tp2.NewBPM {
	// 				next = tp2.Time
	// 				break
	// 			}
	// 		}
	// 		next = c.TransPoints[i+1].Time
	// 	}
	// 	unit := float64(tp.Meter) * 60000 / tp.BPM
	// 	for t := float64(tp.Time); t < float64(next); t += unit {
	// 		ts = append(ts, int64(t))
	// 	}
	// }
	return ts
}
