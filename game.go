package gosu

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
	"github.com/hndada/gosu/scene/selects"
)

// Avoid embedding game.Options directly.
// Pass options as pointers for syncing and saving memory.
type Game struct {
	ctx    *scene.Context
	scn    scene.Scene
	events chan any // from webserver via websocket
}

func NewGame(fsys fs.FS) (*Game, error) {
	ctx, err := scene.NewContext(fsys)
	if err != nil {
		return nil, fmt.Errorf("failed to create scene context: %w", err)
	}

	scn, err := selects.NewScene()
	if err != nil {
		return nil, fmt.Errorf("failed to create selects scene: %w", err)
	}

	g := &Game{
		ctx:    ctx,
		scn:    scn,
		events: make(chan any, 10),
	}

	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(g.ctx.Options.Resolution.IntValues())
	ebiten.SetWindowTitle("gosu")
	// issue: It jitters when Vsync is enabled.
	// ebiten.SetVsyncEnabled(false)

	go g.openWebServer()
	// wait for server to start
	time.Sleep(3 * time.Second)

	return g, nil
}

// TODO: result page
func (g *Game) Update() error {
	// 1. Handle WebSocket events before scene update
	select {
	case args := <-g.events: // assume chan scene.PlayArgs
		switch args := args.(type) {
		case scene.PlayArgs:
			scn, err := play.NewScene(g.ctx, args)
			if err != nil {
				return fmt.Errorf("play scene error: %w", err)
			}
			g.scn = scn
			ebiten.SetWindowTitle(g.scn.WindowTitle())
			return nil
		case error:
			fmt.Println("WS message error:", args)
			return args
		}
	default:
		// no WS message
	}

	// 2. Update current scene
	if g.scn == nil {
		return nil
	}

	switch args := g.scn.Update().(type) {
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
	if g.ctx.Options.DebugPrint {
		str += "\n" + g.ctx.Options.DebugString()
	}
	ebitenutil.DebugPrint(screen, str)
}

func (s Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scene.ScreenSizeX, scene.ScreenSizeY
}
