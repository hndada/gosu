package gosu

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/ebitenui"
)

// todo: ui
type SceneResult struct {
	ScreenResult *ebiten.Image
	Buttons      []ebitenui.Button
}

func (g *Game) NewSceneResult(sp ScenePlay) *SceneResult {
	sr := &SceneResult{}
	// 리절트 이미지 render
	// 리트라이 -> ScenePlay 에 있는 차트 with 모드 다시 실행
	// 리플레이 -> ScenePlay 에 있는 리플레이 with 모드 실행
	return sr
}

// select를 매번 새로 New 해야할까?
func (s *SceneResult) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		g.ChangeScene(NewSceneSelect())
	}
	for _, b := range s.Buttons {
		b.Update()
	}
	return nil
}

func (s *SceneResult) Draw(screen *ebiten.Image) {
	for _, b := range s.Buttons {
		b.Draw(screen)
	}
}
