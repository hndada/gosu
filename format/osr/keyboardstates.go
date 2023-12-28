package osr

import (
	"time"

	"github.com/hndada/gosu/format/osu"
	"github.com/hndada/gosu/input"
)

func (f Format) KeyboardStates(keyCount int) []input.KeyboardState {
	switch f.GameMode {
	case osu.ModeMania:
		return f.maniaKeyboardStates(keyCount)
	case osu.ModeTaiko:
		// Key count of taiko mode is always 4.
		return f.taikoKeyboardStates(4)
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

// Replay format itself has no information about key count.
func (f Format) maniaKeyboardStates(keyCount int) []input.KeyboardState {
	// Need to clean first two replay actions.
	for i := 0; i < 2; i++ {
		if i >= len(f.ReplayData) {
			break
		}
		if f.ReplayData[i].Y == -500 {
			f.ReplayData[i].X = 0
		}
	}

	// Remove a data for RNG seed, which looks like this:
	// -12345|0|0|seed
	f.ReplayData = f.ReplayData[:len(f.ReplayData)-1]

	var t time.Duration
	states := make([]input.KeyboardState, len(f.ReplayData))
	for i, a := range f.ReplayData {
		t += time.Duration(a.W) * time.Millisecond

		// From least significant bit to most significant bit.
		// Example: 13 = 1+0+4+8; all but 2nd key are pressed.
		var k int
		ps := make([]bool, keyCount)
		for x := int(a.X); x > 0; x /= 2 {
			if x%2 == 1 {
				ps[k] = true
			}
			k++
		}
		states[i] = input.KeyboardState{Time: t, PressedList: ps}
	}
	return states
}

// Reference 1: testdata/taiko.osr (Wizdomiot, koyomi's Oni)
// The first three replay actions of normal taiko play:
// {W:0 X:256 Y:-500 Z:0}
// {W:-1 X:256 Y:-500 Z:0}
// {W:1 X:320 Y:9999 Z:0}

// Reference 2: testdata/taiko.osr (Spelunker (ghm12) Oni)
// The first three replay actions of auto taiko play:
// {W:-100000 X:-150 Y:-150 Z:0}
// {W:99133 X:-150 Y:-150 Z:0}
// {W:1000 X:-150 Y:-150 Z:1}

// Reference 3: - Soleily - Renatus [don DON] (2022-09-16) Taiko.osr
// Idle: {W:13 X:320 Y:9999 Z:0}
// Left don: {W:16 X:0 Y:9999 Z:1}
// Right don: {W:15 X:640 Y:9999 Z:20}
// Left kat: {W:12 X:0 Y:9999 Z:2}
// Right kat: {W:3 X:640 Y:9999 Z:8}
// Z value for [K, D, D, K]: [2, 1, 4+16, 8]

// X = 0: left hand hits.
// X = 320: no hands hit.
// X = 640: right hand hits.
// X = 0: both hands hit. (beware: X = 0 again)

func (f Format) taikoKeyboardStates(keyCount int) []input.KeyboardState {
	// Unlike mania, taiko doesn't need to clean first two replay actions.

	// Remove a data for RNG seed, which looks like this:
	// -12345|0|0|seed
	f.ReplayData = f.ReplayData[:len(f.ReplayData)-1]

	var t time.Duration
	states := make([]input.KeyboardState, len(f.ReplayData))
	for i, a := range f.ReplayData {
		t += time.Duration(a.W) * time.Millisecond

		// Format uses bit mask system.
		// Sometimes, Z value is either 4 or 16 for same press.
		// It is better to use []int{2, 1, 4, 8},
		// instead of []int{2, 1, 20, 8}.
		z := int(a.Z)
		ps := make([]bool, keyCount)
		for k, v := range []int{2, 1, 4, 8} {
			if z&v != 0 {
				ps[k] = true
			}
		}
		states[i] = input.KeyboardState{Time: t, PressedList: ps}
	}
	return states
}
