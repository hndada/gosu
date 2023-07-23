package scene

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
)

// Todo: SoundVolume -> SoundVolumeScale?

// *piano.Config wins piano.Config
// 1. Easier to pass around.
// 2. Easier to check whether config is defined or not.
type Config struct {
	MusicRoots []string
	ScreenSize draws.Vector2

	MusicVolume          float64
	SoundVolume          float64
	MusicOffset          int32
	BackgroundBrightness float64
	DebugPrint           bool
	Replay               bool

	CursorSpriteScale float64
	ListItemWidth     float64 // For folder, chart list
	ListItemHeight    float64
	ListItemShrink    float64 // Items are not focused will be shrinked.
	SearchBoxWidth    float64
	SearchBoxHeight   float64
	ClearSpriteScale  float64

	Mode        int
	SubMode     int
	PianoConfig *piano.Config
}

func NewConfig() *Config {
	screenSize := draws.Vector2{X: 1600, Y: 900}
	cfg := &Config{
		MusicRoots: []string{"music"},
		ScreenSize: screenSize,

		MusicVolume:          0.30,
		SoundVolume:          0.50,
		MusicOffset:          0,
		BackgroundBrightness: 0.6,
		DebugPrint:           true,
		Replay:               false,

		CursorSpriteScale: 0.1,
		ListItemWidth:     550, // 400(card) + 150(list)
		ListItemHeight:    40,
		ListItemShrink:    0.05 * 550,
		SearchBoxWidth:    250,
		SearchBoxHeight:   30,
		ClearSpriteScale:  0.5,

		Mode:    mode.ModePiano,
		SubMode: 4,
	}
	cfg.loadPianoConfig(screenSize)
	return cfg
}

func (cfg *Config) loadPianoConfig(screenSize draws.Vector2) {
	cfg.PianoConfig = piano.NewConfig(screenSize)
	cfg.PianoConfig.ScreenSize = &cfg.ScreenSize
	cfg.PianoConfig.MusicVolume = &cfg.MusicVolume
	cfg.PianoConfig.SoundVolume = &cfg.SoundVolume
	cfg.PianoConfig.MusicOffset = &cfg.MusicOffset
}

func (cfg Config) ListItemCount() int {
	return int(cfg.ScreenSize.Y/cfg.ListItemHeight) + 1
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

func (cfg Config) DebugString() string {
	f := fmt.Fprintf
	var b strings.Builder

	var speedScale float64
	switch cfg.Mode {
	case mode.ModePiano:
		speedScale = cfg.PianoConfig.SpeedScale
	}

	f(&b, "FPS: %.2f\n", ebiten.ActualFPS())
	f(&b, "TPS: %.2f\n", ebiten.ActualTPS())
	f(&b, "\n")
	// issue: percent literal (%%) does not work.
	f(&b, "Music volume (Ctrl+ Left/Right): %.0f\n", cfg.MusicVolume*100)
	f(&b, "Sound volume (Alt+ Left/Right): %.0f\n", cfg.SoundVolume*100)
	f(&b, "Music offset (Shift+ Left/Right): %dms\n", cfg.MusicOffset)
	f(&b, "Background brightness: (Ctrl+ O/P): %.0f\n", cfg.BackgroundBrightness*100)
	f(&b, "Debug print (F12): %v\n", cfg.DebugPrint)
	f(&b, "Replay (F11): %v\n", cfg.Replay)
	f(&b, "\n")
	// f(&b, "Mode (F1): %d\n", cfg.Mode)
	// f(&b, "Sub mode (F2/F3): %d\n", cfg.SubMode)
	f(&b, "Speed scale: (Page Down/Up): %.2f\n", speedScale)
	return b.String()
}
