package game

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/resources"
	"github.com/hndada/gosu/ui"
)

// Resources, Options, States.
// Resources are loaded from file system.
// Options are set by user and saved to file system.
// States are generated when runtime, and not saved.

// I would keep os package not to be in scene package.

type States struct {
	*ui.KeyboardState
	// ws        *websocket.Conn

	Resources *Resources
	Options   *Options
	Handlers  *Handlers
	Databases *Databases

	Scene       Scene
	SceneSelect Scene
}

func NewStates(fsys fs.FS) (*States, error) {
	s := &States{
		KeyboardState: &ui.KeyboardState{},
	}

	if _, err := fs.Stat(fsys, "resources"); err != nil {
		s.Resources = NewResources(resources.DefaultFS)
	} else {
		s.Resources = NewResources(fsys)
	}

	// NewOptions is always called, as there
	// might be omitted fields on a local option file.
	s.Options = NewOptions()
	if _, err := fs.Stat(fsys, "options.json"); err == nil {
		data, err := fs.ReadFile(fsys, "options.json")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, s.Options)
		if err != nil {
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

func (s States) Draw(screen *ebiten.Image) {
	s.Scene.Draw(draws.Image{Image: screen})
	str := s.Scene.DebugString()
	if s.Options.DebugPrint {
		str += "\n" + s.Options.DebugString()
	}
	ebitenutil.DebugPrint(screen, str)
}

func (s States) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenSizeX, ScreenSizeY
}
