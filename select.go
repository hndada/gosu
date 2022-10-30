package gosu

import (
	"fmt"
	"image/color"
	"io"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hndada/gosu/audios"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/format/osr"
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

	BackgroundDrawer BackgroundDrawer
	MusicPlayer      *audio.Player // Todo: Rewind after preview has finished.
	MusicCloser      io.Closer
}

func NewSceneSelect() *SceneSelect {
	s := &SceneSelect{}
	s.BackgroundDrawer.Brightness = &BackgroundBrightness
	s.UpdateMode()
	return s
}
func (s *SceneSelect) Update() any {
	if set := ModeKeyHandler.Update() || SortKeyHandler.Update(); set {
		s.UpdateMode()
	}
	if set := s.CursorKeyHandler.Update() || BrightKeyHandler.Update(); set {
		s.UpdateBackground()
	}
	if set := SizeKeyHandler.Update(); set {
		s.UpdateWindowSize()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyNumpadEnter) {
		audios.PlayEffect(SelectSound, EffectVolume)
		prop := modeProps[currentMode]
		info := prop.ChartInfos[s.Cursor]
		var replay *osr.Format
		// b, err := os.ReadFile("replay/MuangMuangE - Doubutsu Biscuits x PPP - Youkoso JAPARI PARK e (TV size ver.) [Muzukashii] (2019-05-15) Taiko.osr")
		// if err != nil {
		// 	panic(err)
		// }
		// replay, err := osr.Parse(b)
		// if err != nil {
		// 	panic(err)
		// }
		return SelectToPlayArgs{
			Path:   info.Path,
			Replay: replay,
		}
	}
	return nil
}

func (s *SceneSelect) UpdateMode() {
	SpeedScaleKeyHandler.Handler = speedScaleHandlers[currentMode]
	s.View = modeProps[currentMode].ChartInfos
	switch currentSort {
	case SortByName:
		sort.Slice(s.View, func(i, j int) bool {
			if s.View[i].MusicName == s.View[j].MusicName {
				return s.View[i].Level < s.View[j].Level
			}
			return s.View[i].MusicName < s.View[j].MusicName
		})
	case SortByLevel:
		sort.Slice(s.View, func(i, j int) bool {
			if s.View[i].Level == s.View[j].Level {
				return s.View[i].MusicName < s.View[j].MusicName
			}
			return s.View[i].Level < s.View[j].Level
		})
	}
	s.Cursor = 0
	s.CursorKeyHandler = NewCursorKeyHandler(&s.Cursor, len(s.View))
	s.UpdateBackground()
}
func (s *SceneSelect) UpdateWindowSize() {
	ebiten.SetWindowSize(WindowSizeX[currentSize], WindowSizeY[currentSize])
}
func NewCursorKeyHandler(cursor *int, len int) ctrl.KeyHandler {
	return ctrl.KeyHandler{
		Handler: &ctrl.IntHandler{
			Value: cursor,
			Min:   0,
			Max:   len - 1,
			Loop:  true,
		},
		Modifiers: []ebiten.Key{},
		Keys:      [2]ebiten.Key{ebiten.KeyUp, ebiten.KeyDown},
		Sounds:    [2][]byte{SwipeSound, SwipeSound},
		Volume:    &EffectVolume,
	}
}
func (s *SceneSelect) UpdateBackground() {
	s.BackgroundDrawer.Sprite = DefaultBackground
	if len(s.View) == 0 {
		return
	}
	if s.Cursor >= len(s.View) {
		return
	}
	info := s.View[s.Cursor]
	sprite := NewBackground(info.BackgroundPath())
	if sprite.IsValid() {
		s.BackgroundDrawer.Sprite = sprite
	}
}

// Currently Chart infos are not in loop.
// May add extra effect to box arrangement. e.g., x -= y / 5
func (s SceneSelect) Draw(screen *ebiten.Image) {
	s.BackgroundDrawer.Draw(screen)
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
	sizeX, sizeY := ebiten.WindowSize()
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf(
			"Mode (F1): %s\n"+
				"Sort (F2): %s\n"+
				"Window size (F3): %dx%d\n"+
				"\n"+
				"Music volume (Alt+ Left/Right): %.0f%%\n"+
				"Effect volume (Ctrl+ Left/Right): %.0f%%\n"+
				"Brightness (Ctrl+ O/P): %.0f%%\n"+
				"\n"+
				"Speed (PageUp/Down): %.0f (Exposure time: %.0fms)\n"+
				"Offset (Shift+ Left/Right): %dms\n",
			prop.Name,
			[]string{"by name", "by level"}[currentSort],
			sizeX,
			sizeY,

			MusicVolume*100,
			EffectVolume*100,
			BackgroundBrightness*100,

			speed*100, prop.ExposureTime(speed),
			Offset))
}
