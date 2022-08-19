package gosu

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hndada/gosu/audioutil"
	"github.com/hndada/gosu/ctrl"
	"github.com/hndada/gosu/db"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/render"
	"golang.org/x/image/font/basicfont"
)

// var ChartInfoSprites []render.Sprite

// 1. Load data from local db (It may be skipped since no local db)
// 2. Find new music, then add to SceneSelect (and also to local db)
// 3. NewSelectScene
// 4. NewPlayScene, based on mode. Args: path, mods, replay, play
// Todo: PlaySoundMove. map[string]func()
func (s *SceneSelect) HandleMove() {}

type SceneSelect struct {
	SelectHandler ctrl.IntHandler
	// Todo: Delayed at Cursor
	Mode int
	// ViewMode int
	View []db.ChartInfo // Todo: should it be []*db.ChartInfo ?

	// ChartBoxs []db.ChartBox
	Cursor     int
	Background render.Sprite
	// Hold       int
	// HoldKey ebiten.Key

	// ReplayMode    bool
	// IndexToMD5Map map[int][md5.Size]byte
	// MD5ToIndexMap map[[md5.Size]byte]int
	// Replays       []*osr.Format

	MusicPlayer *audio.Player // May rewind after preview has finished.
	MusicCloser io.Closer
	SoundMap    audioutil.SoundMap
	// SoundBytes   map[string][]byte
	// SoundClosers []io.Closer
	// PlaySoundMove   func()
	// PlaySoundSelect func()
	ChartInfoBoxSprite render.Sprite
}

const (
	BoxWidth  = 450 // Box width
	BoxHeight = 50  // Box height
)
const count = 20
const pop = BoxWidth / 10

