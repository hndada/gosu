package mode

import (
	"fmt"
	"io/fs"
	"math"

	"github.com/hndada/gosu/draws"
)

// const (
//	ScoreDot = iota + 10
//	ScoreComma
//	ScorePercent
// )

type Score struct {
	Score     float64
	lastScore float64 // to reset tween
	sprites   [13]draws.Sprite
	w         float64 // Score's width is fixed.
	tween     draws.Tween
	easing    draws.TweenFunc
}

type ScoreConfig struct {
	ScreenSize    *draws.Vector2
	FieldPosition *float64

	Position float64 // x
	Scale    float64
	DigitGap float64
}

func LoadScoreImages(fsys fs.FS) [13]draws.Image {
	var imgs [13]draws.Image
	for i := 0; i < 10; i++ {
		imgs[i] = draws.NewImageFromFile(fsys, fmt.Sprintf("score/%d.png", i))
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		imgs[i+10] = draws.NewImageFromFile(fsys, fmt.Sprintf("score/%s.png", name))
	}
	return imgs
}

// Name of a function which returns closure ends with "-er".
// 'score' is used as a name instead of 's' to avoid confusion with 's' in 'sprites'.
func NewScore(imgs [13]draws.Image, cfg ScoreConfig) (score Score) {
	// h0 is the height of number 0. Other numbers are located at h0 - h.
	// Score needs to set same base line,
	// since each number might have different height.
	var h0 float64
	{
		s0 := draws.NewSprite(imgs[0])
		s0.MultiplyScale(cfg.Scale)
		h0 = s0.Height()
		score.w = s0.Width() + cfg.DigitGap
	}

	for i, img := range imgs {
		sprite := draws.NewSprite(img)
		sprite.MultiplyScale(cfg.Scale)
		sprite.Locate(cfg.ScreenSize.X, h0-sprite.Height(), draws.RightTop)
		score.sprites[i] = sprite
	}

	score.easing = draws.EaseOutExponential
	score.setTween()
	return
}

func (score *Score) setTween() {
	begin := score.tween.Current()
	change := score.Score - begin
	score.tween = draws.NewTween(begin, change, ToTick(400), score.easing)
}

func (score *Score) Update() {
	score.tween.Tick()
	if score.lastScore != score.Score {
		score.lastScore = score.Score
		score.setTween()
	}
}

func (score Score) Draw(screen draws.Image) {
	tweenScore := int(math.Ceil(score.tween.Current()))
	digits := make([]int, 0)
	for v := tweenScore; v > 0; v /= 10 {
		digits = append(digits, v%10) // Little endian.
	}

	// Append zero if digits are not enough.
	const zeroFill = 1
	for i := len(digits); i < zeroFill; i++ {
		digits = append(digits, 0)
	}

	var tx float64
	for _, d := range digits {
		sprite := score.sprites[d]
		sprite.Move(tx, 0)
		// Need to set at center since anchor is RightTop.
		sprite.Move(-score.w/2+sprite.Width()/2, 0)
		sprite.Draw(screen, draws.Op{})
		tx -= score.w
	}
}
