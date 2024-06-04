package gosu

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/play"
	"github.com/hndada/gosu/scene/selects"
)

type Game struct {
	resources *scene.Resources
	options   *scene.Options

	currentScene scene.Scene
	sceneSelect  *selects.Scene
	scenePlay    *play.Scene
}

func NewGame(fsys fs.FS) *Game {
	g := &Game{}
	g.resources = scene.NewResources(fsys)
	g.loadOptions()

	sceneSelect, err := selects.NewScene(g.resources, g.options)
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
		scenePlay, err := play.NewScene(g.resources, g.options, args)
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
	if g.options.DebugPrint {
		ebitenutil.DebugPrint(screen, g.options.DebugString())
	}
}

func (g Game) screenSize() draws.XY {
	return g.options.Resolution
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenSize().IntValues()
}
