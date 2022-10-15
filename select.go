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

// SceneSelect might be created after one play at multiplayer.
// Todo: fetch Score with Replay
// Todo: preview music. Start at PreviewTime, keeps playing until end.
type SceneSelect struct {
	// Query     string
	View             []ChartInfo // Todo: ChartInfo -> *ChartInfo?
	Cursor           int
	CursorKeyHandler ctrl.KeyHandler
	// board       draws.Box

	Background  draws.Sprite  // Todo: BackgroundDrawer with some effects
	MusicPlayer *audio.Player // Todo: Rewind after preview has finished.
	MusicCloser io.Closer
}

func NewSceneSelect() *SceneSelect {
	s := &SceneSelect{}
	s.UpdateMode()
	return s
}
func (s *SceneSelect) Update() any {
	if set := ModeKeyHandler.Update(); set {
		s.UpdateMode()
	}
	if set := s.CursorKeyHandler.Update(); set {
		s.UpdateBackground()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyNumpadEnter) {
		// Sounds.Play("restart")
		prop := modeProps[currentMode]
		info := prop.ChartInfos[s.Cursor]
		return SelectToPlayArgs{
			Path:   info.Path,
			Replay: nil, // replay,
		}
	}
	return nil
}

func (s *SceneSelect) UpdateMode() {
	SpeedScaleKeyHandler.Handler = speedScaleHandlers[currentMode]
	s.View = modeProps[currentMode].ChartInfos
	s.Cursor = 0
	s.CursorKeyHandler = NewCursorKeyHandler(&s.Cursor, len(s.View))
	s.UpdateBackground()
}
func NewCursorKeyHandler(cursor *int, len int) ctrl.KeyHandler {
	return ctrl.KeyHandler{
		Handler: &ctrl.IntHandler{
			Value: cursor,
			Min:   0,
			Max:   len,
			Loop:  true,
		},
		Modifiers: []ebiten.Key{},
		Keys:      [2]ebiten.Key{ebiten.KeyUp, ebiten.KeyDown},
		Sounds:    [2][]byte{SwipeSound, SwipeSound},
	}
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
	// s.View[cursor].NewChartBoard().Draw(screen, ebiten.DrawImageOptions{}, draws.Point{})
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
	prop := modeProps[currentMode]
	speed := *prop.SpeedScale
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"Music volume (Q/W): %.0f%%\n"+
				"Effect volume (A/S): %.0f%%\n"+
				"Speed (Z/X): %.0f (Exposure time: %.0fms)\n\n"+
				"Mode (Space): %s\n",
			MusicVolume*100, EffectVolume*100,
			speed*100, prop.ExposureTime(speed),
			prop.Name))
	// "Music volume (Alt+↑/↓): %.0f%%\n"+"Effect volume (Ctrl+↑/↓): %.0f%%\n"+
	// "Speed (Ctrl+PageUp/PageDown): %.0f\n"+"(Exposure time: %.0fms)\n\n"+
	// "Mode (Ctrl+Alt+Shift+←/→): %s\n",
}
