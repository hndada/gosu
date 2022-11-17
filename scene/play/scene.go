package play

import (
	"io/fs"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/mode"
)

type Scene struct{}

func NewScene(fsys fs.FS, name string,
	mode int, mods interface{}, replay *osr.Format) (*Scene, error) {
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	debug.SetGCPercent(0)
	ebiten.SetWindowTitle("music name")
	return &Scene{}, nil
}
func (s *Scene) Update() any {
	return Return{}
}
func (s Scene) Draw(screen draws.Image) {
}

type Return struct {
	mode.Result
}
