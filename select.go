package gosu

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hndada/gosu/audioutil"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/render"
	"golang.org/x/image/font/basicfont"
)

type SceneSelect struct {
	ChartInfoBoxSprite render.Sprite
	SoundMap           audioutil.SoundMap

	Mode int
	// Mods Mods
	// ReplayMode    bool
	View          []data.ChartInfo // Todo: should it be []*data.ChartInfo ?
	SelectHandler ctrl.IntHandler
	// ViewMode int

	Cursor      int // Todo: Delayed at Cursor
	Background  render.Sprite
	MusicPlayer *audio.Player // May rewind after preview has finished.
	MusicCloser io.Closer
}

const (
	BoxWidth  = 450
	BoxHeight = 50
	pop       = BoxWidth / 10
	count     = screenSizeY/BoxHeight/2*2 + 2 // Gives count some margin
)

// Todo: Score / Replay fetch
// Todo: preview music
func NewSceneSelect() *SceneSelect {
	s := new(SceneSelect)
	s.View = data.ChartInfos // View must be set before SelectHandler is set. // Todo: temp
	{
		b, err := audioutil.NewBytes("skin/default-hover.wav")
		if err != nil {
			fmt.Println(err)
		}
		play := audioutil.Context.NewPlayerFromBytes(b).Play
		s.SelectHandler = ctrl.IntHandler{
			Handler: ctrl.Handler{
				Keys:       []ebiten.Key{ebiten.KeyDown, ebiten.KeyUp},
				PlaySounds: []func(){play, play},
				HoldKey:    -1,
			},
			Min:    0,
			Max:    len(s.View) - 1,
			Unit:   1,
			Target: &s.Cursor,
		}
	}
	s.Mode = mode.ModePiano7 // Todo: temp

	purple := color.RGBA{172, 49, 174, 255}
	white := color.RGBA{255, 255, 255, 128}
	const border = 3
	{
		img := image.NewRGBA(image.Rect(0, 0, BoxWidth, BoxHeight))
		draw.Draw(img, img.Bounds(), &image.Uniform{purple}, image.Point{}, draw.Src)
		inRect := image.Rect(border, border, BoxWidth-border, BoxHeight-border)
		draw.Draw(img, inRect, &image.Uniform{white}, image.Point{}, draw.Src)
		s.ChartInfoBoxSprite = render.Sprite{
			I: ebiten.NewImageFromImage(img),
			W: BoxWidth,
			H: BoxHeight,
			X: screenSizeX - BoxWidth + pop,
			// Y is not fixed.
		}
	}
	s.SoundMap.Bytes = make(map[string][]byte)
	err := s.SoundMap.Register("skin/restart.wav", "select")
	if err != nil {
		fmt.Println(err)
	}
	s.UpdateBackground()
	return s
}
func (s *SceneSelect) UpdateBackground() {
	s.Background = mode.DefaultBackground
	if len(s.View) == 0 {
		return
	}
	if s.Cursor >= len(s.View) {
		return
	}
	info := s.View[s.Cursor]
	img := render.NewImage(info.Header.BackgroundPath(info.Path))
	if img != nil {
		s.Background.I = img
	}
	s.Background.SetWidth(screenSizeX)
	s.Background.SetCenterY(screenSizeY / 2)
}

// Default HoldKey value is 0, which is Key0.
func (s *SceneSelect) Update() any {
	moved := s.SelectHandler.Update()
	if moved {
		s.UpdateBackground()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyNumpadEnter) {
		s.SoundMap.Play("select")
		info := s.View[s.Cursor]
		return SelectToPlayArgs{
			Path:   info.Path,
			Mode:   s.Mode,
			Replay: nil,
			Play:   true,
		}
	}
	return nil
}

// Currently Chart infos are not in loop.
// May add extra effect to box arrangement.
// x -= y / 5, for example.
func (s SceneSelect) Draw(screen *ebiten.Image) {
	s.Background.Draw(screen)

	const (
		dx = 20
		dy = 30
	)
	var viewport []data.ChartInfo
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
	for i, info := range viewport {
		sprite := s.ChartInfoBoxSprite
		if i == cursor {
			sprite.X -= pop
		}
		offset := i - cursor
		sprite.SetCenterY(screenSizeY/2 + float64(offset)*BoxHeight)
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

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Volume (Press 1/2): %.0f%%\n"+
			"SpeedBase (Press 3/4): %.0f\n"+
			"(Exposure time: %.0fms)\n\n:"+"Handler:%+v (Target: %d)\n",
			mode.Volume*100, piano.SpeedBase*100, piano.ExposureTime(piano.SpeedBase), s.SelectHandler, *s.SelectHandler.Target))
}
