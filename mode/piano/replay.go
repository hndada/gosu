package piano

import "github.com/hndada/gosu/format/osr"

// Todo: Make sure to ReplayListener time is independent of Game's update tick
// ReplayListener supposes closure function is called every 1 ms.
func NewReplayListener(f *osr.Format, keyCount int, waitBefore int64) func() []bool {
	// actions := f.TrimmedActions()
	actions := append(f.ReplayData, osr.Action{W: 2e9})
	for i := 0; i < 2; i++ {
		if i < len(actions) {
			break
		}
		if a := actions[i]; a.Y == -500 {
			a.X = 0
		}
	}

	var i int // Index of current replay action
	var t = waitBefore
	// var next = 0 + 1 + actions[0].W + actions[1].W
	var next int64
	return func() []bool {
		// if t >= next {
		// 	i++
		// 	next += actions[i+1].W
		// }
		for t >= next { // There might be negative values on actions in a row.
			i++
			next += actions[i+1].W
		}
		pressed := make([]bool, keyCount)
		var k int
		for x := int(actions[i].X); x > 0; x /= 2 {
			if x%2 == 1 {
				pressed[k] = true
			}
			k++
		}
		t++
		return pressed
	}
}
