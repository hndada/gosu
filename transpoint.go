package gosu

import (
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/hndada/gosu/format/osu"
)

type TransPoint struct {
	Time         int64
	BPM          float64
	BeatScale    float64
	Meter        uint8
	Volume       float64 // Range is 0 to 1.
	Highlight    bool
	NewBPM       bool
	Prev         *TransPoint
	Next         *TransPoint
	NextBPMPoint *TransPoint // For performance
	Position     float64
}

// Uninherited point is base point. whereas, Inherited point 'inherits'
// values from previous Uninherited point first, then after.
// All osu format have at least one Uninherited timing point.
// Uninherited: BPM, Inherited: beat scale
// Initial BPM is derived from the first timing point's BPM.
// When there are multiple timing points with same time, the last one will overwrites all precedings.
func NewTransPoints(f any) []*TransPoint {
	var tps []*TransPoint
	switch f := f.(type) {
	case *osu.Format:
		sort.SliceStable(f.TimingPoints, func(i int, j int) bool {
			if f.TimingPoints[i].Time == f.TimingPoints[j].Time {
				return f.TimingPoints[i].Uninherited
			}
			return f.TimingPoints[i].Time < f.TimingPoints[j].Time
		})

		// Drop inherited points preceding to the first uninherited point.
		// Todo: need actual test
		for len(f.TimingPoints) > 0 && !f.TimingPoints[0].Uninherited {
			f.TimingPoints = f.TimingPoints[1:]
		}
		if len(f.TimingPoints) == 0 {
			return tps
		}
		tps = make([]*TransPoint, 0, len(f.TimingPoints))
		lastBPM, _ := f.TimingPoints[0].BPM()
		var prev *TransPoint
		var prevBPMPoint *TransPoint
		for _, timingPoint := range f.TimingPoints {
			if timingPoint.Uninherited {
				tp := &TransPoint{
					Time: int64(timingPoint.Time),
					// BPM:,
					BeatScale: 1,
					Meter:     uint8(timingPoint.Meter),
					Volume:    float64(timingPoint.Volume) / 100,
					Highlight: timingPoint.IsKiai(),
					NewBPM:    true,
					Prev:      prev,
				}
				tp.BPM, _ = timingPoint.BPM()
				if prev != nil {
					beatLength := prev.BPM * prev.BeatScale
					duration := float64(tp.Time - prev.Time)
					tp.Position = prev.Position + beatLength*duration
				} else {
					beatLength := tp.BPM * tp.BeatScale
					duration := float64(tp.Time - 0)
					tp.Position = 0 + beatLength*duration
				}
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
					Volume:    float64(timingPoint.Volume) / 100,
					Highlight: timingPoint.IsKiai(),
					NewBPM:    false,
					Prev:      prev,
				}
				tp.BeatScale, _ = timingPoint.BeatScale()
				if prev != nil {
					beatLength := prev.BPM * prev.BeatScale
					duration := float64(tp.Time - prev.Time)
					tp.Position = prev.Position + beatLength*duration
				} else {
					beatLength := tp.BPM * tp.BeatScale
					duration := float64(tp.Time - 0)
					tp.Position = 0 + beatLength*duration
				}
				if prev == nil {
					fmt.Printf("%s - %s: no uninherited point at first.\n", f.Title, f.Version)
					fmt.Printf("%+v\n", f.TimingPoints)
					os.Exit(1)
				}
				prev.Next = tp // Inherited point is never the first.
				prev = tp
				tps = append(tps, tp)
			}
		}
	}
	return tps
}

// BPM with longest duration will be main BPM.
// Suppose when there are multiple BPMs with same duration, larger one will be main.
func BPMs(tps []*TransPoint, duration int64) (main, min, max float64) {
	bpmDurations := make(map[float64]int64)
	for i, tp := range tps {
		if i == 0 {
			bpmDurations[tp.BPM] += tp.Time
		}
		if i < len(tps)-1 {
			bpmDurations[tp.BPM] += tps[i+1].Time - tp.Time
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

// wb, wa stands for buffer times: wait before, wait after.
// Multiply wa with 2 for preventing indexing a time slice over length.
func BarTimes(tps []*TransPoint, endTime int64) []int64 {
	var margin int64 = 5000
	ts := make([]int64, 0)
	tp0 := tps[0]
	for t := float64(tp0.Time); t >= float64(-margin); t -= float64(tp0.Meter) * 60000 / tp0.BPM {
		ts = append([]int64{int64(t)}, ts...)
	}
	ts = ts[:len(ts)-1] // Drop bar line for tp0 for avoiding duplicated
	for tp := tps[0]; tp != nil; tp = tp.NextBPMPoint {
		next := endTime + margin
		if tp.NextBPMPoint != nil {
			next = tp.NextBPMPoint.Time
		}
		unit := float64(tp.Meter) * 60000 / tp.BPM
		for t := float64(tp.Time); t < float64(next); t += unit {
			ts = append(ts, int64(t))
		}
	}
	return ts
}
