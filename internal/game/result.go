package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Result struct {
	// score
	// hp graph
	// hit error deviation
}

// scene을 나누면 안될것 같다
func (s *Result) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.Key3) {
		g.NextScene=&Select{}
		g.TransCountdown = 99
	}
	// 키 입력 받으면 곡선택 scene으로 이동
	return nil
}

func (s *Result) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Result: Press Key 3")
}