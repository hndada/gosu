package mode

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
	"golang.org/x/image/font/basicfont"
)

type SceneSelect struct {
	Modes     []Mode
	*ModeType // Todo: IntHandler

	View []ChartInfo // Todo: ChartInfo -> *ChartInfo?
	// ViewMode int
	Cursor        int // Todo: Delayed at Cursor
	SelectHandler ctrl.IntHandler

	Background  draws.Sprite  // Todo: BackgroundDrawer with some effects
	MusicPlayer *audio.Player // Todo: Rewind after preview has finished.
	MusicCloser io.Closer
}

var count = int(screenSizeY/ChartInfoBoxHeight/2*2) + 2 // Gives count some margin

// Todo: Score / Replay fetch
// Todo: preview music
func NewSceneSelect(modes []Mode, modeType *ModeType) *SceneSelect {
	s := new(SceneSelect)
	s.Modes = modes
	s.ModeType = modeType

	s.View = modes[*modeType].ChartInfos
	s.SelectHandler = NewSelectHandler(&s.Cursor, len(s.View))
	s.UpdateBackground()
	return s
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
	img := draws.NewImage(info.Header.BackgroundPath(info.Path))
	if img != nil {
		s.Background.I = img
	}
	s.Background.SetWidth(screenSizeX)
	s.Background.SetCenterY(screenSizeY / 2)
	ebiten.SetWindowTitle("gosu")
}

// Default HoldKey value is 0, which is Key0.
func (s *SceneSelect) Update() any {
	moved := s.SelectHandler.Update()
	if moved {
		s.UpdateBackground()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyNumpadEnter) {
		Sounds.Play("restart")
		info := s.View[s.Cursor]
		return SelectToPlayArgs{
			Path:     info.Path,
			ModeType: *s.ModeType, // Todo: duplicated. Should it be removed?
			Replay:   nil,
		}
	}
	return nil
}

// Currently Chart infos are not in loop.
// May add extra effect to box arrangement.
// x -= y / 5, for example.
func (s SceneSelect) Draw(screen *ebiten.Image) {
	const (
		dx = 20
		dy = 30
	)
	s.Background.Draw(screen)
	viewport, cursor := s.Viewport()
	for i, info := range viewport {
		sprite := ChartInfoBoxSprite
		if i == cursor {
			sprite.X -= chartInfoBoxshrink
		}
		offset := i - cursor
		sprite.SetCenterY(screenSizeY/2 + float64(offset)*ChartInfoBoxHeight)
		sprite.Draw(screen)

		t := info.Text()
		// rect := text.BoundString(basicfont.Face7x13, t)
		x := int(sprite.X) + dx //+ rect.Dx()
		y := int(sprite.Y) + dy //+ rect.Dy()
		text.Draw(screen, t, basicfont.Face7x13, x, y, color.Black)
	}
	// Code of drawing cursor
	// {
	// 	sprite := GeneralSkin.CursorSprites[0]
	// 	x, y := ebiten.CursorPosition()
	// 	sprite.X, sprite.Y = float64(x), float64(y)
	// 	sprite.Draw(screen)
	// }
	s.DebugPrint(screen)
}

func (s SceneSelect) Viewport() ([]ChartInfo, int) {
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
	mode := s.Modes[*s.ModeType]
	speedBase := *mode.SpeedHandler.Target * 100
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Volume (Press 1/2): %.0f%%\n"+
			"SpeedBase (Press 3/4): %.0f\n"+"(Exposure time: %.0fms)\n\n:"+
			"Handler:%+v (Target: %d)\n",
			Volume*100,
			speedBase, mode.ExposureTime(speedBase),
			s.SelectHandler, *s.SelectHandler.Target))
}
