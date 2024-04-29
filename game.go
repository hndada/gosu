package gosu

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/play"
	"github.com/hndada/gosu/scene/selects"
)

type Game struct {
	resources *scene.Resources
	options   *scene.Options
	// scenes       []scene.Scene
	sceneSelect  *selects.Scene
	scenePlay    *play.Scene
	currentScene scene.Scene
}

func NewGame(root Root) *Game {
	g := &Game{}
	g.loadResources()
	g.loadOptions()

	sceneSelect, err := selects.NewScene(g.Config, g.Asset, root)
	if err != nil {
		panic(err)
	}
	g.sceneSelect = sceneSelect
	g.currentScene = sceneSelect

	// issue: It jitters when Vsync is enabled.
	// ebiten.SetVsyncEnabled(false)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(g.screenSize().IntValues())
	ebiten.SetWindowTitle("gosu")
	return g
}

func (g *Game) Update() error {
	sc := g.currentScene
	switch args := sc.Update().(type) {
	case error:
		fmt.Println("play scene error:", args)
		g.currentScene = g.sceneSelect

	case piano.Scorer:
		ebiten.SetWindowTitle("gosu")
		g.currentScene = g.sceneSelect
		// debug.SetGCPercent(100)

	case scene.PlayArgs:
		fsys := args.MusicFS
		name := args.ChartFilename
		replay := args.Replay
		scenePlay, err := play.NewScene(g.Config, g.Asset, fsys, name, replay)
		if err != nil {
			fmt.Println("play scene error:", args)
			g.currentScene = g.sceneSelect
			return nil
		}
		ebiten.SetWindowTitle(scenePlay.WindowTitle())
		g.scenePlay = scenePlay
		g.currentScene = scenePlay
		// debug.SetGCPercent(0)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.currentScene.Draw(draws.Image{Image: screen})
}

func (g Game) screenSize() draws.XY {
	return g.options.Screen.Resolution
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenSize().IntValues()
}
