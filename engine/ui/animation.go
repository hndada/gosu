package ui

import (
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// each frame has same interval
type Animation struct {
	srcs   []*ebiten.Image
	Sprite // no use src
	// Duration time.Time // temp: using global duration

	Rep int64 // 반복 횟수
	ebiten.CompositeMode
}

const RepInfinite = -1

func NewAnimation(srcs []*ebiten.Image) Animation {
	var a Animation
	a.srcs = make([]*ebiten.Image, len(srcs))
	// temp: only []*ebiten.Image can be passed
	for i, src := range srcs {
		a.srcs[i] = src
	}
	a.BornTime = time.Now()
	a.Saturation = 1
	a.Dimness = 1
	return a
}

const AnimationDuration = 450 // temp: global duration in ms

func (a Animation) Draw(screen *ebiten.Image) {
	timePerFrame := AnimationDuration / float64(len(a.srcs))
	elapsedTime := time.Since(a.BornTime).Milliseconds()
	rep := elapsedTime / AnimationDuration
	if a.Rep != RepInfinite && rep >= a.Rep {
		return
	}
	t := elapsedTime % AnimationDuration
	i := int(float64(t) / timePerFrame)
	// temp: suppose all animation goes not fixed
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(a.scaleW(), a.scaleH())
	op.GeoM.Translate(float64(a.X), float64(a.Y))
	op.ColorM.ChangeHSV(0, a.Saturation, a.Dimness)
	if a.CompositeMode != 0 {
		op.CompositeMode = a.CompositeMode
	}
	screen.DrawImage(a.srcs[i], op)
}

// temp: suppose all frames have same size in an animation
func (a Animation) scaleW() float64 {
	w1, _ := a.srcs[0].Size()
	return float64(a.W) / float64(w1)
}
func (a Animation) scaleH() float64 {
	_, h1 := a.srcs[0].Size()
	return float64(a.H) / float64(h1)
}
