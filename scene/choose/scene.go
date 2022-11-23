package choose

import (
	"io/fs"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
)

// Background brightness at Song select: 60% (153 / 255), confirmed.
// Score box color: Gray128 with 50% transparent
// Hovered Score box color: Gray96 with 50% transparent

// https://osu.ppy.sh/docs/index.html#beatmapsetcompact-covers
// cover, card, list, slimcover

// rnkaed status
// graveyard, wip, pending
// ranked, approved, qualified, loved
type Scene struct{}

func NewScene() *Scene {
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	debug.SetGCPercent(100)
	ebiten.SetWindowTitle("gosu")
	return &Scene{}
}
func (s *Scene) Update() any {
	return Return{}
}
func (s Scene) Draw(screen draws.Image) {
}

type Return struct {
	FS     fs.FS
	Name   string
	Mode   int
	Mods   interface{}
	Replay *osr.Format
}
