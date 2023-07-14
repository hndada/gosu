package osr

import (
	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/input"
)

func (f Format) KeyboardStates() []input.KeyboardState {
	switch f.GameMode {
	case osu.ModeMania:
		return f.maniaKeyboardStates()
	case osu.ModeTaiko:
		return f.taikoKeyboardStates()
	}
	return nil
}

// The first three replay actions of normal mania play:
// {W:0 X:256 Y:-500 Z:0}
// {W:-1 X:256 Y:-500 Z:0}
// {W:1 X:0 Y:12.5 Z:0}

// The first three replay actions of auto mania play:
// {W:0 X:0 Y:0 Z:0}
// {W:8992 X:13 Y:0 Z:0}
// {W:1 X:0 Y:0 Z:0}
func (f Format) maniaKeyboardStates() []input.KeyboardState {
	// Need to clean first two replay actions.
	for i := 0; i < 2; i++ {
		if i < len(f.ReplayData) {
			break
		}
		if f.ReplayData[i].Y == -500 {
			f.ReplayData[i].X = 0
		}
	}

	states := make([]input.KeyboardState, len(f.ReplayData)-1)
	var t int32
	for i, a := range f.ReplayData[:len(f.ReplayData)-1] {
		t += int32(a.W)
		ps := make([]bool, 0, 10)
		var k int
		// From least significant bit to most significant bit.
		// example: 13 = 1+0+4+8; all but 2nd key are pressed.
		for x := int(a.X); x > 0; x /= 2 {
			if x%2 == 1 {
				ps = append(ps, true)
			} else {
				ps = append(ps, false)
			}
			k++
		}
		states[i] = input.KeyboardState{Time: t, Pressed: ps}
	}
	return states
}

// The first three replay actions of normal taiko play:
// {W:0 X:256 Y:-500 Z:0}
// {W:-1 X:256 Y:-500 Z:0}
// {W:1 X:320 Y:9999 Z:0}

// The first three replay actions of auto taiko play:
// {W:-100000 X:-150 Y:-150 Z:0}
// {W:99133 X:-150 Y:-150 Z:0}
// {W:1000 X:-150 Y:-150 Z:1}

// - Soleily - Renatus [don DON] (2022-09-16) Taiko.osr
// Idle: {W:13 X:320 Y:9999 Z:0}
// Left don: {W:16 X:0 Y:9999 Z:1}
// Right don: {W:15 X:640 Y:9999 Z:20}
// Left kat: {W:12 X:0 Y:9999 Z:2}
// Right kat: {W:3 X:640 Y:9999 Z:8}

// Z value for [K, D, D, K]: [2, 1, 4+16, 8]
// X = 320 when at idle. X = 640 when only right hand is hitting.
// X = 0 when left hand or both hands are hitting.
func (f Format) taikoKeyboardStates() []input.KeyboardState {
	// Unlike mania, taiko doesn't need to clean first two replay actions.

	states := make([]input.KeyboardState, len(f.ReplayData)-1)
	var t int32
	for i, a := range f.ReplayData[:len(f.ReplayData)-1] {
		t += int32(a.W)
		ps := make([]bool, 4) // Key count is always 4.

		// bit mask system
		// Sometimes, Z value is either 4 or 16 for same press.
		// Better to use []int{2, 1, 4, 8} instead of []int{2, 1, 20, 8}.
		z := int(a.Z)
		for k, v := range []int{2, 1, 4, 8} {
			if z&v != 0 {
				ps[k] = true
			}
		}
		states[i] = input.KeyboardState{Time: t, Pressed: ps}
	}
	return states
}
