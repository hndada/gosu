package scene

import (
	"strings"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode/piano"
)

type Config struct {
	ScreenSize draws.Vector2
	MusicRoots []string

	CursorScale float64
	ClearScale  float64

	MusicVolume          float64
	SoundVolume          float64
	Offset               int64
	BackgroundBrightness float64
	DebugPrint           bool

	PianoConfig *piano.Config
}

func DefaultConfig() Config {
	return Config{
		ScreenSize: draws.Vector2{X: 1600, Y: 900},
		MusicRoots: []string{"music"},

		CursorScale: 0.1,
		ClearScale:  0.5,

		MusicVolume:          0.50,
		SoundVolume:          0.50,
		Offset:               -20,
		BackgroundBrightness: 0.6,
		DebugPrint:           true,

		PianoConfig: piano.DefaultConfig(),
	}
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
