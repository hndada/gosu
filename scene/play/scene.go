package play

import (
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/scene"
)

// ScenePlay: struct, PlayScene: function
// Interface declares at 'user' package.
type Scene struct {
	cfg   *scene.Config
	asset *scene.Asset
	*scene.BaseScene

	mode int
	mode.ScenePlay

	drawBackground func(draws.Image)
}

func NewScene(cfg *scene.Config, asset *scene.Asset, fsys fs.FS, name string, _mode int, mods any, replay mode.Replay) (s *Scene, err error) {
	s = &Scene{
		cfg:       cfg,
		asset:     asset,
		BaseScene: scene.TheBaseScene,
	}

	switch _mode {
	case mode.ModePiano:
		s.mode = _mode
		mods := mods.(piano.Mods)
		s.ScenePlay, err = piano.NewScenePlay(cfg.PianoConfig, asset.PianoAssets, fsys, name, mods, replay)
	case mode.ModeDrum:
	}

	ch := s.ScenePlay.ChartHeader()
	bgFilename := ch.BackgroundFilename
	bgSprite := scene.NewBackgroundSprite(fsys, bgFilename, cfg.ScreenSize)
	if bgSprite.IsEmpty() {
		bgSprite = asset.DefaultBackgroundSprite
	}
	s.drawBackground = scene.NewBackgroundDrawer(bgSprite, cfg.ScreenSize, &cfg.BackgroundBrightness)

	ebiten.SetWindowTitle(s.WindowTitle())
	// debug.SetGCPercent(0)
	return
}

// The order of function calls may not consistent with
// the order of methods of mode.ScenePlay.

// Changed speed might not be applied after positions are calculated.
// But this is not tested.
func (s *Scene) Update() any {
	// set
	if s.MusicVolumeKeyHandler.Update() {
		s.SetMusicVolume(s.cfg.MusicVolume)
	}
	s.SoundVolumeKeyHandler.Update()
	if s.SpeedScaleKeyHandlers[s.mode].Update() {
		s.SetSpeedScale()
	}
	if s.MusicOffsetKeyHandler.Update() {
		s.SetMusicOffset(s.cfg.MusicOffset)
	}

	// life cycle
	args := s.ScenePlay.Update()

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

	// draw
	s.BackgroundBrightnessKeyHandler.Update()
	s.DebugPrintKeyHandler.Update()

	return args
}

func (s Scene) Draw(screen draws.Image) {
	s.drawBackground(screen)
	s.ScenePlay.Draw(screen)
	if s.cfg.DebugPrint {
		s.ScenePlay.DebugPrint(screen)
	}
}
