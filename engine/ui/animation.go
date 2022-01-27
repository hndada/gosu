package ui

import (
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const RepInfinite = -1
const AnimationDuration = 450 // temp: global duration in ms

// each frame has same interval
type Animation struct {
	imgs   []*ebiten.Image
	Sprite // no use src
	// Duration time.Time // temp: using global duration

	Rep int64 // repetition
	// ebiten.CompositeMode
}

// TEMP: all imgs should be *ebiten.Image
// TODO: Should I make animation.SetImage?
// TODO: FixedAnimation
func NewAnimation(imgs []*ebiten.Image) Animation {
	var a Animation
	a.imgs = make([]*ebiten.Image, len(imgs))
	for i, img := range imgs {
		a.imgs[i] = img
	}
	a.Saturation = 1
	a.Dimness = 1
	a.BornTime = time.Now()
	return a
}

func (a Animation) Draw(screen *ebiten.Image) {
	timePerFrame := AnimationDuration / float64(len(a.imgs))
	elapsedTime := time.Since(a.BornTime).Milliseconds()
	rep := elapsedTime / AnimationDuration
	if a.Rep != RepInfinite && rep >= a.Rep {
		return
	}
	t := elapsedTime % AnimationDuration
	i := int(float64(t) / timePerFrame)
	// TODO: Sprite's Op() uses its member variable, i, which is nil in animation.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(a.scaleW(), a.scaleH())
	op.GeoM.Translate(float64(a.X), float64(a.Y))
	op.ColorM.ChangeHSV(0, a.Saturation, a.Dimness)
	if a.CompositeMode != 0 {
		op.CompositeMode = a.CompositeMode
	}
	screen.DrawImage(a.imgs[i], op)
}

// TEMP: suppose all frames have same size in an animation
func (a Animation) scaleW() float64 {
	w1, _ := a.imgs[0].Size()
	return float64(a.W) / float64(w1)
}
func (a Animation) scaleH() float64 {
	_, h1 := a.imgs[0].Size()
	return float64(a.H) / float64(h1)
}
