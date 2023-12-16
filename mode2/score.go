package mode

import (
	"fmt"
	"io/fs"
	"math"

	"github.com/hndada/gosu/draws"
)

type Score struct {
	Score float64
}

type ScoreConfig struct {
	Position float64 // x
	Scale    float64
}

func NewScoreSprites(fsys fs.FS, ScreenSize draws.Vector2, scale float64) [13]draws.Sprite {
	var sprites [13]draws.Sprite
	for i := 0; i < 10; i++ {
		s := draws.NewSpriteFromFile(fsys, fmt.Sprintf("score/%d.png", i))
		s.MultiplyScale(scale)
		// Score needs to set same base line,
		// since each number might have different height.
		if i == 0 {
			s.Locate(ScreenSize.X, 0, draws.RightTop)
		} else {
			s.Locate(ScreenSize.X, sprites[0].Height()-s.Height(), draws.RightTop)
		}
		sprites[i] = s
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		s := draws.NewSpriteFromFile(fsys, fmt.Sprintf("score/%s.png", name))
		s.MultiplyScale(scale)
		s.Locate(ScreenSize.X, 0, draws.RightTop)
		sprites[i+10] = s
	}
	return sprites
}

const (
	ScoreDot = iota + 10
	ScoreComma
	ScorePercent
)

func (s Score) Draw(screen draws.Image, sprites [13]draws.Sprite, digitGap float64) {

}

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
