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
	MusicRoots []string
	// WindowSizeIndex int
	WindowSizeX int
	WindowSizeY int
	VolumeMusic float64
	VolumeSound float64
	CursorScale float64
}

var (
	settings *Settings
	// Default settings should not be directly exported.
	// It may be modified by others.
	defaultSettings *Settings
)

func (Settings) init() {
	// //go:embed settings.toml
	// var data string
	// s, err := Settings{}.Load(data)
	// if err != nil {
	// 	panic(err)
	// }
	defaultSettings = &Settings{
		MusicRoots:  []string{"music"},
		WindowSizeX: 1600,
		WindowSizeY: 900,
		VolumeMusic: 0.25,
		VolumeSound: 0.25,
		CursorScale: 0.1,
	}
}

//	func (Settings) Default() Setter {
//		return Settings{
//			MusicRoots:   []string{"music"},
//			WindowSizeX:  1600,
//			WindowSizeY:  900,
//			VolumeMusic:  0.25,
//			VolumeSound: 0.25,
//			CursorScale:  0.1,
//		}
//	}
func (Settings) Default() Setter { return *defaultSettings }
func (Settings) Current() Setter { return *settings }
func (Settings) Set(s Setter) {
	settings = s.(*Settings)
	ebiten.SetWindowTitle("gosu")
	ebiten.SetWindowSize(settings.WindowSizeX, settings.WindowSizeY)
	ebiten.SetTPS(TPS)
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
}

// Unmatched fields will not be touched, feel free to pre-fill default values.
// Todo: alert warning message to user when some lines are failed to be decoded
func (Settings) Load(data any) (Setter, error) {
	s := Settings{}.Default()
	_, err := toml.Decode(data.(string), &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
