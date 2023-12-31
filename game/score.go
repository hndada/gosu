package game

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
)

const (
	ScoreDot = iota + 10
	ScoreComma
	ScorePercent
)

type ScoreResources struct {
	imgs [13]draws.Image // numbers with sign (. , %)
}

func (res *ScoreResources) Load(fsys fs.FS) {
	for i := 0; i < 10; i++ {
		res.imgs[i] = draws.NewImageFromFile(fsys, fmt.Sprintf("score/%d.png", i))
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		res.imgs[i+10] = draws.NewImageFromFile(fsys, fmt.Sprintf("score/%s.png", name))
	}
}

type ScoreOptions struct {
	Scale    float64
	DigitGap float64
}

func NewScoreOptions() ScoreOptions {
	return ScoreOptions{
		Scale:    0.65,
		DigitGap: 0,
	}
}

type ScoreComponent struct {
	sprites [13]draws.Sprite
	score   float64
	w       float64 // Score's width is fixed.
	tween   draws.Tween
}

// Name of a function which returns closure ends with "-er".
func NewScoreComponent(res ScoreResources, opts ScoreOptions) (cmp ScoreComponent) {
	// h0 is the height of number 0. Other numbers are located at h0 - h.
	// Score needs to set same base line, since
	// each number might have different height.
	var h0 float64
	s0 := draws.NewSprite(res.imgs[0])
	s0.MultiplyScale(opts.Scale)
	h0 = s0.H()
	cmp.w = s0.W() + opts.DigitGap

	for i, img := range res.imgs {
		sprite := draws.NewSprite(img)
		sprite.MultiplyScale(opts.Scale)
		sprite.Locate(ScreenW, h0-sprite.H(), draws.RightTop)
		cmp.sprites[i] = sprite
	}
	// cmp.tween = cmp.newTween()
	return
}

func (cmp *ScoreComponent) newTween() draws.Tween {
	begin := cmp.tween.Current()
	change := cmp.score - begin
	return draws.NewTween(begin, change, 400, draws.EaseOutExponential)
}

func (cmp *ScoreComponent) Update(newScore float64) {
	if old := cmp.score; old != newScore {
		cmp.score = newScore
		cmp.tween = cmp.newTween()
	}
}

func (cmp ScoreComponent) Draw(screen draws.Image) {
	score := int(cmp.tween.Current())
	digits := make([]int, 0)
	for v := score; v > 0; v /= 10 {
		digits = append(digits, v%10) // Little endian.
	}

	// Append zero if digits are not enough.
	const zeroFill = 1
	for i := len(digits); i < zeroFill; i++ {
		digits = append(digits, 0)
	}

	var tx float64
	for _, d := range digits {
		s := cmp.sprites[d]
		s.Move(tx, 0)
		// Need to set at center since anchor is RightTop.
		s.Move(-cmp.w/2+s.W()/2, 0)
		s.Draw(screen)
		tx -= cmp.w
	}
}
