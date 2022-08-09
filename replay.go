package main

import "github.com/hndada/gosu/parse/osr"

func NewReplayListener(f *osr.Format, keyCount int) func(int64) KeysState {
	var i int // Index of current replay action
	var state KeysState
	var actions = f.ReplayData[2 : len(f.ReplayData)-1] // Drop non-playing data
	return func(now int64) KeysState {
		var next int64 = 2e9 // Next replay action time in millisecond.
		if i < len(actions)-1 {
			next = state.Time + actions[i+1].W
		}
		if now >= next {
			pressed := make([]bool, keyCount)
			var k int
			for x := int(actions[i+1].X); x > 0; x /= 2 {
				if x%2 == 1 {
					pressed[k] = true
				}
				k++
			}
			state = KeysState{next, pressed}
		}
		return state
	}
}
