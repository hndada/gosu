package gosu

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene/choose"
)

type Resources struct {
}

type Options struct {
}

type Game struct {
	Root
	Resources
	Options
	Scenes
}

func NewGame(root Root) (g *Game) {
	// issue: It jitters when Vsync is enabled.
	// ebiten.SetVsyncEnabled(false)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(g.ScreenSize.XYInts())
	ebiten.SetWindowTitle("gosu")

	g.Scenes["choose"], err = choose.NewScene(g.Config, g.Asset, root)
	if err != nil {
		panic(err)
	}

	return g
}

func (g *Game) Update() error {
	return g.Scenes.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Scenes.Draw(draws.Image{Image: screen})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenSize.XYInts()
}
