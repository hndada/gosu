package ebitenui

import (
	"github.com/hajimehoshi/ebiten"
	"image"
	"image/color"
)

const (
	checkboxBorder  = 2
	checkboxRadius  = 8
	checkboxPadding = 8
)

var (
	ImgUncheckedBox *ebiten.Image
	ImgCheckedBox   *ebiten.Image
)

func init() {
	bw := 2 * checkboxRadius
	bh := bw
	fw := bw - 2*checkboxBorder
	fh := fw
	w := bw + 2*checkboxPadding
	h := bh + checkboxPadding

	fop := &ebiten.DrawImageOptions{}
	fop.GeoM.Translate(checkboxBorder, checkboxBorder)
	bop := &ebiten.DrawImageOptions{}
	bop.GeoM.Translate(checkboxPadding, checkboxPadding/2)

	clrBorder := color.RGBA{102, 0, 95, 255}
	clrUncheck := color.Transparent
	clrCheck := color.RGBA{95, 240, 252, 255}
	{
		fill, _ := ebiten.NewImage(fw, fh, ebiten.FilterDefault)
		fill.Fill(clrUncheck)
		box, _ := ebiten.NewImage(bw, bh, ebiten.FilterDefault)
		box.Fill(clrBorder)
		box.DrawImage(fill, fop)
		ImgUncheckedBox, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
		ImgUncheckedBox.DrawImage(box, bop)
	}
	{
		fill, _ := ebiten.NewImage(fw, fh, ebiten.FilterDefault)
		fill.Fill(clrCheck)
		box, _ := ebiten.NewImage(bw, bh, ebiten.FilterDefault)
		box.Fill(clrBorder)
		box.DrawImage(fill, fop)
		ImgCheckedBox, _ = ebiten.NewImage(w, h, ebiten.FilterDefault)
		ImgCheckedBox.DrawImage(box, bop)
	}
}

// Checkbox is a struct for implementing a checkbox.
// Checkbox consists of padded checkbox image and text image.
// The font of Text is fixed, including text size.
type Checkbox struct {
	MinPt     image.Point
	Text      *ebiten.Image
	checked   bool
	mouseDown bool
	onChanged func(c *Checkbox)
}

func NewCheckbox(s string, p image.Point) *Checkbox {
	c := &Checkbox{}
	c.MinPt = p
	c.Text = RenderText(s, MplusNormalFont, color.Black)
	return c
}

// todo: could it be abstracted as button?
func (c *Checkbox) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p := image.Pt(ebiten.CursorPosition())
		if p.In(c.rect()) {
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

func (c *Checkbox) Draw(screen *ebiten.Image) {
	var box *ebiten.Image
	if c.checked {
		box = ImgCheckedBox
	} else {
		box = ImgUncheckedBox
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.MinPt.X), float64(c.MinPt.Y))
	screen.DrawImage(box, op)

	bx, _ := ImgUncheckedBox.Size()
	op.GeoM.Translate(float64(bx), 0)
	screen.DrawImage(c.Text, op)
}

func (c *Checkbox) Checked() bool {
	return c.checked
}

func (c *Checkbox) SetOnChanged(f func(c *Checkbox)) {
	c.onChanged = f
}

func (c *Checkbox) rect() image.Rectangle {
	bw, bh := ImgUncheckedBox.Size()
	tw, th := c.Text.Size()
	w := bw + tw
	h := th
	if h < bh {
		h = bh
	}
	maxPt := c.MinPt.Add(image.Pt(w, h))
	return image.Rectangle{c.MinPt, maxPt}
}

// it will not use widely
func f64Pt(p image.Point) (float64, float64) {
	return float64(p.X), float64(p.Y)
}
