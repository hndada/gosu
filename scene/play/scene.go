package play

import (
	"io/fs"

	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"

	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
)

type Scene struct {
	mode int
	scene.Scene
}

func NewScene(fsys fs.FS, cname string, mode int, mods interface{}, rf *osr.Format) (*Scene, error) {
	var (
		scene scene.Scene
		err   error
	)
	switch mode {
	case piano.Mode:
		scene, err = piano.NewScenePlay(fsys, cname, mods, rf)
	case drum.Mode:
		scene, err = drum.NewScenePlay(fsys, cname, mods, rf)
	}
	// ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	// debug.SetGCPercent(0)
	return &Scene{mode: mode, Scene: scene}, err
}
func (s *Scene) Update() any {
	scene.VolumeMusic.Update()
	scene.VolumeSound.Update()
	scene.Brightness.Update()
	scene.Offset.Update()
	scene.SpeedScales[s.mode].Update()
	return s.Scene.Update()
}
func (s Scene) Draw(screen draws.Image) {
	s.Scene.Draw(screen)
}

type Return struct {
	FS     fs.FS
	Name   string
	Mode   int
	Mods   interface{}
	Replay *osr.Format
}
