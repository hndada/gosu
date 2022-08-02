package main

import "github.com/hndada/gosu/parse/osr"

type ReplayState struct {
	Time    int64
	Pressed []bool
	Prev    *ReplayState
	Next    *ReplayState
}

func ExtractReplayState(f *osr.Format, keyCount int) []*ReplayState {
	var t int64
	var prev *ReplayState
	rss := make([]*ReplayState, len(f.ReplayData))
	for i, action := range f.ReplayData {
		t += action.W
		rs := &ReplayState{
			Time: t,
			Prev: prev,
		}

		var k int
		pressed := make([]bool, keyCount)
		for x := int(action.X); x > 0; x /= 2 {
			if x%2 == 1 {
				pressed[k] = true
			}
			k++
		}
		rs.Pressed = pressed
		prev.Next = rs
		prev = rs
		rss[i] = rs
	}
	return rss
}
