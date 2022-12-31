package play

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/format/osu"
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
type Chart struct {
	ParentSetId  int
	OsuMode      int
	CS           int
	OsuFile      string
	ChartName    string
	DownloadPath string
}

func (c Chart) Download() []byte {
	u := fmt.Sprintf("https://api.chimu.moe/v1/%s", c.DownloadPath)
	resp, err := http.Get(u)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return body
}

func NewScene(c Chart) (*Scene, error) {
	var mode int
	switch c.OsuMode {
	case osu.ModeMania:
		mode = piano.Mode
	case osu.ModeTaiko:
		mode = drum.Mode
	}
	body := c.Download()
	fsys, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, err
	}
	return newScene(fsys, c.ChartName, mode, nil, nil)

}
func newScene(fsys fs.FS, cname string, mode int, mods interface{}, rf *osr.Format) (*Scene, error) {
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
	// if inpututil.IsKeyJustPressed(input.KeyEscape) {
	// 	// fmt.Println("finish play")
	// 	return s.ScenePlay.Finish()
	// }
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
