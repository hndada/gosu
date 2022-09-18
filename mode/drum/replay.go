package drum

import (
	"github.com/hndada/gosu/format/osr"
)

// Todo: Make sure to ReplayListener time is independent of Game's update tick
// ReplayListener supposes closure function is called every 1 ms.
// Suppose the first the time of replay data is 0ms and no any inputs.
func NewReplayListener(f *osr.Format, waitBefore int64) func() []bool {
	// for i, a := range f.ReplayData[:15] {
	// 	fmt.Printf("%d: %v\n", i, a)
	// }
	// actions := f.TrimmedActions()
	actions := append(f.ReplayData, osr.Action{W: 2e9})
	// for i, a := range actions[:500] {
	// 	if a.Z != 0 {
	// 		fmt.Printf("%d: %v\n", i, a)
	// 	}
	// }

	var i int // Index of current replay action
	var t = waitBefore
	// var next = 0 + actions[0].W + actions[1].W // + 1
	var next int64 = actions[0].W + actions[1].W // 0 + actions[1].W
	return func() []bool {
		for t >= next { // There might be negative values on actions in a row.
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
		t++
		return pressed
	}
}
