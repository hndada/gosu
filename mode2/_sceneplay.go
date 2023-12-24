package mode

import (
	"github.com/hndada/gosu/draws"
)

// Interface is also used when it uses the unknown struct.
type ScenePlay interface {
	ChartHeader() ChartHeader
	WindowTitle() string
	// int32 is enough for dealing with scene time in millisecond.
	// Maximum duration with int32 is around 24 days.
	Now() int32
	Speed() float64
	IsPaused() bool
	DebugString() string

	SetMusicVolume(float64)
	SetSpeedScale()
	SetMusicOffset(int32)

	Update() any
	Pause()
	Resume()
	Finish() any
	Draw(screen draws.Image)
}
