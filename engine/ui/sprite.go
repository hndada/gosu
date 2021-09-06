package ui

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	src  *ebiten.Image
	W, H int // desired w, h
	X, Y int

	Color      color.Color
	Saturation float64
	Dimness    float64

	BornTime time.Time
	ebiten.CompositeMode
	// LifeTime time.Time // zero value goes eternal
}

func NewSprite(src image.Image) Sprite {
	var sprite Sprite
	sprite.SetImage(src)

	sprite.Saturation = 1
	sprite.Dimness = 1

	sprite.BornTime = time.Now()
	return sprite
}

func (s Sprite) isOut(screenSize image.Point) bool {
	return (s.X+s.W < 0 || s.X > screenSize.X ||
		s.Y+s.H < 0 || s.Y > screenSize.Y)
}

func (s *Sprite) SetImage(i image.Image) {
	switch i.(type) {
	case *ebiten.Image:
		s.src = i.(*ebiten.Image)
	default:
		i2 := ebiten.NewImageFromImage(i)
		s.src = i2
	}
}

func (s Sprite) scaleW() float64 {
	w1, _ := s.src.Size()
	return float64(s.W) / float64(w1)
}
func (s Sprite) scaleH() float64 {
	_, h1 := s.src.Size()
	return float64(s.H) / float64(h1)
}
func (s Sprite) Draw(screen *ebiten.Image) {
	if s.src == nil {
		panic("s.src is nil")
	}
	if s.isOut(screen.Bounds().Max) {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(s.scaleW(), s.scaleH())
	op.GeoM.Translate(float64(s.X), float64(s.Y))
	op.ColorM.ChangeHSV(0, s.Saturation, s.Dimness)
	if s.CompositeMode != 0 {
		op.CompositeMode = s.CompositeMode
	}
	screen.DrawImage(s.src, op)
}

// for debugging
func (s Sprite) PrintWHXY(comment string) {
	fmt.Println(comment, s.W, s.H, s.X, s.Y)
}

type FixedSprite struct { // a sprite that never moves once appears
	Sprite
	op *ebiten.DrawImageOptions
}

func NewFixedSprite(src image.Image) FixedSprite {
	return FixedSprite{
		Sprite: NewSprite(src),
	}
}
func (s FixedSprite) Draw(screen *ebiten.Image) {
	if s.src == nil {
		panic("s.src is nil")
	}
	if s.isOut(screen.Bounds().Max) {
		return
	}
	screen.DrawImage(s.src, s.op)
}

// minor parameter should already been set
func (s *FixedSprite) Fix() {
	op := &ebiten.DrawImageOptions{}
	if s.CompositeMode != 0 {
		op.CompositeMode = s.CompositeMode
	}
	op.GeoM.Scale(s.scaleW(), s.scaleH())
	op.GeoM.Translate(float64(s.X), float64(s.Y))

	if s.Color != nil {
		r, g, b, _ := s.Color.RGBA()
		op.ColorM.Scale(0, 0, 0, 1) // reset
		op.ColorM.Translate(
			float64(r)/0xff,
			float64(g)/0xff,
			float64(b)/0xff,
			0, // temp
		)
	}
	op.ColorM.ChangeHSV(0, s.Saturation, s.Dimness)
	s.op = op
}
