package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/play"
)

var scenes = make(map[string]scene.Scene)

// All structs and variables in game package should be unexported
// since the game package is for being called at main via NewGame.
type game struct {
	scene.Scene
}

func main() {
	// Todo: input fs.FS by os.Getwd()
	scenes["choose"] = choose.NewScene()
	g := &game{
		Scene: scenes["choose"],
	}
	// scene, err := piano.NewScenePlay(ZipFS(filepath.Join(dir, "test.osz")),
	// "nekodex - circles! (MuangMuangE) [Hard].osu", nil, nil)
	// scene, err := drum.NewScenePlay(os.DirFS(path.Join(dir, "asdf - 1223")),
	// 	"asdf - 1223 (MuangMuangE) [Oni].osu", nil, nil)
	ebiten.SetWindowSize(scene.ScreenSizeX, scene.ScreenSizeY)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

func (g *game) Update() error {
	switch r := g.Scene.Update().(type) {
	case error:
		return r
	case scene.Return:
		switch r.To {
		case "choose":
			g.Scene = scenes["choose"]
		case "play":
			scene, err := play.NewScene(r.Args.(scene.ScenePlayArgs))
			if err != nil {
				return err
			}
			g.Scene = scene
		}
	}
	return nil
}

// Todo: print error using ebitenutil.DebugPrintAt
func (g *game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(draws.Image{Image: screen})
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scene.ScreenSizeX, scene.ScreenSizeY
}
