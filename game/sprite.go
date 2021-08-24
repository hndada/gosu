package game

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
)

type Sprite struct {
	src   *ebiten.Image
	W, H  int // desired w, h
	X, Y  int
	fixed bool // a sprite that never moves once appears
	op    *ebiten.DrawImageOptions
	// BornTime int64
	// LifeTime int64

	Saturation float64
	Dimness    float64
}

func (s Sprite) IsOut(screenSize image.Point) bool {
	return (s.X+s.W < 0 || s.X > screenSize.X ||
		s.Y+s.H < 0 || s.Y > screenSize.Y)
}

func (s Sprite) Draw(screen *ebiten.Image) {
	if s.src == nil {
		log.Fatal("s.src is nil")
	}
	if s.IsOut(screen.Bounds().Max) {
		return
	}
	if s.fixed {
		screen.DrawImage(s.src, s.op)
		// fmt.Println(s.W, s.H, s.X, s.Y)
	} else {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(s.ScaleW(), s.ScaleH())
		op.GeoM.Translate(float64(s.X), float64(s.Y))
		op.ColorM.ChangeHSV(0, s.Saturation, s.Dimness)
		screen.DrawImage(s.src, op)
	}
}

func (s *Sprite) SetImage(i image.Image) {
	switch i.(type) {
	case *ebiten.Image:
		s.src = i.(*ebiten.Image)
	default:
		i2, err := ebiten.NewImageFromImage(i, ebiten.FilterDefault)
		if err != nil {
			log.Fatal(err)
		}
		s.src = i2
	}
}

func (s Sprite) ScaleW() float64 {
	w1, _ := s.src.Size()
	return float64(s.W) / float64(w1)
}
func (s Sprite) ScaleH() float64 {
	_, h1 := s.src.Size()
	return float64(s.H) / float64(h1)
}

func (s *Sprite) SetFixedOp(w, h, x, y int) {
	s.W = w
	s.H = h
	s.X = x
	s.Y = y
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(s.ScaleW(), s.ScaleH())
	op.GeoM.Translate(float64(x), float64(y))
	s.op = op
	s.fixed = true
}

func NewSprite(src *ebiten.Image) Sprite {
	var sprite Sprite
	sprite.src = src

	sprite.Saturation = 1
	sprite.Dimness = 1
	return sprite
}

type LongSprite struct {
	Sprite
	Vertical bool
}

// temp: no need to be method of LongSprite, to make sure only LongSprite uses this
func (s LongSprite) IsOut(w, h, x, y int, screenSize image.Point) bool {
	return x+w < 0 || x > screenSize.X || y+h < 0 || y > screenSize.Y
}

// 사이즈 제한 있어서 *ebiten.Image로 직접 그리면 X
func (s LongSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w1, h1 := s.src.Size()
	switch s.Vertical {
	case true:
		op.GeoM.Scale(s.ScaleW(), 1) // height 쪽은 굳이 scale 하지 않는다
		// important: op is not AB = BA
		x, y := s.X, s.Y
		op.GeoM.Translate(float64(x), float64(y))
		q, r := s.H/h1, s.H%h1+1 // quotient, remainder // temp: +1

		first := s.src.Bounds()
		w, h := s.W, r
		first.Min = image.Pt(0, h1-r)
		if !s.IsOut(w, h, x, y, screen.Bounds().Size()) {
			screen.DrawImage(s.src.SubImage(first).(*ebiten.Image), op)
		}
		op.GeoM.Translate(0, float64(h))
		y += h
		h = h1
		for i := 0; i < q; i++ {
			if !s.IsOut(w, h, x, y, screen.Bounds().Size()) {
				screen.DrawImage(s.src, op)
			}
			op.GeoM.Translate(0, float64(h))
			y += h
		}

	default:
		op.GeoM.Scale(1, s.ScaleH())
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
