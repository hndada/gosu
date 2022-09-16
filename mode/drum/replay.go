package drum

import "github.com/hndada/gosu/format/osr"

// Todo: Make sure to ReplayListener time is independent of Game's update tick
// ReplayListener supposes closure function is called every 1 ms.
func NewReplayListener(f *osr.Format, waitBefore int64) func() []bool {
	actions := f.TrimmedActions()
	actions = append(actions, osr.Action{W: 2e9})
	var i int // Index of current replay action
	var t = waitBefore
	var next = 0 + actions[0].W + actions[1].W // + 1
	return func() []bool {
		if t >= next {
			i++
			next += actions[i+1].W
		}
		pressed := make([]bool, 4)
		z := int(actions[i].Z)
		for k, v := range []int{2, 1, 20, 8} {
			if z&v != 0 {
				pressed[k] = true
			}
		}
		t++
		return pressed
	}
}
