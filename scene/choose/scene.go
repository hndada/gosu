package choose

import (
	"io/fs"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
)

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
