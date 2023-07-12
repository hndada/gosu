package main

import (
	"os"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
)

var scenes = make(map[string]scene.Scene)

// All structs and variables in cmd/* package should be unexported.
type game struct {
	screenSize *draws.Vector2
	scene.Scene
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	musicName := "asdf - 1223"
	fsys := os.DirFS(path.Join(dir, musicName))
	name := "nekodex - circles! (MuangMuangE) [Hard].osu"

	// oszName := "asdf - 1223.osz"
	// fsys:= ZipFS(filepath.Join(dir, oszName))
	// name := "nekodex - circles! (MuangMuangE) [Hard].osu"

	// In 240Hz monitor, TPS is 60 and FPS is 240 at start.
	// SetTPS will set TPS to FPS, hence TPS will be 240 too.
	// I guess ebiten.SetTPS should be called before scene.DefaultConfig is called.
	ebiten.SetTPS(ebiten.SyncWithFPS)

	cfg := scene.DefaultConfig()
	asset := scene.NewAsset(cfg, fsys)
	mods := piano.Mods{}
	var rf *osr.Format = nil
	scenePlay, err := piano.NewScenePlay(cfg.PianoConfig, asset.PianoAssets, fsys, name, mods, rf)
	if err != nil {
		panic(err)
	}
	g := &game{
		screenSize: &cfg.ScreenSize,
		Scene:      scenePlay,
	}

	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(int(g.screenSize.X), int(g.screenSize.Y))
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}

// Todo: implement
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
