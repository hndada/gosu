package game

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

type LongSprite struct {
	Sprite
	Vertical bool
}

// temp: no need to be method of LongSprite, to make sure only LongSprite uses this
func (s LongSprite) isOut(w, h, x, y int, screenSize image.Point) bool {
	return x+w < 0 || x > screenSize.X || y+h < 0 || y > screenSize.Y
}

// 사이즈 제한 있어서 *ebiten.Image로 직접 그리면 X
func (s LongSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w1, h1 := s.src.Size()
	switch s.Vertical {
	case true:
		op.GeoM.Scale(s.scaleW(), 1) // height 쪽은 굳이 scale 하지 않는다
		// important: op is not AB = BA
		x, y := s.X, s.Y
		op.GeoM.Translate(float64(x), float64(y))
		q, r := s.H/h1, s.H%h1+1 // quotient, remainder // temp: +1

		first := s.src.Bounds()
		w, h := s.W, r
		first.Min = image.Pt(0, h1-r)
		if !s.isOut(w, h, x, y, screen.Bounds().Size()) {
			screen.DrawImage(s.src.SubImage(first).(*ebiten.Image), op)
		}
		op.GeoM.Translate(0, float64(h))
		y += h
		h = h1
		for i := 0; i < q; i++ {
			if !s.isOut(w, h, x, y, screen.Bounds().Size()) {
				screen.DrawImage(s.src, op)
			}
			op.GeoM.Translate(0, float64(h))
			y += h
		}

	default:
		op.GeoM.Scale(1, s.scaleH())
		op.GeoM.Translate(float64(s.X), float64(s.Y))
		q, r := s.W/w1, s.W%w1+1 // temp: +1

		for i := 0; i < q; i++ {
			screen.DrawImage(s.src, op)
			op.GeoM.Translate(float64(w1), 0)
		}

		last := s.src.Bounds()
		last.Max = image.Pt(r, h1)
		screen.DrawImage(s.src.SubImage(last).(*ebiten.Image), op)
	}
}
