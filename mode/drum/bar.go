package drum

import "github.com/hndada/gosu"

type Bar struct {
	Position float64
	Speed    float64
}

func NewBars(tp0 *gosu.TransPoint, duration int64) (bs []Bar) {
	tp0.FetchLatest()
	var margin int64 = 5000
	if margin > tp0.Time {
		margin = tp0.Time
	}
	// Bar positions before first TransPoint.
	// Start with one step before for avoiding duplication.
	speed := tp0.Speed()
	step := tp0.BeatDuration()
	for t := float64(tp0.Time) - step; t >= float64(-margin); t -= step {
		b := Bar{
			Position: speed * t,
			Speed:    speed,
		}
		bs = append([]Bar{b}, bs...)
	}
	// Bar positions for first TransPoint and after it.
	for tp := tp0; tp != nil; tp = tp.NextBPMPoint.FetchLatest() {
		nextTime := duration + margin
		if tp.NextBPMPoint != nil {
			nextTime = tp.NextBPMPoint.Time
		}
		speed := tp.Speed()
		step = tp.BeatDuration()
		for t := float64(tp.Time); t < float64(nextTime); t += step {
			b := Bar{
				Position: speed * t,
				Speed:    speed,
			}
			bs = append(bs, b)
		}
	}
	return
}
