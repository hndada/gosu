package play

import (
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	*scene.Config
	*scene.Asset
	*scene.BaseScene

	mode int
	mode.ScenePlay

	drawBackground func(draws.Image)
}

func NewScene(cfg *scene.Config, asset *scene.Asset, fsys fs.FS, name string, _mode int, mods any, replay mode.Replay) (s *Scene, err error) {
	s = &Scene{
		Config:    cfg,
		Asset:     asset,
		BaseScene: scene.TheBaseScene,
	}

	switch _mode {
	case mode.ModePiano:
		s.mode = _mode
		mods := mods.(piano.Mods)
		s.ScenePlay, err = piano.NewScenePlay(s.PianoConfig, s.PianoAssets, fsys, name, mods, replay)
	case mode.ModeDrum:
	}

	ch := s.ScenePlay.ChartHeader()
	bgFilename := ch.BackgroundFilename
	bgSprite := scene.NewBackgroundSprite(fsys, bgFilename, s.ScreenSize)
	if bgSprite.IsEmpty() {
		bgSprite = asset.DefaultBackgroundSprite
	}
	s.drawBackground = scene.NewBackgroundDrawer(bgSprite, s.ScreenSize, &s.BackgroundBrightness)

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
		s.SetMusicVolume(s.MusicVolume)
	}
	s.SoundVolumeKeyHandler.Update()
	if s.SpeedScaleKeyHandlers[s.mode].Update() {
		s.SetSpeedScale()
	}
	if s.MusicOffsetKeyHandler.Update() {
		s.SetMusicOffset(s.MusicOffset)
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
		// return s.ScenePlay.Finish()
		os.Exit(1)
	}

	// draw
	s.BackgroundBrightnessKeyHandler.Update()
	s.DebugPrintKeyHandler.Update()

	return args
}

func (s Scene) Draw(screen draws.Image) {
	s.drawBackground(screen)
	s.ScenePlay.Draw(screen)
	if s.DebugPrint {
		var b strings.Builder
		f := fmt.Fprintf

		f(&b, "FPS: %.2f\n", ebiten.ActualFPS())
		f(&b, "TPS: %.2f\n", ebiten.ActualTPS())
		f(&b, "%s\n", s.ScenePlay.DebugString()) // interpolated ScenePlay debug string
		f(&b, "Music volume (Ctrl+ Left/Right): %.0f%%\n", s.MusicVolume*100)
		f(&b, "Sound volume (Alt+ Left/Right): %.0f%%\n", s.SoundVolume*100)
		f(&b, "MusicOffset (Shift+ Left/Right): %dms\n", s.MusicOffset)
		f(&b, "\n")
		f(&b, "Press ESC to back to choose a song.\n")
		f(&b, "Press TAB to pause.\n")
		f(&b, "Press Ctrl+ O/P to change background brightness\n")
		f(&b, "Press F12 to print debug.\n")
		ebitenutil.DebugPrint(screen.Image, b.String())
	}
}
