package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/mania"
)

type Game struct {
	jm *common.JudgmentMeter
}

func (g *Game) Update() error {
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{128, 128, 128, 255})
	g.jm.Sprite.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.Settings.ScreenSizeX, common.Settings.ScreenSizeY
}
func main() {
	g := &Game{}
	common.Settings.ScreenSizeX = 800
	common.Settings.ScreenSizeX = 600
	g.jm = common.NewJudgmentMeter(mania.Judgments[:])
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
