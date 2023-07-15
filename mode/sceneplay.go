package mode

import (
	"github.com/hndada/gosu/draws"
)

// Interface is also used when it uses the unknown struct.
type ScenePlay interface {
	// get
	ChartHeader() ChartHeader
	WindowTitle() string
	Now() int32 // int32: Maximum duration is around 24 days.
	Speed() float64
	IsPaused() bool

	// set
	SetMusicVolume(float64)
	SetSpeedScale()
	SetMusicOffset(int32)

	// life cycle
	Update() any
	Pause()
	Resume()
	Finish() any

	// draw
	Draw(screen draws.Image)
	DebugPrint(screen draws.Image)
}

func NextDynamics(d *Dynamic, now int32) *Dynamic {
	for d.Next != nil && now >= d.Next.Time {
		d = d.Next
	}
	return d
}
