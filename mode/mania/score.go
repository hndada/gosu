package mania

import "github.com/hndada/gosu/mode"

var Judgements = [5]mode.Judgement{
	{"KOOL", 16 / 16, 0, 0.75, 16},
	{"COOL", 15 / 16, 0, 0.5, 40},
	{"GOOD", 10 / 16, 4, 0.25, 70},
	{"BAD", 4 / 16, 10, 0, 100},
	{"MISS", 0, 25, -3, 150},
}

type Score struct {
	mode.BaseScore
	Counts [len(Judgements)]int
}

func (s Score) JudgeCounts() []int { return s.Counts[:] }
func (s Score) IsFullCombo() bool  { return s.Counts[4] == 0 }
func (s Score) IsPerfect() bool {
	for _, c := range s.Counts[2:] {
		if c != 0 {
			return false
		}
	}
	return true
}
