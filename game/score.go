package game

import (
	"fmt"
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/tween"
)

const (
	ScoreDot = iota + 10
	ScoreComma
	ScorePercent
)

func LoadScoreImages(fsys fs.FS) []draws.Image {
	imgs := make([]draws.Image, 13)
	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("score/%d.png", i)
		imgs[i] = draws.NewImageFromFile(fsys, fname)
	}
	for i, name := range []string{"dot", "comma", "percent"} {
		fname := fmt.Sprintf("score/%s.png", name)
		imgs[i+10] = draws.NewImageFromFile(fsys, fname)
	}
	return imgs
}

type ScoreOptions struct {
	ImageScale float64
	DigitGap   float64
}

type ScoreComponent struct {
	sprites []draws.Sprite
	score   float64
	w       float64 // Score's width is fixed.
	tween   tween.Tween
}

// Name of a function which returns closure ends with "-er".
func NewScoreComponent(imgs []draws.Image, opts *ScoreOptions) (cmp ScoreComponent) {
	cmp.sprites = make([]draws.Sprite, 13)
	// h0 is the height of number 0. Other numbers are located at h0 - h.
	// Score needs to set same base line, since
	// each number might have different height.
	var h0 float64
	s0 := draws.NewSprite(imgs[0])
	s0.Scale(opts.ImageScale)
	h0 = s0.H()
	cmp.w = s0.W() + opts.DigitGap
	for i, img := range imgs {
		sprite := draws.NewSprite(img)
		sprite.Scale(opts.ImageScale)
		sprite.Locate(ScreenSizeX, h0-sprite.H(), draws.RightTop)
		cmp.sprites[i] = sprite
	}
	return
}

func (cmp *ScoreComponent) Update(newScore float64) {
	if old := cmp.score; old != newScore {
		cmp.score = newScore

		tw := tween.Tween{MaxLoop: 1}
		begin := cmp.tween.Value()
		change := cmp.score - begin
		tw.Add(begin, change, 400, tween.EaseOutExponential)
		cmp.tween = tw
		cmp.tween.Start()
	}
}

func (cmp ScoreComponent) Draw(screen draws.Image) {
	score := int(cmp.tween.Value())
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
