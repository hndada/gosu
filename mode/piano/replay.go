package piano

import (
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

// ReplayListener supposes closure function is called every 1 ms.
// ReplayListener supposes the first the time of replay data is 0ms and no any inputs.
// Todo: Make sure to ReplayListener time is independent of Game's update tick

// Fetch() ([]KeyPressedLog, []KeyActionLog)
func NewReplayListener(f *osr.Format, keyCount int, timer *mode.Timer) func() []bool {
	actions := append(f.ReplayData, osr.Action{W: 2e9})
	for i := 0; i < 2; i++ {
		if i < len(actions) {
			break
		}
		if a := actions[i]; a.Y == -500 {
			a.X = 0
		}
	}

	var i int                                    // Index of current replay action
	var next int64 = actions[0].W + actions[1].W // +1
	return func() []bool {
		for timer.Now >= next { // There might be negative values on actions in a row.
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
		return pressed
	}
}
func loadReplayFromOsu(f *osr.Format, keyCount int) *input.KeyListener {
	actions := append(f.ReplayData, osr.Action{W: 2e9})

	for i := 0; i < 2; i++ {
		if i < len(actions) {
			break
		}
		if a := actions[i]; a.Y == -500 {
			a.X = 0
		}
	}

	var (
		i   int // Index of current replay action
		now int64
	)
	now = actions[0].W + actions[1].W // +1
	return &input.KeyListener{
		PollingRate: 1000,
		Listen:      func() input.KeyPressedLog {},
	}
}
