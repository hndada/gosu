package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// intro 도 1회용으로 넣기
type SceneTitle struct {
}

func (s *SceneTitle) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.Key0) {
		g.NextScene = &SceneSelect{}
		g.TransCountdown = 99
	}
	return nil
}

func (s *SceneTitle) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "SceneTitle: Press Key 0")
}
