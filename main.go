package main

import (
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/assets"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/play"
)

// scenes should be out of main(), because it will be used in other methods of game.
// var scenes = make(map[string]scene.Scene)

func init() {
	// In 240Hz monitor, TPS is 60 and FPS is 240 at start.
	// SetTPS will set TPS to FPS, hence TPS will be 240 too.
	// I guess ebiten.SetTPS should be called before scene.NewConfig is called.

	// issue: TPS becomes literally -1 when setting with ebiten.SyncWithFPS.
	// ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetTPS(240)
	ctrl.UpdateTPS()
}

// All structs and variables in cmd/* package should be unexported.
type game struct {
	musicRoots []string
	screenSize *draws.Vector2
	scene.Scene
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fsys := os.DirFS(dir)
	cfg := scene.NewConfig()
	asset := scene.NewAsset(cfg, assets.FS)

	g := &game{
		musicRoots: cfg.MusicRoots,
		screenSize: &cfg.ScreenSize,
	}
	musicRoot := g.musicRoots[0]
	musicsFS, err := fs.Sub(fsys, musicRoot)
	if err != nil {
		panic(err)
	}

	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(int(g.screenSize.X), int(g.screenSize.Y))

	scene.TheBaseScene = scene.NewBaseScene(cfg, asset)

	g.loadTestPiano(cfg, asset, musicsFS)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
func (g *game) loadTestPiano(cfg *scene.Config, asset *scene.Asset, musicsFS fs.FS) {
	musicName := "test"
	musicFS, err := fs.Sub(musicsFS, musicName)
	// musicFS := ZipFS(filepath.Join(dir, musicName+".osz"))
	if err != nil {
		panic(err)
	}
	name := "nekodex - circles! (MuangMuangE) [Hard].osu"

	scenePlay, err := play.NewScene(cfg, asset, musicFS, name, mode.ModePiano, piano.Mods{}, nil)
	if err != nil {
		panic(err)
	}
	// fmt.Println(len(scenePlay.ScenePlay.(*piano.ScenePlay).Chart.Notes))
	g.Scene = scenePlay
}

// Todo: implement
// pp.Print(g.Scene.(*piano.ScenePlay).Now())
func (g *game) Update() error {
	switch r := g.Scene.Update().(type) {
	case error:
		return r
		// case "choose":
		// 	g.Scene = scenes["choose"]
		// case "play":
		// 	scene, err := play.NewScene()
		// 	if err != nil {
		// 		return err
		// 	}
		// 	g.Scene = scene
	}
	return nil
}

// Todo: print error using ebitenutil.DebugPrintAt
func (g *game) Draw(screen *ebiten.Image) {
	g.Scene.Draw(draws.Image{Image: screen})
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.screenSize.X), int(g.screenSize.Y)
}