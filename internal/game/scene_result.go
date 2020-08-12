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
		g.TransCountdown = g.MaxTransCountDown() // todo: 함수 하나로 묶기
	}
	return nil
}

// input, score, hp 다 구현되고 나서 ui 등 고민
func (s *SceneResult) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "SceneResult: Press Key 3")
}
