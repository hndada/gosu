package gosu

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/coder/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/plays/piano"
	"github.com/hndada/gosu/scene/play"
	"github.com/hndada/gosu/scene/selects3"
)

type Game struct {
	resources *scene.Resources
	options   *scene.Options
	dbs       *scene.Databases
	ws        *websocket.Conn

	scene scene.Scene
	// currentScene scene.Scene
	// scenes       map[int]scene.Scene
	// sceneSelect  *selects.Scene
	// scenePlay    *play.Scene
}

func NewGame(fsys fs.FS) *Game {
	g := &Game{}
	g.resources = scene.NewResources(fsys)
	g.loadOptions()
	dbs, err := NewDatabases(fsys)
	if err != nil {
		panic(err)
	}
	g.dbs = dbs

	// Connect to the WebSocket server
	ctx := context.Background()
	u := fmt.Sprintf("ws://%s%s", HostName, Port) // "ws://localhost:5488"
	c, _, err := websocket.Dial(ctx, u, nil)
	if err != nil {
		panic(err)
	}
	g.ws = c

	sceneSelect, err := selects3.NewScene(g.resources, g.options, g.dbs, g.ws)
	if err != nil {
		panic(err)
	}
	g.scene = sceneSelect

	RunWebServer()

	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(g.screenSize().IntValues())
	ebiten.SetWindowTitle("gosu")
	// issue: It jitters when Vsync is enabled.
	// ebiten.SetVsyncEnabled(false)

	return g
}

func (g *Game) Update() error {
	switch args := g.scene.Update().(type) {
	case scene.PlayArgs:
		scenePlay, err := play.NewScene(g.resources, g.options, args)
		if err != nil {
			fmt.Println("play scene error:", args)
			// g.currentScene = g.sceneSelect
			return nil
		}
		g.scene = scenePlay

		ebiten.SetWindowTitle(scenePlay.WindowTitle())
		// debug.SetGCPercent(0)
	case piano.Scorer:
		// TODO: result page
		sceneSelect, err := selects3.NewScene(g.resources, g.options, g.dbs, g.ws)
		if err != nil {
			panic(err)
		}
		g.scene = sceneSelect

		ebiten.SetWindowTitle("gosu")
		// debug.SetGCPercent(100)
	case error:
		fmt.Println("play scene error:", args)
		panic(args)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(draws.Image{Image: screen})
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
