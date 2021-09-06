package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/common"
	"github.com/hndada/gosu/engine/ui"
)

type Game struct {
	i *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(g.i, op)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return common.Settings.ScreenSize.X, common.Settings.ScreenSize.Y
}
func main() {
	g := &Game{}
	ri, err := ui.LoadImageImage("test.png")
	if err != nil {
		panic(err)
	}
	// ri := imaging.Rotate90(i)
	// ri := imaging.FlipV(i)
	ei := ebiten.NewImageFromImage(ri)
	g.i = ei
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

// this doesn't work
// func main() {
// 	g := &Game{}
// 	i, err := ui.LoadImageHD("test.png")
// 	if err != nil {
// 		panic(err)
// 	}
// 	ri := imaging.Rotate90(i)
// 	ei := ebiten.NewImageFromImage(ri)
// 	g.i = ei
// 	if err := ebiten.RunGame(g); err != nil {
// 		panic(err)
// 	}
// }
