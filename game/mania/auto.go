package mania

import (
	"math/rand"
	"sort"
)

func (c Chart) GenAutoKeyEvents(instability float64) func(int64) []keyEvent {
	i := 0 // closure
	keyEvents := make([]keyEvent, 0, len(c.Notes)*2)
	deviation := func(v float64) int64 {
		d := int64(rand.NormFloat64() * v * 2)
		if d > 20 { // temp: KOOL 더 띄우기
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
			e1 := keyEvent{
				Time:    n.Time + d,
				Key:     n.Key,
				Pressed: true,
			}
			e2 := keyEvent{
				Time:    n.Time + 30 + d,
				Key:     n.Key,
				Pressed: false,
			}
			keyEvents = append(keyEvents, e1, e2)
		case TypeLNHead:
			e := keyEvent{
				Time:    n.Time + d,
				Key:     n.Key,
				Pressed: true,
			}
			keyEvents = append(keyEvents, e)
		case TypeLNTail:
			e := keyEvent{
				Time:    n.Time + d, // Time2: opposite time
				Key:     n.Key,
				Pressed: false,
			}
			keyEvents = append(keyEvents, e)
		}
	}
	sort.Slice(keyEvents, func(i, j int) bool { return keyEvents[i].Time < keyEvents[j].Time })
	return func(time int64) []keyEvent {
		var c int
		window := make([]keyEvent, 0, 10)
		for _, e := range keyEvents[i:] {
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
