package game

import (
	"image"
	"os"

	"github.com/hajimehoshi/ebiten"
)

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

// Skin2 도 image.Image가 아닌 *ebiten.Image로 해야한다
// 그래야 이미지 자체가 한 번만 로드 됨
func LoadImage2(path string) (*ebiten.Image, error) {
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

// todo: 판정 오차 막대기 그리기
type Sprite22 struct {
	src  *ebiten.Image
	W, H int // desired w, h
	x, y float64
}

func (s Sprite22) Op() *ebiten.DrawImageOptions {
	op := &ebiten.DrawImageOptions{}
	w1, h1 := s.src.Size()
	sw := float64(s.w) / float64(w1)
	sh := float64(s.h) / float64(h1)
	op.GeoM.Scale(sw, sh)
	return op
}
