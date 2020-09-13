package mania

import (
	"github.com/hndada/gosu/mode"
)

var (
	kool = mode.Judgment{16 / 16, 0, 0.75, 16}
	cool = mode.Judgment{15 / 16, 0, 0.5, 40}
	good = mode.Judgment{10 / 16, 4, 0.25, 70}
	bad  = mode.Judgment{4 / 16, 10, 0, 100}
	miss = mode.Judgment{0, 25, -3, 150}
)
var judgments = [5]mode.Judgment{kool, cool, good, bad, miss}

// func judge(time int64) mode.Judgment {
// 	if time < 0 {
// 		time *= -1
// 	}
// 	for _, judge := range judgments {
// 		if time <= judge.Window {
// 			return judge
// 		}
// 	}
// 	return miss // todo: 미스 범위보다 멀면 그냥 무시
// }

const (
	idle = iota
	press
	release
	hold
)

func KeyAction(last, now bool) int { // action are realized with 2 snapshots
	switch {
	case !last && !now:
		return idle
	case !last && now:
		return press
	case last && !now:
		return release
	case last && now:
		return hold
	default:
		panic("not reach")
	}
}
