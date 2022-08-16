package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
)

type Game struct {
	Scene
}
type Scene interface {
	Update() *mode.ScoreResult
	Draw(screen *ebiten.Image)
}

var selectScene *SceneSelect

func NewGame() *Game {
	piano.LoadSkin()
	selectScene = NewSceneSelect()
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(WindowSizeX, WindowSizeY)
	ebiten.SetMaxTPS(mode.MaxTPS)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	g := &Game{
		Scene: selectScene,
	}
	return g
}
func (g *Game) Update() error {
	result := g.Scene.Update()
	if result != nil {
		g.Scene = selectScene // Todo: selectResult
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(screen)
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenSizeX, screenSizeY
}
