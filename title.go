package gosu

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/ebitenui"
)

// todo: 화면 비율 조정
type SceneTitle struct {
	Buttons []ebitenui.Button
	// background
	// music player
}

func (g *Game) NewSceneTitle() *SceneTitle {
	s := &SceneTitle{}
	s.Buttons = make([]ebitenui.Button, 0, 6)
	var center, play, multi, edit, option, exit ebitenui.Button


	play := ebitenui.Button{
		MinPt: image.Pt(500, 100),
		Image: ebiten.NewImage(100, 50, ebiten.FilterDefault),
	}
	play.SetOnPressed(func(b *ebitenui.Button) {
		g.ChangeScene(NewSceneSelect()) // 얘 때문에 g가 필요해서 g의 메소드로 만듦
	})
	// 버튼 생성
	s.Buttons = append(s.Buttons, play)

	return s
}
func (s *SceneTitle) Update(g *Game) error {
	for _, b := range s.Buttons {
		b.Update()
	}
	return nil
}

func (s *SceneTitle) Draw(screen *ebiten.Image) {
	for _, b := range s.Buttons {
		b.Draw(screen)
	}
}