// Todo: Score / Replay fetch
func NewSceneSelect() *SceneSelect {
	s := new(SceneSelect)
	s.UpdateBackground()
	{
		b, err := audioutil.NewBytes("skin/default-hover.wav")
		if err != nil {
			fmt.Println(err)
		}
		play := audioutil.Context.NewPlayerFromBytes(b).Play
		s.SelectHandler = ctrl.IntHandler{
			Handler: ctrl.Handler{
				Keys:       []ebiten.Key{ebiten.KeyUp, ebiten.KeyDown},
				PlaySounds: []func(){play, play},
				HoldKey:    -1,
			},
			Min:    0,
			Max:    len(s.View),
			Unit:   1,
			Target: &s.Cursor,
		}
	}
	s.Mode = mode.ModePiano7 // Todo: temp
	s.View = db.ChartInfos   // Todo: temp
	// var err error
	// err = s.SoundMap.Register("skin/default-hover.wav", "move")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	s.SoundMap.Bytes = make(map[string][]byte)
	err := s.SoundMap.Register("skin/restart.wav", "select")
	if err != nil {
		fmt.Println(err)
	}
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
	return s
}
func (s *SceneSelect) UpdateBackground() {
	s.Background = mode.DefaultBackground
	if len(s.View) == 0 {
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

// const HoldKeyNone = -1

// Require holding for a while to move a cursor
// var (
// threshold1 = mode.TimeToTick(100)
// threshold2 = mode.TimeToTick(80)
// )

// // FetchReplay returns first MD5-matching replay format.
// func (s SceneSelect) FetchReplay(i int) *osr.Format {
// 	md5 := s.IndexToMD5Map[i]
// 	for _, r := range s.Replays {
// 		if md5 == r.MD5() {
// 			return r
// 		}
// 	}
// 	return nil
// }

// Default HoldKey value is 0, which is Key0.
func (s *SceneSelect) Update() any {
	if moved := s.SelectHandler.Update(); moved {
		s.UpdateBackground()
	}
	// if s.HoldKey == HoldKeyNone {
	// 	s.Hold++
	// 	if s.Hold > threshold1 {
	// 		s.Hold = threshold1
	// 	}
	// } else {
	// 	if ebiten.IsKeyPressed(s.HoldKey) {
	// 		s.Hold++
	// 	} else {
	// 		s.Hold = 0
	// 	}
	// }
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
	// switch {
	// case
	// switch s.Mode {
	// case db.ModePiano4, db.ModePiano7, db.ModePiano8:

	// 	g.Scene = piano.NewScenePlay(info.Path, nil, true)
	// }
	// if s.ReplayMode {
	// 	return PlayChartArgs{
	// 		Mode:   mode.ModePiano,
	// 		Path:   info.Path,
	// 		Replay: s.FetchReplay(s.Cursor),
	// 		Play:   true,
	// 	}
	// 	// g.Scene = NewScenePlay(info.Chart, info.Path, s.FetchReplay(s.Cursor), true)
	// } else {
	// 	return PlayChartArgs{
	// 		Mode:   mode.ModePiano,
	// 		Path:   info.Path,
	// 		Replay: nil,
	// 		Play:   true,
	// 	}
	// 	// g.Scene = NewScenePlay(info.Chart, info.Path, nil, true)
	// }
	// case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
	// 	s.HoldKey = ebiten.KeyArrowDown
	// 	if s.Hold < threshold1 {
	// 		break
	// 	}
	// 	s.SoundMap.Play("move")
	// 	s.Hold = 0
	// 	s.Cursor++
	// 	s.Cursor %= len(s.View)
	// 	s.UpdateBackground()
	// case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
	// 	s.HoldKey = ebiten.KeyArrowUp
	// 	if s.Hold < threshold1 {
	// 		break
	// 	}
	// 	s.SoundMap.Play("move")
	// 	s.Hold = 0
	// 	s.Cursor--
	// 	if s.Cursor < 0 {
	// 		s.Cursor += len(s.View)
	// 	}
	// 	s.UpdateBackground()
	// case ebiten.IsKeyPressed(ebiten.KeyQ):
	// 	s.HoldKey = ebiten.KeyQ
	// 	if s.Hold < threshold2 {
	// 		break
	// 	}
	// 	s.Hold = 0
	// 	BaseSpeed -= 0.1
	// 	if BaseSpeed < 0.1 {
	// 		BaseSpeed = 0.1
	// 	}
	// case ebiten.IsKeyPressed(ebiten.KeyW):
	// 	s.HoldKey = ebiten.KeyW
	// 	if s.Hold < threshold2 {
	// 		break
	// 	}
	// 	s.Hold = 0
	// 	BaseSpeed += 0.1
	// 	if BaseSpeed > 2 {
	// 		BaseSpeed = 2
	// 	}
	// case ebiten.IsKeyPressed(ebiten.KeyA):
	// 	s.HoldKey = ebiten.KeyA
	// 	if s.Hold < threshold2 {
	// 		break
	// 	}
	// 	s.Hold = 0
	// 	Volume -= 0.05
	// 	if Volume < 0 {
	// 		Volume = 0
	// 	}
	// case ebiten.IsKeyPressed(ebiten.KeyS):
	// 	s.HoldKey = ebiten.KeyS
	// 	if s.Hold < threshold2 {
	// 		break
	// 	}
	// 	s.Hold = 0
	// 	Volume += 0.05
	// 	if Volume > 1 {
	// 		Volume = 1
	// 	}
	// case ebiten.IsKeyPressed(ebiten.KeyZ):
	// 	s.HoldKey = ebiten.KeyZ
	// 	if s.Hold < threshold1 {
	// 		break
	// 	}
	// 	s.Hold = 0
	// 	s.ReplayMode = !s.ReplayMode
	// default:
	// 	s.HoldKey = HoldKeyNone
	// }
}

// Currently topmost and bottommost boxes are not adjoined.
// May add extra effect to box arrangement.
// x -= y / 5, for example.
func (s SceneSelect) Draw(screen *ebiten.Image) {
	s.Background.Draw(screen)
	const (
		dx = 20
		dy = 30
	)

	var viewport []db.ChartInfo
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
	// fmt.Println(len(viewport))
	for i, info := range viewport {
		sprite := s.ChartInfoBoxSprite
		if i == cursor {
			sprite.X -= pop
		}
		offset := i - cursor
		sprite.SetCenterY(screenSizeY/2 + float64(offset)*BoxHeight)
		sprite.Draw(screen)
		// sprite := box.Box

		// if ChartInfoSprites[s.Cursor+offset].I == nil {
		// 	// info := db.ChartInfos[s.Cursor+offset]
		// 	ChartInfoSprites[s.Cursor+offset] = db.NewChartInfoSprite(info)
		// }
		// sprite := ChartInfoSprites[s.Cursor+offset]
		t := info.Text()
		// rect := text.BoundString(basicfont.Face7x13, t)
		x := int(sprite.X) + dx //+ rect.Dx()
		y := int(sprite.Y) + dy //+ rect.Dy()
		text.Draw(screen, t, basicfont.Face7x13, x, y, color.Black)

		// y := (i-s.Cursor)*bh + screenSizeY/2 - bh/2
		// if y > screenSizeY || y+bh < 0 {
		// 	continue
		// }
		// x := screenSizeX - bw + pop
		// if i == s.Cursor {
		// 	x = screenSizeX - bw
		// }
		// op := &ebiten.DrawImageOptions{}
		// op.GeoM.Translate(float64(x), float64(y))
		// screen.DrawImage(s.View[i].Box.I, op)
	}
	// Code of drawing cursor
	// {
	// 	sprite := GeneralSkin.CursorSprites[0]
	// 	x, y := ebiten.CursorPosition()
	// 	sprite.X, sprite.Y = float64(x), float64(y)
	// 	sprite.Draw(screen)
	// }

	// Todo: BaseSpeed
	// ebitenutil.DebugPrint(screen,
	// 	fmt.Sprintf("BaseSpeed (Press Q/W): %.0f\n(Exposure time: %.0fms)\n\nVolume (Press A/S): %d%%\nHold:%d\nReplay mode (Press Z): %v\n", // %.1f
	// 		BaseSpeed*100, ExposureTime(BaseSpeed), int(mode.Volume*100), s.Hold, s.ReplayMode))
}

// Box: render.Sprite{
// 	I: NewBox(c, mode.Level(c.Difficulties())),
// 	W: bw,
// 	H: bh,
// },
// box's x value is not fixed.
// box's y value is not fixed.
