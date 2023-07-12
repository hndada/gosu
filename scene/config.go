package scene

import (
	"strings"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode/piano"
)

type Config struct {
	MusicRoots []string

	ScreenSize  draws.Vector2
	MusicVolume float64
	SoundVolume float64
	Offset      int32

	BackgroundBrightness float64
	DebugPrint           bool

	CursorSpriteScale float64
	ClearSpriteScale  float64

	PianoConfig *piano.Config
}

func DefaultConfig() Config {
	cfg := Config{
		MusicRoots: []string{"music"},

		ScreenSize:  draws.Vector2{X: 1600, Y: 900},
		MusicVolume: 0.50,
		SoundVolume: 0.50,
		Offset:      -20,

		BackgroundBrightness: 0.6,
		DebugPrint:           true,

		CursorSpriteScale: 0.1,
		ClearSpriteScale:  0.5,

		PianoConfig: piano.DefaultConfig(),
	}

	cfg.loadPianoConfig()

	return cfg
}
func (cfg *Config) loadPianoConfig() {
	cfg.PianoConfig = piano.DefaultConfig()
	cfg.PianoConfig.ScreenSize = &cfg.ScreenSize
	cfg.PianoConfig.MusicVolume = &cfg.MusicVolume
	cfg.PianoConfig.SoundVolume = &cfg.SoundVolume
	cfg.PianoConfig.Offset = &cfg.Offset
}

func (c *Config) NormalizeMusicRoots() {
	if len(c.MusicRoots) == 0 {
		c.MusicRoots = []string{"music"}
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
