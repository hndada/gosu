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
	var lastUninherited TransPoint
	var prev *TransPoint
	for _, timingPoint := range o.TimingPoints {
		var tp TransPoint
		if timingPoint.Uninherited {
			tp = TransPoint{
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
				prev.Next = &tp
			}
			lastUninherited = tp
			prev = &tp
		} else {
			tp = lastUninherited
			tp.Time = int64(timingPoint.Time)
			tp.SpeedFactor, _ = timingPoint.SpeedFactor()
			tp.Meter = uint8(timingPoint.Meter)
			tp.Volume = uint8(timingPoint.Volume)
			tp.Meter = uint8(timingPoint.Meter)
			tp.Highlight = timingPoint.IsKiai()
			if prev != nil {
				prev.Next = &tp
			}
			prev = &tp
		}
		tps = append(tps, &tp)
	}
	// for _, tp := range tps {
	// 	fmt.Printf("%+v\n", *tp)
	// }
	return tps
}
