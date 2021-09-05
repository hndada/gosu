package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/mania"
)

type Game struct {
	jm *game.JudgmentMeter
}

func (g *Game) Update(screen *ebiten.Image) error {
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{128, 128, 128, 255})
	g.jm.Sprite.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return game.Settings.ScreenSize.X, game.Settings.ScreenSize.Y
}
func main() {
	g := &Game{}
	game.Settings.ScreenSize = image.Pt(800, 600)
	g.jm = game.NewJudgmentMeter(mania.Judgments[:])
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
