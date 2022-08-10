package gosu

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hndada/gosu/parse/osr"
	"github.com/hndada/gosu/parse/osu"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Todo: should be ChartHeader instead of Chart
// Todo: might integrate database here
type SceneSelect struct {
	Charts      []*Chart
	ChartPaths  []string
	ChartLevels []float64
	ChartBoxes  []Sprite
	ChartCursor int
	Background  Sprite

	Replays      []*osr.Format
	ReplayBoxes  []Sprite
	ReplayCursor int

	Hold    int
	HoldKey ebiten.Key

	PlaySoundMove   func()
	PlaySoundSelect func()
}

const (
	bw = 300 // Box width
	bh = 40  // Box height
)

func NewSceneSelect() *SceneSelect {
	s := &SceneSelect{
		Charts:      make([]*Chart, 0, 50),
		ChartPaths:  make([]string, 0, 50),
		ChartLevels: make([]float64, 0, 50),
		ChartBoxes:  make([]Sprite, 0, 50),
	}
	dirs, err := os.ReadDir("music")
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dpath := filepath.Join("music", dir.Name())
		fs, err := os.ReadDir(dpath)
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			fpath := filepath.Join(dpath, f.Name())
			b, err := os.ReadFile(fpath)
			if err != nil {
				panic(err)
			}
			switch strings.ToLower(filepath.Ext(fpath)) {
			case ".osu":
				o, err := osu.Parse(b)
				if err != nil {
					panic(err)
				}
				if o.Mode == ModeMania {
					c, err := NewChartFromOsu(o)
					if err != nil {
						panic(err)
					}
					s.Charts = append(s.Charts, c)
					s.ChartPaths = append(s.ChartPaths, fpath)
					s.ChartLevels = append(s.ChartLevels, c.Level())
					s.ChartBoxes = append(s.ChartBoxes, Sprite{
						I: NewBox(c, c.Level()),
						W: bw,
						H: bh,
					})
					// box's x value is not fixed.
					// box's y value is not fixed.
				}
			}
		}
	}
	s.UpdateBackground()
	_, apMove := NewAudioPlayer("skin/default-hover.wav")
	s.PlaySoundMove = apMove.PlaySoundEffect
	_, apSelect := NewAudioPlayer("skin/restart.wav")
	s.PlaySoundSelect = apSelect.PlaySoundEffect
	s.HoldKey = HoldKeyNone
	return s
}
func (s *SceneSelect) UpdateBackground() {
	s.Background = RandomDefaultBackground
	if len(s.Charts) == 0 {
		return
	}
	img := NewImage(s.Charts[s.ChartCursor].BackgroundPath(s.ChartPaths[s.ChartCursor]))
	if img != nil {
		s.Background.I = img
	}
}

const (
	border = 3
)
const (
	dx = 20 // dot x
	dy = 30 // dot y
)

var borderColor = color.RGBA{172, 49, 174, 255} // Purple

func NewBox(c *Chart, lv float64) *ebiten.Image {
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
	threshold1 = MsecToTick(150)
	threshold2 = MsecToTick(100)
)

// Todo: enable to pass replay format pointer to NewScenePlay
// Default HoldKey value is 0, which is Key0.
func (s *SceneSelect) Update(g *Game) {
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
		s.PlaySoundSelect()
		g.Scene = NewScenePlay(s.Charts[s.ChartCursor], s.ChartPaths[s.ChartCursor], nil, true)
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		s.HoldKey = ebiten.KeyArrowDown
		if s.Hold < threshold1 {
			break
		}
		s.PlaySoundMove()
		s.Hold = 0
		s.ChartCursor++
		s.ChartCursor %= len(s.Charts)
		s.UpdateBackground()
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		s.HoldKey = ebiten.KeyArrowUp
		if s.Hold < threshold1 {
			break
		}
		s.PlaySoundMove()
		s.Hold = 0
		s.ChartCursor--
		if s.ChartCursor < 0 {
			s.ChartCursor += len(s.Charts)
		}
		s.UpdateBackground()
	case ebiten.IsKeyPressed(ebiten.KeyQ):
		s.HoldKey = ebiten.KeyQ
		if s.Hold < threshold2 {
			break
		}
		s.Hold = 0
		Speed -= 0.05
		if Speed < 0.05 {
			Speed = 0.05
		}
	case ebiten.IsKeyPressed(ebiten.KeyW):
		s.HoldKey = ebiten.KeyW
		if s.Hold < threshold2 {
			break
		}
		s.Hold = 0
		Speed += 0.05
		if Speed > 1.5 {
			Speed = 1.5
		}
	case ebiten.IsKeyPressed(ebiten.KeyA):
		s.HoldKey = ebiten.KeyA
		if s.Hold < threshold2 {
			break
		}
		s.Hold = 0
		Volume -= 0.05
		if Volume < 0 {
			Volume = 0
		}
	case ebiten.IsKeyPressed(ebiten.KeyS):
		s.HoldKey = ebiten.KeyS
		if s.Hold < threshold2 {
			break
		}
		s.Hold = 0
		Volume += 0.05
		if Volume > 1 {
			Volume = 1
		}
	default:
		s.HoldKey = HoldKeyNone
	}
}

// Currently topmost and bottommost boxes are not adjoined.
// May add extra effect to box arrangement.
// x -= y / 5, for example.
func (s SceneSelect) Draw(screen *ebiten.Image) {
	const pop = 25
	s.Background.Draw(screen)
	for i := range s.Charts {
		y := (i-s.ChartCursor)*bh + screenSizeY/2 - bh/2
		if y > screenSizeY || y+bh < 0 {
			continue
		}
		x := screenSizeX - bw + pop
		if i == s.ChartCursor {
			x = screenSizeX - bw
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(s.ChartBoxes[i].I, op)
	}
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Speed (Press Q/W): %d\nVolume (Press A/S): %d%%\nHold:%d\n", // %.1f
			int(Speed*20), int(Volume*100), s.Hold))
}
