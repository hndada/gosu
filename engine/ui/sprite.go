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

// Y-axis towards bottom

// concept: Depth ~= layer; drawing order
type Sprite struct {
	i          *ebiten.Image
	W, H, X, Y int
	Theta      float64

	Color      color.Color
	Saturation float64
	Dimness    float64
	ebiten.CompositeMode

	BornTime time.Time
	// LifeTime time.Time // zero value goes eternal
}

func NewSprite(i image.Image) Sprite {
	var s Sprite
	s.SetImage(i)
	s.Saturation = 1
	s.Dimness = 1
	s.BornTime = time.Now()
	return s
}

func (s Sprite) isOut(screenSize image.Point) bool {
	return (s.X+s.W < 0 || s.X > screenSize.X ||
		s.Y+s.H < 0 || s.Y > screenSize.Y)
}

func (s Sprite) scaleW() float64 {
	w1, _ := s.i.Size()
	return float64(s.W) / float64(w1)
}
func (s Sprite) scaleH() float64 {
	_, h1 := s.i.Size()
	return float64(s.H) / float64(h1)
}

func (s Sprite) Op() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Rotate(s.Theta)
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

	if s.CompositeMode != 0 {
		op.CompositeMode = s.CompositeMode
	}

	return op
}

func (s *Sprite) SetImage(i image.Image) {
	switch t := i.(type) {
	case *ebiten.Image:
		s.i = t
	default:
		s.i = ebiten.NewImageFromImage(i)
	}
}

// Rotate -> Scale -> Translate
func (s Sprite) Draw(screen *ebiten.Image) {
	if s.i == nil {
		panic("no image")
	}
	if s.isOut(screen.Bounds().Max) {
		return
	}
	screen.DrawImage(s.i, s.Op())
}

// for debugging
func (s Sprite) PrintWHXY(comment string) {
	fmt.Println(comment, s.W, s.H, s.X, s.Y)
}
