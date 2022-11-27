package play

import (
	"fmt"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"

	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
)

type ScenePlay interface {
	scene.Scene
	// Pause()
	// Resume()
	IsDone() bool
	Finish() any
}
type Scene struct {
	mode int
	ScenePlay
}

func NewScene(fsys fs.FS, cname string, mode int, mods interface{}, rf *osr.Format) (*Scene, error) {
	var (
		play ScenePlay
		err  error
	)
	switch mode {
	case piano.Mode:
		play, err = piano.NewScenePlay(fsys, cname, mods, rf)
	case drum.Mode:
		play, err = drum.NewScenePlay(fsys, cname, mods, rf)
	}
	// ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	// debug.SetGCPercent(0)
	if err != nil {
		return nil, err
	}
	return &Scene{mode: mode, ScenePlay: play}, err
}
func (s *Scene) Update() any {
	scene.VolumeMusic.Update()
	scene.VolumeSound.Update()
	scene.Brightness.Update()
	scene.Offset.Update()
	scene.SpeedScales[s.mode].Update()
	if inpututil.IsKeyJustPressed(input.KeyEscape) {
		fmt.Println("end the song")
		s.ScenePlay.Finish()
	}
	return s.ScenePlay.Update()
}
func (s Scene) Draw(screen draws.Image) {
	s.ScenePlay.Draw(screen)
}

// type Return struct {
// 	FS     fs.FS
// 	Name   string
// 	Mode   int
// 	Mods   interface{}
// 	Replay *osr.Format
// }
