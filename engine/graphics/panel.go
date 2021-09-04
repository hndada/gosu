package graphics

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/hndada/gosu/game"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

type Panel struct {
	BodyText game.Sprite
	Body     game.LongSprite
	Left     game.Sprite
	Right    game.Sprite
}

// X와 Y는 update에서 매번 새로 설정
func NewPanel(t string) Panel {
	const PanelHeight = 40
	var p Panel
	p.Body.Sprite = game.NewSprite(game.Skin.BoxMiddle)
	p.Body.Vertical = false
	p.Body.W = 450
	p.Body.H = PanelHeight
	{
		src := game.Skin.BoxLeft
		sprite := game.NewSprite(src)
		sprite.H = PanelHeight
		scale := float64(sprite.H) / float64(src.Bounds().Dy())
		sprite.W = int(float64(src.Bounds().Dx()) * scale)
		p.Left = sprite
	}
	{
		src := game.Skin.BoxRight
		sprite := game.NewSprite(src)
		sprite.H = PanelHeight
		scale := float64(sprite.H) / float64(src.Bounds().Dy())
		sprite.W = int(float64(src.Bounds().Dx()) * scale)
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

		src, _ := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		sprite := game.NewSprite(src)
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
}

// todo: SEPlayer
func NewPanelHandler() PanelHandler {
	h := PanelHandler{}
	// h.playSE = mania.SEPlayer(cwd)
	return h
}

func (h *PanelHandler) Update() int {
	i := -1
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		i = h.cursor
		h.holdCount = 0
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if h.holdCount >= 2 { // todo: MaxTPS가 변하여도 체감 시간은 그대로이게 설정
			// h.playSE()
			h.cursor++
			if h.cursor >= len(h.panels) {
				h.cursor = 0
			}
			h.holdCount = 0
		} else {
			h.holdCount++
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if h.holdCount >= 2 {
			// h.playSE()
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
		mid := game.Settings.ScreenSize.Y / 2 // 현재 선택된 차트 focus 틀 위치 고정
		x := game.Settings.ScreenSize.X - 400
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
