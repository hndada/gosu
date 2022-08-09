package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Scene
}
type Scene interface {
	Update(g *Game)
	Draw(screen *ebiten.Image)
}

var selectScene = NewSceneSelect()

func NewGame() *Game {
	LoadSkin()
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(WindowSizeX, WindowSizeY)
	ebiten.SetMaxTPS(MaxTPS)
	ebiten.SetRunnableOnUnfocused(true)
	g := &Game{
		Scene: selectScene,
	}
	return g
}
func (g *Game) Update() error {
	g.Scene.Update(g)
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
