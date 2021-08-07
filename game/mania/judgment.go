package mania

import (
	"github.com/hndada/gosu/game"
)

type Judgment struct {
	Value   float64
	Penalty float64
	HP      float64
	Window  int64
}

var (
	kool  = game.Judgment{Value: 1.0000, Penalty: 0, HP: 0.75, Window: 16}
	cool  = game.Judgment{Value: 0.9375, Penalty: 0, HP: 0.5, Window: 40} // 15/16
	good  = game.Judgment{Value: 0.625, Penalty: 4, HP: 0.25, Window: 70} // 10/16
	bad   = game.Judgment{Value: 0.25, Penalty: 10, HP: 0, Window: 100}
	miss  = game.Judgment{Value: 0, Penalty: 25, HP: -3, Window: 150}
	empty = game.Judgment{}
)
var judgments = [5]game.Judgment{kool, cool, good, bad, miss}

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
