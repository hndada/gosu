package piano

import "github.com/hndada/gosu"

type Bar struct {
	Position float64
	Next     *Bar
	Prev     *Bar
}

func NewBars(transPoints []*gosu.TransPoint, duration int64) (bs []*Bar) {
	var start, end, step float64 // Next time.
	// Bar positions before first TransPoint.
	start = float64(transPoints[0].Time)
	end = start
	if end > -5000 {
		end = -5000
	}
	step = transPoints[0].BeatDuration()
	for t := start; t >= end; t -= step {
		b := Bar{Position: t * transPoints[0].Speed}
		bs = append([]*Bar{&b}, bs...)
	}
	bs = bs[:len(bs)-1] // Drop for avoiding duplicattion

	// Bar positions after first TransPoint.
	newBeatPoints := make([]*gosu.TransPoint, 0)
	for _, tp := range transPoints {
		if tp.NewBeat {
			newBeatPoints = append(newBeatPoints, tp)
		}
	}
	for i, tp := range newBeatPoints {
		start = float64(tp.Time)
		if i == len(newBeatPoints)-1 {
			end = float64(duration + 5000)
		} else {
			end = float64(newBeatPoints[i+1].Time)
		}
		step = tp.BeatDuration()
		for t := start; t < end; t += step {
			b := Bar{Position: tp.Position + (t-start)*tp.Speed}
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