package gosu

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/plays/piano"
	"github.com/hndada/gosu/scene/play"
)

// json.MarshalIndent(options, "", "  ")
// os.WriteFile(fname, data, 0644)

type Game struct {
	scene.States
}

var playArgs = scene.PlayArgs{
	// ChartFS:       os.DirFS("C:/Users/hndada/Documents/GitHub/gosu/cmd/gosu/music/nekodex - circles!"),
	// ChartFilename: "nekodex - circles! (MuangMuangE) [Hard].osu",
	ChartFS:       os.DirFS("C:/Users/hndada/Documents/GitHub/gosu/cmd/gosu/music/cYsmix - triangles"),
	ChartFilename: "cYsmix - triangles (MuangMuangE) [Easy].osu",
	Mods:          piano.Mods{},
	// ReplayFS       fs.FS
	// ReplayFilename string
}

func NewGame(fsys fs.FS) *Game {
	g := &Game{}
	// scn, err := selects.NewScene(g.res, g.opts, g.states, g.hds, g.dbs)
	scn, err := play.NewScene(g.Resources, g.Options, playArgs)
	if err != nil {
		panic(err)
	}
	g.Scene = scn
	return g
}

func (g *Game) Update() error {
	switch args := g.Scene.Update().(type) {
	case scene.PlayArgs:
		scenePlay, err := play.NewScene(g.Resources, g.Options, args)
		if err != nil {
			fmt.Println("play scene error:", args)
			// g.currentScene = g.SceneSelect
			return nil
		}
		g.Scene = scenePlay

		ebiten.SetWindowTitle(scenePlay.WindowTitle())
		// debug.SetGCPercent(0)
	case piano.Scorer:
		// TODO: result page
		g.Scene = g.SceneSelect
		ebiten.SetWindowTitle("gosu")
		// debug.SetGCPercent(100)
	case error:
		fmt.Println("play scene error:", args)
		panic(args)
	}
	return nil
}
