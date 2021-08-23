package game

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten"
)

// todo: Life, fade-in
type Sprite struct {
	src      *ebiten.Image
	W, H     int // desired w, h
	X, Y     int
	Op       *ebiten.DrawImageOptions
	BornTime int64
	LifeTime int64
}

// todo: IsFixed bool field로 대체하기
func (s Sprite) Fixed() bool { return s.Op != nil } // A sprite that never moves once appears

func (s Sprite) Draw(screen *ebiten.Image) {
	if s.src == nil {
		log.Fatal("s.src is nil")
	}
	if s.Fixed() {
		screen.DrawImage(s.src, s.Op)
	} else {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(s.ScaleW(), s.ScaleH())
		op.GeoM.Translate(float64(s.X), float64(s.Y))
		screen.DrawImage(s.src, op)
	}
}

func (s *Sprite) SetImage(i image.Image) {
	switch i.(type) { // doesn't work
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

// // todo: SetImage 자체를 *ebiten.Image만 받게.
// func (s *Sprite) SetImage(i *ebiten.Image) {
// 	s.src = i
// }

func (s Sprite) ScaleW() float64 {
	w1, _ := s.src.Size()
	return float64(s.W) / float64(w1)
}
func (s Sprite) ScaleH() float64 {
	_, h1 := s.src.Size()
	return float64(s.H) / float64(h1)
}

type LongSprite struct {
	Sprite
	Vertical bool
}

// 사이즈 제한 있어서 *ebiten.Image로 직접 그리면 X
func (s LongSprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w1, h1 := s.src.Size()
	switch s.Vertical {
	case true:
		op.GeoM.Scale(s.ScaleW(), 1)                  // height 쪽은 굳이 scale 하지 않는다
		op.GeoM.Translate(float64(s.X), float64(s.Y)) // important: op는 AB != BA
		q, r := s.H/h1, s.H%h1+1                      // quotient, remainder // temp: +1

		first := s.src.Bounds()
		first.Min = image.Pt(0, h1-r)
		screen.DrawImage(s.src.SubImage(first).(*ebiten.Image), op)
		op.GeoM.Translate(0, float64(r))

		for i := 0; i < q; i++ {
			screen.DrawImage(s.src, op)
			op.GeoM.Translate(0, float64(h1))
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
