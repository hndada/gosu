package main

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	defaultasset "github.com/hndada/gosu/asset"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/choose"
	"github.com/hndada/gosu/scene/play"
)

func init() {
	// In 240Hz monitor, TPS is 60 and FPS is 240 at start.
	// SetTPS will set TPS to FPS, hence TPS will be 240 too.
	// I guess ebiten.SetTPS should be called before scene.NewConfig is called.

	// issue: It jitters when Vsync is enabled.
	ebiten.SetVsyncEnabled(false)
	// issue: TPS becomes literally -1 when setting with ebiten.SyncWithFPS.
	// ebiten.SetTPS(ebiten.SyncWithFPS)
	scene.SetTPS(480)
}

// All structs and variables in cmd/* package should be unexported.
type game struct {
	root fs.FS
	*scene.Config
	*scene.Asset
	scenes map[string]scene.Scene
	scene  scene.Scene
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	root := os.DirFS(dir)
	g := newGame(root)
	ebiten.SetWindowSize(g.ScreenSize.XYInt())
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

func newGame(root fs.FS) *game {
	cfg := scene.NewConfig()

	assetFS, err := fs.Sub(root, "asset")
	if err != nil {
		assetFS = defaultasset.FS
	}
	asset := scene.NewAsset(cfg, assetFS)

	return &game{
		root:   root,
		Config: cfg,
		Asset:  asset,
		scenes: make(map[string]scene.Scene),
	}
}

func (g *game) Update() error {
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

// Memo: print error using ebitenutil.DebugPrintAt?
func (g *game) Draw(screen *ebiten.Image) {
	g.scene.Draw(draws.Image{Image: screen})
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.ScreenSize.XYInt()
}
