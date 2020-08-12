package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// intro, close 1회용으로 넣기
// asset/logo 이미지
type SceneTitle struct {
}

func (s *SceneTitle) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.Key0) {
		g.NextScene = &SceneSelect{}
		g.TransCountdown = g.MaxTransCountDown()
	}
	return nil
}

func (s *SceneTitle) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "SceneTitle: Press Key 0")
}
