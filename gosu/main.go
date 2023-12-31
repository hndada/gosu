package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/choose"
	"github.com/hndada/gosu/scene/play"
)

//go:embed resources/*
var defaultResourcesFS embed.FS

type Game struct {
	root fs.FS
	*scene.Config
	*scene.Asset
	scenes map[string]scene.Scene
	scene  scene.Scene
}

func NewGame(root fs.FS) *Game {
	resFS, err := fs.Sub(root, "resources")
	if err != nil {
		resFS = defaultResourcesFS
	}

	cfg := scene.NewConfig()

	asset := scene.NewAsset(cfg, resFS)

	return &Game{
		root:   root,
		Config: cfg,
		Asset:  asset,
		scenes: make(map[string]scene.Scene),
	}
}

func (g *Game) Update() error {
	if g.scene == nil {
		g.scene = g.scenes["choose"]
	}

	switch args := g.scene.Update().(type) {
	case error:
		fmt.Println("play scene error:", args)
		g.scene = g.scenes["choose"]
		// err := args
		// return fmt.Errorf("game update error: %w", err)
	case piano.Scorer:
		ebiten.SetWindowTitle("gosu")
		// debug.SetGCPercent(100)
		g.scene = g.scenes["choose"]
	case scene.PlayArgs:
		fsys := args.MusicFS
		name := args.ChartFilename
		replay := args.Replay
		scene, err := play.NewScene(g.Config, g.Asset, fsys, name, replay)
		if err != nil {
			fmt.Println("play scene error:", args)
			g.scene = g.scenes["choose"]
		} else {
			// debug.SetGCPercent(0)
			g.scene = scene
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(draws.Image{Image: screen})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenSize.XYInts()
}

func main() {
	// issue: It jitters when Vsync is enabled.
	// ebiten.SetVsyncEnabled(false)
	ebiten.SetTPS(ebiten.SyncWithFPS)

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	root := os.DirFS(dir)
	g := NewGame(root)
	ebiten.SetWindowSize(g.ScreenSize.XYInts())
	ebiten.SetWindowTitle("gosu")

	g.scenes["choose"], err = choose.NewScene(g.Config, g.Asset, root)
	if err != nil {
		panic(err)
	}
	g.scene = g.scenes["choose"]

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
