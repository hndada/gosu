package scene

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/mode"
)

const (
	TPS         = mode.TPS
	ScreenSizeX = mode.ScreenSizeX
	ScreenSizeY = mode.ScreenSizeY
)

// Todo: settings -> settingsType?
// struct as a type of settings value is discouraged.
type settings struct {
	MusicRoots  []string
	WindowSize  int
	NumberScale float64
	CursorScale float64

	backgroundBrightness *float64
}

const (
	WindowSizeStandard = iota
	WindowSizeFull
)

// Default settings should not be directly exported.
// It may be modified by others.
var defaultSettings = settings{
	MusicRoots:           []string{"music"},
	WindowSize:           WindowSizeStandard,
	VolumeMusic:          0.25,
	VolumeSound:          0.25,
	BackgroundBrightness: 0.6,
	NumberScale:          0.65,
	CursorScale:          0.1,
}
var Settings = defaultSettings

func ResetSettings() { Settings = defaultSettings }

// Unmatched fields will not be touched, feel free to pre-fill default values.
// Todo: alert warning message to user when some lines are failed to be decoded
func LoadSettings(data string) {
	_, err := toml.Decode(data, &Settings)
	if err != nil {
		fmt.Println(err)
	}
	if len(Settings.MusicRoots) == 0 {
		Settings.MusicRoots = append(Settings.MusicRoots, "music")
	}
	Normalize(&Settings.VolumeMusic, 0, 1)
	switch Settings.WindowSize {
	case WindowSizeStandard:
		ebiten.SetWindowSize(1600, 900)
	case WindowSizeFull:
		ebiten.SetFullscreen(true)
	}
	Normalize(&Settings.VolumeMusic, 0, 1)
	Normalize(&Settings.VolumeSound, 0, 1)
	Normalize(&Settings.BackgroundBrightness, 0, 1)
	ebiten.SetTPS(TPS)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}
