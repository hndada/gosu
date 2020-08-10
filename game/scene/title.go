package scene

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/game"
)

type Title struct {

}
func (s *Title) Update(state *game.State) error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		// state.NextScene=NewSceneMania()
	}
	// 키 입력 받으면 곡선택 scene으로 이동
	return nil
}

func (s *Title) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Title")
}