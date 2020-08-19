package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/ebitenui"
	"image"
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
	play := ebitenui.Button{
		MinPt: image.Pt(900, 100),
		// Image: ,
	}
	play.SetOnPressed(func(b *ebitenui.Button) {
		g.ChangeScene(NewSceneSelect())
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
