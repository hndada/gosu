package main

import (
	"fmt"
	"math"

	"github.com/hndada/gosu/parse/osr"
)

type ReplayState struct {
	Time    int64
	Pressed []bool
	// Prev    *ReplayState
	// Next    *ReplayState
}

func ExtractReplayState(f *osr.Format, keyCount int) []ReplayState {
	var t int64
	// var prev *ReplayState
	maxX := int(math.Pow(2, float64(keyCount))) - 1
	rss := make([]ReplayState, 0, len(f.ReplayData))
	for _, action := range f.ReplayData[:len(f.ReplayData)-1] { // Drop last data: RNG seed
		t += action.W
		rs := ReplayState{
			Time: t,
			// Prev: prev,
		}

		var k int
		pressed := make([]bool, keyCount)
		if int(action.X) > maxX {
			fmt.Printf("skip replay action %v\n", action)
			continue
		}
		for x := int(action.X); x > 0; x /= 2 {
			if x%2 == 1 {
				pressed[k] = true
			}
			k++
		}
		rs.Pressed = pressed
		// prev.Next = rs
		// prev = rs
		rss = append(rss, rs)
	}
	return rss
}
