package main

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/play"
)

// Avoid embedding game.Options directly.
// Pass options as pointers for syncing and saving memory.
type Game struct {
	*scene.Context
	scn        scene.Scene
	WSMessages chan scene.PlayArgs
}

func NewGame(fsys fs.FS) (*Game, error) {
	ctx, err := scene.NewContext(fsys)
	if err != nil {
		return nil, err
	}
	g := &Game{
		Context:    ctx,
		WSMessages: make(chan scene.PlayArgs, 10),
	}

	// issue: It jitters when Vsync is enabled.
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(g.Options.Resolution.IntValues())
	ebiten.SetWindowTitle("gosu")
	// ebiten.SetVsyncEnabled(false)

	go OpenWebServer(g)
	time.Sleep(3 * time.Second) // wait for server to start
	OpenBrowser("http://localhost:8080/")

	return g, nil
}

// TODO: result page
func (g *Game) Update() error {
	// 1. Handle WebSocket events before scene update
	select {
	case args := <-g.WSMessages: // assume chan scene.PlayArgs
		scn, err := play.NewScene(g.Context, args)
		if err != nil {
			return fmt.Errorf("play scene error: %w", err)
		}
		g.scn = scn
		ebiten.SetWindowTitle(g.scn.WindowTitle())
		return nil
	default:
		// no WS message
	}

	if g.scn == nil {
		return nil
	}

	switch args := g.scn.Update().(type) {
	case scene.PlayArgs:
	case piano.Scorer:
		// g.CurrentScene = g.SceneSelect
		ebiten.SetWindowTitle(g.scn.WindowTitle())
		// debug.SetGCPercent(100)
	case error:
		fmt.Println("play scene error:", args)
		panic(args)
	}
	return nil
}

func (g Game) Draw(screen *ebiten.Image) {
	if g.scn == nil {
		return
	}
	g.scn.Draw(draws.Image{Image: screen})
	str := g.scn.DebugString()
	if g.Options.DebugPrint {
		str += "\n" + g.Options.DebugString()
	}
	ebitenutil.DebugPrint(screen, str)
}

func (s Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scene.ScreenSizeX, scene.ScreenSizeY
}

// I would keep os package not to be in scene package.
// json.MarshalIndent(options, "", "  ")
// os.WriteFile(fname, data, 0644)
