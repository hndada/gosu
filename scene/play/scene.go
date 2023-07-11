package play

import (
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"

	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
)

// ScenePlay: struct, PlayScene: function
// Interface declares at 'user' package.
type ScenePlay struct {
	SceneModePlay
	SpeedScaleKeyHandler *ctrl.KeyHandler

	Background           scene.BackgroundDrawer
	backgroundBrightness *float64
	debugPrint           *bool
}
type SceneModePlay interface {
	scene.Scene

	PlayPause()
	Finish() any // Return Scene.

	SetMusicVolume(float64)
	SetSoundVolume(float64)
	SetOffset(int64)
	SetSpeedScale() // Each mode has its own variable for speed scale.
	DebugPrint(draws.Image)
}

func init() {
	skins.loadSkin(4)
	skins.loadSkin(7)
}
func NewScene(m int, args mode.ScenePlayArgs) (*ScenePlay, error) {
	var (
		play SceneModePlay
		err  error
	)
	switch m {
	case mode.ModePiano:
		play, err = piano.NewSceneModePlay(args)
		// case mode.ModeDrum:
		// play, err = drum.NewSceneModePlay(args)
	}

	s := &ScenePlay{
		SceneModePlay:        play,
		SpeedScaleKeyHandler: &scene.SpeedScaleKeyHandlers[m],

		backgroundBrightness: &scene.TheSettings.BackgroundBrightness,
		debugPrint:           &scene.TheSettings.DebugPrint,
	}
	s.SetMusicVolume(scene.TheSettings.MusicVolume)
	s.SetSoundVolume(scene.TheSettings.SoundVolume)
	s.SetOffset(scene.TheSettings.Offset)
	s.SetSpeedScale()
	// Todo: pass background drawer from choose scene
	// s.Background = scene.BackgroundDrawer{
	// 	Sprite: scene.NewBackground(args.FS, c.ImageFilename),
	// }
	// if s.Background.Sprite.IsEmpty() {
	// 	s.Background.Sprite = skin.DefaultBackground
	// }
	// ebiten.SetWindowTitle(c.WindowTitle())
	// debug.SetGCPercent(0)
	return s, err
}

func (s *ScenePlay) Update() any {
	scene.BackgroundBrightnessKeyHandler.Update()

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

	scene.DebugPrintKeyHandler.Update()
	return r
}
func (s ScenePlay) Draw(screen draws.Image) {
	s.Background.Draw(screen)
	s.SceneModePlay.Draw(screen)
	if *s.debugPrint {
		s.SceneModePlay.DebugPrint(screen)
	}
}
