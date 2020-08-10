package scene

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/game"
)

type Result struct {
	// score
	// hp graph
	// hit error deviation
}

func (s *Result) Update(state *game.State) error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		panic("game over!")
	}
	// 키 입력 받으면 곡선택 scene으로 이동
	return nil
}

func (s *Result) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Result")
}