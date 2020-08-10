package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// intro 도 1회용으로 넣기
type Title struct {
}

func (s *Title) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.Key0) {
		g.NextScene = &Select{}
		g.TransCountdown = 99
	}
	return nil
}

func (s *Title) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Title: Press Key 0")
}
