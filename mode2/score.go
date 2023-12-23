package mode

import (
	"fmt"
	"io/fs"
	"math"

	"github.com/hndada/gosu/draws"
)

const (
	ScoreDot = iota + 10
	ScoreComma
	ScorePercent
)

type ScoreRes struct {
	imgs [13]draws.Image
}

func (res *ScoreRes) Load(fsys fs.FS) {
	for i := 0; i < 10; i++ {
		res.imgs[i] = draws.NewImageFromFile(fsys, fmt.Sprintf("score/%d.png", i))
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		res.imgs[i+10] = draws.NewImageFromFile(fsys, fmt.Sprintf("score/%s.png", name))
	}
}

type ScoreOpts struct {
	screen   draws.Box
	Scale    float64
	DigitGap float64
}

func NewScoreOpts(screen draws.Box) ScoreOpts {
	return ScoreOpts{
		screen:   screen,
		Scale:    0.65,
		DigitGap: 0,
	}
}

type ScoreComp struct {
	Score     float64
	lastScore float64 // to reset tween
	w         float64 // Score's width is fixed.
	sprites   [13]draws.Sprite
	easing    draws.TweenFunc
	tween     draws.Tween
}

// Name of a function which returns closure ends with "-er".
func NewScoreComp(res ScoreRes, opts ScoreOpts) (comp ScoreComp) {
	// h0 is the height of number 0. Other numbers are located at h0 - h.
	// Score needs to set same base line, since
	// each number might have different height.
	var h0 float64
	{
		s0 := draws.NewSprite(res.imgs[0])
		s0.MultiplyScale(opts.Scale)
		h0 = s0.H()
		comp.w = s0.W() + opts.DigitGap
	}

	for i, img := range res.imgs {
		sprite := draws.NewSprite(img)
		sprite.MultiplyScale(opts.Scale)
		sprite.Locate(opts.screen.X, h0-sprite.H(), draws.RightTop)
		comp.sprites[i] = sprite
	}

	comp.easing = draws.EaseOutExponential
	comp.setTween()
	return
}

func (comp *ScoreComp) setTween() {
	begin := comp.tween.Current()
	change := comp.Score - begin
	comp.tween = draws.NewTween(begin, change, ToTick(400), comp.easing)
}

func (comp *ScoreComp) Update() {
	comp.tween.Tick()
	if comp.lastScore != comp.Score {
		comp.lastScore = comp.Score
		comp.setTween()
	}
}

func (comp ScoreComp) Draw(screen draws.Image) {
	tweenScore := int(math.Ceil(comp.tween.Current()))
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
		sprite := comp.sprites[d]
		sprite.Move(tx, 0)
		// Need to set at center since anchor is RightTop.
		sprite.Move(-comp.w/2+sprite.W()/2, 0)
		sprite.Draw(screen, draws.Op{})
		tx -= comp.w
	}
}
