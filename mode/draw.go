package mode

import (
	"math"

	"github.com/hndada/gosu/draws"
)

// Name of a function which returns closure ends with "-er".
func NewScoreDrawer(sprites [13]draws.Sprite, score *float64, digitGap float64) func(draws.Image) {
	const zeroFill = 1

	numbers := sprites[:10]
	digitWidth := sprites[0].Width() // Use number 0's width.
	delayedScore := NewDelayed(score)

	return func(dst draws.Image) {
		delayedScore.Update()
		score := int(math.Floor(delayedScore.Delayed))
		digits := make([]int, 0)
		for v := score; v > 0; v /= 10 {
			digits = append(digits, v%10) // Little endian.
		}
		for i := len(digits); i < zeroFill; i++ {
			digits = append(digits, 0)
		}

		w := digitWidth + digitGap
		var tx float64
		for _, d := range digits {
			sprite := numbers[d]
			sprite.Move(tx, 0)
			// Need to set at center since anchor is RightTop.
			sprite.Move(-w/2+sprite.Width()/2, 0)
			sprite.Draw(dst, draws.Op{})
			tx -= w
		}
	}
}

// Each number has different width. Number 0's width is used as standard.
// ComboDrawer's Draw draws each number at constant x regardless of their widths.
func NewComboDrawer(sprites [10]draws.Sprite, src *int, timer *draws.Timer, digitGap float64, bounce float64) func(draws.Image) {
	digitWidth := sprites[0].Width() // Use number 0's width.
	combo := *src
	return func(dst draws.Image) {
		timer.Ticker()
		if combo != *src {
			combo = *src
			timer.Reset()
		}
		if timer.IsDone() {
			return
		}
		if combo == 0 {
			return
		}
		vs := make([]int, 0)
		for v := combo; v > 0; v /= 10 {
			vs = append(vs, v%10) // Little endian.
		}

		// Size of the whole image is 0.5w + (n-1)(w+gap) + 0.5w.
		// Since sprites are already at anchor, no need to care of two 0.5w.
		w := digitWidth + digitGap
		tx := float64(len(vs)-1) * w / 2
		const (
			boundary1 = 0.05
			boundary2 = 0.1
		)
		for _, v := range vs {
			sprite := sprites[v]
			sprite.Move(tx, 0)
			age := timer.Age()
			if age < boundary1 {
				scale := 0.1 * timer.Progress(0, boundary1)
				sprite.Move(0, bounce*sprite.Height()*scale)
			}
			if age >= boundary1 && age < boundary2 {
				scale := 0.1 - 0.1*timer.Progress(boundary1, boundary2)
				sprite.Move(0, bounce*sprite.Height()*scale)
			}
			sprite.Draw(dst, draws.Op{})
			tx -= w
		}
	}
}
