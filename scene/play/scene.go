package play

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"

	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
)

// Interface declares at 'user' package.
type ScenePlay struct {
	SceneModePlay
	SpeedScaleKeyHandler *ctrl.KeyHandler
}
type SceneModePlay interface {
	scene.Scene

	PlayPause()
	Finish() any // Return Scene.
	// IsDone() bool // Draw clear mark.

	SetMusicVolume(float64)
	SetSoundVolume(float64)
	SetOffset(int64)
	SetSpeedScale() // Each mode has its own variable for speed scale.
}

func NewScene(m int, args mode.ScenePlayArgs) (*ScenePlay, error) {
	var (
		play SceneModePlay
		err  error
	)
	switch m {
	case mode.ModePiano:
		play, err = piano.NewSceneModePlay(args)
	case mode.ModeDrum:
		play, err = drum.NewSceneModePlay(args)
	}
	s := &ScenePlay{
		SceneModePlay:        play,
		SpeedScaleKeyHandler: &scene.SpeedScaleKeyHandlers[m],
	}
	// debug.SetGCPercent(0)
	return s, err
}

func (s *ScenePlay) Update() any {
	r := s.SceneModePlay.Update()

	// Settings which affect scene flow.
	if inpututil.IsKeyJustPressed(input.KeyTab) {
		s.PlayPause()
	}
	if inpututil.IsKeyJustPressed(input.KeyEscape) {
		return s.SceneModePlay.Finish()
	}

	// Settings which affect SceneModePlay.
	if fired := scene.MusicVolumeKeyHandler.Update(); fired {
		s.SetMusicVolume(scene.TheSettings.MusicVolume)
	}
	if fired := scene.SoundVolumeKeyHandler.Update(); fired {
		s.SetSoundVolume(scene.TheSettings.SoundVolume)
	}
	if fired := scene.OffsetKeyHandler.Update(); fired {
		s.SetOffset(scene.TheSettings.Offset)
	}
	if fired := s.SpeedScaleKeyHandler.Update(); fired {
		s.SetSpeedScale()
	}

	// Settings which don't affect SceneModePlay.
	scene.BackgroundBrightnessKeyHandler.Update()
	scene.DebugPrintKeyHandler.Update()

	return r
}
func (s ScenePlay) Draw(screen draws.Image) {
	s.SceneModePlay.Draw(screen)
}
