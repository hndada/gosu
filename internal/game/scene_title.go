package game

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hndada/gosu/mode/mania"
	"image"
)

// intro, close 1회용으로 넣기
// asset/logo 이미지
type SceneTitle struct {
}

func (s *SceneTitle) Update(g *Game) error {
	if ebiten.IsKeyPressed(ebiten.Key0) {
		c := mania.NewChart(`C:\Users\hndada\Documents\GitHub\hndada\gosu\mode\mania\test\test_ln.osu`)
		cp := NewChartPanel(c, image.Pt(200, 200))
		cp.Render()
		cps:=[]ChartPanel{*cp}
		g.NextScene = &SceneSelect{cps, 0}
		g.TransCountdown = g.MaxTransCountDown()
	}
	return nil
}

func (s *SceneTitle) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "SceneTitle: Press Key 0")
}
