package game

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/plays"
	"github.com/hndada/gosu/plays/piano"
)

// Options passed to each scene.
// Other scenes should not add any additional options or resources.

// Options vs Settings:
// Options consider both game-specific and user-specific options.
// While Settings are for user-specific options only.

// Roots is a list of root directories to search for files
// such as music, resources, and replays. Each root directory contains
// a set of files which is a directory or a zip file.

// Load server first, then local.
// In web mode, server is the only option.
const (
	ScreenSizeX = plays.ScreenSizeX
	ScreenSizeY = plays.ScreenSizeY
)

type Options struct {
	// These are likely to be modified by the user manually,
	// hence it is at the top of the struct.
	ResourcesPaths []string
	MusicPaths     []string
	ReplaysPaths   []string

	// screenSize is the logical size of the screen, and
	// Resolution is the physical size of the screen.
	screenSize           draws.XY
	Resolution           draws.XY
	IsFullscreen         bool
	BackgroundBrightness float32
	DebugPrint           bool

	MusicVolume      float64
	SoundVolumeScale float64
	MusicOffset      int32

	MouseCursorImageScale float64

	Mode            int
	SubMode         int
	ErrorMeterScale float64
	ScoreImageScale float64
	Piano           *piano.Options
}

// Todo: *Options vs Options
// But I think, to use pointer, *Options is inevitable.
func NewOptions() *Options {
	opts := &Options{
		ResourcesPaths: []string{"resources"},
		MusicPaths:     []string{"music"},
		ReplaysPaths:   []string{"replays"},

		screenSize:           draws.NewXY(plays.ScreenSizeX, plays.ScreenSizeY),
		Resolution:           draws.NewXY(1600, 900),
		IsFullscreen:         false,
		BackgroundBrightness: 0.6,
		DebugPrint:           true,

		MusicVolume:      0.60,
		SoundVolumeScale: 0.60,
		MusicOffset:      -20,

		MouseCursorImageScale: 1.0,

		Mode:            plays.ModePiano,
		SubMode:         4,
		ErrorMeterScale: 1.0,
		ScoreImageScale: 1.0,
		Piano:           piano.NewOptions(),
	}
	opts.Piano.SetDerived()
	return opts
}

func (opts *Options) Normalize() {
	// Leading dot and slash is not allowed in fs.
	for i, path := range opts.MusicPaths {
		path = strings.TrimPrefix(path, "..")
		path = strings.TrimPrefix(path, ".")
		path = strings.TrimPrefix(path, "/")
		path = strings.TrimPrefix(path, "\\")
		opts.MusicPaths[i] = path
	}
}

func (opts Options) DebugString() string {
	f := fmt.Fprintf
	var b strings.Builder

	var speedScale float64
	switch opts.Mode {
	case plays.ModePiano:
		speedScale = opts.Piano.SpeedScale
	}

	f(&b, "FPS: %.2f\n", ebiten.ActualFPS())
	f(&b, "TPS: %.2f\n", ebiten.ActualTPS())
	f(&b, "\n")
	// issue: percent literal (%%) does not work.
	f(&b, "Music volume (Ctrl+ Left/Right): %.0f\n", opts.MusicVolume*100)
	f(&b, "Sound volume (Alt+ Left/Right): %.0f\n", opts.SoundVolumeScale*100)
	f(&b, "Music offset (Shift+ Left/Right): %dms\n", opts.MusicOffset)
	f(&b, "Background brightness: (Ctrl+ O/P): %.0f\n", opts.BackgroundBrightness*100)
	f(&b, "Debug print (F12): %v\n", opts.DebugPrint)
	// f(&b, "Replay (F11): %v\n", opts.Replay)
	f(&b, "\n")
	// f(&b, "Mode (F1): %d\n", opts.Mode)
	// f(&b, "Sub mode (F2/F3): %d\n", opts.SubMode)
	f(&b, "Speed scale: (Page Down/Up): %.2f\n", speedScale)
	return b.String()
}
