package scene

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
)

// *piano.Config wins piano.Config
// 1. Easier to pass around.
// 2. Easier to check whether config is defined or not.
type Config struct {
	MusicRoots []string

	ScreenSize  draws.Vector2
	MusicVolume float64
	SoundVolume float64
	MusicOffset int32

	BackgroundBrightness float64
	DebugPrint           bool

	CursorSpriteScale float64
	ClearSpriteScale  float64

	PianoConfig *piano.Config
}

func NewConfig() *Config {
	cfg := &Config{
		MusicRoots: []string{"musics"},

		ScreenSize:  draws.Vector2{X: 1600, Y: 900},
		MusicVolume: 0.30,
		SoundVolume: 0.50,
		MusicOffset: 0,

		BackgroundBrightness: 0.6,
		DebugPrint:           true,

		CursorSpriteScale: 0.1,
		ClearSpriteScale:  0.5,

		PianoConfig: piano.NewConfig(),
	}

	cfg.loadPianoConfig()

	return cfg
}
func (cfg *Config) loadPianoConfig() {
	cfg.PianoConfig = piano.NewConfig()
	cfg.PianoConfig.ScreenSize = &cfg.ScreenSize
	cfg.PianoConfig.MusicVolume = &cfg.MusicVolume
	cfg.PianoConfig.SoundVolume = &cfg.SoundVolume
	cfg.PianoConfig.MusicOffset = &cfg.MusicOffset
}

func (c *Config) NormalizeMusicRoots() {
	if len(c.MusicRoots) == 0 {
		c.MusicRoots = []string{"musics"}
	}

	// Leading dot and slash is not allowed in fs.
	for i, name := range c.MusicRoots {
		name = strings.TrimPrefix(name, "..")
		name = strings.TrimPrefix(name, ".")
		name = strings.TrimPrefix(name, "/")
		name = strings.TrimPrefix(name, "\\")
		c.MusicRoots[i] = name
	}
}

func SetTPS(tps float64) {
	ebiten.SetTPS(int(tps))
	ctrl.UpdateTPS(tps)
	mode.SetTPS(tps)
}
