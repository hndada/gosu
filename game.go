package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

type Game struct {
	Scene
}
type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

const SampleRate = 44100

var Context *audio.Context = audio.NewContext(SampleRate)

func NewGame() *Game {
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(ScreenSizeX, ScreenSizeY)
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
	return ScreenSizeX, ScreenSizeY
}
