package gosu

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hndada/gosu/audioutil"
	"github.com/hndada/gosu/db"
	"github.com/hndada/gosu/mode"
	"github.com/hndada/gosu/mode/piano"
	"github.com/hndada/gosu/render"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// 1. Load data from local db (It may be skipped since no local db)
// 2. Find new music, then add to SceneSelect (and also to local db)
// 3. NewSelectScene
// 4. NewPlayScene, based on mode. Args: path, mods, replay, play
// Todo: PlaySoundMove. map[string]func()
func (s *SceneSelect) HandleMove() {}

type SceneSelect struct {
	ViewMode int
	View     []db.ChartBox
	// ChartBoxs []db.ChartBox
	Cursor     int
	Background render.Sprite
	Hold       int
	HoldKey    ebiten.Key

	// ReplayMode    bool
	// IndexToMD5Map map[int][md5.Size]byte
	// MD5ToIndexMap map[[md5.Size]byte]int
	// Replays       []*osr.Format

	MusicPlayer *audio.Player // May rewind after preview has finished.
	MusicCloser io.Closer
	SoundPad    audioutil.SoundPad
	// SoundBytes   map[string][]byte
	// SoundClosers []io.Closer
	// PlaySoundMove   func()
	// PlaySoundSelect func()

	Mode int
}

const (
	bw  = 450 // Box width
	bh  = 50  // Box height
	pop = bw / 10
)

func NewSceneSelect() *SceneSelect {
	// s := &SceneSelect{
	// 	View:          make([]db.ChartBox, 0, 50),
	// 	IndexToMD5Map: make(map[int][16]byte),
	// 	MD5ToIndexMap: map[[16]byte]int{},
	// 	Replays:       make([]*osr.Format, 0, 10),
	// }
	s := new(SceneSelect)
	s.UpdateBackground()
	s.HoldKey = HoldKeyNone
	s.Hold = threshold1
	_ = s.SoundPad.Register("skin/default-hover.wav", "move")
	_ = s.SoundPad.Register("skin/restart.wav", "select")
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

const (
	border = 3
)
const (
	dx = 20 // dot x
	dy = 30 // dot y
)

var borderColor = color.RGBA{172, 49, 174, 255} // Purple

func NewBox(c *piano.Chart, lv float64) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, bw, bh))
	draw.Draw(img, img.Bounds(), &image.Uniform{borderColor}, image.Point{}, draw.Src)
	inRect := image.Rect(border, border, bw-border, bh-border)
	draw.Draw(img, inRect, &image.Uniform{color.White}, image.Point{}, draw.Src)
	t := fmt.Sprintf("(%dK Lv %.1f) %s [%s]", c.KeyCount, lv, c.MusicName, c.ChartName)
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(dx * 64), Y: fixed.Int26_6(dy * 64)},
	}
	d.DrawString(t)
	return ebiten.NewImageFromImage(img)
}

const HoldKeyNone = -1

// Require holding for a while to move a cursor
var (
	threshold1 = mode.TimeToTick(100)
	// threshold2 = mode.TimeToTick(80)
)

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
	if s.HoldKey == HoldKeyNone {
		s.Hold++
		if s.Hold > threshold1 {
			s.Hold = threshold1
		}
	} else {
		if ebiten.IsKeyPressed(s.HoldKey) {
			s.Hold++
		} else {
			s.Hold = 0
		}
	}
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyEnter), ebiten.IsKeyPressed(ebiten.KeyNumpadEnter):
		s.SoundPad.Play("select")
		info := s.View[s.Cursor]
		return SelectToPlayArgs{
			Path:   info.Path,
			Mode:   s.Mode,
			Replay: nil,
			Play:   true,
		}
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
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		s.HoldKey = ebiten.KeyArrowDown
		if s.Hold < threshold1 {
			break
		}
		s.SoundPad.Play("move")
		s.Hold = 0
		s.Cursor++
		s.Cursor %= len(s.View)
		s.UpdateBackground()
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		s.HoldKey = ebiten.KeyArrowUp
		if s.Hold < threshold1 {
			break
		}
		s.SoundPad.Play("move")
		s.Hold = 0
		s.Cursor--
		if s.Cursor < 0 {
			s.Cursor += len(s.View)
		}
		s.UpdateBackground()
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
	default:
		s.HoldKey = HoldKeyNone
	}
	return nil
}

// Currently topmost and bottommost boxes are not adjoined.
// May add extra effect to box arrangement.
// x -= y / 5, for example.
func (s SceneSelect) Draw(screen *ebiten.Image) {
	s.Background.Draw(screen)
	for i := range s.View {
		y := (i-s.Cursor)*bh + screenSizeY/2 - bh/2
		if y > screenSizeY || y+bh < 0 {
			continue
		}
		x := screenSizeX - bw + pop
		if i == s.Cursor {
			x = screenSizeX - bw
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
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
