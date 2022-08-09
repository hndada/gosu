package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Scene
}
type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

func NewGame() *Game {
	LoadSkin()
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(WindowSizeX, WindowSizeY)
	ebiten.SetMaxTPS(MaxTPS)
	ebiten.SetRunnableOnUnfocused(true)
	g := &Game{
		Scene: NewSceneSelect(),
	}
	return g
}
func (g *Game) Update() error {
	g.Scene.Update()
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
