package mode

import "github.com/hndada/gosu/draws"

type GameOpts struct {
	TPS int
	FPS int
}

// Todo: fix screen ratio with 4:3 so that in-game
// elements are not stretched regardless of the screen ratio?
type ScreenOpts struct {
	W, H float64
}

func (opts ScreenOpts) WHXY() draws.WHXY {
	return draws.WHXY{W: opts.W, H: opts.H}
}

type MusicOpts struct {
	Volume float64
	Offset int32
}

type SoundOpts struct {
	Volume float64
}
