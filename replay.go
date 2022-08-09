package main

import "github.com/hndada/gosu/parse/osr"

// ReplayListener supposes closure function is called every 1 ms, and it should be.
func NewReplayListener(f *osr.Format, keyCount int) func() []bool {
	var actions = f.ReplayData[2 : len(f.ReplayData)-1] // Drop non-playing data
	// var t int64
	// ts := make([]int64, len(actions))
	// for i, a := range actions {
	// 	t += a.W
	// 	ts[i] = t
	// }
	var now = actions[0].W
	var i int // Index of current replay action
	return func() []bool {
		var next int64 = 2e9 // Next replay action time in millisecond.
		if i < len(actions)-1 {
			next = now + actions[i+1].W
		}
		if now >= next {
			i++
		}
		pressed := make([]bool, keyCount)
		var k int
		for x := int(actions[i+1].X); x > 0; x /= 2 {
			if x%2 == 1 {
				pressed[k] = true
			}
			k++
		}
		return pressed
	}
}
