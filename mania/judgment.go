package mania

import "github.com/hndada/gosu/common"

var (
	Kool  = common.Judgment{Value: 1.0000, Penalty: 0, HP: 0.75, Window: 16}
	Cool  = common.Judgment{Value: 0.9375, Penalty: 0, HP: 0.5, Window: 40} // 15/16
	Good  = common.Judgment{Value: 0.625, Penalty: 4, HP: 0.25, Window: 70} // 10/16
	Bad   = common.Judgment{Value: 0.25, Penalty: 10, HP: 0, Window: 100}
	Miss  = common.Judgment{Value: 0, Penalty: 25, HP: -3, Window: 150}
	empty = common.Judgment{}
)
var Judgments = [5]common.Judgment{Kool, Cool, Good, Bad, Miss}

// var JudgmentMeter common.JudgmentMeter

func init() {
	for i := range Judgments {
		Judgments[i].Window = int64(float64(Judgments[i].Window) * 1.5)
	}
}

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
