package scene

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/defaultskin"
	"github.com/hndada/gosu/mode"
)

const (
	TPS         = mode.TPS
	ScreenSizeX = mode.ScreenSizeX
	ScreenSizeY = mode.ScreenSizeY
)

type Settings struct {
	MusicRoots  []string
	WindowSize  int
	CursorScale float64
}

const (
	WindowSizeStandard = iota
	WindowSizeFull
)

var (
	DefaultSettings = Settings{
		MusicRoots:  []string{"music"},
		WindowSize:  WindowSizeStandard,
		CursorScale: 0.1,
	}
	UserSettings = DefaultSettings
	S            = &UserSettings
)

func init() {
	DefaultSettings.process()
	UserSettings.process()
	DefaultSkin.Load(defaultskin.FS)
}

func (settings *Settings) Load(data string) {
	_, err := toml.Decode(data, settings)
	if err != nil {
		fmt.Println(err)
	}
	defer settings.process()

	// Leading dot and slash is not allowed in fs.
	for i, name := range settings.MusicRoots {
		name = strings.TrimPrefix(name, ".")
		name = strings.TrimPrefix(name, ".") // There might be two dots
		name = strings.TrimPrefix(name, "/")
		name = strings.TrimPrefix(name, "\\")
		settings.MusicRoots[i] = name
	}
	if len(settings.MusicRoots) == 0 {
		settings.MusicRoots = []string{"music"}
	}
	mode.Normalize(&settings.WindowSize, 0, WindowSizeFull)
	mode.Normalize(&settings.CursorScale, 0, 2)
}
func (settings *Settings) process() {
	switch settings.WindowSize {
	case WindowSizeStandard:
		ebiten.SetWindowSize(1600, 900)
	case WindowSizeFull:
		ebiten.SetFullscreen(true)
	}
	ebiten.SetTPS(TPS)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}
