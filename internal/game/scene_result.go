package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type SceneResult struct {
	// score
	// hp graph
	// hit error deviation
}

func (s *SceneResult) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.Key3) {
		g.NextScene = &SceneSelect{}
		g.TransLifetime = g.MaxTransLifetime()
	}
	// 키 입력 받으면 곡선택 scene으로 이동
	return nil
}

func (s *SceneResult) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "SceneResult: Press Key 3")
}
