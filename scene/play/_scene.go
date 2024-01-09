package play

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/game/piano"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/scene"
)

// ScenePlay: struct, PlayScene: function
// Interface declares at 'user' package.
type Scene struct {
	*scene.Config
	*scene.Asset

	mode.ScenePlay
	drawBackground func(draws.Image)

	// KeyHandleXxx looks slightly better than HandleXxxByKey.
	KeyHandleMusicVolume          func() bool
	KeyHandleSoundVolume          func() bool
	KeyHandleMusicOffset          func() bool
	KeyHandleBackgroundBrightness func() bool
	KeyHandleDebugPrint           func() bool
	KeyHandleSpeedScale           func() bool
}

func NewScene(cfg *scene.Config, asset *scene.Asset, fsys fs.FS, name string, replay *osr.Format) (s *Scene, err error) {
	s = &Scene{Config: cfg, Asset: asset}

	switch s.Mode {
	case mode.ModePiano:
		s.ScenePlay, err = piano.NewScenePlay(s.PianoConfig, s.PianoAssets, fsys, name, replay)
	}
	if err != nil {
		return
	}
	ch := s.ScenePlay.ChartHeader()
	s.drawBackground = scene.NewBackgroundDrawer(s.Config, s.Asset, fsys, ch.BackgroundFilename)

	s.KeyHandleMusicVolume = scene.NewMusicVolumeKeyHandler(s.Config, s.Asset)
	s.KeyHandleSoundVolume = scene.NewSoundVolumeKeyHandler(s.Config, s.Asset)
	s.KeyHandleMusicOffset = scene.NewMusicOffsetKeyHandler(s.Config, s.Asset)
	s.KeyHandleBackgroundBrightness = scene.NewBackgroundBrightnessKeyHandler(s.Config, s.Asset)
	s.KeyHandleDebugPrint = scene.NewDebugPrintKeyHandler(s.Config, s.Asset)
	s.KeyHandleSpeedScale = scene.NewSpeedScaleKeyHandler(s.Config, s.Asset, s.Mode)

	ebiten.SetWindowTitle(s.WindowTitle())
	return
}

// Changing speed might not be applied after positions are calculated.
// But this is not tested.
func (s *Scene) Update() any {
	if s.KeyHandleMusicVolume() {
		s.SetMusicVolume(s.MusicVolume)
	}
	s.KeyHandleSoundVolume()
	if s.KeyHandleMusicOffset() {
		s.SetMusicOffset(s.MusicOffset)
	}
	s.KeyHandleBackgroundBrightness()
	s.KeyHandleDebugPrint()
	if s.KeyHandleSpeedScale() {
		s.SetSpeedScale()
	}

	if input.IsKeyJustPressed(input.KeyTab) {
		if s.IsPaused() {
			s.Resume()
		} else {
			s.Pause()
		}
	}
	if input.IsKeyJustPressed(input.KeyEscape) {
		return s.ScenePlay.Finish()
		// os.Exit(1)
	}

	return s.ScenePlay.Update()
}

func (s Scene) Draw(screen draws.Image) {
	s.drawBackground(screen)
	s.ScenePlay.Draw(screen)
	if s.DebugPrint {
		f := fmt.Fprintf
		var b strings.Builder
		f(&b, s.Config.DebugString())
		f(&b, "\n")
		f(&b, s.ScenePlay.DebugString()) // interpolated ScenePlay debug string
		f(&b, "\n")
		f(&b, "Press TAB to pause.\n")
		f(&b, "Press ESC to back to choose a song.\n")
		ebitenutil.DebugPrint(screen.Image, b.String())
	}
}
