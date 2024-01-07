package gosu

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/scene/choose"
	"github.com/hndada/gosu/scene/play"
)

const (
	SceneChoose = iota
	ScenePlay
)

type Game struct {
	scene.Resources
	scene.Options
	scenes []scene.Scene
	idx    int
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
	if scs.Scene == nil {
		scs.scenes = scs.scenes["choose"]
	}

	sc := scs.Scene()
	switch args := sc.Update().(type) {
	case error:
		fmt.Println("play scene error:", args)
		scs.scene = scs.scenes["choose"]
	case piano.Scorer:
		ebiten.SetWindowTitle("gosu")
		// debug.SetGCPercent(100)
		scs.scene = scs.scenes["choose"]
	case scene.PlayArgs:
		fsys := args.MusicFS
		name := args.ChartFilename
		replay := args.Replay
		scene, err := play.NewScene(g.Config, g.Asset, fsys, name, replay)
		if err != nil {
			fmt.Println("play scene error:", args)
			scs.scene = scs.scenes["choose"]
		} else {
			// debug.SetGCPercent(0)
			scs.scene = scene
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scenes[g.idx].Draw(draws.Image{Image: screen})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenSize.XYInts()
}
