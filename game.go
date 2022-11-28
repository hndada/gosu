package gosu

import (
	"fmt"
	"io/fs"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
	"github.com/hndada/gosu/scene/choose"
	"github.com/hndada/gosu/scene/play"
)

const (
	TPS         = scene.TPS
	ScreenSizeX = scene.ScreenSizeX
	ScreenSizeY = scene.ScreenSizeY
)

// All structs and variables in game package should be unexported
// since the game package is for being called at main via NewGame.
type game struct {
	fs.FS
	scene.Scene
	choose    scene.Scene
	err       error
	countdown int
}
type Settings struct {
	General mode.Settings
	Piano   piano.Settings
	Drum    drum.Settings
	Scene   scene.Settings
}

var (
	UserSettings Settings
	S            = &UserSettings
)

// Todo: tidy NewSettings() and *Settings.Load()?
func Load(fsys fs.FS) {
	data, err := fs.ReadFile(fsys, "settings.toml")
	if err != nil {
		fmt.Println("no settings.toml detected")
	} else {
		S.General = mode.NewSettings()
		S.Piano = piano.NewSettings()
		S.Drum = drum.NewSettings()
		_, err := toml.Decode(string(data), &UserSettings)
		if err != nil {
			fmt.Println(err)
		} else {
			mode.UserSettings.Load(S.General)
			piano.UserSettings.Load(S.Piano)
			drum.UserSettings.Load(S.Drum)
			scene.UserSettings.Load(S.Scene)
		}
	}

	skinFS, err := fs.Sub(fsys, "skin")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = skinFS.Open(".")
	if err != nil {
		fmt.Println("no /skin detected")
	} else {
		mode.UserSkin.Load(skinFS)
		piano.UserSkins.Load(skinFS)
		drum.UserSkin.Load(skinFS)
		scene.UserSkin.Load(skinFS)
	}
}

func NewGame() *game {
	// load(fsys)
	// dir, err := os.Getwd()
	// if err != nil {
	// log.Fatal(err)
	// }
	// scene, err := piano.NewScenePlay(ZipFS(filepath.Join(dir, "test.osz")),
	// "nekodex - circles! (MuangMuangE) [Hard].osu", nil, nil)
	// scene, err := drum.NewScenePlay(os.DirFS(path.Join(dir, "asdf - 1223")),
	// 	"asdf - 1223 (MuangMuangE) [Oni].osu", nil, nil)
	// if err != nil {
	// 	panic(err)
	// }
	g := &game{
		// FS:     fsys,
		Scene:  nil,
		choose: choose.NewScene(),
	}
	return g
}

func (g *game) Update() error {
	if g.countdown > 0 {
		g.countdown--
	}
	if g.Scene == nil {
		g.Scene = g.choose
	}
	switch r := g.Scene.Update().(type) {
	case error:
		if r != nil {
			fmt.Println(r)
		}
	case choose.Return:
		scene, err := play.NewScene(r.FS, r.Name, r.Mode, r.Mods, r.Replay)
		if err != nil {
			fmt.Println(err)
		} else {
			g.Scene = scene
		}
	case mode.Result:
		g.Scene = g.choose
	}
	// g.Scene.Update()
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	if g.Scene == nil {
		return
	}
	g.Scene.Draw(draws.Image{Image: screen})
	if g.err != nil && g.countdown > 0 {
		ebitenutil.DebugPrintAt(screen, g.err.Error(), ScreenSizeX/2, ScreenSizeY/2)
	}
}
func (g *game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return mode.ScreenSizeX, mode.ScreenSizeY
}
