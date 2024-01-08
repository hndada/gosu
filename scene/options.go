package scene

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/game"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/mode"
)

// Options passed to each scene.
// Other scenes should not add any additional options or resources.

// Options vs Settings:
// Options consider both game-specific and user-specific options.
// While Settings are for user-specific options only.
type Options struct {
	// These are likely to be modified by the user manually,
	// hence it is at the top of the struct.
	Root RootOptions

	Screen ScreenOptions
	Audio  AudioOptions
	UI     UIOptions
	Game   GameOptions
}

// Roots is a list of root directories to search for files
// such as music, resources, and replays. Each root directory contains
// a set of files which is a directory or a zip file.

// Load server first, then local.
// In web mode, server is the only option.
type RootOptions struct {
	ResourcesPaths []string
	MusicPaths     []string
	ReplaysPaths   []string
}

func (opts *RootOptions) Normalize() {
	// Leading dot and slash is not allowed in fs.
	for i, path := range opts.MusicPaths {
		path = strings.TrimPrefix(path, "..")
		path = strings.TrimPrefix(path, ".")
		path = strings.TrimPrefix(path, "/")
		path = strings.TrimPrefix(path, "\\")
		opts.MusicPaths[i] = path
	}
}

type ScreenOptions struct {
	// Resolution is the physical size of the screen,
	// whereas ScreenSize is the logical size of the screen.
	Resolution           draws.Vector2
	Fullscreen           bool
	BackgroundBrightness float64
	DebugPrint           bool
}

type AudioOptions struct {
	MusicVolume      float64
	SoundVolumeScale float64
	MusicOffset      int32
}

type UIOptions struct {
	MouseCursorScale float64
}

type GameOptions struct {
	Mode            int
	SubMode         int
	ErrorMeterScale float64
	ScoreImageScale float64
	Piano           piano.Options
}

func NewOptions() *Options {
	return &Options{
		Root: RootOptions{
			ResourcesPaths: []string{"resources"},
			MusicPaths:     []string{"music"},
			ReplaysPaths:   []string{"replays"},
		},
		Screen: ScreenOptions{
			Resolution:           draws.Vec2(1600, 900),
			Fullscreen:           false,
			BackgroundBrightness: 0.6,
			DebugPrint:           true,
		},
		Audio: AudioOptions{
			MusicVolume:      0.60,
			SoundVolumeScale: 0.60,
			MusicOffset:      -20,
		},
		UI: UIOptions{
			MouseCursorScale: 1.0,
		},
		Game: GameOptions{
			Mode:            game.ModePiano,
			SubMode:         4,
			ErrorMeterScale: 1.0,
			ScoreImageScale: 1.0,
		},
	}
}

func (opts Options) DebugString() string {
	f := fmt.Fprintf
	var b strings.Builder

	var speedScale float64
	switch opts.Mode {
	case mode.ModePiano:
		speedScale = opts.PianoOptions.SpeedScale
	}

	f(&b, "FPS: %.2f\n", ebiten.ActualFPS())
	f(&b, "TPS: %.2f\n", ebiten.ActualTPS())
	f(&b, "\n")
	// issue: percent literal (%%) does not work.
	f(&b, "Music volume (Ctrl+ Left/Right): %.0f\n", opts.MusicVolume*100)
	f(&b, "Sound volume (Alt+ Left/Right): %.0f\n", opts.SoundVolume*100)
	f(&b, "Music offset (Shift+ Left/Right): %dms\n", opts.MusicOffset)
	f(&b, "Background brightness: (Ctrl+ O/P): %.0f\n", opts.BackgroundBrightness*100)
	f(&b, "Debug print (F12): %v\n", opts.DebugPrint)
	f(&b, "Replay (F11): %v\n", opts.Replay)
	f(&b, "\n")
	// f(&b, "Mode (F1): %d\n", opts.Mode)
	// f(&b, "Sub mode (F2/F3): %d\n", opts.SubMode)
	f(&b, "Speed scale: (Page Down/Up): %.2f\n", speedScale)
	return b.String()
}
