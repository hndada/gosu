package drum

import (
	"image/color"

	"github.com/hndada/gosu"
)

// Todo: Tick judgment should be bound to MaxScaledBPM: 280
// Todo: let them put custom window
var (
	Cool = gosu.Judgment{Flow: 0.01, Acc: 1, Window: 25}
	Good = gosu.Judgment{Flow: 0.01, Acc: 0.25, Window: 60}
	Miss = gosu.Judgment{Flow: -1, Acc: 0, Window: 100}
)

// var Judgments = []gosu.Judgment{Cool, Good, Miss}
var Judgments = [2][3]gosu.Judgment{
	{Cool, Good, Miss},
	{Cool, Good, Miss},
}

// Todo: match the order betwen colors and judgments
var JudgmentColors = []color.NRGBA{
	{109, 120, 134, 255}, // Gray
	{51, 255, 40, 255},   // Lime
	{85, 251, 255, 255},  // Skyblue
}
var (
	white  = color.NRGBA{255, 255, 255, 192}
	purple = color.NRGBA{213, 0, 242, 192}
)

// When hit big notes only with one press, the note gives half the score only.
// For example, when hit a Big note by one press with Good, it will gives 0.25 * 0.5 = 0.125.
// No Flow decrease for hitting Big note by one press.
// When one side of judgment is Cool, Good on the other hand, overall judgment of Big note goes Good.
// In other word, to get Cool at Big note, you have to hit it Cool with both sides.

// Roll / Shake note does not affect on Flow / Acc scores.
// For example, a Roll / Shake only chart has extra score only: max score is 100k.
