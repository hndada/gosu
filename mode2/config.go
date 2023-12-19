package mode

import "github.com/hndada/gosu/draws"

type ScreenConfig struct {
	ScreenSize draws.Vector2
	TPS        int
	FPS        int
}

type MusicConfig struct {
	MusicVolume float64
	SoundVolume float64
	Offset      int32
}
