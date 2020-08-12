package graphics

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
)

const (
	checkboxRadius  = 8
	checkboxPadding = 8
)

// todo: 글자 크기가 안 맞는 체크박스가 나올 수 있음, 현재 수동으로 해야함
type Checkbox struct {
	Text      *ebiten.Image
	Rect      image.Rectangle
	checked   bool
	mouseDown bool
	onChanged func(c *Checkbox)
}

func NewCheckbox(text *ebiten.Image, x, y int) *Checkbox {
	c := &Checkbox{}
	tw, th := text.Size()
	w := 2*checkboxRadius + checkboxPadding + tw
	h := th
	if h < 2*checkboxRadius {
		h = 2 * checkboxRadius
	}
	c.Rect = image.Rect(x, y, x+w, y+h)
	c.Text = text
	return c
}

// todo: button과 한번에 합칠 수 있을까?
func (c *Checkbox) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if InRect(c.Rect, x, y) {
			c.mouseDown = true
		} else {
			c.mouseDown = false
		}
	} else {
		if c.mouseDown {
			c.checked = !c.checked
			if c.onChanged != nil {
				c.onChanged(c)
			}
		}
		c.mouseDown = false
	}
}

var (
	imgUncheckedBox *ebiten.Image
	imgCheckedBox   *ebiten.Image
)

func init() {
	w := 2*checkboxPadding + 2*checkboxRadius
	h := 2*checkboxRadius + checkboxPadding
	{
		fill, _ := ebiten.NewImage(2*checkboxRadius-2, 2*checkboxRadius-2, ebiten.FilterDefault)
		fill.Fill(color.Transparent)

		box, _ := ebiten.NewImage(2*checkboxRadius, 2*checkboxRadius, ebiten.FilterDefault)
		box.Fill(color.Black)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(1, 1)
		box.DrawImage(fill, op)

		imgUncheckedBox, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(checkboxPadding, checkboxPadding/2)
		imgUncheckedBox.DrawImage(box, op)
	}
	{
		fill, _ := ebiten.NewImage(2*checkboxRadius-2, 2*checkboxRadius-2, ebiten.FilterDefault)
		fill.Fill(color.RGBA{95, 240, 252, 255})

		box, _ := ebiten.NewImage(2*checkboxRadius, 2*checkboxRadius, ebiten.FilterDefault)
		box.Fill(color.Black)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(1, 1)
		box.DrawImage(fill, op)

		imgCheckedBox, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(checkboxPadding, checkboxPadding/2)
		imgCheckedBox.DrawImage(box, op)
	}
}

// 패딩 넣어 그려진 체크박스 그리고 텍스트 이미지 그리기
func (c *Checkbox) Draw(screen *ebiten.Image) {
	p := c.Rect.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(p.X), float64(p.Y))
	var box *ebiten.Image
	if c.checked {
		box = imgCheckedBox
	} else {
		box = imgUncheckedBox
	}
	screen.DrawImage(box, op)
	bx, by := imgUncheckedBox.Size()
	op.GeoM.Translate(float64(bx), float64(by))
	screen.DrawImage(c.Text, op)
}

func (c *Checkbox) Checked() bool {
	return c.checked
}

func (c *Checkbox) SetOnChanged(f func(c *Checkbox)) {
	c.onChanged = f
}

// func Rect(p image.Point, w, h int) image.Rectangle {
// 	return image.Rect(p.X, p.Y, p.X+w, p.Y+h)
// }
