package game

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
)

type Sprite struct {
	i *ebiten.Image
	p image.Point
}

func (s *Sprite) SetImage(i *ebiten.Image)  { s.i = i }
func (s Sprite) Image() *ebiten.Image       { return s.i }
func (s Sprite) Size() image.Point          { return image.Pt(s.i.Size()) }
func (s *Sprite) SetPosition(p image.Point) { s.p = p }
func (s Sprite) Position() image.Point      { return s.p }

// func (s *Sprite) ResetPosition(op *ebiten.DrawImageOptions) {
// 	op.GeoM.Reset()
// 	rw, rh := s.i.Size()
// 	op.GeoM.Scale(float64(s.w)/float64(rw), float64(s.h)/float64(rh))
// 	op.GeoM.Translate(float64(s.x), float64(s.y))
// }

// todo: 판정 오차 막대기 그리기
type Sprite2 struct {
	src  *ebiten.Image
	W, H int // desired w, h
	X, Y float64
}

func (s Sprite2) Op() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	w1, h1 := s.src.Size()
	sw := float64(s.W) / float64(w1)
	sh := float64(s.H) / float64(h1)
	op.GeoM.Scale(sw, sh)
	op.GeoM.Translate(float64(s.X), float64(s.Y))
	return op
}

type LongSprite2 Sprite2

// temp: vertically-long sprite only
// temp: draw directly
// func (s LongSprite2) Draw(screen *ebiten.Image) {
// 	op := s.Op()

// 	w1, h1 := s.src.Size()
// 	sw := float64(s.W) / float64(w1)
// 	sh := float64(s.H) / float64(h1)
// 	q, r := int(s.H)/h1, int(s.H)%h1+1

// 	for i := 0; i < q; i++ {
// 		screen.DrawImage(s.src, op)
// 		op.GeoM.Translate(0, -float64(s.H))
// 	}
// }

// expandible sprite
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

// spritesheet
// 마지막으로 불러온 스킨 불러오기: 처음 불러오는 등 err != nil 일 경우 defaultSkin
const (
	ScoreComma = iota + 10
	ScoreDot
	ScorePercent
)

// todo: image.Image로 바꾸기
var Skin struct {
	name       string
	score      [10]*ebiten.Image
	combo      [10]*ebiten.Image
	hpBarFrame *ebiten.Image
	hpBarColor *ebiten.Image
	// boxLeft         *ebiten.Image
	// boxMiddle       *ebiten.Image
	// boxRight        *ebiten.Image
	// chartPanelFrame *ebiten.Image
}

func LoadSkin(skinPath string) {
	var filename, path string
	var err error
	for i := 0; i < 10; i++ {
		filename = fmt.Sprintf("score-%d.png", i)
		path = filepath.Join(skinPath, filename)
		if Skin.score[i], err = LoadImage(path); err != nil {
			panic("failed to load images")
		}
	}
	for i := 0; i < 10; i++ {
		filename = fmt.Sprintf("combo-%d.png", i)
		path = filepath.Join(skinPath, filename)
		if Skin.combo[i], err = LoadImage(path); err != nil {
			panic("failed to load images")
		}
	}
}

// 길이/4, 높이에 해당하는 거 만큼 롱노트 SubImage

// Skin2 도 image.Image가 아닌 *ebiten.Image로 해야한다
// 그래야 이미지 자체가 한 번만 로드 됨
func LoadImage(path string) (*ebiten.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	ei, _ := ebiten.NewImageFromImage(i, ebiten.FilterDefault)
	return ei, nil
}

var Skin2 struct {
	Number [13]*ebiten.Image // including dot, comma, percent
	// Combo     [10]*ebiten.Image
	BoxLeft   *ebiten.Image
	BoxRight  *ebiten.Image
	BoxMiddle *ebiten.Image

	Cursor      *ebiten.Image
	CursorSmoke *ebiten.Image

	DefaultBG *ebiten.Image
}
