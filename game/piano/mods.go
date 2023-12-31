package piano

import "github.com/hndada/gosu/game"

// Alternative names of Mods:
// Modifiers, Parameters
// Occupied: Options, Settings, Configs
// If Mods is gonna be used, it might be good to change "Mode".

// the ideal number of Judgments is: 3 + 1
func (Scorer) DefaultJudgments() []game.Judgment {
	return []game.Judgment{
		{Window: 20, Weight: 1},
		{Window: 40, Weight: 1},
		{Window: 80, Weight: 0.5},
		{Window: 120, Weight: 0},
	}
}
