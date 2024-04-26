package gosu

import (
	"encoding/json"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	draws "github.com/hndada/gosu/draws5"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/choose"
	"github.com/hndada/gosu/scene/play"
)

const (
	SceneChoose = iota
	ScenePlay
)

type Game struct {
	resources *scene.Resources
	options   *scene.Options
	scenes    []scene.Scene
	idx       int
}

func NewGame(root Root) (g *Game) {
	// issue: It jitters when Vsync is enabled.
	// ebiten.SetVsyncEnabled(false)
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(g.ScreenSize.XYInts())
	ebiten.SetWindowTitle("gosu")

	g.loadOptions()
	g.scenes = make([]scene.Scene, 2)
	g.idx = SceneChoose

	g.Scenes["choose"], err = choose.NewScene(g.Config, g.Asset, root)
	if err != nil {
		panic(err)
	}

	return g
}

func (g *Game) loadOptions() {
	jsonData := g.loadOptionsData()
	err := json.Unmarshal(jsonData, g.options)
	if err != nil {
		panic(err)
	}
	g.options.Normalize()
}

func (g *Game) Update() error {
	sc := g.scenes[g.idx]
	switch args := sc.Update().(type) {
	case error:
		fmt.Println("play scene error:", args)
		g.idx = SceneChoose
	case piano.Scorer:
		ebiten.SetWindowTitle("gosu")
		g.idx = SceneChoose
		// debug.SetGCPercent(100)
	case scene.PlayArgs:
		fsys := args.MusicFS
		name := args.ChartFilename
		replay := args.Replay
		scene, err := play.NewScene(g.Config, g.Asset, fsys, name, replay)
		if err != nil {
			fmt.Println("play scene error:", args)
			g.idx = SceneChoose
		} else {
			g.idx = ScenePlay
			g.scenes[g.idx] = scene
			// debug.SetGCPercent(0)
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scenes[g.idx].Draw(draws.Image{Image: screen})
}
func (g Game) ScreenSize() draws.XY {
	return g.options.Screen.Resolution
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.options.Screen.Layout()
}
