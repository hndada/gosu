package drum

import (
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/game"
)

// ReplayListener supposes closure function is called every 1 ms.
// ReplayListener supposes the first the time of replay data is 0ms and no any inputs.
// Todo: Make sure to ReplayListener time is independent of Game's update tick
func NewReplayListener(f *osr.Format, timer *game.Timer) func() []bool {
	actions := append(f.ReplayData, osr.Action{W: 2e9})

	var i int
	var next int64 = actions[0].W + actions[1].W // +1
	return func() []bool {
		for timer.Now >= next { // There might be negative values on actions in a row.
			i++
			next += actions[i+1].W
		}
		pressed := make([]bool, 4)
		z := int(actions[i].Z)
		for k, v := range []int{2, 1, 4, 8} { // []int{2, 1, 20, 8}
			if z&v != 0 {
				pressed[k] = true
			}
		}
		return pressed
	}
}
