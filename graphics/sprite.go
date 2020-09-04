package graphics

import "github.com/hajimehoshi/ebiten"

// 포인터 여부는 상관 없게. 메소드로 한번 감싸다
// 메소드로 option 생성; 이미지 원래 크기도 알 수 있음
// 100 스케일로 그리고 확대하면 깨지는 문제도 해결 가능
type Sprite struct {
	i    *ebiten.Image
	w, h int
	x, y int
}

func (s *Sprite) Image() *ebiten.Image { return s.i }

// field값들은 이미 값이 맞춰져있다고 가정
func (s *Sprite) ResetPosition(op *ebiten.DrawImageOptions) {
	op.GeoM.Reset()
	rw, rh := s.i.Size()
	op.GeoM.Scale(float64(s.w)/float64(rw), float64(s.h)/float64(rh))
	op.GeoM.Translate(float64(s.x), float64(s.y))
}

// Stands for Expandible sprite
type ExpSprite struct {
	vertical bool
	i        *ebiten.Image
	wh       int
	x, y     int
}

func (s *ExpSprite) Image(length float64) *ebiten.Image {
	var i *ebiten.Image
	var ratio float64 // only need to consider either one of w or h when scaling
	var count int
	rw, rh := s.i.Size()
	op := &ebiten.DrawImageOptions{}
	if s.vertical {
		i, _ = ebiten.NewImage(int(s.wh), int(length), ebiten.FilterDefault)
		ratio = float64(s.wh) / float64(rw)
		op.GeoM.Scale(ratio, ratio)
		count = int(length / (float64(rh) * ratio))
		for c := 0; c <= count; c++ {
			i.DrawImage(s.i, op)
			op.GeoM.Translate(0, float64(rh)*ratio)
		}
	} else {
		i, _ = ebiten.NewImage(int(length), int(s.wh), ebiten.FilterDefault)
		ratio = float64(s.wh) / float64(rh)
		op.GeoM.Scale(ratio, ratio)
		count = int(length / (float64(rw) * ratio))
		for c := 0; c <= count; c++ {
			i.DrawImage(s.i, op)
			op.GeoM.Translate(float64(rw)*ratio, 0)
		}
	}
	return i
}

func (s *ExpSprite) ResetPosition(op *ebiten.DrawImageOptions) {
	op.GeoM.Reset()
	op.GeoM.Translate(float64(s.x), float64(s.y))
}
