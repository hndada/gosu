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
	"github.com/hndada/gosu/parse/osu"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Todo: change Boxes type to Sprite
// Todo: load background of each chart.
type SceneSelect struct {
	Charts     []*Chart // Todo: Should be ChartHeader instead of Chart
	ChartPaths []string
	Boxes      []Sprite

	Cursor     int
	Hold       int
	Pressed    bool
	Background Sprite
}

const (
	bw = 300 // Box width
	bh = 40  // Box height
)

// Todo: play sound effect when moving a cursor
func NewSceneSelect() *SceneSelect {
	charts := make([]*Chart, 0, 50)
	paths := make([]string, 0, 50)
	boxes := make([]Sprite, 0, 50)
	dirs, err := os.ReadDir("music")
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		dpath := filepath.Join("music", dir.Name())
		files, err := os.ReadDir(dpath)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
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
					charts = append(charts, c)
					paths = append(paths, fpath)
					boxes = append(boxes, Sprite{
						I: NewBox(c),
						W: bw,
						H: bh,
					})
					// box's x value is not fixed.
					// box's y value is not fixed.
				}
			}
		}
	}
	bg := Sprite{
		I: NewImage("skin/bg.jpg"),
		W: screenSizeX,
		H: screenSizeY,
	}
	return &SceneSelect{
		Charts:     charts,
		ChartPaths: paths,
		Boxes:      boxes,
		Background: bg,
	}
}

func NewBox(c *Chart) *ebiten.Image {
	const (
		bx = 20
		by = 30
	)
	img := image.NewRGBA(image.Rect(0, 0, bw, bh))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	t := fmt.Sprintf("(%dKey Lv %.2f) %s [%s]", c.KeyCount, c.Level, c.MusicName, c.ChartName)
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: basicfont.Face7x13,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(bx * 64), Y: fixed.Int26_6(by * 64)},
	}
	d.DrawString(t)
	return ebiten.NewImageFromImage(img)
}

// Todo: should hold goes reset when different key.
// Todo: enable to pass replay format pointer to NewScenePlay.
// Todo: map to *Game g
func (s *SceneSelect) Update(g *Game) {
	const threshold = 80 // Require holding for 80ms to move a cursor
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyEnter), ebiten.IsKeyPressed(ebiten.KeyNumpadEnter):
		g.Scene = NewScenePlay(s.Charts[s.Cursor], s.ChartPaths[s.Cursor], nil)
	case ebiten.IsKeyPressed(ebiten.KeyArrowDown):
		if TickToMsec(s.Hold) > threshold {
			s.Hold = 0
			s.Cursor++
			s.Cursor %= len(s.Charts)
		} else {
			s.Hold++
		}
	case ebiten.IsKeyPressed(ebiten.KeyArrowUp):
		if TickToMsec(s.Hold) > threshold {
			s.Hold = 0
			s.Cursor--
			if s.Cursor < 0 {
				s.Cursor += len(s.Charts)
			}
		} else {
			s.Hold++
		}
	default:
		s.Hold = 0
		if !ebiten.IsKeyPressed(ebiten.KeyO) && !ebiten.IsKeyPressed(ebiten.KeyP) ||
			!ebiten.IsKeyPressed(ebiten.KeyQ) && !ebiten.IsKeyPressed(ebiten.KeyW) {
			s.Pressed = false
			break
		}
		if s.Pressed {
			break
		}
		s.Pressed = true
		switch {
		case ebiten.IsKeyPressed(ebiten.KeyO):
			Speed -= 0.01
			if Speed < 0.01 {
				Speed = 0.01
			}
		case ebiten.IsKeyPressed(ebiten.KeyP):
			Speed += 0.01
			if Speed > 0.4 {
				Speed = 0.4
			}
		case ebiten.IsKeyPressed(ebiten.KeyQ):
			Volume -= 0.05
			if Volume < 0 {
				Volume = 0
			}
		case ebiten.IsKeyPressed(ebiten.KeyW):
			Volume += 0.05
			if Volume > 1 {
				Volume = 1
			}
		}
	}
}

// Currently topmost and bottommost boxes are not adjoined.
// May add extra effect to box arrangement,
// x -= y / 5, for example.
func (s SceneSelect) Draw(screen *ebiten.Image) {
	const pop = 25
	s.Background.Draw(screen)
	for i := range s.Charts {
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
		screen.DrawImage(s.Boxes[i].I, op)
	}
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Speed (Press O/P): %.1f\nVolume (Press Q/W): %.0f",
			Speed*100, Volume*100))
}

// func (s *SceneSelect) IsHoldEnough(threshold int64) {}
