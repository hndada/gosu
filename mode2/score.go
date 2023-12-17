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

type ScoreComponent struct {
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
func NewScoreComponent(imgs [13]draws.Image, cfg ScoreConfig) (sc ScoreComponent) {
	// h0 is the height of number 0. Other numbers are located at h0 - h.
	// Score needs to set same base line,
	// since each number might have different height.
	var h0 float64
	{
		s0 := draws.NewSprite(imgs[0])
		s0.MultiplyScale(cfg.Scale)
		h0 = s0.Height()
		sc.w = s0.Width() + cfg.DigitGap
	}

	for i, img := range imgs {
		sprite := draws.NewSprite(img)
		sprite.MultiplyScale(cfg.Scale)
		sprite.Locate(cfg.ScreenSize.X, h0-sprite.Height(), draws.RightTop)
		sc.sprites[i] = sprite
	}

	sc.easing = draws.EaseOutExponential
	sc.setTween()
	return
}

func (sc *ScoreComponent) setTween() {
	begin := sc.tween.Current()
	change := sc.Score - begin
	sc.tween = draws.NewTween(begin, change, ToTick(400), sc.easing)
}

func (sc *ScoreComponent) Update() {
	sc.tween.Tick()
	if sc.lastScore != sc.Score {
		sc.lastScore = sc.Score
		sc.setTween()
	}
}

func (sc ScoreComponent) Draw(screen draws.Image) {
	tweenScore := int(math.Ceil(sc.tween.Current()))
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
		sprite := sc.sprites[d]
		sprite.Move(tx, 0)
		// Need to set at center since anchor is RightTop.
		sprite.Move(-sc.w/2+sprite.Width()/2, 0)
		sprite.Draw(screen, draws.Op{})
		tx -= sc.w
	}
}
