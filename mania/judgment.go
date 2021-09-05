package mania

import (
	"github.com/hndada/gosu/game"
)

var (
	Kool  = game.Judgment{Value: 1.0000, Penalty: 0, HP: 0.75, Window: 16}
	Cool  = game.Judgment{Value: 0.9375, Penalty: 0, HP: 0.5, Window: 40} // 15/16
	Good  = game.Judgment{Value: 0.625, Penalty: 4, HP: 0.25, Window: 70} // 10/16
	Bad   = game.Judgment{Value: 0.25, Penalty: 10, HP: 0, Window: 100}
	Miss  = game.Judgment{Value: 0, Penalty: 25, HP: -3, Window: 150}
	empty = game.Judgment{}
)
var Judgments = [5]game.Judgment{Kool, Cool, Good, Bad, Miss}

// var JudgmentMeter game.JudgmentMeter

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
