package game

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/plays/piano"
	"github.com/hndada/gosu/resources"
	"github.com/hndada/gosu/ui"
)

// Resources are loaded from file system.
// Options are set by user and saved to file system.
type Scene interface {
	New(*Game, Args) (Scene, error)
	Update() any
	Draw(screen draws.Image)
	WindowTitle() string
	DebugString() string
}

type Game struct {
	*ui.KeyboardState

	Resources *Resources
	Options   *Options
	Handlers  *Handlers
	Databases *Databases

	SceneSelect  Scene
	ScenePlay    Scene
	CurrentScene Scene
}

func NewGame(fsys fs.FS) (*Game, error) {
	s := &Game{
		KeyboardState: &ui.KeyboardState{},
	}

	if resFS, err := fs.Sub(fsys, "resources"); err == nil {
		s.Resources = NewResources(resFS)
	} else {
		s.Resources = NewResources(resources.DefaultFS)
	}

	// NewOptions is always called, as there
	// might be omitted fields on a local option file.
	s.Options = NewOptions()
	if data, err := fs.ReadFile(fsys, "options.json"); err == nil {
		if err := json.Unmarshal(data, s.Options); err != nil {
			fmt.Printf("Failed to unmarshal options.json: %v\n", err)
		}
	}
	// It is always necessary to set derived values.
	s.Options.Normalize()
	s.Options.Piano.SetDerived()

	s.Handlers = NewHandlers(s.Options, s.KeyboardState)

	dbs, err := NewDatabases(fsys)
	if err != nil {
		panic(err)
	}
	s.Databases = dbs

	// issue: It jitters when Vsync is enabled.
	ebiten.SetTPS(ebiten.SyncWithFPS)
	ebiten.SetWindowSize(s.Options.Resolution.IntValues())
	ebiten.SetWindowTitle("gosu")
	// ebiten.SetVsyncEnabled(false)
	return s, nil
}

// TODO: result page
func (g *Game) Update() error {
	var err error
	switch args := g.CurrentScene.Update().(type) {
	case PlayArgs:
		g.ScenePlay, err = g.ScenePlay.New(g, args)
		if err != nil {
			fmt.Println("play scene error:", args)
			return nil
		}
		g.CurrentScene = g.ScenePlay
		ebiten.SetWindowTitle(g.CurrentScene.WindowTitle())
		// debug.SetGCPercent(0)
	case piano.Scorer:
		g.CurrentScene = g.SceneSelect
		ebiten.SetWindowTitle(g.CurrentScene.WindowTitle())
		// debug.SetGCPercent(100)
	case error:
		fmt.Println("play scene error:", args)
		panic(args)
	}
	return nil
}

func (s Game) Draw(screen *ebiten.Image) {
	s.CurrentScene.Draw(draws.Image{Image: screen})
	str := s.CurrentScene.DebugString()
	if s.Options.DebugPrint {
		str += "\n" + s.Options.DebugString()
	}
	ebitenutil.DebugPrint(screen, str)
}

func (s Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenSizeX, ScreenSizeY
}

// I would keep os package not to be in scene package.
// json.MarshalIndent(options, "", "  ")
// os.WriteFile(fname, data, 0644)
