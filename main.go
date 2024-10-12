package main

import (
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/play"
	"github.com/hndada/gosu/game/selects"
	"github.com/hndada/gosu/plays/piano"
)

var testPlayArgs = game.PlayArgs{
	// ChartFS:       os.DirFS("C:/Users/hndada/Documents/GitHub/gosu/music/nekodex - circles!"),
	// ChartFilename: "nekodex - circles! (MuangMuangE) [Hard].osu",
	ChartFS:       os.DirFS("C:/Users/hndada/Documents/GitHub/gosu/music/cYsmix - triangles"),
	ChartFilename: "cYsmix - triangles (MuangMuangE) [Easy].osu",
	Mods:          piano.Mods{},
	// ReplayFS       fs.FS
	// ReplayFilename string
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	root := os.DirFS(dir)

	g, err := game.NewGame(root)
	if err != nil {
		panic(err)
	}

	{
		scn, err := selects.Scene{}.New(g, nil)
		if err != nil {
			panic(err)
		}
		g.SceneSelect = scn
	}
	{
		scn, err := play.Scene{}.New(g, testPlayArgs)
		if err != nil {
			panic(err)
		}
		g.ScenePlay = scn
	}
	g.CurrentScene = g.ScenePlay

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
