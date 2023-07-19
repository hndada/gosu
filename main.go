package main

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/assets"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/play"
)

// scenes should be out of main(), because it will be used in other methods of game.
// Todo: save the scene to gosu.scenes
// var scenes = make(map[string]scene.Scene)

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
	musicRoots []string
	screenSize *draws.Vector2
	scene      interface {
		Update() any
		Draw(screen draws.Image)
	}
}

var print = func(args ...any) { fmt.Printf("%+v\n", args...) }

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

	g.loadTestPiano(cfg, asset, musicsFS)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
func (g *game) loadTestPiano(cfg *scene.Config, asset *scene.Asset, musicsFS fs.FS) {
	musicName := "nekodex - circles!"
	musicFS, err := fs.Sub(musicsFS, musicName)
	// musicFS := ZipFS(filepath.Join(dir, musicName+".osz"))
	if err != nil {
		panic(err)
	}
	name := "nekodex - circles! (MuangMuangE) [Hard].osu"

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fsys := os.DirFS(dir)

	replay, err := mode.NewReplay(fsys, "format/osr/testdata/circles(7k).osr", 7)
	if err != nil {
		panic(err)
	}

	scenePlay, err := play.NewScene(cfg, asset, musicFS, name, replay)
	if err != nil {
		panic(err)
	}
	// fmt.Println(len(scenePlay.ScenePlay.(*piano.ScenePlay).Chart.Notes))
	g.scene = scenePlay
}

// Todo: implement
// pp.Print(g.Scene.(*piano.ScenePlay).Now())
func (g *game) Update() error {
	switch r := g.scene.Update().(type) {
	case error:
		return fmt.Errorf("game update error: %w", r)
		// case "choose":
		// ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
		// debug.SetGCPercent(100)
		// ebiten.SetWindowTitle("gosu")
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

// Todo: print error using ebitenutil.DebugPrintAt?
func (g *game) Draw(screen *ebiten.Image) {
	g.scene.Draw(draws.Image{Image: screen})
}

func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(g.screenSize.X), int(g.screenSize.Y)
}
