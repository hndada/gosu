package piano

import "github.com/hndada/gosu"

type Bar struct {
	Position float64
	Next     *Bar
	Prev     *Bar
}

func NewBars(transPoints []*gosu.TransPoint, duration int64) (bs []*Bar) {
	tp := transPoints[0]
	tp.FetchPresent()
	var margin int64 = 5000
	if margin > tp.Time {
		margin = tp.Time
	}
	// Bar positions before first TransPoint.
	// Start with one step before for avoiding duplication.
	step := tp.BeatDuration()
	for t := float64(tp.Time) - step; t >= float64(-margin); t -= step {
		b := Bar{Position: tp.Speed() * t}
		bs = append([]*Bar{&b}, bs...)
	}
	// Bar positions for first TransPoint and after it.
	for ; tp != nil; tp = tp.NextBPMPoint {
		tp = tp.FetchPresent()
		nextTime := duration + margin
		if tp.NextBPMPoint != nil {
			nextTime = tp.NextBPMPoint.Time
		}
		step = tp.BeatDuration()
		for t := float64(tp.Time); t < float64(nextTime); t += step {
			pos := tp.Position
			pos += tp.Speed() * (t - float64(tp.Time))
			b := Bar{Position: pos}
			bs = append(bs, &b)
		}
	}
	var prev *Bar
	for _, b := range bs {
		b.Prev = prev
		if prev != nil {
			prev.Next = b
		}
		prev = b
	}
	return
}
