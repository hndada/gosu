package scene

import (
	"github.com/BurntSushi/toml"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// Currently, TPS should be 1000 or greater.
	// TPS supposed to be multiple of 1000, since only one speed value
	// goes passed per Update, while unit of TransPoint's time is 1ms.
	// TPS affects only on Update(), not on Draw().
	// Todo: add lower TPS support
	TPS = 1000

	// ScreenSize is a logical size of in-game screen.
	ScreenSizeX = 1600
	ScreenSizeY = 900
)

type Settings struct {
	empty      bool
	MusicRoots []string
	// WindowSizeIndex int
	WindowSizeX int
	WindowSizeY int
	VolumeMusic float64
	VolumeSound float64
	CursorScale float64
}

// Default settings should not be directly exported.
// It may be modified by others.
var (
	defaultSettings Settings
	settings        Settings
)

func initSettings() {
	defaultSettings = Settings{
		MusicRoots:  []string{"music"},
		WindowSizeX: 1600,
		WindowSizeY: 900,
		VolumeMusic: 0.25,
		VolumeSound: 0.25,
		CursorScale: 0.1,
	}
	settings = defaultSettings
}
func DefaultSettings() Settings { return defaultSettings }
func CurrentSettings() Settings { return settings }

// Unmatched fields will not be touched, feel free to pre-fill default values.
// Todo: alert warning message to user when some lines are failed to be decoded
func LoadSettings(data string, base Settings) {
	_, err := toml.Decode(data, &settings)
	if err != nil && base.empty {
		panic(err)
	}
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(settings.WindowSizeX, settings.WindowSizeY)
	ebiten.SetTPS(TPS)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
}
