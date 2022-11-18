package drum

import (
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/format/osr"
)

// - Soleily - Renatus [don DON] (2022-09-16) Taiko.osr
// Idle: {W:13 X:320 Y:9999 Z:0}
// Left don: {W:16 X:0 Y:9999 Z:1}
// Right don: {W:15 X:640 Y:9999 Z:20}
// Left kat: {W:12 X:0 Y:9999 Z:2}
// Right kat: {W:3 X:640 Y:9999 Z:8}

// Z value for [K, D, D, K]: [2, 1, 4+16, 8]
// X = 320 when at idle. X = 640 when only right hand is hitting.
// X = 0 when left hand or both hands are hitting.

// ReplayListener supposes closure function is called every 1 ms.
// ReplayListener supposes the first the time of replay data is 0ms and no any inputs.
// Todo: Make sure to ReplayListener time is independent of Game's update tick
func NewReplayListener(f *osr.Format, timer *audios.Timer) func() []bool {
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
