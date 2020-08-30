package config

import (
	"github.com/hajimehoshi/ebiten"
)

// 포인터 여부는 상관 없게. 메소드로 한번 감싸다
// 메소드로 option 생성; 이미지 원래 크기도 알 수 있음
// 100 스케일로 그리고 확대하면 깨지는 문제도 해결 가능
type Sprite struct {
	i    *ebiten.Image
	w, h float64
	x, y float64
}

// field값들은 이미 값이 맞춰져있다고 가정
func (s *Sprite) ResetPosition(op *ebiten.DrawImageOptions) {
	op.GeoM.Reset()
	rw, rh := s.i.Size()
	op.GeoM.Scale(s.w/float64(rw), s.h/float64(rh))
	op.GeoM.Translate(s.x, s.y)
}

func (s *Sprite) Image() *ebiten.Image { return s.i }

type Sprites struct {
	Name            string
	skin            skin
	BoxLeft         Sprite // unscaled
	BoxMiddle       Sprite // unscaled
	BoxRight        Sprite // unscaled
	ChartPanelFrame Sprite // unscaled

	Score            [10]Sprite // unscaled
	ManiaCombo       [10]Sprite // unscaled
	ManiaHitResults  [5]Sprite  // unscaled
	ManiaStages      []ManiaStage // 키별로 option 다름
}

func (s *Sprites) SkinName() string { return s.skin.name }

func (s *Sprites) Render(settings *Settings) {
	// screenSize := settings.screenSize

	for key := range maniaNoteKinds {
		s.ManiaStages[key].Render(settings)
	}
}
