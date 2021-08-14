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
	var ok bool
	for i := 0; i < 10; i++ {
		filename = fmt.Sprintf("score-%d.png", i)
		path = filepath.Join(skinPath, filename)
		if Skin.score[i], ok = LoadImage(path); !ok {
			panic("failed to load images")
		}
	}
	for i := 0; i < 10; i++ {
		filename = fmt.Sprintf("combo-%d.png", i)
		path = filepath.Join(skinPath, filename)
		if Skin.combo[i], ok = LoadImage(path); !ok {
			panic("failed to load images")
		}
	}
}

// loadSkinImage로 한번에 표시하려면 reflect 써야함
func LoadImage(path string) (*ebiten.Image, bool) {
	empty, _ := ebiten.NewImage(0, 0, ebiten.FilterDefault)
	f, err := os.Open(path)
	if err != nil {
		return empty, false
	}
	defer f.Close()
	src, _, err := image.Decode(f)
	if err != nil {
		return empty, false
	}
	img, _ := ebiten.NewImageFromImage(src, ebiten.FilterDefault)
	return img, true
}

// 길이/4, 높이에 해당하는 거 만큼 롱노트 SubImage
