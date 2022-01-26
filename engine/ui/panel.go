package ui

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hndada/gosu/engine/audio"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Panel struct {
	BodyText Sprite
	Body     LongSprite
	Left     Sprite
	Right    Sprite
}

type BoxSkin struct {
	Left   *ebiten.Image
	Middle *ebiten.Image
	Right  *ebiten.Image
}

// X and Y position values are updated on every Update()
func NewPanel(t string, skin BoxSkin) Panel {
	const PanelHeight = 40
	var p Panel
	p.Body.Sprite = NewSprite(skin.Middle)
	p.Body.Vertical = false
	p.Body.W = 450
	p.Body.H = PanelHeight
	{
		i := skin.Left
		sprite := NewSprite(i)
		sprite.H = PanelHeight
		scale := float64(sprite.H) / float64(i.Bounds().Dy())
		sprite.W = int(float64(i.Bounds().Dx()) * scale)
		p.Left = sprite
	}
	{
		i := skin.Right
		sprite := NewSprite(i)
		sprite.H = PanelHeight
		scale := float64(sprite.H) / float64(i.Bounds().Dy())
		sprite.W = int(float64(i.Bounds().Dx()) * scale)
		p.Right = sprite
	}
	{
		rect := image.Rect(0, 0, 450, 40)
		img := image.NewRGBA(rect)
		PanelTextcolor := color.Black
		x, y := 20, 30
		point := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(PanelTextcolor),
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		d.DrawString(t)

		i := ebiten.NewImageFromImage(img)
		sprite := NewSprite(i)
		sprite.W = 450
		sprite.H = 40
		p.BodyText = sprite
	}
	return p
}

func (p *Panel) SetXY(x, y int) {
	p.Body.X = x
	p.BodyText.X = x
	p.Left.X = x - p.Left.W
	p.Right.X = x + p.Body.W

	p.Body.Y = y
	p.BodyText.Y = y
	p.Left.Y = y
	p.Right.Y = y
}

func (p Panel) Draw(screen *ebiten.Image) {
	p.Left.Draw(screen)
	p.Body.Draw(screen)
	p.BodyText.Draw(screen)
	p.Right.Draw(screen)
}

type PanelHandler struct {
	panels    []Panel
	cursor    int
	holdCount int
	playSE    func()
	size      image.Point
}

func NewPanelHandler(screenSize image.Point, sePath string) PanelHandler {
	h := PanelHandler{}
	h.size = screenSize
	h.playSE = audio.NewSEPlayer(sePath, 25) // temp: volume
	return h
}

func (h *PanelHandler) Update() int {
	i := -1
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		i = h.cursor
		h.holdCount = 0
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if h.holdCount >= 2 { // TODO: actual duration should be consistent independent of maxTPS
			h.playSE()
			h.cursor++
			if h.cursor >= len(h.panels) {
				h.cursor = 0
			}
			h.holdCount = 0
		} else {
			h.holdCount++
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if h.holdCount >= 2 {
			h.playSE()
			h.cursor--
			if h.cursor < 0 {
				h.cursor = len(h.panels) - 1
			}
			h.holdCount = 0
		} else {
			h.holdCount++
		}
	} else {
		h.holdCount = 0
	}

	for i := range h.panels {
		mid := h.size.Y / 2 // A position of 'Currently selected chart' is fixed.
		x := h.size.X - 400
		d := i - h.cursor
		if d < 0 {
			d = -d
		}
		x += int(math.Pow(1.55, float64(d)))
		y := mid + 40*(i-h.cursor)
		if d == 0 {
			x -= 40
		}
		h.panels[i].SetXY(x, y)
	}
	return i
}

func (h *PanelHandler) Draw(screen *ebiten.Image) {
	for _, p := range h.panels {
		p.Draw(screen)
	}
}

func (h *PanelHandler) Append(p Panel) { h.panels = append(h.panels, p) }
