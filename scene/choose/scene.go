package choose

import (
	"fmt"
	"io/fs"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/drum"
	"github.com/hndada/gosu/mode/piano"
	scene "github.com/hndada/gosu/scene"
)

// s.UpdateBackground()
const (
	TPS         = scene.TPS
	ScreenSizeX = scene.ScreenSizeX
	ScreenSizeY = scene.ScreenSizeY
)

// Background brightness at Song select: 60% (153 / 255), confirmed.
// Score box color: Gray128 with 50% transparent
// Hovered Score box color: Gray96 with 50% transparent

// View             []ChartInfo // Todo: ChartInfo -> *ChartInfo?

// Todo: fetch Score with Replay
// Todo: preview music. Start at PreviewTime, keeps playing until end.
// Todo: Rewind after preview has finished.
// Group1, Group2, Sort, Filter int
type Scene struct {
	volumeMusic   *float64
	volumeSound   *float64
	brightness    *float64
	offset        *int64
	speedFactors  []*float64
	exposureTimes []func(float64) float64

	choose     audios.Sound
	Music      audios.MusicPlayer
	Background mode.BackgroundDrawer

	mode        int
	subMode     int
	Mode        ctrl.KeyHandler
	SubMode     ctrl.KeyHandler
	TypeWriter  TypeWriter // Query here
	SetSelected bool
	// ChartSetPanel *ChartSetPanel
	ChartSetList List
	// ChartPanel    *ChartPanel
	ChartList List
}

func NewScene() *Scene {
	s := &Scene{}
	s.volumeMusic = &mode.S.VolumeMusic
	s.volumeSound = &mode.S.VolumeSound
	s.brightness = &mode.S.BackgroundBrightness
	s.offset = &mode.S.Offset
	s.speedFactors = []*float64{
		&piano.S.SpeedScale, &drum.S.SpeedScale}
	s.exposureTimes = []func(float64) float64{
		piano.ExposureTime, drum.ExposureTime}
	s.Mode = ctrl.KeyHandler{
		Handler: ctrl.IntHandler{
			Value: &s.mode,
			Min:   0,
			Max:   2 - 1, // There are two modes.
			Loop:  true,
		},
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{-1, input.KeyF1},
		Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
		Volume:    &mode.S.VolumeSound,
	}
	s.SubMode = ctrl.KeyHandler{
		Handler: ctrl.IntHandler{
			Value: &s.subMode,
			Min:   4,
			Max:   9,
			Loop:  true,
		},
		Modifiers: []input.Key{},
		Keys:      [2]input.Key{input.KeyF2, input.KeyF3},
		Sounds:    [2]audios.Sounder{scene.UserSkin.Swipe, scene.UserSkin.Swipe},
		Volume:    &mode.S.VolumeSound,
	}
	s.Background.Sprite = mode.NewBackground()
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	debug.SetGCPercent(100)
	ebiten.SetWindowTitle("gosu")
	return &Scene{}
}
func isEnter() bool {
	return ebiten.IsKeyPressed(input.KeyEnter) ||
		ebiten.IsKeyPressed(input.KeyNumpadEnter)
}
func isBack() bool {
	return ebiten.IsKeyPressed(input.KeyEscape)
}
func (s *Scene) Update() any {
	scene.VolumeMusic.Update()
	scene.VolumeSound.Update()

	scene.Offset.Update()
	scene.SpeedScales[s.mode].Update()

	if isEnter() {
		if s.SetSelected {
			s.choose.Play(*s.volumeSound)
			var c Chart
			fs, name, err := c.Select()
			if err != nil {
				return err
			}
			return Return{
				FS:     fs,
				Name:   name,
				Mode:   s.mode,
				Mods:   nil,
				Replay: nil,
			}
		} else {
			s.SetSelected = true
		}
	}
	if isBack() {
		if s.SetSelected {
			s.SetSelected = false
		} else {
			s.Query = ""
		}
	}
	if s.Mode.Update() || s.SubMode.Update() {
		s.UpdateMode()
	}
	if s.SetSelected {
		s.ChartList.Update()
	} else {
		s.ChartSetList.Update()
	}
	return nil
}
func (s Scene) Draw(screen draws.Image) {
	s.Background.Draw(screen)
	if s.SetSelected {
		s.ChartPanel.Draw(screen)
		s.ChartList.Draw(screen)
	} else {
		s.ChartSetPanel.Draw(screen)
		s.ChartSetList.Draw(screen)
	}
	s.DebugPrint(screen)
}
func (s Scene) DebugPrint(screen draws.Image) {
	speed := *s.speedFactors[s.mode]
	ebitenutil.DebugPrint(screen.Image,
		fmt.Sprintf(
			"Mode (F1): %s\n"+
				"Sub mode (F2/F3): %s\n"+
				"\n"+
				"Music volume (Alt+ Left/Right): %.0f%%\n"+
				"Sound volume (Ctrl+ Left/Right): %.0f%%\n"+
				"Offset (Shift+ Left/Right): %dms\n"+
				"Brightness (Ctrl+ O/P): %.0f%%\n"+
				"\n"+
				"Speed (PageUp/Down): %.0f (Exposure time: %.0fms)\n"+
				[]string{"Piano", "Drum"}[s.mode],
			fmt.Sprintf("%d Key", s.SubMode),

			*s.volumeMusic*100,
			*s.volumeSound*100,
			*s.offset,
			*s.brightness*100,

			speed*100, s.exposureTimes[s.mode](speed)))
}

type Return struct {
	FS     fs.FS
	Name   string
	Mode   int
	Mods   interface{}
	Replay *osr.Format
}
