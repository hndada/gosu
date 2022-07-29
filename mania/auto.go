package mania

import (
	"math/rand"
	"sort"

	"github.com/hndada/gosu/common"
)

func (c Chart) GenAutoKeyEvents(instability float64) func(int64) []common.PlayKeyEvent {
	i := 0 // closure
	es := make([]common.PlayKeyEvent, 0, len(c.Notes)*2)
	deviation := func(v float64) int64 {
		d := int64(rand.NormFloat64() * v * 2)
		if d > 20 { // temp: more biased to KOOL
			d = int64(rand.NormFloat64() * v * 2)
		}
		return d
	}
	var d int64
	for _, n := range c.Notes {
		d = deviation(instability)
		if d > Miss.Window { // lost
			continue
		}
		switch n.Type {
		case TypeNote:
			e1 := common.PlayKeyEvent{
				Time:    n.Time + d,
				Pressed: true,
				Key:     n.Key,
			}
			e2 := common.PlayKeyEvent{
				Time:    n.Time + 30 + d,
				Pressed: false,
				Key:     n.Key,
			}
			es = append(es, e1, e2)
		case TypeLNHead:
			e := common.PlayKeyEvent{
				Time:    n.Time + d,
				Pressed: true,
				Key:     n.Key,
			}
			es = append(es, e)
		case TypeLNTail:
			e := common.PlayKeyEvent{
				Time:    n.Time + d, // Time2: opposite time
				Pressed: false,
				Key:     n.Key,
			}
			es = append(es, e)
		}
	}
	sort.Slice(es, func(i, j int) bool { return es[i].Time < es[j].Time })
	return func(time int64) []common.PlayKeyEvent {
		var c int
		window := make([]common.PlayKeyEvent, 0, 10)
		for _, e := range es[i:] {
			if e.Time <= time {
				window = append(window, e)
				c++
			} else {
				break
			}
		}
		i += c
		return window
	}
}
