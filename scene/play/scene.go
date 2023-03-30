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
	PlayPause()
	IsDone() bool // Draw clear mark.
	Finish() any  // Return Scene.
}
type Scene struct {
	mode int
	ScenePlay
}

func NewScene(fsys fs.FS, cname string, mode int, mods interface{}, rf *osr.Format) (*Scene, error) {
	fmt.Println("NewScene", fsys, cname)
	// fmt.Println(cname, mode)
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
		// fmt.Println("finish play")
		return s.ScenePlay.Finish()
	}
	if inpututil.IsKeyJustPressed(input.KeyTab) {
		s.PlayPause()
	}
	return s.ScenePlay.Update()
}
func (s Scene) Draw(screen draws.Image) {
	s.ScenePlay.Draw(screen)
	if s.IsDone() {
		scene.UserSkin.Clear.Draw(screen, draws.Op{})
	}
}

// type Return struct {
// 	FS     fs.FS
// 	Name   string
// 	Mode   int
// 	Mods   interface{}
// 	Replay *osr.Format
// }
