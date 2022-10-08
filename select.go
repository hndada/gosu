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

// var SelectedChart ChartInfo
// var (
//
//	Query string
//	View  []ChartInfo
//
// )
// Currently SceneSelect's Cursor has reset every play has finished.
type SceneSelect struct {
	// ModeProps []ModeProp
	// Mode      *int
	// ViewQuery     string
	View []ChartInfo // Todo: ChartInfo -> *ChartInfo?
	// Cursor        *int
	Cursor int
	// CursorHandler ctrl.IntHandler
	CursorKeyHandler ctrl.KeyHandler
	// ModeHandler   ctrl.IntHandler

	Background  draws.Sprite  // Todo: BackgroundDrawer with some effects
	MusicPlayer *audio.Player // Todo: Rewind after preview has finished.
	MusicCloser io.Closer
	// MultiMode   int
}

// const (
// 	MultiModeNone = iota
// 	MultiMode
// )

// Todo: Score / Replay fetch
// Todo: preview music. Start at PreviewTime, keeps playing until end.
// func NewSceneSelect(modes []ModeProp, mode *int) *SceneSelect {
// SceneSelect should be created every play has done considering multiplay.
func NewSceneSelect() *SceneSelect {
	s := new(SceneSelect)
	// s.ModeHandler = NewModeHandler(s.Mode, len(modes))
	s.UpdateMode()
	// ebiten.SetWindowTitle("gosu")
	return s
}
func (s *SceneSelect) Update() any {
	// Todo: refactor it
	// if act := VsyncSwitchHandler.Update(); act {
	// 	if *VsyncSwitchHandler.Target {
	// 		ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	// 	} else {
	// 		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	// 	}
	// }
	if act := ModeKeyHandler.Update(); act {
		s.UpdateMode()
	}
	if act := s.CursorKeyHandler.Update(); act {
		s.UpdateBackground()
	}
	prop := ModeProps[CurrentMode]
	prop.SpeedKeyHandler.Update()
	// s.ModeProps[*s.Mode].SpeedHandler.Update()
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyNumpadEnter) {
		// Sounds.Play("restart")
		// cursor := s.ModeProps[*s.Mode].Cursor
		// info := s.ModeProps[*s.Mode].ChartInfos[cursor]
		info := prop.ChartInfos[s.Cursor]

		// b, err := os.ReadFile("replay/MuangMuangE - cillia - Ringo Uri no Utakata Shoujo [Ringo Oni] (2019-06-14) Taiko-1.osr")
		// if err != nil {
		// 	panic(err)
		// }
		// replay, err := osr.Parse(b)
		// if err != nil {
		// 	panic(err)
		// }
		return SelectToPlayArgs{
			Path: info.Path,
			// Mode: *s.Mode, // Todo: duplicated. Should it be removed?
			// Replay:       replay,
			Replay: nil,
			// SpeedHandler: s.ModeProps[*s.Mode].SpeedHandler,
		}
	}
	return nil
}

func NewCursorKeyHandler(cursor *int, len int) ctrl.KeyHandler {
	h := ctrl.KeyHandler{
		Handler: &ctrl.IntHandler{
			Value:  cursor,
			Unit:   1,
			Min:    0,
			Max:    len,
			Loop:   true,
			Sounds: [2][]byte{SwipeSound, SwipeSound},
		},
	}
	h.SetKeys([]ebiten.Key{}, [2]ebiten.Key{ebiten.KeyDown, ebiten.KeyUp})
	return h
}

//	func (s *SceneSelect) UpdateMode() {
//		s.View = s.ModeProps[*s.Mode].ChartInfos
//		s.Cursor = &s.ModeProps[*s.Mode].Cursor
//		s.CursorHandler = NewCursorHandler(&s.ModeProps[*s.Mode].Cursor, len(s.View))
//		s.UpdateBackground()
//	}
func (s *SceneSelect) UpdateMode() {
	s.View = ModeProps[CurrentMode].ChartInfos
	// s.Cursor = &ModeProps[CurrentMode].Cursor
	s.CursorKeyHandler = NewCursorKeyHandler(&s.Cursor, len(s.View))
	s.UpdateBackground()
}

func (s *SceneSelect) UpdateBackground() {
	s.Background = DefaultBackground
	if len(s.View) == 0 {
		return
	}
	if s.Cursor >= len(s.View) {
		return
	}
	info := s.View[s.Cursor]
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
		text.Draw(screen, t, Face12, x, y+int(offset), color.Black)
		// text.Draw(screen, t, basicfont.Face7x13, x, y+int(offset), color.Black)
	}
	s.DebugPrint(screen)
}
func (s SceneSelect) Viewport() ([]ChartInfo, int) {
	count := chartItemBoxCount
	var viewport []ChartInfo
	var cursor int
	if s.Cursor <= count/2 {
		viewport = append(viewport, s.View[0:s.Cursor]...)
		cursor = s.Cursor
	} else {
		bound := s.Cursor - count/2
		viewport = append(viewport, s.View[bound:s.Cursor]...)
		cursor = count / 2
	}
	if s.Cursor >= len(s.View)-count/2 {
		viewport = append(viewport, s.View[s.Cursor:len(s.View)]...)
	} else {
		bound := s.Cursor + count/2
		viewport = append(viewport, s.View[s.Cursor:bound]...)
	}
	return viewport, cursor
}
func (s SceneSelect) DebugPrint(screen *ebiten.Image) {
	// mode := s.ModeProps[*s.Mode]
	prop := ModeProps[CurrentMode]
	speed := *prop.SpeedScale
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Music volume (Alt+↑/↓): %.0f%%\n"+"Effect volume (Ctrl+↑/↓): %.0f%%\n"+
			// "Vsync enabled (Press 5): %v\n"+
			"Speed (Ctrl+PageUp/PageDown): %.0f\n"+"(Exposure time: %.0fms)\n\n"+
			"Mode (Ctrl+Alt+Shift+←/→): %s\n",
			// "Chart info index: %d\n",
			MusicVolume*100, EffectVolume*100,
			// VsyncSwitch,
			speed*100, prop.ExposureTime(speed),
			prop.Name))
	// s.Cursor))
}
