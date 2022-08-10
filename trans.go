package gosu

import (
	"sort"

	"github.com/hndada/gosu/parse/osu"
)

type TransPoint struct {
	Time        int64
	SpeedFactor float64
	BPM         float64
	Meter       uint8
	Volume      uint8
	Highlight   bool
	Prev        *TransPoint
	Next        *TransPoint
}

// Uninherited point is base point. whereas, Inherited point 'inherits'
// values from previous Uninherited point first, then after.
// All osu format have at least one Uninherited timing point.
// Uninherited: BPM, Inherited: speed factor
func NewTransPointsFromOsu(o *osu.Format) []*TransPoint {
	sort.Slice(o.TimingPoints, func(i int, j int) bool {
		return o.TimingPoints[i].Time < o.TimingPoints[j].Time
	})
	tps := make([]*TransPoint, 0, len(o.TimingPoints))
	var lastUninherited TransPoint // Need to deep-copied
	var prev *TransPoint
	for _, timingPoint := range o.TimingPoints {
		if timingPoint.Uninherited {
			tp := &TransPoint{
				Time:        int64(timingPoint.Time),
				SpeedFactor: 1,
				// BPM:,
				Meter:     uint8(timingPoint.Meter),
				Volume:    uint8(timingPoint.Volume),
				Highlight: timingPoint.IsKiai(),
				Prev:      prev,
			}
			tp.BPM, _ = timingPoint.BPM()

			if prev != nil {
				prev.Next = tp
			}
			prev = tp
			lastUninherited = *tp
			tps = append(tps, tp)
		} else {
			tp := &TransPoint{
				Time: int64(timingPoint.Time),
				// SpeedFactor: 1,
				BPM:       lastUninherited.BPM,
				Meter:     uint8(timingPoint.Meter),
				Volume:    uint8(timingPoint.Volume),
				Highlight: timingPoint.IsKiai(),
				Prev:      prev,
			}
			tp.SpeedFactor, _ = timingPoint.SpeedFactor()

			prev.Next = tp // Inherited point is never the first.
			prev = tp
			tps = append(tps, tp)
		}
	}
	return tps
}
