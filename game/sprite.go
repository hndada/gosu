package game

import (
	"github.com/hajimehoshi/ebiten"
	"image"
)

// 포인터 여부는 상관 없게. 메소드로 한번 감싸다
// 메소드로 option 생성; 이미지 원래 크기도 알 수 있음
// 100 스케일로 그리고 확대하면 깨지는 문제도 해결 가능
type Sprite struct {
	i *ebiten.Image
	// w, h int
	// x, y int
	p image.Point
}

func (s *Sprite) SetImage(i *ebiten.Image)  { s.i = i }
func (s Sprite) Image() *ebiten.Image       { return s.i }
func (s Sprite) Size() image.Point          { return image.Pt(s.i.Size()) }
func (s *Sprite) SetPosition(p image.Point) { s.p = p }
func (s Sprite) Position() image.Point      { return s.p }

type Sprite2 struct {
	Size     image.Point
	Position image.Point
	i        *ebiten.Image
	Op       *ebiten.DrawImageOptions
}

func (s *Sprite2) SetImage(i *ebiten.Image) { s.i = i }
func (s *Sprite2) Image() *ebiten.Image     { return s.i }

type LongSprite struct {
	Size  image.Point // 여기에 가변 길이 값 들어감
	Start Sprite2
	End   Sprite2 // Position
	i     *ebiten.Image
	Op    *ebiten.DrawImageOptions
}

// // field값들은 이미 값이 맞춰져있다고 가정
// func (s *Sprite) ResetPosition(op *ebiten.DrawImageOptions) {
// 	op.GeoM.Reset()
// 	rw, rh := s.i.Size()
// 	op.GeoM.Scale(float64(s.w)/float64(rw), float64(s.h)/float64(rh))
// 	op.GeoM.Translate(float64(s.x), float64(s.y))
// }

// Stands for Expandible sprite
// type ExpSprite struct {
// 	vertical bool
// 	i        *ebiten.Image
// 	wh       int
// 	x, y     int
// }

// func (s *ExpSprite) Image(length float64) *ebiten.Image {
// 	var i *ebiten.Image
// 	var ratio float64 // only need to consider either one of w or h when scaling
// 	var count int
// 	rw, rh := s.i.Size()
// 	op := &ebiten.DrawImageOptions{}
// 	if s.vertical {
// 		i, _ = ebiten.NewImage(int(s.wh), int(length), ebiten.FilterDefault)
// 		ratio = float64(s.wh) / float64(rw)
// 		op.GeoM.Scale(ratio, ratio)
// 		count = int(length / (float64(rh) * ratio))
// 		for c := 0; c <= count; c++ {
// 			i.DrawImage(s.i, op)
// 			op.GeoM.Translate(0, float64(rh)*ratio)
// 		}
// 	} else {
// 		i, _ = ebiten.NewImage(int(length), int(s.wh), ebiten.FilterDefault)
// 		ratio = float64(s.wh) / float64(rh)
// 		op.GeoM.Scale(ratio, ratio)
// 		count = int(length / (float64(rw) * ratio))
// 		for c := 0; c <= count; c++ {
// 			i.DrawImage(s.i, op)
// 			op.GeoM.Translate(float64(rw)*ratio, 0)
// 		}
// 	}
// 	return i
// }
//
// func (s *ExpSprite) ResetPosition(op *ebiten.DrawImageOptions) {
// 	op.GeoM.Reset()
// 	op.GeoM.Translate(float64(s.x), float64(s.y))
// }
