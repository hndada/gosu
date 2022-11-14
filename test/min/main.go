package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct{}

func main() {
	g := Game{}
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
func (g Game) Update() error { return nil }
func (g Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %4.2f\nTPS: %4.2f",
		ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
