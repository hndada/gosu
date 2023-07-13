package play

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
)

// ScenePlay: struct, PlayScene: function
// Interface declares at 'user' package.
type Scene struct {
	cfg *scene.Config
	// asset *scene.Asset

	*scene.BaseScene
	mode int
	mode.ScenePlay
	drawBackground func(draws.Image)
}

func NewScene(cfg *scene.Config, asset *scene.Asset, fsys fs.FS, name string,
	_mode int, mods any, rf *osr.Format) (s *Scene, err error) {

	s = new(Scene)
	s.cfg = cfg
	s.BaseScene = scene.TheBaseScene
	s.mode = _mode
	switch s.mode {
	case mode.ModePiano:
		mods := mods.(piano.Mods)
		s.ScenePlay, err = piano.NewScenePlay(cfg.PianoConfig, asset.PianoAssets,
			fsys, name, mods, rf)
	}

	bgFilename := s.ScenePlay.BackgroundFilename()
	bgSprite := scene.NewBackgroundSprite(fsys, bgFilename, cfg.ScreenSize)
	if bgSprite.IsEmpty() {
		bgSprite = asset.DefaultBackgroundSprite
	}
	s.drawBackground = scene.NewDrawBackgroundFunc(bgSprite,
		cfg.ScreenSize, &cfg.BackgroundBrightness)

	ebiten.SetWindowTitle(s.WindowTitle())
	// debug.SetGCPercent(0)
	return
}

// The order of function calls may not consistent with
// the order of methods of mode.ScenePlay.
func (s *Scene) Update() any {
	if inpututil.IsKeyJustPressed(input.KeyTab) {
		if s.IsPaused() {
			s.Resume()
		} else {
			s.Pause()
		}
	}
	if inpututil.IsKeyJustPressed(input.KeyEscape) {
		return s.ScenePlay.Finish()
	}

	if s.MusicVolumeKeyHandler.Update() {
		s.SetMusicVolume(s.cfg.MusicVolume)
	}
	s.SoundVolumeKeyHandler.Update()
	s.BackgroundBrightnessKeyHandler.Update()
	if s.OffsetKeyHandler.Update() {
		s.SetOffset(s.cfg.Offset)
	}
	s.DebugPrintKeyHandler.Update()
	if s.SpeedScaleKeyHandlers[s.mode].Update() {
		s.SetSpeedScale()
	}

	return s.ScenePlay.Update()
}

func (s Scene) Draw(screen draws.Image) {
	s.drawBackground(screen)
	s.ScenePlay.Draw(screen)
	if s.cfg.DebugPrint {
		s.ScenePlay.DebugPrint(screen)
	}
}
