package choose

import (
	"fmt"
	"io/fs"
	"runtime/debug"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/draws"
	"github.com/hndada/gosu/format/osr"
	"github.com/hndada/gosu/input"
	"github.com/hndada/gosu/mode"
)

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
	Mode    int
	SubMode int
	Sort    int

	volumeMusic          *float64
	volumeSound          *float64
	offset               *int64
	backgroundBrightness *float64
	speedFactors         []*float64
	exposureTimes        []func(float64) float64

	choose     audios.Sound
	Music      audios.MusicPlayer
	Background mode.BackgroundDrawer
	Panel      Panel
	List
}

func NewScene() *Scene {
	ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	debug.SetGCPercent(100)
	ebiten.SetWindowTitle("gosu")
	return &Scene{}
}
func isEnter() bool {
	return ebiten.IsKeyPressed(input.KeyEnter) ||
		ebiten.IsKeyPressed(input.KeyNumpadEnter)
}
func (s *Scene) Update() any {
	if ModeKeyHandler.Update() {
		s.UpdateMode()
	}
	if s.CursorKeyHandler.Update() {
		s.UpdateBackground()
	}
	s.Background.Sprite = mode.NewBackground()
	if isEnter() {
		s.choose.Play(*s.volumeSound)
		var c Chart
		fs, name, err := c.Select()
		if err != nil {
			return err
		}
		return Return{
			FS:     fs,
			Name:   name,
			Mode:   s.Mode,
			Mods:   nil,
			Replay: nil,
		}
	}
	return nil
}
func (s Scene) Draw(screen draws.Image) {
	s.Background.Draw(screen)
	s.Panel.Draw(screen)
	s.List.Draw(screen)
	s.DebugPrint(screen)
}
func (s Scene) DebugPrint(screen draws.Image) {
	speed := *s.speedFactors[s.Mode]
	ebitenutil.DebugPrint(screen.Image,
		fmt.Sprintf(
			"Mode (F1): %s\n"+
				"Sub mode (F2): %s\n"+
				"Sort (F3): %s\n"+
				"\n"+
				"Music volume (Alt+ Left/Right): %.0f%%\n"+
				"Sound volume (Ctrl+ Left/Right): %.0f%%\n"+
				"Offset (Shift+ Left/Right): %dms\n"+
				"Brightness (Ctrl+ O/P): %.0f%%\n"+
				"\n"+
				"Speed (PageUp/Down): %.0f (Exposure time: %.0fms)\n"+
				[]string{"Piano", "Drum"}[s.Mode],
			fmt.Sprintf("%d Key", s.SubMode),
			[]string{"by name", "by level"}[s.Sort],

			*s.volumeMusic*100,
			*s.volumeSound*100,
			*s.offset,
			*s.backgroundBrightness*100,

			speed*100, s.exposureTimes[s.Mode](speed)))
}

type Return struct {
	FS     fs.FS
	Name   string
	Mode   int
	Mods   interface{}
	Replay *osr.Format
}
