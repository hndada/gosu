package mode

import "github.com/hndada/gosu/draws"

type ScreenConfig struct {
	Size draws.Vector2
	TPS  int
	FPS  int
}

type MusicOffset struct {
	MusicVolume float64
	SoundVolume float64
	Offset      int32
}
