package drum

import "github.com/hndada/gosu/mode"

type Bar struct {
	Floater
	Next *Bar
	Prev *Bar
}

func NewBars(transPoints []*mode.Dynamic, duration int64) (bs []*Bar) {
	var start, end, step float64 // Next time.

	// Bar positions before first Dynamic.
	start = float64(transPoints[0].Time)
	end = start
	if end > -5000 {
		end = -5000
	}
	step = transPoints[0].BeatDuration()
	for t := start; t >= end; t -= step {
		b := Bar{
			Floater: Floater{
				Time: int64(t),
			},
		}
		bs = append([]*Bar{&b}, bs...)
	}

	// Bar positions after first Dynamic.
	bs = bs[:len(bs)-1] // Drop for avoiding duplicattion
	newBeatPoints := make([]*mode.Dynamic, 0)
	for _, dy := range transPoints {
		if dy.NewBeat {
			newBeatPoints = append(newBeatPoints, tp)
		}
	}
	for i, dy := range newBeatPoints {
		start = float64(dy.Time)
		if i == len(newBeatPoints)-1 {
			end = float64(duration + 5000)
		} else {
			end = float64(newBeatPoints[i+1].Time)
		}
		step = dy.BeatDuration()
		for t := start; t < end; t += step {
			b := Bar{
				Floater: Floater{
					Time: int64(t),
				},
			}
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
