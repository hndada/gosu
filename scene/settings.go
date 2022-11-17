package scene

import (
	"fmt"
	"io/fs"

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
const (
	WindowSizeStandard = iota
	WindowSizeFull
)

// struct type is discouraged for settings.
type Settings struct {
	MusicRoots           []string
	WindowSize           int
	VolumeMusic          float64
	VolumeSound          float64
	BackgroundBrightness float64
	CursorScale          float64
	ScoreScale           float64
	ScoreDigitGap        float64
}

// Default settings should not be directly exported.
// It may be modified by others.
var (
	defaultSettings Settings
	settings        Settings
)

func init() {
	initSettings()
	initSkin()
}

func initSettings() {
	defaultSettings = Settings{
		MusicRoots:           []string{"music"},
		WindowSize:           WindowSizeStandard,
		VolumeMusic:          0.25,
		VolumeSound:          0.25,
		BackgroundBrightness: 0.6,
		CursorScale:          0.1,
		ScoreScale:           0.65,
		ScoreDigitGap:        0,
	}
	settings = defaultSettings
}
func DefaultSettings() Settings { return defaultSettings }
func CurrentSettings() Settings { return settings }

// Unmatched fields will not be touched, feel free to pre-fill default values.
// Todo: alert warning message to user when some lines are failed to be decoded
func LoadSettings(data string) {
	_, err := toml.Decode(data, &settings)
	if err != nil {
		fmt.Println(err)
	}
	if len(settings.MusicRoots) == 0 {
		settings.MusicRoots = append(settings.MusicRoots, "music")
	}
	switch settings.WindowSize {
	case WindowSizeStandard:
		ebiten.SetWindowSize(1600, 900)
	case WindowSizeFull:
		ebiten.SetFullscreen(true)
	}
	if settings.VolumeMusic > 1 {
		settings.VolumeMusic = 1
	}
	if settings.VolumeMusic < 0 {
		settings.VolumeMusic = 0
	}
	if settings.VolumeSound > 1 {
		settings.VolumeSound = 1
	}
	if settings.VolumeSound < 0 {
		settings.VolumeSound = 0
	}
	if settings.BackgroundBrightness > 1 {
		settings.BackgroundBrightness = 1
	}
	if settings.BackgroundBrightness < 0 {
		settings.BackgroundBrightness = 0
	}
	ebiten.SetWindowTitle("gosu")
	ebiten.SetTPS(TPS)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
}

// Load is called at game.go.
func Load(fsys fs.FS) {
	data, _ := fs.ReadFile(fsys, "settings.toml")
	LoadSettings(string(data))
	LoadSkin(fsys, LoadSkinPlay)
}
