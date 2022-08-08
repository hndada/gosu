package main

import (
	"fmt"
	"image"
	"image/color"
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

type SceneSelect struct {
	Charts     []*Chart // Todo: Should be ChartHeader instead of Chart
	ChartPaths []string
	Boxes      []*ebiten.Image
	Cursor     int
	Hold       int
	Pressed    bool
	Bg         *ebiten.Image
}

// Chart, Mods are fixed once enters to scene
// Speed should be mutable during playing
// Todo: load background
// Todo: play sound effect when moving a cursor
func NewSceneSelect() *SceneSelect {
	cs := make([]*Chart, 0, 50)
	cps := make([]string, 0, 50)
	bs := make([]*ebiten.Image, 0, 50)
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
					cs = append(cs, c)
					cps = append(cps, fpath)
					bs = append(bs, NewBox(c))
				}
			}
		}
	}
	return &SceneSelect{
		Charts:     cs,
		ChartPaths: cps,
		Boxes:      bs,
		Bg:         NewImage("skin/bg.jpg"),
	}
}

func NewBox(c *Chart) *ebiten.Image {
	const (
		w  = 350
		h  = 100
		bx = 20
		by = 30
	)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
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

// Todo: should hold goes reset when different key
func (s *SceneSelect) Update() {
	const threshold = 50 // Require holding for 50ms to move a cursor
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyEnter):
		NewScenePlay(s.Charts[s.Cursor], s.ChartPaths[s.Cursor]) // Todo: map to *Game g
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
			s.Cursor %= len(s.Charts)
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

// Currently Topmost and bottommost are not adjoined
func (s SceneSelect) Draw(screen *ebiten.Image) {
	const (
		w = 350
		h = 100
	)
	{
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(ScreenSizeX)/float64(s.Bg.Bounds().Dx()),
			float64(ScreenSizeY)/float64(s.Bg.Bounds().Dy()))
		screen.DrawImage(s.Bg, op)
	}
	for i := range s.Charts {
		y := (i-s.Cursor)*h + ScreenSizeY/2 - h/2
		if y > ScreenSizeY || y+h < 0 {
			continue
		}
		x := ScreenSizeX - w + 25
		if i == s.Cursor {
			x = ScreenSizeX - w
			// May add extra arrangement effect
			// Ex. x -= y / 5
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(s.Boxes[i], op)
	}
	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Speed (Press O/P): %.1f\nVolume (Press Q/W): %.0f",
			Speed*100, Volume*100))
}

// func (s *SceneSelect) IsHoldEnough(threshold int64) {}
