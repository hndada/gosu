package gosu

import (
	"fmt"
	"image/color"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/draws"
)

type SceneSelect struct {
	ModeProps []ModeProp
	Mode      *int
	// ViewQuery     string
	View          []ChartInfo // Todo: ChartInfo -> *ChartInfo?
	Cursor        *int
	ModeHandler   ctrl.IntHandler
	CursorHandler ctrl.IntHandler

	Background  draws.Sprite  // Todo: BackgroundDrawer with some effects
	MusicPlayer *audio.Player // Todo: Rewind after preview has finished.
	MusicCloser io.Closer
}

var count = int(screenSizeY/ChartInfoBoxHeight/2*2) + 2 // Gives count some margin

// Todo: Score / Replay fetch
// Todo: preview music
func NewSceneSelect(modes []ModeProp, mode *int) *SceneSelect {
	s := new(SceneSelect)
	s.ModeProps = modes
	s.Mode = mode
	s.ModeHandler = NewModeHandler(s.Mode, len(modes))
	s.UpdateMode()
	ebiten.SetWindowTitle("gosu")
	return s
}
func (s *SceneSelect) Update() any {
	// Todo: refactor it
	if fired := VsyncSwitchHandler.Update(); fired {
		if *VsyncSwitchHandler.Target {
			ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
		} else {
			ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
		}
	}
	if moved := s.CursorHandler.Update(); moved {
		s.UpdateBackground()
	}
	s.ModeProps[*s.Mode].SpeedHandler.Update()
	if fired := s.ModeHandler.Update(); fired {
		s.UpdateMode()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyNumpadEnter) {
		Sounds.Play("restart")
		cursor := s.ModeProps[*s.Mode].Cursor
		info := s.ModeProps[*s.Mode].ChartInfos[cursor]
		return SelectToPlayArgs{
			Path:         info.Path,
			Mode:         *s.Mode, // Todo: duplicated. Should it be removed?
			Replay:       nil,
			SpeedHandler: s.ModeProps[*s.Mode].SpeedHandler,
		}
	}
	return nil
}
func (s *SceneSelect) UpdateMode() {
	s.View = s.ModeProps[*s.Mode].ChartInfos
	s.Cursor = &s.ModeProps[*s.Mode].Cursor
	s.CursorHandler = NewCursorHandler(&s.ModeProps[*s.Mode].Cursor, len(s.View))
	s.UpdateBackground()
}

func (s *SceneSelect) UpdateBackground() {
	s.Background = DefaultBackground
	if len(s.View) == 0 {
		return
	}
	if *s.Cursor >= len(s.View) {
		return
	}
	info := s.View[*s.Cursor]
	sprite := NewBackground(info.BackgroundPath())
	if sprite.IsValid() {
		s.Background = sprite
	}
}

// Currently Chart infos are not in loop.
// May add extra effect to box arrangement. e.g., x -= y / 5
func (s SceneSelect) Draw(screen *ebiten.Image) {
	s.Background.Draw(screen, nil)
	viewport, cursor := s.Viewport()
	for i := range viewport {
		sprite := ChartItemBoxSprite
		var tx float64
		if i == cursor {
			tx -= chartInfoBoxshrink
		}
		ty := float64(i-cursor) * ChartInfoBoxHeight
		sprite.Move(tx, ty)
		sprite.Draw(screen, nil)
	}

	const (
		dx = 20 // Padding left.
		dy = 30 // Padding bottom.
	)
	for i, info := range viewport {
		sprite := ChartItemBoxSprite
		t := info.Text()
		offset := float64(i-cursor) * ChartInfoBoxHeight
		// rect := text.BoundString(draws.Face24, t)
		x := int(sprite.X()-sprite.W()) + dx   //+ rect.Dx()
		y := int(sprite.Y()-sprite.H()/2) + dy //+ rect.Dy()
		if i == cursor {
			x -= int(chartInfoBoxshrink)
		}
		text.Draw(screen, t, draws.Face20, x, y+int(offset), color.Black)
		// text.Draw(screen, t, basicfont.Face7x13, x, y+int(offset), color.Black)
	}
	s.DebugPrint(screen)
}
func (s SceneSelect) Viewport() ([]ChartInfo, int) {
	var viewport []ChartInfo
	var cursor int
	if *s.Cursor <= count/2 {
		viewport = append(viewport, s.View[0:*s.Cursor]...)
		cursor = *s.Cursor
	} else {
		bound := *s.Cursor - count/2
		viewport = append(viewport, s.View[bound:*s.Cursor]...)
		cursor = count / 2
	}
	if *s.Cursor >= len(s.View)-count/2 {
		viewport = append(viewport, s.View[*s.Cursor:len(s.View)]...)
	} else {
		bound := *s.Cursor + count/2
		viewport = append(viewport, s.View[*s.Cursor:bound]...)
	}
	return viewport, cursor
}
func (s SceneSelect) DebugPrint(screen *ebiten.Image) {
	mode := s.ModeProps[*s.Mode]
	speed := *mode.SpeedHandler.Target
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Music volume (Press 1/2): %.0f%%\n"+"Effect volume (Press 3/4): %.0f%%\n"+
			"VsyncOn (Press 5): %v\n"+
			"SpeedScale (Press 8/9): %.0f\n"+"(Exposure time: %.0fms)\n\n"+
			"Mode (Press 0): %s\n"+
			"Chart info index: %d\n",
			MusicVolume*100, EffectVolume*100,
			VsyncSwitch,
			speed*100, mode.ExposureTime(speed),
			ModeNames[*s.Mode],
			*s.Cursor))
}
