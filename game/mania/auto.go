package mania

import "sort"

func (c Chart) GenAutoKeyEvents() func(int64) []keyEvent {
	i := 0
	keyEvents := make([]keyEvent, 0, len(c.Notes)*2)
	for _, n := range c.Notes {
		switch n.Type {
		case TypeNote:
			e1 := keyEvent{
				time:    n.Time,
				key:     n.Key,
				pressed: true,
			}
			e2 := keyEvent{
				time:    n.Time + 1,
				key:     n.Key,
				pressed: false,
			}
			keyEvents = append(keyEvents, e1, e2)
		case TypeLNHead:
			e := keyEvent{
				time:    n.Time,
				key:     n.Key,
				pressed: true,
			}
			keyEvents = append(keyEvents, e)
		case TypeLNTail:
			e := keyEvent{
				time:    n.Time, // Time2: opposite time
				key:     n.Key,
				pressed: false,
			}
			keyEvents = append(keyEvents, e)
		}
	}
	sort.Slice(keyEvents, func(i, j int) bool { return keyEvents[i].time < keyEvents[j].time })
	return func(time int64) []keyEvent {
		var c int
		window := make([]keyEvent, 0, 10)
		for _, e := range keyEvents[i:] {
			if e.time <= time {
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
